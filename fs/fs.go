package fs

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/AspieSoft/go-regex-re2/v2"
	"github.com/alphadose/haxmap"
	"github.com/fsnotify/fsnotify"
)

// JoinPath joins multiple file types with safety from backtracking
func JoinPath(path ...string) (string, error) {
	resPath, err := filepath.Abs(string(path[0]))
	if err != nil {
		return "", err
	}
	for i := 1; i < len(path); i++ {
		p := filepath.Join(resPath, string(path[i]))
		if p == resPath || !strings.HasPrefix(p, resPath) {
			return "", errors.New("path leaked outside of root")
		}
		resPath = p
	}
	return resPath, nil
}

// Copy lets you copy files from the src to the dst
func Copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}


var regDirEndSlash *regex.Regexp = regex.Comp(`[\\/][^\\/]*$`)

// GetFileFromParent checks if the parent (or sub parent) directory of a file contains a specific file or folder
//
// @root is the highest grandparent to check before quitting
//
// @start is the lowest level to start searching from (if a directory is passed, it will not be included in your search)
//
// @search is what file you want to search from
func GetFileFromParent(root string, start string, search string) (string, bool) {
	dir := string(regDirEndSlash.RepStrLit([]byte(start), []byte{}))
	if len(dir) == 0 || dir == root || !strings.HasPrefix(dir, root) {
		return "", false
	}

	if dirList, err := os.ReadDir(string(dir)); err == nil {
		for _, file := range dirList {
			name := file.Name()
			if name == search {
				if path, err := JoinPath(string(dir), name); err == nil {
					return path, true
				}
				return "", false
			}
		}
	}

	return GetFileFromParent(root, dir, search)
}


// A watcher instance for the `FS.FileWatcher` method
type FileWatcher struct {
	watcherList *haxmap.Map[string, *watcherObj]

	// when a file changes
	//
	// @path: the file path the change happened to
	//
	// @op: the change operation
	OnFileChange func(path string, op string)

	// when a directory is added
	//
	// @path: the file path the change happened to
	//
	// @op: the change operation
	//
	// return false to prevent that directory from being watched
	OnDirAdd func(path string, op string) (addWatcher bool)

	// when a file or directory is removed
	//
	// @path: the file path the change happened to
	//
	// @op: the change operation
	//
	// return false to prevent that directory from no longer being watched
	OnRemove func(path string, op string) (removeWatcher bool)

	// every time something happenes
	//
	// @path: the file path the change happened to
	//
	// @op: the change operation
	OnAny func(path string, op string)
}

type watcherObj struct {
	watcher *fsnotify.Watcher
	close *bool
}

func Watcher() *FileWatcher {
	return &FileWatcher{watcherList: haxmap.New[string, *watcherObj]()}
}

// WatchDir watches the files in a directory and its subdirectories for changes
func (fw *FileWatcher) WatchDir(root string) error {
	var err error
	if root, err = filepath.Abs(root); err != nil {
		return err
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	runClose := false

	fw.watcherList.Set(root, &watcherObj{watcher: watcher, close: &runClose})

	go func() {
		defer watcher.Close()
		for {
			if runClose {
				break
			}

			if event, ok := <-watcher.Events; ok {
				filePath := event.Name

				stat, err := os.Stat(filePath)
				if err != nil {
					if fw.OnRemove == nil || fw.OnRemove(filePath, event.Op.String()){
						watcher.Remove(filePath)
					}
				}else if stat.IsDir() {
					if fw.OnDirAdd == nil || fw.OnDirAdd(filePath, event.Op.String()){
						watcher.Add(filePath)
					}
				}else{
					if fw.OnFileChange != nil {
						fw.OnFileChange(filePath, event.Op.String())
					}
				}

				if fw.OnAny != nil {
					fw.OnAny(filePath, event.Op.String())
				}
			}
		}
	}()

	err = watcher.Add(root)
	if err != nil {
		return err
	}

	fw.watchDirSub(watcher, root)

	return nil
}

func (fw *FileWatcher) watchDirSub(watcher *fsnotify.Watcher, dir string){
	files, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	for _, file := range files {
		if file.IsDir() {
			if path, err := JoinPath(dir, file.Name()); err == nil {
				watcher.Add(path)
				fw.watchDirSub(watcher, path)
			}
		}
	}
}

// CloseWatcher will close the watcher by the root name you used 
//
// @root pass a file path for a specific watcher or "*" for all watchers that exist
func (fw *FileWatcher) CloseWatcher(root string) error {
	if root == "" || root == "*" {
		rList := []string{}
		fw.watcherList.ForEach(func(r string, w *watcherObj) bool {
			rList = append(rList, r)
			*w.close = true
			return true
		})
		fw.watcherList.Del(rList...)
	}else{
		var err error
		if root, err = filepath.Abs(root); err != nil {
			return err
		}

		if w, ok := fw.watcherList.Get(root); ok {
			*w.close = true
			fw.watcherList.Del(root)
		}
	}

	return nil
}

// Wait for all Watchers to close
func (fw *FileWatcher) Wait(){
	for fw.watcherList.Len() != 0 {
		time.Sleep(100 * time.Nanosecond)
	}
}

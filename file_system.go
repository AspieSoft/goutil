package goutil

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/AspieSoft/go-regex/v4"
	"github.com/alphadose/haxmap"
	"github.com/fsnotify/fsnotify"
)

type fileSystem struct {}
var FS *fileSystem = &fileSystem{}

// linux package installer
type LinuxPKG struct {
	sudo bool
}

// JoinPath joins multiple file types with safety from backtracking
func (fs *fileSystem) JoinPath(path ...string) (string, error) {
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

var regDirEndSlash *regex.Regexp = regex.Compile(`[\\/][^\\/]*$`)

// GetFileFromParent checks if the parent (or sub parent) directory of a file contains a specific file or folder
//
// @root is the highest grandparent to check before quitting
//
// @start is the lowest level to start searching from (if a directory is passed, it will not be included in your search)
//
// @search is what file you want to search from
func (fs *fileSystem) GetFileFromParent(root string, start string, search string) (string, bool) {
	dir := string(regDirEndSlash.RepStr([]byte(start), []byte{}))
	if len(dir) == 0 || dir == root || !strings.HasPrefix(dir, root) {
		return "", false
	}

	if dirList, err := os.ReadDir(string(dir)); err == nil {
		for _, file := range dirList {
			name := file.Name()
			if name == search {
				if path, err := fs.JoinPath(string(dir), name); err == nil {
					return path, true
				}
				return "", false
			}
		}
	}

	return fs.GetFileFromParent(root, dir, search)
}


// A watcher instance for the `FS.FileWatcher` method
type FileWatcher struct {
	fs *fileSystem
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

func (fs *fileSystem) FileWatcher() *FileWatcher {
	return &FileWatcher{fs: fs, watcherList: haxmap.New[string, *watcherObj]()}
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
			if path, err := fw.fs.JoinPath(dir, file.Name()); err == nil {
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


// InstallLinuxPkg attempts to install a linux package
//
// this method will also resolve the sudo command and ask for a user password if needed
//
// this method will not attempt to run an install, if it finds the package is already installed
func (linuxPKG *LinuxPKG) InstallLinuxPkg(pkg []string, man ...string){
	if !linuxPKG.HasLinuxPkg(pkg) {
		var pkgMan string
		if len(man) != 0 {
			pkgMan = man[0]
		}else{
			pkgMan = linuxPKG.GetLinuxInstaller([]string{`apt-get`, `apt`, `yum`})
		}

		var cmd *exec.Cmd
		if linuxPKG.sudo {
			cmd = exec.Command(`sudo`, append([]string{pkgMan, `install`, `-y`}, pkg...)...)
		}else{
			cmd = exec.Command(pkgMan, append([]string{`install`, `-y`}, pkg...)...)
		}

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return
		}

		go (func() {
			out := bufio.NewReader(stdout)
			for {
				s, err := out.ReadString('\n')
				if err == nil {
					fmt.Println(s)
				}
			}
		})()

		stderr, err := cmd.StderrPipe()
		if err != nil {
			return
		}

		go (func() {
			out := bufio.NewReader(stderr)
			for {
				s, err := out.ReadString('\n')
				if err == nil {
					fmt.Println(s)
				}
			}
		})()

		cmd.Run()
	}
}

// HasLinuxPkg attempt to check if a linux package is installed
func (linuxPKG *LinuxPKG) HasLinuxPkg(pkg []string) bool {
	for _, name := range pkg {
		hasPackage := false

		var cmd *exec.Cmd
		if linuxPKG.sudo {
			cmd = exec.Command(`sudo`, `dpkg`, `-s`, name)
		}else{
			cmd = exec.Command(`dpkg`, `-s`, name)
		}

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return true
		}
		go (func() {
			out := bufio.NewReader(stdout)
			for {
				_, err := out.ReadString('\n')
				if err == nil {
					hasPackage = true
				}
			}
		})()
		for i := 0; i < 3; i++ {
			cmd.Run()
			if hasPackage {
				break
			}
		}
		if !hasPackage {
			return false
		}
	}

	return true
}

// GetLinuxInstaller attempt to find out what package manager a linux distro is using or has available
func (linuxPKG *LinuxPKG) GetLinuxInstaller(man []string) string {
	hasInstaller := ""

	for _, m := range man {

		var cmd *exec.Cmd
		if linuxPKG.sudo {
			cmd = exec.Command(`sudo`, `dpkg`, `-s`, m)
		}else{
			cmd = exec.Command(`dpkg`, `-s`, m)
		}

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			continue
		}
		go (func() {
			out := bufio.NewReader(stdout)
			for {
				_, err := out.Peek(1)
				if err == nil {
					hasInstaller = m
				}
			}
		})()

		for i := 0; i < 3; i++ {
			cmd.Run()
			if hasInstaller != "" {
				break
			}
		}

		if hasInstaller != "" {
			break
		}
	}

	return hasInstaller
}

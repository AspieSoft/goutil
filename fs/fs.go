package fs

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/AspieSoft/go-regex-re2/v2"
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

// ReplaceText replaces a []byte with a new one in a file
//
// @all: if true, will replace all text matching @search,
// if false, will only replace the first occurrence
func ReplaceText(name string, search []byte, rep []byte, all bool) error {
	stat, err := os.Stat(name)
	if err != nil || stat.IsDir() {
		return err
	}

	file, err := os.OpenFile(name, os.O_RDWR, stat.Mode().Perm())
	if err != nil {
		return err
	}
	defer file.Close()

	var found bool

	l := int64(len(search))
	i := int64(0)

	buf := make([]byte, l)
	size, err := file.ReadAt(buf, i)
	buf = buf[:size]
	for err == nil {
		if bytes.Equal(buf, search) {
			found = true

			rl := int64(len(rep))
			if rl == l {
				file.WriteAt(rep, i)
				file.Sync()
			}else if rl < l {
				file.WriteAt(rep, i)
				rl = l - rl

				j := i+l

				b := make([]byte, 1024)
				s, e := file.ReadAt(b, j)
				b = b[:s]

				for e == nil {
					file.WriteAt(b, j-rl)
					j += 1024
					b = make([]byte, 1024)
					s, e = file.ReadAt(b, j)
					b = b[:s]
				}

				if s != 0 {
					file.WriteAt(b, j-rl)
					j += int64(s)
				}

				file.Truncate(j-rl)
				file.Sync()
			}else if rl > l {
				rl -= l

				dif := int64(1024)
				if rl > dif {
					dif = rl
				}

				j := i+l

				b := make([]byte, dif)
				s, e := file.ReadAt(b, j)
				bw := b[:s]

				file.WriteAt(rep, i)
				j += rl

				for e == nil {
					b = make([]byte, dif)
					s, e = file.ReadAt(b, j+dif-rl)
				
					file.WriteAt(bw, j)
					bw = b[:s]

					j += dif
				}

				file.WriteAt(bw, j)
				file.Sync()
			}

			if !all {
				file.Sync()
				file.Close()
				return nil
			}

			i += l
		}

		i++
		s := search[:l-1]

		buf = make([]byte, len(s))
		size, err = file.ReadAt(buf, i)
		buf = buf[:size]

		for err == nil && len(s) > 0 {
			if bytes.Equal(buf, s) {
				break
			}

			i++
			s = s[:len(s)-1]

			buf = make([]byte, len(s))
			size, err = file.ReadAt(buf, i)
			buf = buf[:size]
		}

		buf = make([]byte, l)
		size, err = file.ReadAt(buf, i)
		buf = buf[:size]
	}

	file.Sync()
	file.Close()

	if !found {
		return io.EOF
	}
	return nil
}

// ReplaceRegex replaces a regex match with a new []byte in a file
//
// @all: if true, will replace all text matching @re,
// if false, will only replace the first occurrence
func ReplaceRegex(name string, re string, rep []byte, all bool, maxReSize ...int64) error {
	stat, err := os.Stat(name)
	if err != nil || stat.IsDir() {
		return err
	}

	reg, err := regex.CompTry(re)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(name, os.O_RDWR, stat.Mode().Perm())
	if err != nil {
		return err
	}
	defer file.Close()

	var found bool

	l := int64(len(re) * 10)
	if l < 1024 {
		l = 1024
	}
	for _, maxRe := range maxReSize {
		if l < maxRe {
			l = maxRe
		}
	}

	i := int64(0)

	buf := make([]byte, l)
	size, err := file.ReadAt(buf, i)
	buf = buf[:size]
	for err == nil {
		if reg.Match(buf) {
			found = true

			repRes := reg.RepStr(buf, rep)

			rl := int64(len(repRes))
			if rl == l {
				file.WriteAt(repRes, i)
				file.Sync()
			}else if rl < l {
				file.WriteAt(repRes, i)
				rl = l - rl

				j := i+l

				b := make([]byte, 1024)
				s, e := file.ReadAt(b, j)
				b = b[:s]

				for e == nil {
					file.WriteAt(b, j-rl)
					j += 1024
					b = make([]byte, 1024)
					s, e = file.ReadAt(b, j)
					b = b[:s]
				}

				if s != 0 {
					file.WriteAt(b, j-rl)
					j += int64(s)
				}

				file.Truncate(j-rl)
				file.Sync()
			}else if rl > l {
				rl -= l

				dif := int64(1024)
				if rl > dif {
					dif = rl
				}

				j := i+l

				b := make([]byte, dif)
				s, e := file.ReadAt(b, j)
				bw := b[:s]

				file.WriteAt(repRes, i)
				j += rl

				for e == nil {
					b = make([]byte, dif)
					s, e = file.ReadAt(b, j+dif-rl)
				
					file.WriteAt(bw, j)
					bw = b[:s]

					j += dif
				}

				file.WriteAt(bw, j)
				file.Sync()
			}

			if !all {
				file.Sync()
				file.Close()
				return nil
			}

			i += int64(len(repRes))
		}

		i++
		buf = make([]byte, l)
		size, err = file.ReadAt(buf, i)
		buf = buf[:size]
	}

	if reg.Match(buf) {
		found = true

		repRes := reg.RepStr(buf, rep)

		rl := int64(len(repRes))
		if rl == l {
			file.WriteAt(repRes, i)
			file.Sync()
		}else if rl < l {
			file.WriteAt(repRes, i)
			rl = l - rl

			j := i+l

			b := make([]byte, 1024)
			s, e := file.ReadAt(b, j)
			b = b[:s]

			for e == nil {
				file.WriteAt(b, j-rl)
				j += 1024
				b = make([]byte, 1024)
				s, e = file.ReadAt(b, j)
				b = b[:s]
			}

			if s != 0 {
				file.WriteAt(b, j-rl)
				j += int64(s)
			}

			file.Truncate(j-rl)
			file.Sync()
		}else if rl > l {
			rl -= l

			dif := int64(1024)
			if rl > dif {
				dif = rl
			}

			j := i+l

			b := make([]byte, dif)
			s, e := file.ReadAt(b, j)
			bw := b[:s]

			file.WriteAt(repRes, i)
			j += rl

			for e == nil {
				b = make([]byte, dif)
				s, e = file.ReadAt(b, j+dif-rl)
			
				file.WriteAt(bw, j)
				bw = b[:s]

				j += dif
			}

			file.WriteAt(bw, j)
			file.Sync()
		}
	}

	file.Sync()
	file.Close()

	if !found {
		return io.EOF
	}
	return nil
}

// ReplaceRegexFunc replaces a regex match with the result of a callback function in a file
//
// @all: if true, will replace all text matching @re,
// if false, will only replace the first occurrence
func ReplaceRegexFunc(name string, re string, rep func(data func(int) []byte) []byte, all bool, maxReSize ...int64) error {
	stat, err := os.Stat(name)
	if err != nil || stat.IsDir() {
		return err
	}

	reg, err := regex.CompTry(re)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(name, os.O_RDWR, stat.Mode().Perm())
	if err != nil {
		return err
	}
	defer file.Close()

	var found bool

	l := int64(len(re) * 10)
	if l < 1024 {
		l = 1024
	}
	for _, maxRe := range maxReSize {
		if l < maxRe {
			l = maxRe
		}
	}

	i := int64(0)

	buf := make([]byte, l)
	size, err := file.ReadAt(buf, i)
	buf = buf[:size]
	for err == nil {
		if reg.Match(buf) {
			found = true

			repRes := reg.RepFunc(buf, rep)

			rl := int64(len(repRes))
			if rl == l {
				file.WriteAt(repRes, i)
				file.Sync()
			}else if rl < l {
				file.WriteAt(repRes, i)
				rl = l - rl

				j := i+l

				b := make([]byte, 1024)
				s, e := file.ReadAt(b, j)
				b = b[:s]

				for e == nil {
					file.WriteAt(b, j-rl)
					j += 1024
					b = make([]byte, 1024)
					s, e = file.ReadAt(b, j)
					b = b[:s]
				}

				if s != 0 {
					file.WriteAt(b, j-rl)
					j += int64(s)
				}

				file.Truncate(j-rl)
				file.Sync()
			}else if rl > l {
				rl -= l

				dif := int64(1024)
				if rl > dif {
					dif = rl
				}

				j := i+l

				b := make([]byte, dif)
				s, e := file.ReadAt(b, j)
				bw := b[:s]

				file.WriteAt(repRes, i)
				j += rl

				for e == nil {
					b = make([]byte, dif)
					s, e = file.ReadAt(b, j+dif-rl)
				
					file.WriteAt(bw, j)
					bw = b[:s]

					j += dif
				}

				file.WriteAt(bw, j)
				file.Sync()
			}

			if !all {
				file.Sync()
				file.Close()
				return nil
			}

			i += int64(len(repRes))
		}

		i++
		buf = make([]byte, l)
		size, err = file.ReadAt(buf, i)
		buf = buf[:size]
	}

	if reg.Match(buf) {
		found = true

		repRes := reg.RepFunc(buf, rep)

		rl := int64(len(repRes))
		if rl == l {
			file.WriteAt(repRes, i)
			file.Sync()
		}else if rl < l {
			file.WriteAt(repRes, i)
			rl = l - rl

			j := i+l

			b := make([]byte, 1024)
			s, e := file.ReadAt(b, j)
			b = b[:s]

			for e == nil {
				file.WriteAt(b, j-rl)
				j += 1024
				b = make([]byte, 1024)
				s, e = file.ReadAt(b, j)
				b = b[:s]
			}

			if s != 0 {
				file.WriteAt(b, j-rl)
				j += int64(s)
			}

			file.Truncate(j-rl)
			file.Sync()
		}else if rl > l {
			rl -= l

			dif := int64(1024)
			if rl > dif {
				dif = rl
			}

			j := i+l

			b := make([]byte, dif)
			s, e := file.ReadAt(b, j)
			bw := b[:s]

			file.WriteAt(repRes, i)
			j += rl

			for e == nil {
				b = make([]byte, dif)
				s, e = file.ReadAt(b, j+dif-rl)
			
				file.WriteAt(bw, j)
				bw = b[:s]

				j += dif
			}

			file.WriteAt(bw, j)
			file.Sync()
		}
	}

	file.Sync()
	file.Close()

	if !found {
		return io.EOF
	}
	return nil
}



// A watcher instance for the `FS.FileWatcher` method
type FileWatcher struct {
	// watcherList2 *haxmap.Map[string, *watcherObj]
	watcherList *map[string]*watcherObj
	mu sync.Mutex
	size *uint

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

	// every time something happens
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
	size := uint(0)
	return &FileWatcher{watcherList: &map[string]*watcherObj{}, size: &size}
}

// WatchDir watches the files in a directory and its subdirectories for changes
//
// @nosub: do not watch sub directories
func (fw *FileWatcher) WatchDir(root string, nosub ...bool) error {
	var err error
	if root, err = filepath.Abs(root); err != nil {
		return err
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	runClose := false

	fw.mu.Lock()
	(*fw.watcherList)[root] = &watcherObj{watcher: watcher, close: &runClose}
	*fw.size++
	fw.mu.Unlock()

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

	if len(nosub) == 0 || nosub[0] == false {
		fw.watchDirSub(watcher, root)
	}

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
	fw.mu.Lock()
	defer fw.mu.Unlock()

	if root == "" || root == "*" {
		for r, w := range *fw.watcherList {
			*w.close = true
			delete(*fw.watcherList, r)
			*fw.size--
		}
	}else{
		var err error
		if root, err = filepath.Abs(root); err != nil {
			return err
		}

		if w, ok := (*fw.watcherList)[root]; ok {
			*w.close = true
			delete(*fw.watcherList, root)
			*fw.size--
		}
	}

	return nil
}

// Wait for all Watchers to close
func (fw *FileWatcher) Wait(){
	for *fw.size != 0 {
		time.Sleep(100 * time.Millisecond)
	}
}

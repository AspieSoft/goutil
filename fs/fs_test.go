package fs

import (
	"errors"
	"testing"
)

func TestFS(t *testing.T){
	if val, err := JoinPath("test", "1"); err != nil {
		t.Error("[", val, "]\n", errors.New("JoinPath Method Failed"))
	}

	if val, err := JoinPath("test", "../out/of/root"); err == nil {
		t.Error("[", val, "]\n", errors.New("JoinPath Method Leaked Outsite The Root"))
	}

	watcher := Watcher()
	watcher.WatchDir("./fs")
	watcher.CloseWatcher("*")
	watcher.Wait()
}

# Go Util

[![donation link](https://img.shields.io/badge/buy%20me%20a%20coffee-paypal-blue)](https://paypal.me/shaynejrtaylor?country.x=US&locale.x=en_US)

A simple utility package for golang.

This module simply adds a variety of useful functions in an easy to use way.

## Installation

```shell script
  go get github.com/AspieSoft/goutil/v5
```

## Usage

```go

import (
  "github.com/AspieSoft/goutil/v5"
)

func main(){
  goutil.FS.JoinPath("root", "file") // a safer way to join 2 file paths without backtracking

  goutil.Contains([]any, any) // checks if an array contains a value

  // simple AES-CFB Encryption
  encrypted := goutil.Crypt.CFB.Encrypt([]byte("my message"), []byte("password"))
  goutil.Crypt.CFB.Decrypt(encrypted, []byte("password"))

  // simple gzip compression for strings
  // (also supports brotli and smaz)
  compressed := goutil.GZIP.Zip([]byte("my long string"))
  goutil.GZIP.UnZip(compressed)

  // watch a directory recursively
  watcher := goutil.FS.FileWatcher()

  watcher.OnFileChange = func(path, op string) {
    // do something when a file changes
    path // the file path the change happened to
    op // the change operation
  }

  watcher.OnDirAdd = func(path, op string) {
    // do something when a directory is added
    // return false to prevent that directory from being watched
    return true
  }

  watcher.OnRemove = func(path, op string) {
    // do something when a file or directory is removed
    // return false to prevent that directory from no longer being watched
    return true
  }

  watcher.OnAny = func(path, op string) {
    // do something every time something happenes
  }

  // add a directory to the watch list
  watcher.WatchDir("my/folder")

  // close a specific watcher
  watcher.CloseWatcher("my/folder")

  // close all watchers
  watcher.CloseWatcher("*")

  // wait for all watchers to finish closing
  watcher.Wait()
}

```

# Go Utility Methods

[![donation link](https://img.shields.io/badge/buy%20me%20a%20coffee-paypal-blue)](https://paypal.me/shaynejrtaylor?country.x=US&locale.x=en_US)

A simple utility package for golang.

This module simply adds a variety of useful functions in an easy to use way.

## Installation

```shell script
  go get github.com/AspieSoft/goutil/v7
```

## Usage

```go

import (
  "github.com/AspieSoft/goutil/v7"

  // other optional utility files

  // filesystem
  "github.com/AspieSoft/goutil/fs/v3"

  // encryption
  "github.com/AspieSoft/goutil/crypt"

  // compression
  "github.com/AspieSoft/goutil/compress/gzip"
  "github.com/AspieSoft/goutil/compress/brotli"
  "github.com/AspieSoft/goutil/compress/smaz"

  // other
  "github.com/AspieSoft/goutil/bash"
  "github.com/AspieSoft/goutil/cache"
  "github.com/AspieSoft/goutil/syncmap"
  "github.com/AspieSoft/goutil/cputemp"
)

func main(){
  fs.JoinPath("root", "file") // a safer way to join 2 file paths without backtracking

  goutil.Contains([]any, any) // checks if an array contains a value


  // simple AES-CFB Encryption
  encrypted := crypt.CFB.Encrypt([]byte("my message"), []byte("password"))
  crypt.CFB.Decrypt(encrypted, []byte("password"))


  // simple gzip compression for strings
  // (also supports brotli and smaz)
  compressed := gzip.Zip([]byte("my long string"))
  gzip.UnZip(compressed)


  // convert any type to something else
  MyStr := goutil.ToType[string](MyByteArray)
  MyInt := goutil.ToType[int]("1") // this will return `MyInt == 1`
  MyUInt := goutil.ToType[uint32]([]byte{'2'}) // this will return `MyUInt == 2`
  MyBool := goutil.ToType[bool](1) // this will return `MyBool == true`


  // watch a directory recursively
  watcher := fs.FileWatcher()

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


  // loading config files
  type Config struct{}
  config = Config{}

  // this method will automatically search for [.yml, .yaml, .json, etc...] files
  // it allows the user to decide what compatible file type they want to use for their config
  fs.ReadConfig("path/to/config.yml", &config)
}

```

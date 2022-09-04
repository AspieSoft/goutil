# Go Regex

[![donation link](https://img.shields.io/badge/buy%20me%20a%20coffee-square-blue)](https://buymeacoffee.aspiesoft.com)

A simple utility package for golang.

This module simply adds a variety of useful functions in an easy to use way.

## Installation

```shell script

  go get github.com/AspieSoft/goutil

```

## Usage

```go

import (
  "github.com/AspieSoft/goutil"
)

func main(){
  goutil.JoinPath("root", "file") // a safer way to join 2 file paths without backtracking

  goutil.Contains([]any, any) // checks if an array contains a value

  // simple AES-CFB Encryption
  encrypted := goutil.Encrypt("my message", "password")
  goutil.Decrypt(encrypted, "password")

  // simple gzip compression for strings
  compressed := goutil.Compress("my long string")
  goutil.Decompress(compressed)

  // watch a directory recursively
  goutil.WatchDir("my/folder", &goutil.Watcher{
    FileChange: func(path string, op string){
      // do something when a file changes
      path // the file path the change happened to
      op // the change operation
    },

    DirAdd: func(path string, op string){
      // do something when a directory is added
      // return false to prevent that directory from being watched
      return true
    },

    Remove: func(path string, op string){
      // do something when a file or directory is removed
      // return false to prevent that directory from no longer being watched
      return true
    },

    Any: func(path string, op string){
      // do something every time something happenes
    },
  })
}

```

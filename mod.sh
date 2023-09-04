#!/bin/bash

cd $(dirname "$0")
dir="$PWD"

flag1="$1"
shift


started="0"
for line in $(cat go.work); do
  if [ "$started" = "0" ]; then
    if [ "$line" = "use" ]; then
      started="1"
    elif [ "$line" = "use (" -o "$line" = "use(" ]; then
      started="2"
    fi
  elif [ "$started" = "1" ]; then
    if [ "$line" = "(" ]; then
      started="2"
    fi
  elif [ "$started" = "2" ]; then
    if [ "$line" = ")" ]; then
      started="0"
      break
    else
      cd "$line"

      if [ "$flag1" = "-u" -o "$flag1" = "update" ]; then
        go mod tidy
        go get -u $@
      elif [ "$flag1" = "-i" -o "$flag1" = "tidy" ]; then
        go mod tidy $@
        go get -u
      elif [ "$flag1" = "-t" -o "$flag1" = "test" ]; then
        go test $@
      elif [ "$flag1" = "-b" -o "$flag1" = "build" ]; then
        go build $@
      else
        go mod tidy
        go get -u
        go test
      fi

      cd "$dir"
    fi
  fi
done

if [ "$flag1" = "-u" -o "$flag1" = "update" ]; then
  go mod tidy
  go get -u $@
elif [ "$flag1" = "-i" -o "$flag1" = "tidy" ]; then
  go mod tidy $@
  go get -u
elif [ "$flag1" = "-t" -o "$flag1" = "test" ]; then
  go test $@
elif [ "$flag1" = "-b" -o "$flag1" = "build" ]; then
  go build $@
else
  go mod tidy
  go get -u
  go test
fi

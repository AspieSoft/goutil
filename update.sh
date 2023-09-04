#!/bin/bash

go mod tidy
go get -u

cd brotli
go mod tidy
go get -u

cd ..

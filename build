#! /bin/bash
set -ex

env GOOS=darwin GOARCH=amd64 go build -o ./builds/plis-darwin_x86_64 plis.go
env GOOS=windows GOARCH=amd64 go build -o ./builds/plis-windows_x86_64 plis.go
env GOOS=linux GOARCH=amd64 go build -o ./builds/plis-linux_x86_64 plis.go

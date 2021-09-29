#!/bin/sh
set -e

GOOS=windows GOARCH=386 go build -o ./out/ ./cmd/ssh
GOOS=windows GOARCH=386 go build -o ./out/ ./cmd/scp
GOOS=windows GOARCH=386 go build -o ./out/ ./cmd/grun
tar czf mutagen-ssh-wrapper.tar.gz -C ./out scp.exe ssh.exe grun.exe

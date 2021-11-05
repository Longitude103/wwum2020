#!/usr/bin/env zsh

#GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=mingw64-configure go build -o bin/wwum2020-amd64.exe
GOOS=darwin GOARCH=amd64 go build -o bin/wwum2020-amd64-darwin
GOOS=linux GOARCH=amd64 go build -o bin/wwum2020-amd64-linux

echo "Done Compiling wwum2020, look in bin/"

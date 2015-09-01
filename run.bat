@ECHO off

rm go-codeivate.exe

go build -v -o go-codeivate.exe main.go

start "" "go-codeivate.exe" %*
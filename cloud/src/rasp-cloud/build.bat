@echo off
set CGO_ENABLED=0 
set GOOS=linux
set GOPACH=amd64
go build
go install
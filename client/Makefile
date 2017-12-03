SHELL := /bin/bash 
VERSION = 1.0.0
CURDIR = $(shell pwd)

versionDir="grapehttp/pkg/vinfo"
gitTag = $(shell if [ "`git describe --tags --abbrev=0`" != "" ];then git describe --tags --abbrev=0; else git log --pretty=format:'%h' -n 1; fi)
gitCommit = $(shell git log --pretty=format:'%H' -n 1)
gitTreeState = $(shell if git status|grep -q 'clean';then echo clean; else echo dirty; fi)
buildDate=$(shell TZ=Asia/Shanghai date +%FT%T%z)

all:
	@go build -v -ldflags "-w -X ${versionDir}.gitTag=${gitTag} -X ${versionDir}.buildDate=${buildDate} -X ${versionDir}.gitCommit=${gitCommit} -X ${versionDir}.gitTreeState=${gitTreeState}" -o ../fctl

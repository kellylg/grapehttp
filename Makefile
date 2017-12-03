SHELL := /bin/bash 
VERSION = 1.0.0
CURDIR = $(shell pwd)
BASEDIR = ${CURDIR}

# build with verison infos
versionDir = "grapehttp/pkg/vinfo"
gitTag = $(shell if [ "`git describe --tags --abbrev=0`" != "" ];then git describe --tags --abbrev=0; else git log --pretty=format:'%h' -n 1; fi)
buildDate = $(shell TZ=Asia/Shanghai date +%FT%T%z)
gitCommit = $(shell git log --pretty=format:'%H' -n 1)
gitTreeState = $(shell if git status|grep -q 'clean';then echo clean; else echo dirty; fi)
LDFLAGS="-w -X ${versionDir}.gitTag=${gitTag} -X ${versionDir}.buildDate=${buildDate} -X ${versionDir}.gitCommit=${gitCommit} -X ${versionDir}.gitTreeState=${gitTreeState}"

all: gotool fctl server

server:
	@go build -v -ldflags ${LDFLAGS}

fctl:
	@make -C ${BASEDIR}/client
	@echo binary file is: ${BASEDIR}/fctl

gotool:
	@-gofmt -w  .
	@-go tool vet . |& grep -v vendor

clean:
	rm -f grapehttp

.PHONY: gotool clean fctl

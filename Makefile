GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean

VERSION=`git describe --tags`
DATE=`date -u '+%Y-%m-%d_%I:%M:%S%p'`
LDFLAGS=-ldflags "-X main.buildDate=${DATE} -X main.buildVersion=${VERSION}"


all: build

build:
	cp ${GOPATH}/bin/migrate bhlindex
	cd cmd/bhlindex && \
	$(GOCLEAN) && \
	GOOS=linux GOARCH=amd64 $(GOBUILD) ${LDFLAGS} && \
	mv bhlindex ../../bhlindex

release: build
	tar --exclude='bhlindex/development.sh' -zcvf /tmp/bhlindex-${VERSION}-linux.tar.gz bhlindex

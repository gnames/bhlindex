GOCMD=go
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean

VERSION=`git describe --tags`
DATE=`date -u '+%Y-%m-%d_%I:%M:%S%p'`
LDFLAGS=-ldflags "-X main.buildDate=${DATE} -X main.buildVersion=${VERSION}"


all: install

build:
	cp ${GOPATH}/bin/migrate bhlindex
	cd bhlindex && \
	$(GOCLEAN) && \
	GOOS=linux GOARCH=amd64 $(GOBUILD) ${LDFLAGS}

install:
	cd bhlindex && \
	$(GOINSTALL) ${LDFLAGS};

release: build
	tar --exclude='bhlindex/development.sh' --exclude='bhlindex/cmd' -zcvf /tmp/bhlindex-${VERSION}-linux.tar.gz bhlindex

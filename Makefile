GOCMD=go
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean

VERSION=`git describe --tags`
DATE=`date -u '+%Y-%m-%d_%I:%M:%S%p'`
LDFLAGS=-ldflags "-X main.buildDate=${DATE} -X main.buildVersion=${VERSION}"

all: install

build: grpc
	cp ${GOPATH}/bin/migrate bhlindex
	cd bhlindex && \
	$(GOCLEAN) && \
	GO111MODULE=on GOOS=linux GOARCH=amd64 $(GOBUILD) ${LDFLAGS}

install: grpc
	cd bhlindex && \
	GO111MODULE=on $(GOINSTALL) ${LDFLAGS};

release: build
	tar --exclude='bhlindex/development.sh' --exclude='bhlindex/cmd' \
	--exclude='main.go' -zcvf /tmp/bhlindex-${VERSION}-linux.tar.gz bhlindex

grpc:
	cd protob && \
	protoc -I . ./protob.proto --go_out=plugins=grpc:.

GOCMD=go
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean

VERSION=`git describe --tags`
DATE=`date -u '+%Y-%m-%d_%I:%M:%S%p'`
FLAGS_SHARED = CGO_ENABLED=0 GOARCH=amd64
FLAGS_LD=-ldflags "-w -s \
                  -X github.com/gnames/bhlindex.Build=${DATE} \
                  -X github.com/gnames/bhlindex.Version=${VERSION}"

all: install

tools: deps
	@echo Installing tools from tools.go
	@cat gnverifier/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

deps:
	@echo Download go.mod dependencies
	$(GOCMD) mod download; \
	$(GOGENERATE)

test: deps install
	@echo Run tests
	go test -race ./...

build:
	cd bhlindex; \
	$(GOCLEAN); \
	$(FLAGS_SHARED) $(GOBUILD);

install:
	@echo Building and Installing bhlindex
	cd bhlindex; \
	$(FLAGS_SHARED) $(GOINSTALL); \
	$(GOCLEAN); 

release: build
	@echo Building releases for Linux, Mac, Windows
	cd gnverifier; \
	$(GOCLEAN); \
	$(FLAGS_SHARED) GOOS=linux $(GOBUILD); \
	tar zcvf /tmp/bhlindex-${VER}-linux.tar.gz bhlindex; \

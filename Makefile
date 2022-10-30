VERSION=`git describe --tags`
VER = $(shell git describe --tags --abbrev=0)
DATE=`date -u '+%Y-%m-%d_%I:%M:%S%p'`
FLAGS_SHARED = CGO_ENABLED=0 GOARCH=amd64
FLAGS_LD=-ldflags "-w -s \
                  -X github.com/gnames/bhlindex.Build=${DATE} \
                  -X github.com/gnames/bhlindex.Version=${VERSION}"
GOCMD=go
GOBUILD=$(GOCMD) build $(FLAGS_LD)
GOINSTALL=$(GOCMD) install $(FLAGS_LD)
GOCLEAN=$(GOCMD) clean

all: install

tools: deps
	@echo Installing tools from tools.go
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

deps:
	@echo Download go.mod dependencies
	$(GOCMD) mod download; \
	$(GOGENERATE)

test: deps install
	@echo Run tests
	go test -race ./...

build:
	$(GOCLEAN); \
	$(FLAGS_SHARED) $(GOBUILD);

install:
	@echo Building and Installing bhlindex
	$(FLAGS_SHARED) $(GOINSTALL); \
	$(GOCLEAN); 

release:
	@echo Building release for Linux
	$(GOCLEAN); \
	$(FLAGS_SHARED) GOOS=linux $(GOBUILD); \
	tar zcvf /tmp/bhlindex-${VER}-linux.tar.gz bhlindex;
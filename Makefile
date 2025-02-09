PROJ_NAME = bhlindex

VERSION=`git describe --tags`
VER = $(shell git describe --tags --abbrev=0)
DATE=`date -u '+%Y-%m-%d_%I:%M:%S%p'`
FLAGS_LD = -ldflags "-X github.com/gnames/$(PROJ_NAME)/internal.Build=${DATE} \
                     -X github.com/gnames/$(PROJ_NAME)/internal.Version=${VERSION}"
FLAGS_REL = -trimpath -ldflags "-s -w -X github.com/gnames/$(PROJ_NAME)/internal.Build=$(DATE)"

NO_C = CGO_ENABLED=0
FLAGS_LINUX = $(NO_C) GOARCH=amd64 GOOS=linux

GOCMD=go
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install
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
	go test -race -count=1 ./...

build:
	$(GOCLEAN); \
	$(NO_C) $(GOBUILD) $(FLAGS_LD);
	
buildrel:
	$(GOCLEAN); \
	$(NO_C) $(GOBUILD) $(FLAGS_REL);

install:
	@echo Building and Installing bhlindex
	$(NO_C) $(GOINSTALL) $(FLAGS_LD); \
	$(GOCLEAN); 

release:
	@echo Building release for Linux
	$(GOCLEAN); \
	$(FLAGS_LINUX)  $(GOBUILD) $(FLAGS_REL); \
	tar zcvf /tmp/bhlindex-${VER}-linux.tar.gz bhlindex;

include .env

APP=bumper

CGO_ENABLED ?= 0
GOOS ?= linux
GOARCH ?= amd64
GOBUILDFLAGS ?= CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH)

VERSION=$(shell git describe --always --long)
BuildTime=$(shell date +%Y-%m-%d.%H:%M:%S)
GitDirty=$(shell git status --porcelain --untracked-files=no)
GitCommit=$(shell git rev-parse --short HEAD)
ifneq ($(GitDirty),)
	GitCommit:= $(GitCommit)-dirty
endif


DIST=dist
SOURCES = app.go bumper.go

.PHONY: test

compiled: test
	$(GOBUILDFLAGS) go build  $(SOURCES) -ldflags "-s -X main.CommitHash=$(GitCommit) -X main.BuildTime=$(BuildTime)" -a -installsuffix cgo -o $(DIST)/alpine/$(APP) .

build: test
	go build  -ldflags "-s -X main.CommitHash=$(GitCommit) -X main.BuildTime=$(BuildTime)" -a  -o $(DIST)/$(APP) $(SOURCES)

install:
	go install  -ldflags "-s -X main.CommitHash=$(GitCommit) -X main.BuildTime=$(BuildTime)"

run:
	go run $(SOURCES) -i "0.0.1 testTag"
	go run $(SOURCES) -p m -i "0.0.1 testMinor"
	go run $(SOURCES) -p M -i "0.0.1 testMajor"

test:
	 go test -v ./...

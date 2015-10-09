.PHONY: build test clean

BINARY := temperature-monitor
IMPORT_BASE := github.com/alext
IMPORT_PATH := $(IMPORT_BASE)/temperature-monitor

GOPATH := $(CURDIR)/gopath:$(CURDIR)/Godeps/_workspace
export GOPATH

ifdef RELEASE_VERSION
VERSION := $(RELEASE_VERSION)
else
VERSION := $(shell git describe --always | tr -d '\n'; test -z "`git status --porcelain`" || echo '-dirty')
endif

build: Godeps/Godeps.json gopath
	go build -ldflags "-X main.version=$(VERSION)" -o $(BINARY)

test: build gopath
	go test -v ./...
	./$(BINARY) -version

clean:
	rm -rf $(BINARY) gopath

gopath:
	rm -f gopath/src/$(IMPORT_PATH)
	mkdir -p gopath/src/$(IMPORT_BASE)
	ln -s $(CURDIR) gopath/src/$(IMPORT_PATH)

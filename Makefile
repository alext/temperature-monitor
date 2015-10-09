.PHONY: build test clean

BINARY := temperature-monitor
IMPORT_BASE := github.com/alext
IMPORT_PATH := $(IMPORT_BASE)/temperature-monitor

GOPATH := $(CURDIR)/gopath:$(CURDIR)/Godeps/_workspace
export GOPATH

build: Godeps/Godeps.json gopath
	go build -o $(BINARY)

test: gopath
	go test -v ./...

clean:
	rm -rf $(BINARY) gopath

gopath:
	rm -f gopath/src/$(IMPORT_PATH)
	mkdir -p gopath/src/$(IMPORT_BASE)
	ln -s $(CURDIR) gopath/src/$(IMPORT_PATH)

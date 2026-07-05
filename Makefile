BINARY := awm
PREFIX ?= /usr/local/bin
TAG ?= $(shell git describe --tags --exact-match 2>/dev/null)
COMMIT ?= $(shell git rev-parse --short=8 HEAD 2>/dev/null)
DIRTY ?= $(shell test -n "`git status --porcelain 2>/dev/null`" && echo true || echo false)
LDFLAGS := -X awmcli/internal/version.Tag=$(TAG) -X awmcli/internal/version.Commit=$(COMMIT) -X awmcli/internal/version.Dirty=$(DIRTY)

.PHONY: build release test fmt generate install clean

build: generate
	go build -o $(BINARY) ./cmd/awm

release:
	@if [ -z "$(TAG)" ]; then echo "release requires building from an exact git tag"; exit 1; fi
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) ./cmd/awm

test:
	go test ./...

fmt:
	gofmt -w ./cmd/awm/*.go ./internal/**/*.go ./tools/**/*.go

generate:
	go generate ./internal/version

install: build
	cp $(BINARY) $(PREFIX)/$(BINARY)

clean:
	rm -f $(BINARY)

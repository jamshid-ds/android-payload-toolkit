# Makefile for Android Payload Toolkit

BINARY_NAME=android-payload-toolkit
GO=go
GOFLAGS=-ldflags="-s -w"

# Platform specific settings
UNAME := $(shell uname)
ifeq ($(UNAME), Darwin)
    CGO_CFLAGS=-I/opt/homebrew/opt/xz/include -I/usr/local/opt/xz/include
    CGO_LDFLAGS=-L/opt/homebrew/opt/xz/lib -L/usr/local/opt/xz/lib
endif

.PHONY: all build clean test install deps

all: build

deps:
	$(GO) mod download
	$(GO) mod tidy

build: deps
	CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" \
	$(GO) build $(GOFLAGS) -o $(BINARY_NAME) .

build-all: deps
	# Linux AMD64
	GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) -o dist/$(BINARY_NAME)-linux-amd64
	# Linux ARM64
	GOOS=linux GOARCH=arm64 $(GO) build $(GOFLAGS) -o dist/$(BINARY_NAME)-linux-arm64
	# macOS AMD64
	CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" \
	GOOS=darwin GOARCH=amd64 $(GO) build $(GOFLAGS) -o dist/$(BINARY_NAME)-darwin-amd64
	# macOS ARM64 (M1/M2)
	CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" \
	GOOS=darwin GOARCH=arm64 $(GO) build $(GOFLAGS) -o dist/$(BINARY_NAME)-darwin-arm64
	# Windows AMD64
	GOOS=windows GOARCH=amd64 $(GO) build $(GOFLAGS) -o dist/$(BINARY_NAME)-windows-amd64.exe

install: build
	cp $(BINARY_NAME) /usr/local/bin/

test:
	$(GO) test ./...

clean:
	rm -f $(BINARY_NAME)
	rm -rf dist/
	rm -rf extracted_*
	rm -f *.img *.bin
	$(GO) clean

run: build
	./$(BINARY_NAME)

help:
	@echo "Available targets:"
	@echo "  make build      - Build for current platform"
	@echo "  make build-all  - Build for all platforms"
	@echo "  make install    - Install to /usr/local/bin"
	@echo "  make test       - Run tests"
	@echo "  make clean      - Clean build artifacts"
	@echo "  make deps       - Download dependencies"
	@echo "  make help       - Show this help"
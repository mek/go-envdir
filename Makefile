GO ?= go
GOFMT ?= gofmt
PKGS := ./...
BIN_DIR := bin
BIN := $(BIN_DIR)/envdir
FILESERVER_BIN := $(BIN_DIR)/fileserver
GOFILES := $(shell find . -name '*.go' -not -path './bin/*' -print)

.DEFAULT_GOAL := help

.PHONY: all build build-cli build-examples check clean fmt fmt-check help tidy test vet

all: check build # Build after running checks

build: build-cli build-examples # Build the CLI and sample programs

build-cli: # Build the envdir CLI
	@mkdir -p $(BIN_DIR)
	$(GO) build -o $(BIN) ./cmd/envdir

build-examples: # Build the sample programs
	@mkdir -p $(BIN_DIR)
	$(GO) build -o $(FILESERVER_BIN) ./examples/fileserver

check: fmt-check vet test # Run formatting checks, vet, and tests

clean: # Remove build artifacts
	rm -rf $(BIN_DIR)

fmt: # Format Go source files
	$(GOFMT) -w $(GOFILES)

fmt-check: # Fail if Go files are not formatted
	@test -z "$$($(GOFMT) -l $(GOFILES))" || \
		(echo "gofmt found unformatted files"; $(GOFMT) -l $(GOFILES); exit 1)

help: # Show available make targets
	@awk 'BEGIN {FS = ":.*# "}; /^[[:alnum:]_-]+:.*# / {printf "make %-10s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

tidy: # Tidy module dependencies
	$(GO) mod tidy

test: # Run Go tests
	$(GO) test $(PKGS)

vet: # Run go vet
	$(GO) vet $(PKGS)

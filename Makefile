# TypeGo Makefile

BINARY_NAME=typego.exe
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

.PHONY: all build test clean lint deps

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/typego

test:
	$(GOTEST) -v -race ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

lint:
	golangci-lint run

deps:
	$(GOGET) ./...

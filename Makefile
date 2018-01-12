# Current version of the project.
VERSION ?= v0.0.1

# This repo's root import path (under GOPATH).
ROOT := github.com/gaocegege/kubetidb

# Project main package location (can be multiple ones).
CMD_DIR := ./cmd/controller

# Project output directory.
OUTPUT_DIR := ./bin

# Git commit sha.
GitSHA := $(shell git rev-parse --short HEAD)

# Golang standard bin directory.
BIN_DIR := $(GOPATH)/bin

# Golang packages except vendor.
PACKAGES := $(shell go list ./... | grep -v /vendor/ )

build:
	go build -i -v -o $(OUTPUT_DIR)/kubetidb \
	  -ldflags "-s -w -X $(ROOT)/pkg/version.Version=$(VERSION) \
	            -X $(ROOT)/pkg/version.GitSHA=$(GitSHA)" \
	  $(CMD_DIR) \

test:
	go test $(PACKAGES)

clean:
	-rm -vrf ${OUTPUT_DIR}

.PHONY: clean build

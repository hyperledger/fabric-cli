# Copyright State Street Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0

PROJECT_VERSION=dev

GO_CMD		?= go
LINT_CMD	?= gobin -run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.19.1

BIN_DIR := $(CURDIR)/bin
CMD_DIR := $(CURDIR)/cmd

PROJECT_NAME = hyperledger/fabric-cli
PKGNAME = github.com/$(PROJECT_NAME)

EXTRA_VERSION ?= $(shell git rev-parse --short HEAD)

# defined in common/metadata/metadata.go
METADATA_VAR = Version=$(PROJECT_VERSION)
METADATA_VAR += CommitSHA=$(EXTRA_VERSION)

GO_LDFLAGS = $(patsubst %,-X $(PKGNAME)/common/metadata.%,$(METADATA_VAR))

export GO111MODULE := on

all: clean build

.PHONY: clean
clean:
	rm -rf $(BIN_DIR)
	find . -name "mocks" -type d -print0 | xargs -0 /bin/rm -rf

.PHONY: generate
generate:
	$(GO_CMD) generate ./...

.PHONY: lint
lint: generate
	$(LINT_CMD) run

.PHONY: test
test: generate
	$(GO_CMD) test -cover ./...

.PHONY: build
build:
	$(GO_CMD) build -o $(BIN_DIR)/fabric -ldflags "$(GO_LDFLAGS)" $(CMD_DIR)/fabric.go
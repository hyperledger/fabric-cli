# Copyright State Street Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0

GO_CMD		?= go
LINT_CMD	?= gobin -run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.19.1

BIN_DIR := $(CURDIR)/bin
CMD_DIR := $(CURDIR)/cmd

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
	$(GO_CMD) build -o $(BIN_DIR)/fabric $(CMD_DIR)/fabric.go
# Copyright State Street Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0

GO_CMD		?= go
DEP_CMD		?= dep
LINT_CMD	?= gometalinter

BIN_DIR := $(CURDIR)/bin
CMD_DIR := $(CURDIR)/cmd

all: clean dep build

.PHONY: clean
clean:
	rm -rf $(BIN_DIR)
	find . -name "mocks" -type d -print0 | xargs -0 /bin/rm -rf

.PHONY: dep
dep:
	$(DEP_CMD) ensure

.PHONY: lint
lint:
	$(LINT_CMD) --disable=gocyclo --disable=gas --deadline=120s --exclude=vendor ./...

.PHONY: generate
generate: 
	$(GO_CMD) generate ./...

.PHONY: test
test: generate
	$(GO_CMD) test -cover ./...

.PHONY: build
build:
	$(GO_CMD) build -o $(BIN_DIR)/fabric $(CMD_DIR)/fabric.go
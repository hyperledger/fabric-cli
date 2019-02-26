# Copyright State Street Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0

GO_CMD := go

all: clean build

clean:
	rm -rf ./bin

build:
	$(GO_CMD) build -o ./bin/home ./cmd/home.go
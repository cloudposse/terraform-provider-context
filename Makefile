SHELL := /bin/bash

# List of targets the `readme` target should call before generating the readme
export README_DEPS ?= docs/targets.md docs/terraform.md

-include $(shell curl -sSL -o .build-harness "https://cloudposse.tools/build-harness"; echo .build-harness)

.PHONY: build
build:
	go build .

docs: readme

.PHONY: deps
deps:
	go get github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt
	go mod tidy

.PHONY: generate
generate:
	go generate ./...

# Install the binary in $GOPATH/bin
.PHONY: install
install:
	go install .

tfdocs:
	tfplugindocs

.PHONY: test
test:
	go test ./... -json -v $(TESTARGS) -timeout 5m -count 1 | gotestfmt

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -json -v $(TESTARGS) -count 1 -timeout 120m | gotestfmt

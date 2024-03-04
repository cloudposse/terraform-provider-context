default: testacc

.PHONY: build
build:
	go build .

.PHONY: install
install:
	go install .

.PHONY: generate
generate:
	go generate ./...

.PHONY: test
test:
	go test ./... -v $(TESTARGS) -timeout 5m

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

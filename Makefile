default: testacc

# Run acceptance tests.
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./provider/... -v $(TESTARGS) -timeout 120m

# This is used to install the provider locally.
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

# Development version.
VERSION := 0.1.0

# Common settings that help with the build process.
REGISTRY := github.com
ORG := couchbasecloud
NAME := couchbasecapella

# Terraform requieres a specific format.
BINARY := terraform-provider-$(NAME)

.PHONY: build
build:
	@go build -o $(BINARY)

.PHONY: install
install: build
	@mkdir -p ~/.terraform.d/plugins/$(REGISTRY)/$(ORG)/$(NAME)/$(VERSION)/$(GOOS)_$(GOARCH)
	@mv $(BINARY) ~/.terraform.d/plugins/$(REGISTRY)/$(ORG)/$(NAME)/$(VERSION)/$(GOOS)_$(GOARCH)

.PHONY: fmt
fmt:
	@terraform fmt -diff -recursive .

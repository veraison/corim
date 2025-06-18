.DEFAULT_GOAL := help

SHELL := /bin/bash

GO111MODULE := on

GOPKG := github.com/veraison/corim/corim
GOPKG += github.com/veraison/corim/comid
GOPKG += github.com/veraison/corim/cots
GOPKG += github.com/veraison/corim/encoding
GOPKG += github.com/veraison/corim/extensions

GOLINT ?= golangci-lint

GOLINT_ARGS ?= run --timeout=3m -E dupl -E gocritic -E staticcheck -E lll -E prealloc

.PHONY: lint
lint:
	$(GOLINT) $(GOLINT_ARGS)

ifeq ($(MAKECMDGOALS),test)
GOTEST_ARGS ?= -v -race $(GOPKG)
else
  ifeq ($(MAKECMDGOALS),test-cover)
  GOTEST_ARGS ?= -short -cover $(GOPKG)
  endif
endif

COVER_THRESHOLD := $(shell grep '^name: cover' .github/workflows/ci-go-cover.yml | cut -c13-)

.PHONY: test test-cover
test test-cover:
	go test $(GOTEST_ARGS)

realtest:
	go test $(GOTEST_ARGS)
.PHONY: realtest

presubmit:
	@echo
	@echo ">>> Check that the reported coverage figures are $(COVER_THRESHOLD)"
	@echo
	$(MAKE) test-cover
	@echo
	@echo ">>> Fix any lint error"
	@echo
	$(MAKE) lint

.PHONY: licenses
licenses: ; @./scripts/licenses.sh

.PHONY: test-certs
test-certs:
	@echo "Regenerating certificate chain..."
	@$(SHELL) scripts/gen-certs.sh create

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  * test:       run unit tests for $(GOPKG)"
	@echo "  * test-cover: run unit tests and measure coverage for $(GOPKG)"
	@echo "  * lint:       lint sources using default configuration and some extra checkers"
	@echo "  * presubmit:  check you are ready to push your local branch to remote"
	@echo "  * help:       print this menu"
	@echo "  * licenses:   check licenses of dependent packages"
	@echo "  * test-certs: regenerate the certificate chain"

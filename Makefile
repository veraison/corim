.DEFAULT_GOAL := help

SHELL := /bin/bash

GO111MODULE := on

GOPKG := github.com/veraison/corim/corim
GOPKG += github.com/veraison/corim/comid
GOPKG += github.com/veraison/corim/cocli/cmd

GOLINT ?= golangci-lint

ifeq ($(MAKECMDGOALS),lint)
GOLINT_ARGS ?= run --timeout=3m
else
  ifeq ($(MAKECMDGOALS),lint-extra)
  GOLINT_ARGS ?= run --timeout=3m --issues-exit-code=0 -E dupl -E gocritic -E gosimple -E lll -E prealloc
  endif
endif

.PHONY: lint lint-extra
lint lint-extra: ; $(GOLINT) $(GOLINT_ARGS)

ifeq ($(MAKECMDGOALS),test)
GOTEST_ARGS ?= -v -race $(GOPKG)
else
  ifeq ($(MAKECMDGOALS),test-cover)
  GOTEST_ARGS ?= -short -cover $(GOPKG)
  endif
endif

COVER_THRESHOLD := $(shell grep '^name: cover' .github/workflows/ci-go-cover.yml | cut -c13-)

.PHONY: test test-cover
test test-cover: ; go test $(GOTEST_ARGS)

presubmit:
	@echo
	@echo ">>> Check that the reported coverage figures are $(COVER_THRESHOLD)"
	@echo
	$(MAKE) test-cover
	@echo
	@echo ">>> Fix any lint error"
	@echo
	$(MAKE) lint-extra

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  * test:       run unit tests for $(GOPKG)"
	@echo "  * test-cover: run unit tests and measure coverage for $(GOPKG)"
	@echo "  * lint:       lint sources using default configuration"
	@echo "  * lint-extra: lint sources using default configuration and some extra checkers"
	@echo "  * presubmit:  check you are ready to push your local branch to remote"
	@echo "  * help:       print this menu"

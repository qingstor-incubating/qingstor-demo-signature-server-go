SHELL := /bin/bash

DIRS_TO_CHECK=$(shell ls -d */ | grep -vE "vendor|test")
PKGS_TO_CHECK=$(shell go list ./... | grep -v "/vendor/")
FILES_CLIENT=$(shell ls ./client/*.go)
FILES_SERVER=$(shell ls ./server/*.go)

.PHONY: help
help:
	@echo "Please use \`make <target>\` where <target> is one of"
	@echo "  check             to vet and lint the QingStor Demo Signature Server"
	@echo "  run               to run the QingStor Demo Signature Server"
	@echo "  test              to test the QingStor Demo Signature Server"

.PHONY: check
check: vet lint

.PHONY: vet
vet:
	@echo "go tool vet, skipping vendor packages"
	@go tool vet -all ${DIRS_TO_CHECK}
	@echo "ok"

.PHONY: lint
lint:
	@echo "golint, skipping vendor packages"
	@lint=$$(for pkg in ${PKGS_TO_CHECK}; do golint $${pkg}; done); \
	 lint=$$(echo "$${lint}"); \
	 if [[ -n $${lint} ]]; then echo "$${lint}"; exit 1; fi
	@echo "ok"

.PHONY: run
run:
	go run ${FILES_SERVER}
	@echo "ok"

.PHONY: test
test:
	go run ${FILES_CLIENT}
	@echo "ok"

SHELL := /bin/bash

GO_TEST := go test -count=1 -v -p=1

## Run Go unit tests
example:
	${GO_TEST} ./example
.PHONY: example

## Run Go unit tests
test:
	${GO_TEST} ./...
.PHONY: test

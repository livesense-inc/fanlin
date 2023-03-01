MAKEFLAGS += --warn-undefined-variables
SHELL     := /bin/bash -e -u -o pipefail

echo:
	@echo Hi

test:
	@go clean -testcache
	@go test -race ./...

prof: PKG ?= handler
prof: TYPE ?= mem
prof:
	@if [[ -z "${PKG}" ]]; then echo 'empty variable: PKG'; exit 1; fi
	@if [[ -z "${TYPE}" ]]; then echo 'empty variable: TYPE'; exit 1; fi
	@if [[ ! -d "./lib/${PKG}" ]]; then echo 'package not found: ${PKG}'; exit 1; fi
	@go test -bench=. -run=NONE -${TYPE}profile=${TYPE}.out ./lib/${PKG}
	@go tool pprof -text -nodecount=10 ./${PKG}.test ${TYPE}.out

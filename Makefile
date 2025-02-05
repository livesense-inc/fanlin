MAKEFLAGS += --warn-undefined-variables
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
CGO_ENABLED ?= $(shell go env CGO_ENABLED)

build: cmd/fanlin/server

cmd/fanlin/server: cmd/fanlin/main.go
	GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=${CGO_ENABLED} go build -ldflags="-s -w" -trimpath -tags timetzdata -o $@ $^

cmd/fanlin/fanlin.json: cmd/fanlin/sample-conf.json
	@cp $^ $@

fanlin.json: cmd/fanlin/fanlin.json
	@ln -sf $^

run: cmd/fanlin/server fanlin.json
	@$^

test:
	@go clean -testcache
	@go test -race ./...

lint:
	@go vet ./...

bench:
	@go test -bench=. -benchmem -run=NONE ./...

prof: PKG ?= handler
prof: TYPE ?= mem
prof:
	@if [ -z "${PKG}" ]; then echo 'empty variable: PKG'; exit 1; fi
	@if [ -z "${TYPE}" ]; then echo 'empty variable: TYPE'; exit 1; fi
	@if [ ! -d "./lib/${PKG}" ]; then echo 'package not found: ${PKG}'; exit 1; fi
	@go test -bench=. -run=NONE -${TYPE}profile=${TYPE}.out ./lib/${PKG}
	@go tool pprof -text -nodecount=10 ./${PKG}.test ${TYPE}.out

clean:
	@unlink fanlin.json || true
	@rm -f cmd/fanlin/server cmd/fanlin/fanlin.json

.PHONY: build cmd/fanlin/server run test lint bench prof clean

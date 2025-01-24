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

clean:
	@unlink fanlin.json || true
	@rm -f cmd/fanlin/server cmd/fanlin/fanlin.json

.PHONY: build cmd/fanlin/server test clean

MAKEFLAGS += --warn-undefined-variables
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
CGO_ENABLED ?= $(shell go env CGO_ENABLED)

server: cmd/fanlin/main.go
	GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=${CGO_ENABLED} go build -ldflags="-s -w" -trimpath -tags timetzdata -o $@ $^

test:
	@go clean -testcache
	@go test -race ./...

lint:
	@go vet ./...

.PHONY: server test lint

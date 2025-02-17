MAKEFLAGS += --warn-undefined-variables
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
CGO_ENABLED ?= $(shell go env CGO_ENABLED)
AWS_ENDPOINT_URL := http://127.0.0.1:4567
AWS_REGION := ap-northeast-1
AWS_CMD_ENV += AWS_ACCESS_KEY_ID=AAAAAAAAAAAAAAAAAAAA
AWS_CMD_ENV += AWS_SECRET_ACCESS_KEY=AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AWS_CMD_OPT += --endpoint-url=${AWS_ENDPOINT_URL}
AWS_CMD_OPT += --region=${AWS_REGION}
AWS_CMD := ${AWS_CMD_ENV} aws ${AWS_CMD_OPT}
AWS_S3_BUCKET_NAME := local-test

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

create-s3-bucket:
	@${AWS_CMD} s3api create-bucket \
		--bucket=${AWS_S3_BUCKET_NAME} \
		--create-bucket-configuration LocationConstraint=${AWS_REGION}

clean-s3-bucket:
	@${AWS_CMD} s3 rm s3://${AWS_S3_BUCKET_NAME} --include='*' --recursive

list-s3-bucket:
	@${AWS_CMD} s3 ls s3://${AWS_S3_BUCKET_NAME}/${FOLDER}

copy-object:
	@${AWS_CMD} s3 cp ${SRC} s3://${AWS_S3_BUCKET_NAME}/${DEST}

.PHONY: build cmd/fanlin/server run test lint bench prof clean

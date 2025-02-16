FROM golang:1-bookworm AS builder
RUN set -eux; \
  apt-get update -y; \
  apt-get install -y --no-install-recommends ca-certificates build-essential make libaom-dev
WORKDIR /usr/src/app
COPY . .
RUN make build

# https://github.com/GoogleContainerTools/distroless
# https://console.cloud.google.com/gcr/images/distroless/GLOBAL
FROM gcr.io/distroless/cc-debian12:nonroot-amd64
COPY --from=builder /lib/x86_64-linux-gnu/libaom.so.* /lib/x86_64-linux-gnu/
COPY --from=builder /lib/x86_64-linux-gnu/libm.so.* /lib/x86_64-linux-gnu/
COPY --from=builder /usr/src/app/cmd/fanlin/server /usr/local/bin/fanlin
ENTRYPOINT ["/usr/local/bin/fanlin"]

FROM golang:1-bookworm AS builder
RUN set -eux; \
  apt-get update -y; \
  apt-get install -y \
    ca-certificates \
    build-essential \
    make \
    libaom-dev \
    liblcms2-dev \
    libheif-dev
WORKDIR /usr/src/app
COPY . .
RUN make build

# https://github.com/GoogleContainerTools/distroless
# https://console.cloud.google.com/gcr/images/distroless/GLOBAL
FROM gcr.io/distroless/cc-debian12:nonroot-amd64
COPY --from=builder --chmod=644 /lib/x86_64-linux-gnu/libaom.so.*   /lib/x86_64-linux-gnu/
COPY --from=builder --chmod=644 /lib/x86_64-linux-gnu/liblcms2.so.* /lib/x86_64-linux-gnu/
COPY --from=builder --chmod=755 /usr/src/app/cmd/fanlin/server      /usr/local/bin/fanlin
ENTRYPOINT ["/usr/local/bin/fanlin"]

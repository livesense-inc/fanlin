FROM debian:bookworm AS heif
ARG VERSION=1.19.7
RUN set -eux; \
  apt-get update; \
  apt-get install -y --no-install-recommends \
    ca-certificates \
    build-essential \
    libde265-dev \
    wget \
    cmake \
    ; \
  wget https://github.com/strukturag/libheif/releases/download/v$VERSION/libheif-$VERSION.tar.gz; \
  tar zxvf libheif-$VERSION.tar.gz; \
  cd libheif-$VERSION; \
  mkdir build; \
  cd build; \
  cmake --preset=release ..; \
  make; \
  make install

FROM golang:1-bookworm AS builder
RUN set -eux; \
  apt-get update -y; \
  apt-get install -y --no-install-recommends \
    ca-certificates \
    build-essential \
    make \
    libaom-dev \
    liblcms2-dev
WORKDIR /usr/src/app
COPY . .
COPY --from=heif --chmod=644 /usr/local/lib/pkgconfig/libheif.pc /lib/x86_64-linux-gnu/pkgconfig/
COPY --from=heif --chmod=644 /usr/local/include/libheif/*        /usr/local/include/libheif/
COPY --from=heif --chmod=644 /usr/local/lib/libheif.so*          /usr/local/lib/
RUN make build

# https://github.com/GoogleContainerTools/distroless
# https://console.cloud.google.com/gcr/images/distroless/GLOBAL
FROM gcr.io/distroless/cc-debian12:nonroot-amd64
COPY --from=builder --chmod=755 /usr/src/app/cmd/fanlin/server         /usr/local/bin/fanlin
COPY --from=builder --chmod=644 /lib/x86_64-linux-gnu/libaom.so.*      /lib/x86_64-linux-gnu/
COPY --from=builder --chmod=644 /lib/x86_64-linux-gnu/liblcms2.so.*    /lib/x86_64-linux-gnu/
COPY --from=builder --chmod=644 /lib/x86_64-linux-gnu/libsharpyuv.so.* /lib/x86_64-linux-gnu/
COPY --from=heif    --chmod=644 /lib/x86_64-linux-gnu/libde265.so*     /lib/x86_64-linux-gnu/
COPY --from=heif    --chmod=644 /usr/local/lib/libheif.so*             /lib/x86_64-linux-gnu/
COPY --from=heif    --chmod=644 /usr/local/lib/libheif/*               /lib/x86_64-linux-gnu/libheif/plugins/
ENV LIBHEIF_PLUGIN_PATH=/lib/x86_64-linux-gnu/libheif/plugins
ENTRYPOINT ["/usr/local/bin/fanlin"]

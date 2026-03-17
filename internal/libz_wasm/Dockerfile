FROM debian:trixie AS base
ARG LIBZ_REV=4b22c6e013fe0ca21a7e7bc4f0661af42fb463e3
ARG LIBZ_URL=https://gitlab.com/sortix/libz.git

ARG WASI_SDK_VERSION=30
ARG BINARYEN_VERSION=127

ARG WASI_SDK_CHECKSUM=sha256:0507679dff16814b74516cd969a9b16d2ced1347388024bc7966264648c78bfb
ARG BINARYEN_CHECKSUM=sha256:c8ebe5d00c978601086bcad33b2c80fcfe33d6a8b87b754ba4ea86a9a16cc316

ARG WASI_SDK_URL=https://github.com/WebAssembly/wasi-sdk/releases/download/wasi-sdk-${WASI_SDK_VERSION}/wasi-sdk-${WASI_SDK_VERSION}.0-x86_64-linux.tar.gz
ARG BINARYEN_URL=https://github.com/WebAssembly/binaryen/releases/download/version_${BINARYEN_VERSION}/binaryen-version_${BINARYEN_VERSION}-x86_64-linux.tar.gz

FROM base AS wasi-sdk
ADD --checksum=${WASI_SDK_CHECKSUM} ${WASI_SDK_URL} ./wasi-sdk.tar.gz
RUN mkdir /opt/wasi-sdk && tar -C /opt/wasi-sdk --strip-components=1 -xf /wasi-sdk.tar.gz

FROM base AS binaryen
ADD --checksum=${BINARYEN_CHECKSUM} ${BINARYEN_URL} ./binaryen.tar.gz
RUN mkdir /opt/binaryen && tar -C /opt/binaryen --strip-components=1 -xf binaryen.tar.gz

FROM golang:1.25 AS wasm2go
RUN GOBIN=/ go install -trimpath -v github.com/ncruces/wasm2go@d996851b98c1199ffb85db3df9a144d9e96f4d52

FROM base AS build
ADD ${LIBZ_URL}#${LIBZ_REV} /src
COPY --link ./exports.txt /src/
WORKDIR /src

COPY --link --from=wasi-sdk /opt /opt
COPY --link --from=binaryen /opt /opt
ENV PATH="/opt/wasi-sdk/bin:/opt/binaryen/bin:$PATH"

RUN wasm32-wasi-clang -c -Wall -Wextra -O2 \
 -mmutable-globals -mbulk-memory -mmultivalue -mnontrapping-fptoint -mreference-types -msign-ext \
 -std=c11 -I. -DZ_INSIDE_LIBZ -D_GNU_SOURCE -Wno-incompatible-pointer-types-discards-qualifiers \
 adler32.c \
 compress.c \
 crc32.c \
 deflate.c \
 gzclose.c \
 gzlib.c \
 gzread.c \
 gzwrite.c \
 infback.c \
 inffast.c \
 inflate.c \
 inftrees.c \
 trees.c \
 uncompr.c \
 zutil.c

RUN wasm32-wasi-clang -o libz.wasm -mexec-model=reactor -Wl,--soname=libz_wasm -Wl,--global-base=1024 $(printf " -Wl,--export=%s" $(cat exports.txt)) *.o
RUN wasm-ctor-eval -g -c _initialize libz.wasm -o libz.tmp
RUN wasm-opt -g libz.tmp -o libz.wasm -O2 \
 --enable-mutable-globals --enable-bulk-memory --enable-multivalue --enable-nontrapping-float-to-int --enable-reference-types --enable-sign-ext \
 --low-memory-unused

COPY --link --from=wasm2go /wasm2go /usr/local/bin/
RUN wasm2go -endian=little < libz.wasm > libz.go
RUN wasm2go -endian=big < libz.wasm > libz_big.go

FROM scratch
COPY --link --from=build /src/libz.go /
COPY --link --from=build /src/libz_big.go /

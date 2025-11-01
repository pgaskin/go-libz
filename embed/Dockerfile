# output is reproducible
# docker build --platform amd64 --pull --no-cache --progress plain --output . .

FROM debian:trixie AS base
ARG LIBZ_REV=4b22c6e013fe0ca21a7e7bc4f0661af42fb463e3
ARG LIBZ_URL=https://gitlab.com/sortix/libz.git

ARG WASI_SDK_VERSION=27
ARG BINARYEN_VERSION=124

ARG WASI_SDK_CHECKSUM=sha256:b7d4d944c88503e4f21d84af07ac293e3440b1b6210bfd7fe78e0afd92c23bc2
ARG BINARYEN_CHECKSUM=sha256:0290c3779fedf592b8da0ded3032ff55c41a2b7bfa2d6bf7b7bac6f0e6e28963

ARG WASI_SDK_URL=https://github.com/WebAssembly/wasi-sdk/releases/download/wasi-sdk-${WASI_SDK_VERSION}/wasi-sdk-${WASI_SDK_VERSION}.0-x86_64-linux.tar.gz
ARG BINARYEN_URL=https://github.com/WebAssembly/binaryen/releases/download/version_${BINARYEN_VERSION}/binaryen-version_${BINARYEN_VERSION}-x86_64-linux.tar.gz

FROM base AS wasi-sdk
ADD --checksum=${WASI_SDK_CHECKSUM} ${WASI_SDK_URL} ./wasi-sdk.tar.gz
RUN mkdir /opt/wasi-sdk && tar -C /opt/wasi-sdk --strip-components=1 -xf /wasi-sdk.tar.gz

FROM base AS binaryen
ADD --checksum=${BINARYEN_CHECKSUM} ${BINARYEN_URL} ./binaryen.tar.gz
RUN mkdir /opt/binaryen && tar -C /opt/binaryen --strip-components=1 -xf binaryen.tar.gz

FROM base AS build
ADD ${LIBZ_URL}#${LIBZ_REV} /src
COPY --link ./exports.txt /src/
WORKDIR /src

COPY --link --from=wasi-sdk /opt /opt
COPY --link --from=binaryen /opt /opt
ENV PATH="/opt/wasi-sdk/bin:/opt/binaryen/bin:$PATH"

RUN wasm32-wasi-clang -o libz.wasm -mexec-model=reactor -Wall -Wextra -O2 \
 -mmutable-globals -mbulk-memory -mmultivalue -mnontrapping-fptoint -mreference-types -msign-ext -msimd128 \
 -std=c11 -I. -DZ_INSIDE_LIBZ -D_GNU_SOURCE -Wno-incompatible-pointer-types-discards-qualifiers \
 $(printf " -Wl,--export=%s" $(cat exports.txt)) \
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

RUN wasm-ctor-eval -g -c _initialize libz.wasm -o libz.tmp

RUN wasm-opt -g libz.tmp -o libz.wasm -O2 \
 --enable-mutable-globals --enable-bulk-memory --enable-multivalue --enable-nontrapping-float-to-int --enable-reference-types --enable-sign-ext --enable-simd \
 --low-memory-unused --strip-debug --strip-dwarf

RUN chmod 644 libz.wasm

FROM scratch
COPY --link --from=build /src/libz.wasm /

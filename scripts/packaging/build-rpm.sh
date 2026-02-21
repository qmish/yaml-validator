#!/bin/bash
set -e
VERSION=${1:-v1.0.0}
VER=${VERSION#v}
BIN=yaml-validator-${VERSION}-linux-amd64
mkdir -p dist
RELEASE_DIR=docs/release/release-${VERSION}
if [ -f "${RELEASE_DIR}/${BIN}" ]; then cp "${RELEASE_DIR}/${BIN}" /tmp/yaml-validator; else GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X yaml-validator/cmd.version=${VERSION}" -o /tmp/yaml-validator .; fi
chmod +x /tmp/yaml-validator
if command -v fpm &>/dev/null; then fpm -t rpm -s dir -n yaml-validator -v "${VER}" --description "YAML validation tool" --url https://github.com/qmish/yaml-validator --prefix /usr/bin /tmp/yaml-validator; mv yaml-validator-*.rpm dist/ 2>/dev/null; echo "Created dist/yaml-validator-${VER}-1.x86_64.rpm"; else echo "fpm not found"; exit 1; fi

#!/bin/bash
#
# require:
# $ go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
# $ go get -u google.golang.org/grpc

__dirname=$(cd $(dirname $0); pwd)

cd $__dirname
cd ../echo

protoc \
    --plugin=protoc-gen-go=${GOPATH}/bin/protoc-gen-go \
    --plugin=protoc-gen-ts=../node_modules/.bin/protoc-gen-js_service \
    --go_out=plugins=grpc:. \
    --js_out=import_style=commonjs,binary:. \
    --ts_out=service=true:. \
    ./echo.proto

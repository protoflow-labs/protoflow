#!/bin/bash

PROTOC_GEN_ES_PATH="./node_modules/@bufbuild/protoc-gen-es/bin/protoc-gen-es"

protos=$(find ./proto -type f)

for proto in $protos; do
    protoc \
        --go_out=./ \
        --go-grpc_out=./ \
        --plugin=protoc-gen-es=${PROTOC_GEN_ES_PATH} \
        --es_out=./src/rpc \
        -I./proto $proto
done



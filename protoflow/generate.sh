#!/bin/bash

protos=$(find ./proto -type f)

for proto in $protos; do
    protoc --go_out=./ --twirp_out=./ -I./proto $proto
done



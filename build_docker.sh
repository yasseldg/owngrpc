#!/bin/bash

# use 0 for less information or 1 for full information
export DOCKER_BUILDKIT=1

# build docker file:

docker build --rm -t grpc_client:0.0.1 -f client/Dockerfile .

docker build --rm -t grpc_server:0.0.1 -f server/Dockerfile .
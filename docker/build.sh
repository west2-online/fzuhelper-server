#!/bin/bash

# This script will log in to Aliyun Container Registry and push the image up.
# You first need to obtain the password for the account, otherwise you cannot push.

set -e

VERSION="latest"
DIR=$(pwd)

docker_image="registry.cn-hangzhou.aliyuncs.com/west2-online/fzuhelper-server:$VERSION"

# login, only need to login once, then you can comment out this sentence.
docker login registry.cn-hangzhou.aliyuncs.com --username=west2gold@aliyun.com

# push image
docker buildx build --platform linux/amd64 -t $docker_image -f Dockerfile ../. --push
# Copyright 2024 The west2-online Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

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
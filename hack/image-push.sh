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
# 该脚本将在 github/workflow 中加载并推送镜像，可以传入多个服务名
# 镜像来自 build job 导出的 tar，由 workflow 下载为 artifact
set -e

services=("$@")
if [ ${#services[@]} -eq 0 ]; then
    echo "Usage: $0 <service> [service...]"
    exit 1
fi

# 先加载全部镜像，避免推送到一半才发现镜像有问题
for service in "${services[@]}"; do
    make load-image-$service
done

# 再逐个推送，同时 push 会出问题，见 Makefile 中 push-% 的说明
for service in "${services[@]}"; do
    echo "Pushing image for service $service..."
    make push-image-$service
done

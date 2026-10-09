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
# 该脚本将在 github/workflow 中构建项目镜像，并导出为 tar
# 导出的 tar 由 workflow 上传为 artifact，交由后续的 push job 推送
set -e

service="$1"
echo "Building image for service $service..."
make build-image-$service
make save-image-$service

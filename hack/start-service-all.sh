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
IMAGE_NAME="registry.cn-hangzhou.aliyuncs.com/west2-online/fzuhelper-server"

#该脚本负责启动服务，而其他相关的组件（如etcd）则在docker-compose.yml中启动
#只适用于单机部署
SERVICES=(api classroom user)

for service in "${SERVICES[@]}"; do
    echo "Starting service $service..."
    # 根据 service 的名称做不同处理
    case $service in
        "api")
            # 启动 api 服务，暴露端口 20001，使用指定网络 fzu-helper
            docker run -d --name api --network fzu-helper -p 20001:20001 --restart always api_image
            ;;
        "classroom")
            docker run -d --name classroom --network fzu-helper -p 20002:20002 --restart always classroom_image
            ;;
        "user")
            docker run -d --name user --network fzu-helper -p 20003:20003  --restart always user_image
            ;;
    esac
done

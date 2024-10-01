#!/bin/bash

#该脚本负责启动服务，而其他相关的组件（如etcd）则在docker-compose.yml中启动
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

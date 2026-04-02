# 第一阶段：构建应用程序
FROM golang:1.26.1 AS builder
# 定义构建参数，在构建镜像的时候需要传入 SERVICE 来表示构建的模块
ARG SERVICE
# CI=true 时（如 GitHub Actions）不替换软件源，否则使用国内镜像加速
ARG CI

# 设置环境变量
ENV TZ=Asia/Shanghai
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
# 创建工作目录
RUN mkdir -p /app
RUN mkdir -p /app/output

WORKDIR /app

# 非 CI 环境下使用国内 GOPROXY 加速
RUN if [ -z "$CI" ]; then go env -w GOPROXY=https://goproxy.cn,direct; fi

# 复制所有文件到工作目录
COPY . .

# 下载依赖
RUN go mod tidy

# 编译应用程序
RUN bash ./docker/script/build.sh ${SERVICE}

# 第二阶段：创建最终运行环境
FROM alpine

ARG SERVICE
# CI=true 时（如 GitHub Actions）不替换软件源，否则使用国内镜像加速
ARG CI

# 环境变量
ENV TZ=Asia/Shanghai
ENV SERVICE=${SERVICE}
# ETCD_ADDR 默认值在 entrypoint 中设置，不在 Dockerfile 中提供

# 更新软件依赖，非 CI 环境下替换为国内镜像源加速
RUN if [ -z "$CI" ]; then sed -i 's#https\?://dl-cdn.alpinelinux.org/alpine#https://mirrors.tuna.tsinghua.edu.cn/alpine#g' /etc/apk/repositories; fi
RUN apk update --no-cache && apk --no-cache add ca-certificates tzdata bash

# 创建工作目录
WORKDIR /app

# 复制 Binary 及启动脚本，调整权限设置
COPY --from=builder /app/output /app/output
COPY --from=builder /app/docker/script/entrypoint.sh /app/entrypoint.sh
RUN mkdir -p /app/config # 确保有文件路径可以被挂载
RUN chmod +x /app/entrypoint.sh

# 需要注意的是，这个命令是在程序执行时运行，不能直接在这里填入 SERVICE 参，需要绕个弯设置为环境变量
CMD ["sh", "-c", "./entrypoint.sh"]

# 第一阶段：构建应用程序
FROM golang:1.22 AS builder

# 定义构建参数
ARG SERVICE

# 设置环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOPROXY=https://goproxy.cn,direct \
    GOOS=linux \
    GOARCH=amd64

# 创建工作目录
RUN mkdir -p /app


WORKDIR /app

# 复制所有文件到工作目录
COPY . .

# 下载依赖
RUN go mod tidy


# 编译应用程序
RUN cd ./cmd/${SERVICE} && go build -o /app/${SERVICE}

# 第二阶段：创建最终运行环境
FROM alpine

# 安装必要的依赖
RUN apk --no-cache add ca-certificates

# 定义构建参数
ARG SERVICE
# 创建工作目录
WORKDIR /app

# 从构建阶段复制应用程序二进制文件
COPY --from=builder /app/cmd/${SERVICE}/config /app/config
COPY --from=builder /app/${SERVICE} /app/${SERVICE}



# 第一阶段：构建应用程序
FROM golang:1.22 AS builder

# 定义构建参数
ARG SERVICE

# 设置环境变量
ENV TZ Asia/Shanghai
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOPROXY=https://goproxy.cn,direct
ENV GOOS=linux
ENV GOARCH=amd64


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
ENV TZ Asia/Shanghai
# 安装必要的依赖
RUN apk --no-cache add ca-certificates

# 定义构建参数
ARG SERVICE
# 创建工作目录
WORKDIR /app

# 从构建阶段复制应用程序二进制文件和配置文件
COPY --from=builder /app/config /app/config
COPY --from=builder /app/${SERVICE} /app/${SERVICE}
COPY --from=builder /app/docker/entrypoint.sh /app/docker-entrypoint.sh

RUN chmod +x /app/docker-entrypoint.sh

#因为SERVER好像只能在当前环境下被解析出来，如果想要将SERVER作为变量传入脚本，目前的方法是设置一个ENV来接收SERVER的值
ENV SERVER=${SERVICE}

CMD ["sh", "-c", "./docker-entrypoint.sh ${SERVER}"]



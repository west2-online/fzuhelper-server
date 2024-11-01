# Deploy

这份指南将帮助您从零开始部署本项目。只需按照本指南的步骤操作，您将逐步了解如何设置必要的环境、编译并运行服务，以及启动整个系统。本项目包含多个服务组件，通过 Docker 以及简单的脚本来完成构建和启动，让您轻松上手。

在部署过程中，您将使用 Makefile 中的命令来启动依赖的环境（包括 MySQL、etcd、Redis 等数据库和缓存服务），然后编译并运行特定的服务。每个服务都将自动获取并读取配置信息，并注册到 etcd 以实现服务管理。这一切都简化为几个命令操作，极大降低了部署的复杂度。

如果对具体的构建和启动流程感兴趣，请查阅 [Build](./build.md) 页面。

## 先决条件

需要安装的软件包：

- Docker
  - 请参考 [Docker 官方文档](https://docs.docker.com/get-docker/) 安装 Docker。
  - 用于启动 MySQL、etcd、Redis 等服务。
- tmux
  - 便于在终端中管理多个会话，提高运维效率。
  - Makefile 中的命令会默认在 tmux 中执行。

## 本地部署

- 修改 `config/config.yaml` 的配置，将数据库等配置的 ip 修改为 `localhost`（如果没有请新增这个文件）。配置示例请参考 `config.example.yaml`。

### 启动环境

#### 清理本地环境（可选）

```shell
make clean-all
```

#### kafka 环境准备（可选）

```shell
sh docker/script/generate_certs_for_kafka.sh
```

#### 启动环境基础容器（数据库等）

```shell
make env-up
```

### 启动特定服务

```shell
make <target>
```

`<target>` 为服务名称，即 `cmd` 目录下的文件夹名称。

## 服务器部署

_此部分供 west2-online 成员参考_

先决条件: 服务器上已经安装并启动 Docker 服务。

### 登录阿里云容器镜像服务

在开发机和服务器上都要登录阿里云的容器镜像仓库。

```shell
docker login --username=west2gold@aliyun.com registry.cn-hangzhou.aliyuncs.com
```

密码请参阅 [飞书文档](https://west2-online.feishu.cn/wiki/Bnvhw9adGizcOFk5jRccOWnbn1g)。

### 上传镜像

现在切换到开发机上，cd 到项目的根目录下，然后执行以下命令：

```shell
make push-<target> # make push-api
```

开始运行命令后需要再次输入 `<target>` 作为确认，以防误操作。

### 拷贝相关文件

需要拷贝的文件有：

-. `./docker/docker-compose.yml`
-. `./docker/script/etcd-monitor.sh`
-. `./docker/env/`
-. `./config/config.yaml`
-. `./hack/image-refresh.sh`
-. `./hack/docker-run.sh`

需要将以上文件传送到云服务器的用户主目录（`~/`）下，可以使用 rsync 进行这项操作。

完成后的目录结构应该与下面的结构类似：

```shell
├── docker
│   ├── script
│   │   └── etcd-monitor.sh
│   ├── env
│   │   ├── redis.env
│   │   ├── mysql.env
│   │   └── etcd.env
│   ├── docker-compose.yaml
├── config
│   ├── sql
│   │   ├── init.sql
│   └── config.yaml
├── hack
│   ├── image-refresh.sh
│   └── docker-run.sh
```

### 启动容器

启动 `<target>` 服务容器：

```shell
bash ./hack/docker-run.sh <target>
```

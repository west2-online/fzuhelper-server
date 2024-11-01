# Deploy
这份指南将帮助您从零开始部署本项目。只需按照本指南的步骤操作，您将逐步了解如何设置必要的环境、编译并运行服务，以及启动整个系统。本项目包含多个服务组件，通过 Docker 以及简单的脚本来完成构建和启动，让您轻松上手。

在部署过程中，您将使用 Makefile 中的命令来启动依赖的环境（包括 MySQL、etcd、Redis 等数据库和缓存服务），然后编译并运行特定的服务。每个服务都将自动获取并读取配置信息，并注册到 etcd 以实现服务管理。这一切都简化为几个命令操作，极大降低了部署的复杂度。
如果对具体的构建和启动流程感兴趣, 可参考 [build](./build.md)
- 需要你有 **docker**, 用于启动 MySQL、etcd、Redis 等服务, 怎么设置 docker 环境这里不再赘述

## 本地部署

- 修改 `./config/config.yaml` 的配置，将数据库等配置的 ip 修改为 `localhost`（如果没有请新增这个文件）。配置示例请参考 `config.example.yaml`。
- 请在 `tmux` 环境下执行命令
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
- `<target>` 为服务名称，即 `cmd` 目录下的文件夹名称。

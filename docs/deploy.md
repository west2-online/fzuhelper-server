# Deploy

请通过 `make help` 命令来查阅可用的命令列表。

后续的 `<target>` 指代某个服务，例如 `api`。

## 本地部署

修改 `./config/config.yaml` 的配置，将数据库等配置的 ip 修改为 `localhost`（如果没有请新增这个文件）。配置示例请参考 `config.example.yaml`。

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

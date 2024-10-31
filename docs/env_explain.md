# 环境变量解释

## kafka

文件位置: `docker/env/kafka.env`

### 基础配置

- **`BITNAMI_DEBUG=true`**

  > 启用 Bitnami 镜像的调试模式。这会输出更多的调试信息，帮助你跟踪启动和运行过程中的问题。

- **`KAFKA_CFG_PROCESS_ROLES=broker,controller`**

  > 指定 Kafka 节点的角色。在 Kafka 2.8 及以后的版本中，Kafka 引入了单独的控制器角色。此配置表示该节点既是一个 Broker（消息代理），也是一个 Controller（集群管理控制器）。

- **`KAFKA_CFG_NODE_ID=1`**

  > Kafka 节点的唯一标识符，用于区分集群中的不同节点。这里设置为 `1`，表示这是第一个节点。

- **`KAFKA_CFG_CONTROLLER_LISTENER_NAMES=CONTROLLER`**

  > 指定用于控制器之间通信的监听器名称。控制器是管理集群元数据的角色，`CONTROLLER` 是其用于监听控制器间通信的名称。

- **`ALLOW_PLAINTEXT_LISTENER=no`**

  > 禁止 Kafka 通过明文（非加密）通信。这个设置确保所有的监听器都使用加密或安全协议（如 SASL_PLAINTEXT）。

- **`KAFKA_BROKER_ID=1`**

  > Kafka Broker 的 ID，类似于 `KAFKA_CFG_NODE_ID`，用于标识当前的 Kafka Broker。在集群中每个 Broker 都应该有唯一的 ID。

- **`KAFKA_CFG_LISTENERS=INTERNAL://:9094,CLIENT://:9095,CONTROLLER://:9093,EXTERNAL://:9092`**

  > 配置 Kafka Broker 的监听地址。不同的端口用于不同的角色或用途：
  >
  > - `INTERNAL://:9094` 用于内部通信（Broker 间的通信）。
  > - `CLIENT://:9095` 用于客户端连接。
  > - `CONTROLLER://:9093` 用于控制器的通信。
  > - `EXTERNAL://:9092` 用于外部通信（如外部系统连接 Kafka）。

- **`KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP=INTERNAL:PLAINTEXT,CLIENT:SASL_PLAINTEXT,CONTROLLER:PLAINTEXT,EXTERNAL:SASL_PLAINTEXT`**

  > 定义每个监听器的安全协议：
  >
  > - `INTERNAL` 使用 `PLAINTEXT`（不加密的明文通信）。
  > - `CLIENT` 和 `EXTERNAL` 使用 `SASL_PLAINTEXT`，即带有简单认证和安全层的明文通信。
  > - `CONTROLLER` 使用 `PLAINTEXT`。

- **`KAFKA_CFG_ADVERTISED_LISTENERS=INTERNAL://kafka:9094,CLIENT://:9095,EXTERNAL://127.0.0.1:9092`**

  > 定义 Kafka 对外公布的地址和端口。客户端和其他 Broker 将通过这些地址连接：
  >
  > - `INTERNAL` 的地址为 `kafka:9094`，用于内部通信。
  > - `CLIENT` 地址为空，表示在默认网络接口上开放。
  > - `EXTERNAL` 对外公布为 `127.0.0.1:9092`，供外部连接。

- **`KAFKA_CFG_SSL_ENDPOINT_IDENTIFICATION_ALGORITHM=`**

  > 指定 SSL 连接中使用的端点身份验证算法。空值表示禁用该功能。

- **`KAFKA_CFG_INTER_BROKER_LISTENER_NAME=INTERNAL`**

  > 指定用于 Broker 间通信的监听器名称，设置为 `INTERNAL`，对应 `KAFKA_CFG_LISTENERS` 中的 `INTERNAL://:9094`。

- **`KAFKA_CFG_SASL_MECHANISM_INTER_BROKER_PROTOCOL=PLAIN`**

  > 配置 Broker 间使用的 SASL 认证机制，这里使用的是 `PLAIN`（简单用户名/密码认证机制）。

- **`KAFKA_CFG_SASL_ENABLED_MECHANISMS=PLAIN`**

  > 允许的 SASL 机制，设置为 `PLAIN`，表示使用用户名/密码认证。

- **`KAFKA_CFG_LOG_SEGMENT_BYTES=536870912`**

  > 配置日志分段的大小限制，这里设置为 `512MB`。当一个日志分段文件达到这个大小时，Kafka 会滚动生成新的日志分段。

- **`KAFKA_CFG_LOG_RETENTION_HOURS=6`**

  > 配置 Kafka 日志保留时间，单位是小时。日志将在 6 小时后被删除。

- **`KAFKA_CFG_LOG_FLUSH_INTERVAL_MS=1000`**

  > 配置日志刷新的间隔时间，单位是毫秒。这里设置为每 1000 毫秒（1 秒）将日志刷写到磁盘。

- **`KAFKA_TLS_TYPE=JKS`**

  > 指定 Kafka TLS/SSL 证书的格式类型，这里使用 `JKS`（Java KeyStore）格式。

- **`KAFKA_CLIENT_USERS=fzuhelper`**

  > 定义 Kafka 客户端用户的名称，设置为 `fzuhelper`。

- **`KAFKA_INTER_BROKER_USER=fzuhelper`**

  > 配置用于 Broker 间通信的用户，设置为 `fzuhelper`。

- **`KAFKA_CFG_CONTROLLER_QUORUM_VOTERS=1@127.0.0.1:9093`**
  > 配置控制器选举中的投票者信息，表示节点 1 位于 `127.0.0.1:9093`，参与控制器的选举。

### 密码相关配置

- **`KAFKA_CLIENT_PASSWORDS=fzuhelper-password`**

  > 定义 Kafka 客户端用户 `fzuhelper` 的密码。

- **`KAFKA_CERTIFICATE_PASSWORD=fzuhelper-password`**

  > 定义 Kafka TLS/SSL 证书的密码。

- **`KAFKA_INTER_BROKER_PASSWORD=fzuhelper-password`**
  > 配置用于 Broker 间通信的密码，和用户 `fzuhelper` 相关。

## kafka-ui

文件位置: `docker/env/kafka-ui.env`

### 基础配置

- **`KAFKA_CLUSTERS_0_NAME=fzuhelper`**

  > kafka 集群的名字

- **`KAFKA_CLUSTERS_0_BOOTSTRAPSERVERS=kafka:9094`**

  > 访问 broker 的地址，这里通过 docker 网络解析到 kafka

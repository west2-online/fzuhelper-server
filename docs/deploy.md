# Deploy
请通过`make help`来显示可用的服务列表，后续的target指代我们的服务，例如api
## 本地部署
修改`./config/config.yaml`的配置，将数据库等配置的ip修改为localhost（如果没有请新增这个文件), 具体参考config.example.yaml
### 启动环境
#### 清理本地环境(optional)
```shell
make clean-all
```
#### kafka环境准备(optional)
```shell
sh docker/script/generate_certs_for_kafka.sh
```
#### 启动环境
```shell
make env-up
```
### 启动特定服务
```shell
make target #e.g. make api
```
使用make help获取更多信息

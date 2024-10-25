# Deploy
target 请通过`make help`来显示可用的服务列表，后续的target指代我们的服务，例如api
## 本地部署
修改`./config/config.yaml`的配置，将数据库等配置的ip修改为localhost（如果没有请新增这个文件）
### 启动环境
#### 清理本地环境(optional)
```shell
make clean-all
```
#### kafka环境准备
```shell
sh docker/script/generate_certs_for_kafka.sh
```
#### 启动环境
```shell
make env-up
```
### 启动服务
#### 启动所有服务
> 可以使用"ctrl+b s"来切换终端
```shell
make local
```
#### 启动特定服务
```shell
make target #e.g. make api
```
使用make help获取更多信息
## 云服务部署
> 请保证已经使用docker login

### 构建镜像
```shell
make push-target 
```
### 云服务器端

#### 环境准备与加密 (不需要kafka可以跳过)
首先你需要更改一些密码来防止泄露。最好更改为同样的密码
> `config/config.yaml` 中kafka的password部分  
> `docker/env/kafka.env` 中的三个密码，已经特意标注出来了  
> `docker/script/generate_certs_for_kafka.sh` 中的密码，已经标注出来  

其次你需要执行这个脚本来生成kafka所需要的证书与密钥
```shell
sh docker/script/generate_certs_for_kafka.sh
```

#### 环境搭建
```shell
docker compose up -d
```

#### 部署服务
```shell
sh image-refresh.sh target #更新镜像
sh docker-run.sh target #运行容器
```

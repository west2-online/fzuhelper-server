# fzuhelper-server

- target 统一为api、classroom、user
### 本地部署
- 修改./config/config.yaml的配置，将数据库等配置的ip修改为localhost
1. 启动环境
- 清理本地环境(optional)
```shell
make clean-all
```
- 启动环境
```shell
make env-up
```
2. 启动服务
- 启动所有服务
- 可以使用"ctrl+b s"来切换终端
```shell
make local
```
- 启动特定服务
```shell
make target #e.g. make api
```

使用make help获取更多信息

### 云服务部署
##### 请保证已经使用docker login
1. 构建镜像
```shell
make push-target 
```

2. 云服务器端
- 环境搭建
```shell
docker compose up -d
```
- 部署服务
```shell
sh image-refresh.sh target #更新镜像
sh docker-run.sh target #运行容器
```



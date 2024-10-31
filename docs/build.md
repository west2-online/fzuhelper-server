# Build

## 定义
- target： 具体的服务名，例如：api，classroom，course 等
## 大致流程
1. 使用命令```make env-up```启动环境（etcd，redis e.g.)
2. ```make target```编译并运行具体的 target
3. target 从 etcd 中获取 config.yaml
4. 读取 config.yaml 中的配置，完成初始化

## 构建和启动
### 目录结构

项目的关键目录如下：

- `cmd/`：包含各服务模块的启动入口。
- `output/`：构建产物的输出目录。


### 构建流程
构建过程主要通过 `build.sh` 脚本完成，用于编译指定服务模块的二进制文件或进行系统测试。
1. 进入到 cmd 中的对应 target
2. go build 编译对应的二进制文件，并输出到 output 中
### output 结构
```shell
 output
 └── target
        └── binary
```

### 启动流程 
启动过程主要通过 `entrypoint.sh` 脚本完成。
1. 设置 etcd 的地址，为后续程序在***运行时***能够获取到 etcd 的地址并获取 config.yaml
2. cd 到构建阶段生成的 output 目录，执行对应的二进制文件


<div align="center">
  <h1 style="display: inline-block; vertical-align: middle;">fzuhelper-server</h1>
</div>

<div align="center">
  <a href="/README.md">English</a> | <a href="#overview">中文</a>
</div>

## <a id="overview"></a>概述

fzuhelper-server 是基于分布式架构的 fzuhelper 服务器应用程序，自 2024 年以来一直在使用，每天为 **超过 23,000 名** 福州大学的学生提供服务（[数据来源及 fzuhelper 介绍](https://west2-online.feishu.cn/wiki/RG3UwWGqPig8lHk0mYsccKWRnrd)）。

该项目侧重于业务实现。如果你想了解我们如何与教务系统对接，可以查看我们开源版本的项目：[west2-online/jwch](https://github.com/west2-online/jwch)。

> fzuhelper 于 2015 年上线，由 west2-online 从零开发并持续运营，尽可能为校内学生提供工业级实践机会，并为学生就业提供有力支持。

## 功能特点

- **云原生**：采用原生 Golang 分布式架构设计，基于字节跳动的最佳实践。
- **高性能**：支持异步 RPC、非阻塞 I/O、共享内存通信和即时编译（JIT）。
- **可扩展性**：模块化、分层的结构设计，代码清晰易读，降低了开发难度。
- **DevOps**：丰富的脚本和工具减少了不必要的手动操作，简化了使用和部署流程。

## 项目结构

```bash
│  .golangci.yml              # GolangCI 配置文件
│  .licenseignore             
│  go.mod                     
│  go.sum                     
│  LICENSE                    
│  Makefile                   # 一些 make 命令
│  README.md
├── api                       # gateway
├── cmd                       # 各个微服务的启动入口
├── config                    # 配置文件和配置示例
├── docker                    # Docker 构建配置
├── docs
├── hack                      # 用于自动化开发、构建和部署任务的工具
├── idl                       # 接口定义
├── internal                  # 各个微服务的实现
├── kitex_gen                 # Kitex 生成的代码
├── pkg                      
│   ├── base                  # 通用基础服务
│   │      ├── client         # 对应组件(redis, mysql e.g.)的客户端
│   ├── cache                 # 缓存服务
│   ├── db                    # 数据库服务
│   ├── constants/            # 存储常量
│   ├── errno/                # 自定义错误
│   ├── eshook                # elasticsearch hook
│   ├── logger/               # 日志系统
│   ├── tracer/               # 用于 Jaeger 的追踪器
│   └── utils/                # 实用函数
```

## 快速启动和部署

由于我们编写的脚本，流程已大大简化。你只需使用以下命令，即可快速启动环境并运行程序。

详情请访问：[部署文档](deploy.md)

## 架构

<img src="/docs/img/architecture.svg">

## 贡献者

<img src="/docs/img/logo(en).svg" width="400">

如果你有兴趣参与 fzuhelper-server 的维护工作，请访问我们的 [官方网站](https://site.west2.online) 联系我们。

## 许可协议
fzuhelper-server 采用 Apache 2.0 许可协议。详情请参阅 LICENSE 文件。

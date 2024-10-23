<div align="center">
  <h1 style="display: inline-block; vertical-align: middle;">fzuhelper-server</h1>
</div>

<div align="center">
  <a href="#overview">English</a> | <a href="docs/README.zh.md">中文</a>
</div>

## <a id="overview"></a>Overview

fzuhelper-server is a distributed architecture-based server application for fzuhelper, which has been in use since 2024, serving **over 23,000 students** from Fuzhou University everyday ([Data source and introduction to fzuhelper](https://west2-online.feishu.cn/wiki/RG3UwWGqPig8lHk0mYsccKWRnrd)).

This project focuses on business implementation. To see how we interface with the academic affairs office, you can check out our open-source version at [west2-online/jwch](https://github.com/west2-online/jwch).

> fzuhelper was launched in 2015, developed from scratch by west2-online and continuously operated, providing students with industrial-grade practice as much as possible on campus and offering robust support for student employment.

## Features

- **Cloud-Native**: native golang distributed architecture design, based on ByteDance's best practices.
- **High Performance**: Supports asynchronous RPC, non-blocking I/O, shared memory communication, and Just-In-Time (JIT) compilation.
- **Scalability**: Features a modular, layered structural design, with clear and readable code, reducing development difficulty.
- **DevOps**：An abundance of scripts and tools reduce unnecessary manual labor, simplifying usage and deployment.

## Project structure

```bash
│  .golangci.yml              # GolangCI configuration
│  .licenseignore             
│  go.mod                     
│  go.sum                     
│  LICENSE                    
│  Makefile                   # some useful commands
│  README.md                  
├── cmd                       # microservices
├── config                    # for run-directly config and config-example
├── docker                    # docker build configuration
├── docs
├── hack                      # tools for automating development, build, and deployment tasks.
├── idl                       # interface definition
├── kitex_gen                 # kitex generated code
├── pkg                      
│   ├── client/               # client side implementations
│   ├── constants/            # store any consts
│   ├── errno/                # custom error
│   ├── logger/               # logging system
│   ├── middleware/           # common middleware
│   ├── tracer/               # for jaeger
│   └── utils/                # useful funcs
```

## Quick start and deploy

Due to the script we have written, the process has been greatly simplified. You just need to use the following command to quickly start the environment and run the program in a containerized manner.

please visit: [deploy](docs/deploy.md)

## Architecture

<img src="/docs/img/architecture.svg">

## Contributors

<img src="/docs/img/logo(en).svg" width="400">

If you are interested in joining the maintenance of fzuhelper-server, please contact us on our [official website](https://site.west2.online)

## License
`fzuhelper-server` is under the Apache 2.0 license. See the LICENSE file for details.

<div align="center">
  <img src="/docs/img/appicon.png" width="80" style="vertical-align: middle; margin-right: 26px; margin-top: 5px;">
  <h1 style="display: inline-block; vertical-align: middle; font-size: 2em;">fzuhelper-server</h1>
</div>

## Overview

fzuhelper-server is a distributed architecture-based server application for fzuhelper, which has been in use since 2024, serving **over 23,000 students** from Fuzhou University everyday ([Data source and introduction to fzuhelper](https://west2-online.feishu.cn/wiki/RG3UwWGqPig8lHk0mYsccKWRnrd)).

This project focuses on business implementation. To see how we interface with the academic affairs office, you can check out our open-source version at [west2-online/jwch](https://github.com/west2-online/jwch).

> fzuhelper was launched in 2015, developed from scratch by west2-online and continuously operated, providing students with industrial-grade practice as much as possible on campus and offering robust support for student employment.

## Features

- **Cloud-Native**: native golang distributed architecture design, based on ByteDance's best practices.
- **High Performance**: Supports asynchronous RPC, non-blocking I/O, shared memory communication, and Just-In-Time (JIT) compilation.
- **Scalability**: Features a modular, layered structural design, with clear and readable code, reducing development difficulty.
- **DevOps**：An abundance of scripts and tools reduce unnecessary manual labor, simplifying usage and deployment.

## Architecture

<img src="/docs/img/architecture.svg">

## Contributors

<img src="/docs/img/logo(en).svg" width="400">

If you are interested in joining the maintenance of fzuhelper-server, please contact us on our [official website](https://site.west2.online)

## License
`fzuhelper-server` is under the Apache 2.0 license. See the LICENSE file for details.
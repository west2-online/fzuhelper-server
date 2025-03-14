# Copyright 2024 The west2-online Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# 该脚本适用于 Docker 及本地调试，作为执行程序前的 presetting（即 entrypoint）
# 请不要直接执行这个脚本，这个脚本应当由 Makefile/Dockerfile 接管

#! /usr/bin/env bash
CURDIR=$(pwd)

# 此处只涉及 Kitex，但是 Hertz 使用这个没有影响，保留即可
export KITEX_RUNTIME_ROOT=$CURDIR
export KITEX_LOG_DIR="$CURDIR/log"

if [ ! -d "$KITEX_LOG_DIR/app" ]; then
    mkdir -p "$KITEX_LOG_DIR/app"
fi

if [ ! -d "$KITEX_LOG_DIR/rpc" ]; then
    mkdir -p "$KITEX_LOG_DIR/rpc"
fi
# 设置是否是 kubernetes 环境, 默认是 k8s
: ${DEPLOY_ENV:="k8s"}
export DEPLOY_ENV

# 参数替换，检查 ETCD_ADDR 是否已经设置，没有将会设置默认值。如果非 k8s 环境，下面 ETCD_ADDR 将生效
: ${ETCD_ADDR:="localhost:2379"}
export ETCD_ADDR

# 这个 SERVICE 环境变量会自动地由 Dockerfile/Makefile 设置
exec "$CURDIR/output/$SERVICE/fzuhelper-$SERVICE"

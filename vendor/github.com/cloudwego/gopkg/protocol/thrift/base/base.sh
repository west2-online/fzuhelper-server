#! /bin/bash

# please install fastcodec version of thriftgo
# if the `feat-fastcodec` branch doesn't exist try to use `main` branch
# the feature is only for internal use now, but it will be released in the future when it's stable.
# go install github.com/cloudwego/thriftgo@feat-fastcodec

set -e

thriftgo -g fastgo:no_default_serdes=true,gen_setter=true -o=.. ./base.thrift
gofmt -w base.go
gofmt -w k-base.go

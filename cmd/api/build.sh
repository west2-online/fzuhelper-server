#!/bin/bash
RUN_NAME=hertz_service
mkdir -p output/bin
cp script/* output 2>/dev/null
chmod +x output/bootstrap.sh
#go build -toolexec="skywalking-go-agent" -o output/bin/${RUN_NAME} -a
go build -o output/bin/${RUN_NAME} -a

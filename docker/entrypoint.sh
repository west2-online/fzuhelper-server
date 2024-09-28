#!/bin/sh

# This file will serve as the entry point for the container.
# During the image building process, this file will be renamed to docker-entrypoint.sh.
# When the container starts, this shell script will be called to initiate the services within the container.

# THIS SCRIPT SHOULD NOT BE MANUALLY EXECUTED.

CONFIG_PATH="./config/config.yaml"

function read_key() {
    local key="$2"
    local flag=0
    while read -r LINE; do
        if [[ $flag == 0 ]]; then
            if [[ "$LINE" == *"$key:"* ]]; then
                if [[ "$LINE" == *" "* ]]; then
                    echo "$LINE" | awk -F " " '{print $2}'
                    return
                else
                    continue
                fi
            fi
        fi
    done < "$1"
}

export ETCD_ADDR=$(read_key $CONFIG_PATH "etcd-addr")

service=$1
./${service}

#!/bin/sh

# After we use "make xxx" to build a certain service, the makefile will automatically execute this script.

# THIS SCRIPT SHOULD NOT BE MANUALLY EXECUTED.

SERVICE=$1
OUTPUT_PATH="./output"
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


sh $OUTPUT_PATH/$SERVICE/bootstrap.sh
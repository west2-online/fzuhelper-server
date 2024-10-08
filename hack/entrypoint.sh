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

#!/bin/sh

# After we use "make xxx" to build a certain service, the Makefile will automatically execute this script.
# In short, this script serves for local.

# THIS SCRIPT SHOULD NOT BE MANUALLY EXECUTED.

SERVICE=$1 # service name
OUTPUT_PATH="./output" # related to project folder
CONFIG_PATH="./config/config.yaml" # related to project folder


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

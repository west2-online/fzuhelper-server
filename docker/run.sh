#!/bin/bash

# create folder
mkdir -p data/kibana
mkdir -p data/elasticsearch
mkdir -p data/jaeger
mkdir -p data/prometheus
mkdir -p data/grafana

mkdir -p data/mysql
mkdir -p data/redis
mkdir -p data/rabbitmq
mkdir -p data/etcd

IMAGE_NAME="fzuhelper"
SERVICE_TO_START=${1:-all} # default start all


DIR=$(cd $(dirname $0); pwd)

SERVICES=(api classroom user)

remove_container() {
    container_status=$(docker inspect -f '{{.State.Status}}' "$1")
    if [ "$container_status" == "running" ]; then
        echo "Stopping container $1..."
        docker stop "$1"
    elif [ "$container_status" == "paused" ]; then
        echo "Unpausing and then stopping container $1..."
        docker unpause "$1"
        docker stop "$1"
    fi
    echo "Remove container $1..."
    docker rm "$1"
}

start_container() {
    echo "Starting container for $1..."
    docker run -d --name "fzuhelper-$1" \
    -e service=$1 \
    --net=host \
    -v $DIR/config:/app/config \
    "$IMAGE_NAME"
}

containers_to_stop=$(docker ps -aq --filter "ancestor=$IMAGE_NAME")
if [ "$SERVICE_TO_START" == "all" ]; then
    for container_id in $containers_to_stop; do
        remove_container $container_id
    done
else
    for container_id in $containers_to_stop; do
        container_id=$(docker inspect -f '{{.Name}}' "$container_id")
        if [ "$container_id" != "/fzuhelper-$SERVICE_TO_START" ]; then
            continue
        fi
        remove_container $container_id
    done
fi

if [ "$SERVICE_TO_START" == "all" ]; then
    for service in "${SERVICES[@]}"; do
        start_container $service
    done
else
    start_container $SERVICE_TO_START
fi

#!/bin/bash

if [ $1 = "up" ]; then
    docker compose $@
fi

if [ $1 = "start" ]; then
    docker compose start
fi

if [ $1 = "stop" ]; then
    docker compose stop
fi

if [ $1 = "down" ]; then
    docker compose down
fi

if [ $1 = "go" ]; then
    docker compose exec app $@
fi

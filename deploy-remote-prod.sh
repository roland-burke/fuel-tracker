#!/bin/bash
REMOTE_CONTEXT_NAME=the-machine
ENV_FILE=docker/.env.prod
COMPOSE_FILE=docker/docker-compose.prod.yml

docker-compose --context $REMOTE_CONTEXT_NAME --env-file $ENV_FILE -f docker-compose.yml -f $COMPOSE_FILE up -d --build --force-recreate
# to "repair" the terminal
stty sane
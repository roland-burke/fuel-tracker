#!/bin/bash
REMOTE_CONTEXT_NAME=the-machine

docker-compose --context $REMOTE_CONTEXT_NAME --env-file docker/.env.prod -f docker-compose.yml -f docker/docker-compose.prod.yml down
# to "repair" the terminal
stty sane
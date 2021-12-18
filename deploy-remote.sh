#!/bin/bash

FILENAME=`basename "$0"`

PROD_REMOTE_CONTEXT_NAME=the-machine
PROD_ENV_FILE=docker/.env.prod
PROD_COMPOSE_FILE=docker/docker-compose.prod.yml

DEV_REMOTE_CONTEXT_NAME=the-machine
DEV_ENV_FILE=docker/.env.dev
DEV_COMPOSE_FILE=docker/docker-compose.dev.yml

run_docker_compose() 
{
    docker-compose --context $1 --env-file $2 -f docker-compose.yml -f $3 up -d --build --force-recreate
}

if [ "$1" = "prod" ]
then
	run_docker_compose() $PROD_REMOTE_CONTEXT_NAME $PROD_ENV_FILE $PROD_COMPOSE_FILE
elif [ "$1" = "dev" ]
then
	run_docker_compose() $DEV_REMOTE_CONTEXT_NAME $DEV_ENV_FILE $DEV_COMPOSE_FILE
else
	echo "Usage: $FILENAME <prod|dev>"
	exit 0
fi

# to "repair" the terminal
stty sane
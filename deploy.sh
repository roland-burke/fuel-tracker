#!/bin/bash

FILENAME=`basename "$0"`

REMOTE_CONTEXT_NAME=the-machine

PROD_ENV_FILE=docker/.env.prod
PROD_COMPOSE_FILE=docker/docker-compose.prod.yml

DEV_ENV_FILE=docker/.env.dev
DEV_COMPOSE_FILE=docker/docker-compose.dev.yml

run_docker_compose_remote() 
{
    docker-compose --context $1 --env-file $2 -f docker-compose.yml -f $3 up -d --build --force-recreate
}

run_docker_compose_local() 
{
	echo $1
    docker-compose --env-file $1 -f docker-compose.yml -f $2 up -d --build --force-recreate
}

if [[ "$2" = "clean" ]]
then
	echo "clean local:"
	rm -rf ../.fuel-tracker-db
fi

# === REMOTE ===
if [[ "$1" = "prod" ]]
then
	run_docker_compose_remote $REMOTE_CONTEXT_NAME $PROD_ENV_FILE $PROD_COMPOSE_FILE
elif [[ "$1" = "dev" ]]
# === LOCAL ===
then
	run_docker_compose_local $DEV_ENV_FILE $DEV_COMPOSE_FILE
else
	echo "Usage: $FILENAME <prod|dev>"
	exit 0
fi

# to "repair" the terminal
stty sane
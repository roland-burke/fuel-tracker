#!/bin/bash

FILENAME=`basename "$0"`

if [ "$1" = "prod" ]
then
	REMOTE_CONTEXT_NAME=the-machine
	ENV_FILE=docker/.env.prod
	COMPOSE_FILE=docker/docker-compose.prod.yml
elif [ "$1" = "dev" ]
then
	REMOTE_CONTEXT_NAME=the-machine
	ENV_FILE=docker/.env.dev
	COMPOSE_FILE=docker/docker-compose.dev.yml
else
	echo "Usage: $FILENAME <prod|dev>"
	exit 0
fi

docker-compose --context $REMOTE_CONTEXT_NAME --env-file $ENV_FILE -f docker-compose.yml -f $COMPOSE_FILE up -d --build --force-recreate
# to "repair" the terminal
stty sane
#!/bin/bash

FILENAME=`basename "$0"`
REMOTE_CONTEXT_NAME="the-machine"

PROD_ENV_FILE="docker/.env.prod"
PROD_BACKEND_CONFIG="conf.prod.json"

DEV_ENV_FILE="docker/.env.dev"
DEV_BACKEND_CONFIG="conf.dev.json"

# === REMOTE ===
if [[ "$1" = "remote" ]]
then
	docker-compose --context $REMOTE_CONTEXT_NAME --env-file $PROD_ENV_FILE build --no-cache
    docker-compose --context $REMOTE_CONTEXT_NAME --env-file $PROD_ENV_FILE up -d --force-recreate

# === LOCAL ===
elif [[ "$1" = "local" ]]
then
	if [[ "$2" = "clean" ]]
	then
		echo "clean local:"
		rm -rf ../.fuel-tracker-db
	fi
	# build new image
	docker-compose --env-file $DEV_ENV_FILE build --no-cache

	# remove builder image
	docker image prune --force --filter label=stage=builder

	# run container
	docker-compose --env-file $DEV_ENV_FILE up -d --force-recreate
else
	echo "Usage: $FILENAME <remote|local>"
	exit 0
fi

# to "repair" the terminal
stty sane
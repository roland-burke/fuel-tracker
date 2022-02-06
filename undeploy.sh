#!/bin/bash
REMOTE_CONTEXT_NAME=the-machine

# === REMOTE ===
if [ "$1" = "prod" ]
then
	docker-compose --context $REMOTE_CONTEXT_NAME -f docker-compose.yml down
	# to "repair" the terminal
	stty sane
# === LOCAL ===
elif [ "$1" = "dev" ]
then
	docker-compose -f docker-compose.yml down
else
	echo "Usage: $FILENAME <prod|dev>"
	exit 0
fi
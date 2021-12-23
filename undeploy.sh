#!/bin/bash
REMOTE_CONTEXT_NAME=the-machine

# === REMOTE ===
if [ "$1" = "remote" ]
then
	docker-compose --context $REMOTE_CONTEXT_NAME -f docker-compose.yml down
	# to "repair" the terminal
	stty sane
elif [ "$1" = "local" ]
then
	docker-compose -f docker-compose.yml down
else
	echo "Usage: $FILENAME <remote|local>"
	exit 0
fi
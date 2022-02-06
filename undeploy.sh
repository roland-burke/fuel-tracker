#!/bin/bash
REMOTE_CONTEXT_NAME="the-machine"

# === REMOTE ===
if [ "$1" = "remote" ]
then
	docker-compose --context $REMOTE_CONTEXT_NAME down
	# to "repair" the terminal
	stty sane
# === LOCAL ===
elif [ "$1" = "local" ]
then
	docker-compose down
else
	echo "Usage: $FILENAME <remote|local>"
	exit 0
fi
#!/bin/bash
docker-compose --context the-machine --env-file docker/.env.dev -f docker-compose.yml -f docker/docker-compose.dev.yml up -d --build --force-recreate
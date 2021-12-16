#!/bin/bash
docker-compose --context the-machine --env-file docker/.env.prod -f docker-compose.yml -f docker/docker-compose.prod.yml up -d --build --force-recreate
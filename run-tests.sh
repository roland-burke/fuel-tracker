#!/bin/bash

DB_NAME=pgx_test
CONTAINER_NAME=fuel-tracker-postgres-test

# setup 
docker stop ${CONTAINER_NAME} 2> /dev/null
docker rm ${CONTAINER_NAME} 2> /dev/null

# start container and setup db
docker run --name ${CONTAINER_NAME} -p 5432:5432 -e POSTGRES_PASSWORD=testpw -d postgres
sleep 1
docker exec -d ${CONTAINER_NAME} psql -U postgres -c "create database ${DB_NAME};"
docker exec -d ${CONTAINER_NAME} psql -U postgres -c "create domain uint64 as numeric(20,0);"

# init data
docker cp database/db-setup-01.sql ${CONTAINER_NAME}:setup.sql
docker exec -d ${CONTAINER_NAME} psql -U postgres -d ${DB_NAME} -f setup.sql
docker cp database/test/init.sql ${CONTAINER_NAME}:init.sql
docker exec -d ${CONTAINER_NAME} psql -U postgres -d ${DB_NAME} -f init.sql

# run tests
PGX_TEST_DATABASE="host=/var/run/postgresql database=${DB_NAME}" DATABASE_PATH="localhost:5432/${DB_NAME}" DATABASE_USERNAME="postgres" DATABASE_PASSWORD="testpw" go test ./... -v -cover

# cleanup 
docker stop ${CONTAINER_NAME}
docker rm ${CONTAINER_NAME}
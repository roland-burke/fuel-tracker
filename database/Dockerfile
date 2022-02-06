# syntax=docker/dockerfile:1

FROM postgres:latest

# Select config with build arg
ARG userInitFilePath=./dev/db-init-user-dev.sql
ARG dataInitFilePath=./dev/db-init-data-dev.sql

COPY db-setup-01.sql /docker-entrypoint-initdb.d/10-init.sql
COPY ${userInitFilePath} /docker-entrypoint-initdb.d/20-init.sql
COPY ${dataInitFilePath} /docker-entrypoint-initdb.d/30-init.sql
# syntax=docker/dockerfile:1

# Select config with build arg
ARG USER_INIT_PATH
ARG DATA_INIT_PATH

FROM postgres:14

ARG USER_INIT_PATH
ARG DATA_INIT_PATH

COPY db-setup-01.sql /docker-entrypoint-initdb.d/10-init.sql
COPY ${USER_INIT_PATH} /docker-entrypoint-initdb.d/20-init.sql
COPY ${DATA_INIT_PATH} /docker-entrypoint-initdb.d/30-init.sql
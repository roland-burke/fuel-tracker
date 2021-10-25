# syntax=docker/dockerfile:1

FROM postgres:latest

COPY db-setup-01.sql /docker-entrypoint-initdb.d/10-init.sql
COPY db-init-dev.sql /docker-entrypoint-initdb.d/20-init.sql

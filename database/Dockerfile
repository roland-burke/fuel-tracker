# syntax=docker/dockerfile:1

FROM postgres:latest

COPY db-setup.sql /docker-entrypoint-initdb.d/10-init.sql

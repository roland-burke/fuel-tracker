# syntax=docker/dockerfile:1

FROM golang:1.16-alpine as builder
LABEL stage=builder

# Select config with build arg
ARG configFilePath=conf.prod.json

WORKDIR /app

# Download necessary Go modules
COPY src/go.mod ./
COPY src/go.sum ./

RUN go mod download

# Copy files to workdir
COPY src/*.go ./
COPY config/${configFilePath} config/conf.json

RUN pwd
RUN ls config

RUN go build

# Generate clean, final image for deployment
FROM alpine:3.11.3
LABEL stage=deploy

COPY --from=builder ./app/fuel-tracker .
COPY --from=builder ./app/config/conf.json .

# Executable
ENTRYPOINT [ "./fuel-tracker" ]
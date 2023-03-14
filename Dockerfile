# syntax=docker/dockerfile:1

FROM golang:1.16-alpine as builder
LABEL stage=builder

# Select config with build arg
ARG configFilePath=conf.prod.json

WORKDIR /app

# Download necessary Go modules
COPY go.mod ./
COPY go.sum ./

RUN go mod download

# Copy files to workdir
COPY cmd/ ./cmd
COPY internal/ ./internal
COPY config/${configFilePath} config/conf.json

RUN go build -o ./fuel-tracker ./cmd/main

# Generate clean, final image for deployment
FROM alpine:3.11.3
LABEL stage=deploy

COPY --from=builder ./app/fuel-tracker .
COPY --from=builder ./app/config/conf.json .

# Executable
ENTRYPOINT [ "./fuel-tracker" ]

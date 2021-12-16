# syntax=docker/dockerfile:1

FROM golang:1.16-alpine as builder

WORKDIR /app

# Download necessary Go modules
COPY src/go.mod ./
COPY src/go.sum ./

RUN go mod download

# Copy files to workdir
COPY src/*.go ./
COPY config/conf.json ./config/conf.json

RUN go build

# Generate clean, final image for deployment
FROM alpine:3.11.3
COPY --from=builder ./app/fuel-tracker .
COPY --from=builder ./app/config/conf.json .

# Executable
ENTRYPOINT [ "./fuel-tracker" ]
# syntax=docker/dockerfile:1

FROM golang:1.16-alpine as builder

WORKDIR /app

# Download necessary Go modules
COPY go.mod ./
COPY go.sum ./

RUN go mod download

# Copy files to workdir
COPY *.go ./
COPY config/conf.json ./

RUN go build

# Generate clean, final image for deployment
FROM alpine:3.11.3
COPY --from=builder ./app/fuel-tracker .
COPY --from=builder ./app/conf.json .

# Executable
ENTRYPOINT [ "./fuel-tracker" ]
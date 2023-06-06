# syntax=docker/dockerfile:1

FROM golang:1.20-alpine as builder
LABEL stage=builder

WORKDIR /app

# Download necessary Go modules
COPY go.mod ./
COPY go.sum ./

RUN go mod download

# Copy files to workdir
COPY cmd/ ./cmd
COPY internal/ ./internal

RUN go build -o ./fuel-tracker ./cmd/main

# Generate clean, final image for deployment
FROM alpine:3
LABEL stage=deploy

COPY --from=builder ./app/fuel-tracker .

# Executable
ENTRYPOINT [ "./fuel-tracker" ]

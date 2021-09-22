# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

WORKDIR /app

# Download necessary Go modules
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy files to workdir
COPY *.go ./
COPY conf.json ./

RUN go build -o /fuel-tracker-backend

CMD [ "/fuel-tracker-backend" ]

EXPOSE 9008 
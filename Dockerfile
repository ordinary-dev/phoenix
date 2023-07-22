FROM golang:1.20.3-alpine AS builder

RUN apk add gcc
RUN apk add musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY backend ./backend
COPY views ./views

RUN go build -o main

FROM alpine:3.17.3

WORKDIR /app
COPY --from=builder /app/main ./main
COPY assets ./assets
COPY templates ./templates

RUN mkdir /var/lib/phoenix
ENV P_DBPATH="/var/lib/phoenix/db.sqlite3"
ENV P_PRODUCTION="true"

EXPOSE 8080

ENTRYPOINT ["/app/main"]

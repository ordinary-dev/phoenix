FROM golang:1.21.3-alpine3.18 AS builder

RUN apk add gcc
RUN apk add musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

ADD . .

RUN go build -o main

FROM alpine:3.18.4

WORKDIR /app
COPY --from=builder /app/main /usr/local/bin/phoenix
COPY assets ./assets
COPY templates ./templates

RUN mkdir /var/lib/phoenix
ENV P_DBPATH="/var/lib/phoenix/db.sqlite3"
ENV P_PRODUCTION="true"

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/phoenix"]

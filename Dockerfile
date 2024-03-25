FROM golang:1.22 AS builder

RUN apt install -y --no-install-recommends gcc

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

ADD . .

RUN go build -o main

FROM debian:bookworm-slim

WORKDIR /app
COPY --from=builder /app/main /usr/local/bin/phoenix
COPY assets ./assets
COPY templates ./templates

RUN mkdir /var/lib/phoenix
ENV P_DBPATH="/var/lib/phoenix/db.sqlite3"
ENV P_PRODUCTION="true"

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/phoenix"]

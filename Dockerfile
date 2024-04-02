FROM golang:1.22-alpine AS builder

RUN apk add --no-cache gcc libc-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

ADD . .

RUN go build -o main

FROM alpine:3.19

WORKDIR /app
COPY --from=builder /app/main /usr/local/bin/phoenix
COPY assets ./assets
COPY templates ./templates

RUN mkdir /var/lib/phoenix
ENV P_DBPATH="/var/lib/phoenix/db.sqlite3"

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/phoenix"]

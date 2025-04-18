FROM golang:1.23-alpine AS builder

RUN apk add --no-cache gcc libc-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

ADD . .

RUN go build -o main

FROM alpine:3.21

WORKDIR /app
COPY --from=builder /app/main /usr/local/bin/phoenix

RUN mkdir web
COPY web/assets ./web/assets
COPY web/views ./web/views

RUN mkdir /var/lib/phoenix
ENV DB_PATH="/var/lib/phoenix/db.sqlite3"

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/phoenix"]

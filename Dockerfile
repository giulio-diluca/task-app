FROM golang:1.25.1 AS builder

RUN set -eux; \
	apt-get update; \
	apt-get install -y --no-install-recommends coreutils

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /task-app-binary ./cmd

FROM alpine:latest

RUN apk add --no-cache bash

WORKDIR /app

COPY --from=builder /task-app-binary /app/task-app-binary

EXPOSE 8080

CMD ["/app/task-app-binary"]

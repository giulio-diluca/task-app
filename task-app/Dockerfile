FROM golang:1.25.1

RUN set -eux; \
	apt-get update; \
	apt-get install -y --no-install-recommends \
        coreutils

WORKDIR /app

COPY . ./

RUN go mod download

EXPOSE 8080

CMD ["go", "run", "./cmd"]
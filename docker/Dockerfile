FROM golang:1.17 as builder

ADD . /app
WORKDIR /app/

RUN DEBIAN_FRONTEND=noninteractive apt-get -y install make git

RUN env ;\
	go mod vendor ;\
	make -d

FROM node:12 as front

ADD ./web /app
WORKDIR /app/web/

RUN npm ci;\
    npm run build:dev

FROM debian:buster-slim

EXPOSE 8080

ENV APP_CONFIG_PATH /app/configs

RUN adduser --home /app --disabled-password --gecos "" app ;\
	apt update ;\
	apt install ca-certificates curl -y --no-install-recommends ;\
	rm -rf /var/lib/apt/lists/*

USER app

COPY --from=builder /app/bin/go-sso /app/service
COPY --from=builder /app/configs /app/configs

CMD ["/app/service"]

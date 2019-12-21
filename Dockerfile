FROM golang:alpine
WORKDIR /go/src/github.com/step/sauron-reporter/

ADD . .
RUN apk update && apk add --no-cache git ca-certificates make && go get ./... && make reporter

FROM alpine
WORKDIR /app
COPY --from=0 /go/src/github.com/step/sauron-reporter/bin/reporter ./reporter
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs
ENTRYPOINT ["sh","-c", "/app/reporter -redis-address $REDIS_ADDRESS -redis-db $REDIS_DB"]
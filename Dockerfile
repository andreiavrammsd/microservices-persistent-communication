FROM golang:1.7.4-alpine

RUN apk add --update git

ENV APPDIR /go/src/github.com/andreiavrammsd/microservices-persistent-communication
ADD . ${APPDIR}
WORKDIR ${APPDIR}

ENTRYPOINT go get && /go/bin/microservices-persistent-communication

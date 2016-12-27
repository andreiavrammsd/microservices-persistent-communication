FROM golang:1.7.4-alpine

RUN apk add --update git

ENV APPDIR /go/src/github.com/andreiavrammsd/microservices-persistent-communication
ADD . ${APPDIR}
WORKDIR ${APPDIR}
RUN chmod +x wait-for-it.sh
RUN go get

RUN mkdir -p /var/log/microservices-persistent-communication

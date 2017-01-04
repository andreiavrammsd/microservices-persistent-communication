FROM golang:1.7.4-alpine

RUN apk update

RUN apk add --update git
RUN apk add openssl

ENV DOCKERIZE_VERSION v0.3.0
RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz

ENV APPDIR /go/src/github.com/andreiavrammsd/microservices-persistent-communication
ADD . ${APPDIR}
WORKDIR ${APPDIR}

RUN go get

RUN mkdir -p /var/log/microservices-persistent-communication

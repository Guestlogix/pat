FROM golang:1.12.7-alpine3.10

ARG WORKDIR=/go/src/github.com/guestlogix/pat

RUN apk update && apk upgrade && \
    apk add --no-cache bash git

RUN mkdir -p ${WORKDIR}
ADD . ${WORKDIR}
WORKDIR ${WORKDIR}

RUN go get -d ./...
RUN go install github.com/guestlogix/pat
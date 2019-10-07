####
# Used when building the tool for a GitHub Action
####
FROM golang:1.12.7-alpine3.10

ARG PATDIR=/go/src/github.com/guestlogix/pat

RUN apk update && apk upgrade && \
    apk add --no-cache bash git

RUN mkdir -p ${PATDIR}
ADD . ${PATDIR}

# Install all the go tools
RUN cd ${PATDIR} && go get -d ./...
RUN go install github.com/guestlogix/pat

COPY ./actions/entry.sh /tmp
ENTRYPOINT ["tmp/entry.sh"]
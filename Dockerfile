FROM golang:1.12.7-alpine3.10

ARG WORKDIR=/go/src/github.com/guestlogix/pat

RUN apk update && apk upgrade && \
    apk add --no-cache bash git

RUN mkdir -p ${WORKDIR}
ADD . ${WORKDIR}
WORKDIR ${WORKDIR}

# Install all the go tools
RUN go get gopkg.in/Netflix-Skunkworks/go-jira.v1/cmd/jira
RUN go get -d ./...
RUN go install github.com/guestlogix/pat

# Initialize Jira
# RUN mkdir ~/.jira.d
ADD .jira.d /root/.jira.d
# RUN echo "endpoint: https://guestlogix.atlassian.net" > ~/.jira.d/config.yml
# RUN echo "user: hgoddard@guestlogix.com" >> ~/.jira.d/config.yml
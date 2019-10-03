####
# Used when building the tool for a GitHub Action
####
FROM golang:1.12.7-alpine3.10

ARG WORKDIR=/go/src/github.com/guestlogix/pat
ARG JIRA_USER=hgoddard@guestlogix.com
ARG JIRA_ENDPOINT=https://guestlogix.atlassian.net

RUN apk update && apk upgrade && \
    apk add --no-cache bash git

RUN mkdir -p ${WORKDIR}
ADD . ${WORKDIR}
WORKDIR ${WORKDIR}

COPY entry.sh /entry.sh

# Install all the go tools
RUN go get gopkg.in/Netflix-Skunkworks/go-jira.v1/cmd/jira
RUN go get -d ./...
RUN go install github.com/guestlogix/pat

# Initialize Jira
ENV JIRA_USER=${JIRA_USER}
ENV JIRA_ENDPOINT=${JIRA_ENDPOINT}

ENTRYPOINT ["./entry.sh"]
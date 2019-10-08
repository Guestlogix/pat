############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/guestlogix/pat/
COPY . .
# Fetch dependencies.
# Using go get.
RUN go get -d -v
# Build the binary.
RUN go build -o /go/bin/pat
############################
# STEP 2 build a small image
############################
FROM alpine
# Add bash
RUN apk update && apk add bash
# Set tmp as workdir
WORKDIR /tmp
# Copy our static executable.
COPY --from=builder /go/bin/pat /go/bin/pat
# Copy our entry bash to route to proper script
COPY ./actions /tmp
# Entry
ENTRYPOINT ["/tmp/entry.sh"]
# PAT
### Pipeline Automation Tool

## Overview

This maintains a go based CLI tool to perform a number of automated actions in our pipelines.

## Build locally

`go install github.com/guestlogix/pat`

## Docker

```
docker build -t pat .
docker run -it pat
```

## Jira

The PAT image also ships with Netflix's `go-jira` lib for automation involving jira. In order for authentication to work, the env var `JIRA_API_TOKEN` must be set with a vaild API token.

e.g. `docker run -it -e JIRA_API_TOKEN=<TOKEN> JIRA_PROJECT=<PROJECT> pat`

### Create Issue

`jira ctask "Title" -d "Description Body"`
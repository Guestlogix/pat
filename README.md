# PAT
### Pipeline Automation Tool

## Overview

This maintains a go based CLI tool to perform a number of automated actions in our pipelines.

### Build locally

`go install github.com/guestlogix/pat`

### Docker

```
docker build -t pat .
docker run -it pat
```

## Using the CLI
`pat --help`
```
NAME:
   PAT - CLI Tools for pipelines.

USAGE:
   appsettings [global options] command [command options] [arguments...]

VERSION:
   0.0.1

AUTHOR:
   Guestlogix

COMMANDS:
     appsettings, a      Generates a markdown report of altered appsettings and posts a comment on the github pr
     releasenotes, rn    Generates the release notes between two tags
     releaseversion, rv  Obtains the last semver tag in the git history
     help, h             Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```


## Jira

The PAT image also ships with Netflix's `go-jira` lib for automation involving jira. In order for authentication to work, the env var `JIRA_API_TOKEN` must be set with a vaild API token.

e.g. `docker run -it -e JIRA_API_TOKEN=<JIRA_TOKEN> -e JIRA_USER=<JIRA_USER> pat`

### Create Issue
You first need to generate the issue in yml with `pat releasenotes`

`jira --user=$JIRA_USER --endpoint=$JIRA_ENDPOINT create --template ./issue.yml --project RL --noedit`

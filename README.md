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

## Jira

The PAT image also ships with Netflix's `go-jira` lib for automation involving jira. In order for authentication to work, the env var `JIRA_API_TOKEN` must be set with a vaild API token.

e.g. `docker run -it -e JIRA_API_TOKEN=<JIRA_TOKEN> -e JIRA_USER=<JIRA_USER> pat`

### Create Issue
You first need to generate the issue in yml with `pat releasenotes`

`jira --user=$JIRA_USER --endpoint=$JIRA_ENDPOINT create --template ./issue.yml --project RL --noedit`

## Use as a Github Action

You can make use of PAT in a `Github Action` using some bash scripts in the `./actions` folder. The name of the `.sh` script will be the value that is passed in the action `pipeline-command`. In the desired repository of use add a workflow `.yaml` like this, specifying your case statement.

> NOTE: Ensure you give execute permissions to the script (`chmod +x <YOUR_SCRIPT>.sh`)

> NOTE: Be sure to include thr `actions/checkout@master` step if you need access to the actual source code of the calling repo.

```
name: PAT Action
on: [push]

jobs:
  pat:
    runs-on: ubuntu-latest
    name: A job to use pat
    steps:
    - uses: actions/checkout@master
    - name: PAT
      id: pat
      uses: Guestlogix/pat@master
      env:
        GITHUB_ACCESS_TOKEN: ${{ secrets.GITHUB_ACCESS_TOKEN }}
      with:
        pipeline-command: '<YOUR_SCRIPT_NAME>'
```

Finally, update the chart below with the new functionality.

| Name         | Key            | ENV Vars | Notes                                                                                                        |
|--------------|----------------|-----------|--------------------------------------------------------------------------------------------------------------|
| Auto Version | `auto-version` | GITHUB_ACCESS_TOKEN  | Finds the latest semantic version, then increments it according to the appropriate conventional commit name. |
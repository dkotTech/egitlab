# eGitlab

extension for gitlab

## Install

```
go install github.com/dkotTech/egitlab@latest
```

## Usage

### 1. Setup a configuration

You need a gitlab token to request pipelines

```
egitlab set-creds
```

```
> Gitlab token

[ Submit ]

q: exit
```

_By default cli save only one token, so you can not use multiple token at once._

### 2.1 Use it in your repository (with git)

```
egitlab pipelines 
```
or
```
egitlab p 
```

#### Example output

```
master Open pipeline
┏━━━━━━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━━━━━━━━━━━━━┓
┃tests               ┃build               ┃deploy              ┃after-deploy             ┃
┣━━━━━━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━━━━━━━━━━━━┫
┃> tests             ┃> docker            ┃> production        ┃> integration-tests      ┃
┃✅ success          ┃✅ success          ┃✋ manual           ┃⏭ skipped               ┃
┃Open                ┃Open                ┃Open                ┃Open                     ┃
┃2m18s               ┃2m18s               ┃0s                  ┃0s                       ┃
┃alt+5: run again job┃alt+4: run again job┃alt+2: run job      ┃                         ┃
┣━━━━━━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━━━━━━━━━━━━┫
┃                    ┃                    ┃> development       ┃> other-test             ┃
┃                    ┃                    ┃✋ manual           ┃🆕 created               ┃
┃                    ┃                    ┃Open                ┃Open                     ┃
┃                    ┃                    ┃1m15s               ┃0s                       ┃
┃                    ┃                    ┃alt+3: run job      ┃alt+1: run job           ┃
┗━━━━━━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━━━━━━━━━━━━━┛
next update in: 26.38s
last update: 2024-01-01T10:00:00Z
q: exit | u: update now
```

by default cli use a git as executive, so its panic if you have not git installed

### 2.2 Use it anywhere (without git)

```
egitlab p -gitlab-host=https://gitlab.com -gitlab-project=gitlab-org/gitlab -use-git-cli=false
```

```
USAGE:
   egitlab pipelines [command options]

OPTIONS:
   --update-interval value                    pipeline status update interval (default: 30s)
   --create-job value [ --create-job value ]  create a job on update by name, if job in manual state
   --retry-job value [ --retry-job value ]    retry a job on update by name, if job in failed or success state
   --gitlab-host value                        gitlab host
   --ref value                                git ref (default: "master")
   --gitlab-project value                     gitlab project name
   --use-git-cli                              using git cli (default: true)
   --help, -h                                 show help
```

### 2.3 Execute a pipeline as soon as it is ready to run

* `--create-job` - create a new jobs from manual state
* `--retry-job` - rerun jobs if they in failed or success states 

#### Example 

Create a _docker_ job and create a _development_ deploy job. 

Typically, deployment happens after the docker build process, so first, execute the _docker_ job, and once it succeeds, initiate the _development_ deploy job.
```
egitlab p -gitlab-host='https://gitlab.com' -gitlab-project=gitlab-org/gitlab -use-git-cli=false --ref=master --create-job=development --create-job=docker
```

In project directory just:

```
egitlab p --create-job=development --create-job=docker
```
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

Example output

```
master Open pipeline
┏━━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━━┳━━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━━━━━━━━━━━━━┓
┃tests           ┃build         ┃deploy          ┃after-deploy             ┃
┣━━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━━━━━━━━━━━━┫
┃> tests         ┃> docker      ┃> production    ┃> integration-tests      ┃
┃✅ success      ┃✅ success    ┃✋ manual       ┃⏭ skipped               ┃
┃Open            ┃Open          ┃Open            ┃Open                     ┃
┃2m18s           ┃1m51s         ┃0s              ┃0s                       ┃
┣━━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━━━━━━━━━━━━┫
┃                ┃              ┃> development   ┃> other-test             ┃
┃                ┃              ┃✋ manual       ┃🆕 created               ┃
┃                ┃              ┃Open            ┃Open                     ┃
┃                ┃              ┃0s              ┃0s                       ┃
┗━━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━━┻━━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━━━━━━━━━━━━━┛
next update in: 26.38s
q: exit
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
   --update-interval value  pipeline status update interval (default: 30s)
   --gitlab-host value      gitlab host
   --ref value              git ref (default: "master")
   --gitlab-project value   gitlab project name
   --use-git-cli            using git cli (default: true)
   --help, -h               show help
```

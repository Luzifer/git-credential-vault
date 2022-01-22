[![Go Report Card](https://goreportcard.com/badge/github.com/Luzifer/git-credential-vault)](https://goreportcard.com/report/github.com/Luzifer/git-credential-vault)
![](https://badges.fyi/github/license/Luzifer/git-credential-vault)
![](https://badges.fyi/github/downloads/Luzifer/git-credential-vault)
![](https://badges.fyi/github/latest-release/Luzifer/git-credential-vault)
![](https://knut.in/project-status/git-credential-vault)

# Luzifer / git-credential-vault

`git-credential-vault` is an implementation of the [Git Credential Storage](https://git-scm.com/book/en/v2/Git-Tools-Credential-Storage) utilizing [Vault](https://www.vaultproject.io/) as storage backend.

The only supported action is `get` as storage is managed through Vault related tools / the web-UI. The tool expects to find Vault keys per host containing `username` / `password` fields in it. Those fields are then combined with the data received from git and sent back for authentication.

## Expected Vault structure

```
secret/git-credentials (pass this to --vault-path-prefix)
 +- github.com
 |   +- username = api
 |   +- password = verysecrettoken
 +- gitlab.com
     +- username = user
     +- password = anothertoken
```

## Usage

```console
# export VAULT_ADDR=http://localhost:8200
# export VAULT_TOKEN=somesecretvaulttoken
# echo -e "protocol=https\nhost=github.com\n\n" | ./git-credential-vault --vault-path-prefix secret/git-credentials get
host=github.com
username=api
password=myverysecrettoken
protocol=https
```

### Vault KV Secrets Engine - Version 2

This tool supports both versions of the Vault KV Secrets Engine. You just need to consider one thing: Version 2 of the KV Secrets Engine does use slightly modified paths for reading secrets. In order to be compatible to both versions of the Secrets Engine you need to adjust the `vault-path-prefix` slightly when using it:

```bash
# Version 1
vault list secret_v1/git-credentials
# Keys
# ----
# github.com
git config --global credential.helper 'vault --vault-path-prefix secret_v1/git-credentials'
```

```bash
# Version 2
vault kv list secret_v2/git-credentials
# Keys
# ----
# github.com
git config --global credential.helper 'vault --vault-path-prefix secret_v2/data/git-credentials'
```

Mind the extra `/data` after the mountpoint for a mountpoint using version 2. If you omit it the tool will not work properly as it will not yield any credentials.

### Dockerfile example (git clone)

In this example the `VAULT_TOKEN` is passed in through a build-arg which means you **MUST** revoke the token before pushing the image, otherwise you will be leaking an active credential!

```Dockerfile
FROM alpine

ARG VAULT_ADDR
ARG VAULT_TOKEN

RUN set -ex \
 && apk --no-cache add curl git \
 && curl -sSfL "https://github.com/Luzifer/git-credential-vault/releases/download/v0.1.0/git-credential-vault_linux_amd64.tar.gz" | tar -xz -C /usr/bin \
 && mv /usr/bin/git-credential-vault_linux_amd64 /usr/bin/git-credential-vault \
 && git config --global credential.helper 'vault --vault-path-prefix secret/git-credentials'

RUN set -ex \
 && git clone https://github.com/myuser/secretrepo.git /src
```

```console
# docker build --build-arg VAULT_ADDR=${VAULT_ADDR} --build-arg VAULT_TOKEN=${VAULT_TOKEN} --no-cache .
```

### Dockerfile example (go get)

In this example the `VAULT_TOKEN` is passed in through a build-arg which means you **MUST** revoke the token before pushing the image, otherwise you will be leaking an active credential!

```Dockerfile
FROM golang:alpine

ARG VAULT_ADDR
ARG VAULT_TOKEN

RUN set -ex \
 && apk --no-cache add git \
 && go get -u -v github.com/Luzifer/git-credential-vault \
 && git config --global credential.helper 'vault --vault-path-prefix secret/git-credentials'

RUN set -ex \
 && go get -v github.com/myuser/secretrepo
```

```console
# docker build --build-arg VAULT_ADDR=${VAULT_ADDR} --build-arg VAULT_TOKEN=${VAULT_TOKEN} --no-cache .
```

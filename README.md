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
secret/git-credentials
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

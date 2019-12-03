# govaultenv

Expose vault secrets as environment variables in newly spawned process.

## Motivation

Secrets are frequently poorly maintained (never rotated, copied across developers and environemtns, etc, etc), so the intent of this utility is to simplify secrets handling using vault
as primary and the only secrets storage.

There are two primary use cases:

1. Securily provide an environment for ops and devs
2. Provide secrets to kubernetes applications

## Usage

```sh
govaultenv program arguments...
```

Complete workflow with interactive shell spawned and vault secret exposed as environment variables:

```sh
export VAULT_ADDR=...
vault login -method=okta username=...
govaultenv -verbose=debug /bin/bash
```

## How to install on Mac

```sh
brew tap jamhed/govaultenv https://github.com/jamhed/govaultenv
brew install govaultenv
```

## How to authenticate

Either do:

* `vault login` with one of supported auth schemas (govaultenv tries to read `~/.vault-token` file)
* export `VAULT_TOKEN` environment variable with valid TOKEN
* provide token value with `-token` command line flag
* provide kubernetes token and authentication path, see below

## How it works

It traverses all environment variables looking for prefix (`VAULT_` by default), and then if environment variable vaule is in format:

* `secret-path#key`, then secret is fetched from vault, and new environment variable without
prefix is set to secret key value, e.g. if you have a secret `team/solr` with key `password`, and you have an environment
variable `VAULT_SOLR_PASS=team/solr#password` defined, then spawned process has new environment variable `SOLR_PASS`
set to the value of the corresponding vault secret.
* `secret-path#key:local-path`, then secret key is written to local file `local-path`, and all written
secrets are deleted upon completion of calling program.
* `secret-path`, then all keys are exposed as generated environment variable named as `variableName_keyName`, e.g. if you have it
as `VAULT_SOLR=team/solr`, and solr secret has keys `username` and `password`, then following environment variables
are generated: `SOLR_USERNAME` and `SOLR_PASSWORD`.

There are command-line flags to control this behaviour:

* `uppercase`, true by default, set it to false to keep generated environment variable name as it is
* `stripname`, false by default, set it to true to strip the original environment variable name from generated one

## Wrapped tokens

Vault has a useful feature called `wrapped tokens` that allows to securely pass secrets (including tokens) around,
and govaultenv has an option `unwrap` to support it.

```
WRAPPED_TOKEN=$(vault token create -field=wrapping_token -wrap-ttl=1h -ttl=1h)
govaultenv -unwrap -token $WRAPPED_TOKEN env
```

Here wrapped token can be used only once, has limited time-to-live (one hour), and underlying token has also limited time-to-live (one hour).

## How to verify it

Have `govaultenv` binary installed locally, have `VAULT_ADDR` and `VAULT_TOKEN` environment variables set, and expose some vault secret, e.g. `VAULT_SOLR_PASS=team/solr#pass`, and then:

```sh
govaultenv -append=false env
```

You should be able to see a secret value as environment variable value.

## Kubernetes

Make sure you have `govaultenv` binary residing in the image.

Start your image in proper namespace with proper service account, e.g.

```sh
kubectl run --generator=run-pod/v1 tmp --rm -i --tty --serviceaccount=vault-auth --image jamhed/govaultenv
```

Inside kubernetes pod it's possible to use service account vault authentication schema:

```sh
export VAULT_SOLR_PASS=team/solr#pass
govaultenv -kubeauth default@kubernetes -append=false env
SOLR_PASS=...
```

## Operations

How to spawn an interactive shell with secret variable keys pulled out of vault:

```sh
export VAULT_GOVC=team/env
govaultenv /bin/bash
```

## Related projects

1. https://github.com/mumoshu/aws-secret-operator
2. https://github.com/hashicorp/envconsul
3. https://github.com/channable/vaultenv
4. https://github.com/hashicorp/vault/issues/7364
5. https://github.com/sethvargo/vault-kubernetes-authenticator
6. https://github.com/tuenti/secrets-manager
7. https://github.com/DaspawnW/vault-crd
8. https://github.com/Talend/vault-sidecar-injector
9. https://github.com/hashicorp/consul-template
10. https://github.com/ricoberger/vault-secrets-operator

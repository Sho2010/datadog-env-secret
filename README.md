# What's this

`datadog-env-secret` is implemented datadog "secret_backend_command" to get secret from environment variable.

https://docs.datadoghq.com/agent/guide/secrets-management/?tab=linux

# Usage

Execute

```sh
# Environment variable key is always upper case of secret key.
$ export SECRET1=secret_value
echo '{ "version": "1.0", "secrets": ["secret1", "secret2"] }' | ./datadog-env-secret
```

Result

```json
{
  "secret1": {
    "value": "secret_value",
    "error": null
  },
  "secret2": {
    "value": null,
    "error": "environment variable [SECRET2] is not set"
  }
}
```

# Installation(for linux)

## Get datadog-env-secret

```sh
$ go get github.com/Sho2010/datadog-env-secret


# On Linux, the executable set as secret_backend_command must:
#
# Belong to the same user running the Agent (dd-agent by default, or root inside a container).
# Have no rights for group or other.
# Have at least exec rights for the owner.
#
$ chown dd-agent:dd-agent ${GOPATH}/bin/datadog-env-secret
$ chmod 700 ${GOPATH}/bin/datadog-env-secret
```

## Update your DD agent config

e.g. `/etc/datadog-agent/datadog.yaml`

```yaml
# e.g.
# secret_backend_command: "/usr/local/bin/datadog-env-secret"
secret_backend_command: ${YOUR_TOOL_PATH}
```

## Confirmation

```sh
$ sudo -u dd-agent -- datadog-agent secret
=== Checking executable rights ===
Executable path: /usr/local/bin/datadog-env-secret
Check Rights: OK, the executable has the correct rights

Rights Detail:
file mode: 100700
Owner username: dd-agent
Group name: dd-agent
```

- - -

# Use in conf.d

Example:
```
instances:
  - server: db_prod
    # two valid secret handles
    user: "ENC[db_prod_user]"
    password: "ENC[db_prod_password]"

    # The `ENC[]` handle must be the entire YAML value, which means that
    # the following is NOT detected as a secret handle:
    password2: "db-ENC[prod_password]"
```

`!!! IMPORTANT both edit`

- /etc/init/datadog-agent-process.conf
- /etc/init/datadog-agent.conf

```
env DB_PROD_USER="xxxxxxxxx"
env DB_PROD_PASSWORD="xxxxxxxxx"
```

### Confirm

```sh
$ sudo -u dd-agent -- datadog-agent secret
=== Secrets stats ===
Number of secrets decrypted: 2
Secrets handle decrypted:
- DB_PROD_USER: from hoge
- DB_PROD_PASSWORD: from hoge
```

# Vault-Login

Vault-Login is a simple command-line tool that authenticates against HashiCorp Vault using various authentication mechanisms. After successful authentication, it writes the received Vault token to various outputs (e.g., standard output, file or Kubernetes secret).

This tool helps automate the authentication process with Vault and can be useful for CI/CD pipelines or any environment that needs programmatic access to Vault.

## Features

- Supports multiple authentication methods
- Outputs Vault token to various locations
- Lightweight and easy to use for automation and scripting

## Supported Authentication Methods

- AppRole
- Kubernetes

## Configuration

| Flag Option                         | Environment Variable                      | Description                                      |
|-------------------------------------|-------------------------------------------|--------------------------------------------------|
| `--auth-type`                       | `VAULT_LOGIN_AUTH_TYPE`                   | Type of the authentication                       |
| `--auth-role`                       | `VAULT_LOGIN_AUTH_ROLE`                   | Role for authentication                          |
| `--auth-mount`                      | `VAULT_LOGIN_AUTH_MOUNT`                  | Mount point for authentication                   |
| `--auth-approle-secret-id`          | `VAULT_LOGIN_AUTH_APPROLE_SECRET_ID`      | Approle Secret ID for authentication             |
| `--auth-approle-secret-id-file`     | `VAULT_LOGIN_AUTH_APPROLE_SECRET_ID_FILE` | Approle Secret ID file for authentication        |
| `--output-type`                     | `VAULT_LOGIN_OUTPUT_TYPE`                 | Type of output                                   |
| `--output-secret-name`              | `VAULT_LOGIN_OUTPUT_SECRET_NAME`          | Output secret name                               |
| `--output-secret-namespace`         | `VAULT_LOGIN_OUTPUT_SECRET_NAMESPACE`     | Output secret namespace                          |
| `--output-secret-key`               | `VAULT_LOGIN_OUTPUT_SECRET_KEY`           | Output secret key                                |

### Notes:
- If both an environment variable and a flag are provided for a particular field, the flag will take precedence over the environment variable.
- You can configure the fields using either environment variables or command-line flags depending on your preference.

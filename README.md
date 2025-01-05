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

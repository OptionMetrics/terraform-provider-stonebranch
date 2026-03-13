# Security Policy

## Reporting a Vulnerability

If you discover a security vulnerability in this project, please report it responsibly.

**Do not open a public GitHub issue for security vulnerabilities.**

Instead, please send an email to the maintainers with:

1. A description of the vulnerability
2. Steps to reproduce the issue
3. Any potential impact

We will acknowledge receipt within 48 hours and provide a timeline for a fix.

## Supported Versions

| Version | Supported |
|---------|-----------|
| Latest  | Yes       |

## Security Best Practices

When using this provider:

- **Never hardcode API tokens** in Terraform configuration files. Use environment variables (`STONEBRANCH_API_TOKEN`) or a secrets manager.
- **Mark sensitive variables** with `sensitive = true` in your Terraform configurations.
- **Restrict API token permissions** to the minimum required for your use case.
- **Use `.gitignore`** to exclude `.env` files, `.tfstate` files, and any files containing credentials.

# Stonebranch Credential Example

This example demonstrates how to create credential resources in Stonebranch Universal Controller.

## Resources Created

- `stonebranch_credential.app_user` - Application service account credentials
- `stonebranch_credential.sftp` - SFTP server credentials
- `stonebranch_task_unix.with_creds` - Task using credentials
- `stonebranch_task_file_transfer.with_creds` - File transfer using credentials

## Usage

```bash
# Set your API token
export STONEBRANCH_API_TOKEN="your-token"

# Set credential passwords (sensitive - use environment variables)
export TF_VAR_app_password="your-app-password"
export TF_VAR_sftp_password="your-sftp-password"

# Initialize and apply
terraform init
terraform plan
terraform apply
```

## Variables

| Name | Description | Default | Required |
|------|-------------|---------|----------|
| `agent_var` | Variable name for the agent | `agent_name` | No |
| `app_user` | Application username | `appuser` | No |
| `app_password` | Application password | - | Yes |
| `sftp_user` | SFTP username | `sftpuser` | No |
| `sftp_password` | SFTP password | - | Yes |

## Credential Attributes

| Attribute | Description |
|-----------|-------------|
| `name` | Unique name for the credential |
| `description` | Optional description |
| `runtime_user` | Username for authentication |
| `runtime_password` | Password for authentication (sensitive) |

## Using Credentials

### In Tasks
```hcl
resource "stonebranch_task_unix" "example" {
  credentials = stonebranch_credential.app_user.name
  # ...
}
```

### In File Transfers
```hcl
resource "stonebranch_task_file_transfer" "example" {
  remote_credentials = stonebranch_credential.sftp.name
  # ...
}
```

## Security Best Practices

- **Never commit passwords to version control**
- Use environment variables (`TF_VAR_*`) for sensitive values
- Consider using HashiCorp Vault or similar for secrets management
- Use `.tfvars` files (excluded from git) for local development

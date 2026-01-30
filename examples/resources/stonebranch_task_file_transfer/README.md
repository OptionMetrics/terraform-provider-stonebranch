# Stonebranch File Transfer Task Example

This example demonstrates how to create file transfer tasks (SFTP/FTP) in Stonebranch Universal Controller.

## Resources Created

- `stonebranch_task_file_transfer.sftp_download` - SFTP download task
- `stonebranch_task_file_transfer.sftp_upload` - SFTP upload task
- `stonebranch_credential.sftp` - Credentials for SFTP authentication

## Usage

```bash
# Set your API token
export STONEBRANCH_API_TOKEN="your-token"

# Set SFTP password (sensitive)
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
| `sftp_server` | SFTP server hostname | `sftp.example.com` | No |
| `sftp_user` | SFTP username | `sftpuser` | No |
| `sftp_password` | SFTP password | - | Yes |

## Supported Server Types

- `SFTP` - SSH File Transfer Protocol (recommended)
- `FTP` - File Transfer Protocol

## Notes

- The credential resource stores authentication details securely
- File paths should be absolute paths on the respective systems
- Ensure the agent has network access to the remote server

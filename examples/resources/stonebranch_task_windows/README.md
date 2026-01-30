# Stonebranch Windows Task Example

This example demonstrates how to create Windows tasks in Stonebranch Universal Controller.

## Resources Created

- `stonebranch_task_windows.hello` - Simple task that runs an echo command
- `stonebranch_task_windows.with_script` - Task that executes a batch script resource
- `stonebranch_task_windows.powershell` - Task that runs PowerShell commands
- `stonebranch_script.windows_batch` - Supporting batch script resource

## Usage

```bash
# Set your API token
export STONEBRANCH_API_TOKEN="your-token"

# Optionally set the base URL
export STONEBRANCH_BASE_URL="https://your-instance.stonebranch.cloud"

# Initialize and apply
terraform init
terraform plan
terraform apply
```

## Variables

| Name | Description | Default |
|------|-------------|---------|
| `agent_var` | Variable name containing the Windows agent (resolved at runtime) | `windows_agent_name` |

## Windows-Specific Attributes

The Windows task resource supports these platform-specific attributes:

- `elevate_user` - Run with administrator privileges
- `desktop_interact` - Allow interaction with the desktop
- `create_console` - Create a console window

## Notes

- Ensure you have a Windows agent registered in Stonebranch
- The `agent_var` should resolve to a valid Windows agent at runtime
- PowerShell execution policy may need to be configured on the agent

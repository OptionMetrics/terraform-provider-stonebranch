# Stonebranch Unix Task Example

This example demonstrates how to create Unix/Linux tasks in Stonebranch Universal Controller.

## Resources Created

- `stonebranch_task_unix.hello` - Simple task that runs an echo command
- `stonebranch_task_unix.with_script` - Task that executes a script resource
- `stonebranch_task_unix.with_retry` - Task with retry configuration
- `stonebranch_script.backup` - Supporting script resource

## Usage

```bash
# Set your API token
export STONEBRANCH_API_TOKEN="your-token"

# Set the base URL for your StoneBranch instance
export STONEBRANCH_BASE_URL="https://your-instance.stonebranch.cloud"

# Initialize and apply
terraform init
terraform plan
terraform apply
```

## Variables

| Name | Description | Default |
|------|-------------|---------|
| `agent_var` | Variable name containing the agent (resolved at runtime) | `agent_name` |

## Notes

- The `agent_var` variable should match a Stonebranch variable that resolves to a valid agent name at runtime
- Tasks are created but not executed - use triggers or manual execution to run them

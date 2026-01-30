# Stonebranch Script Example

This example demonstrates how to create reusable script resources in Stonebranch Universal Controller.

## Resources Created

- `stonebranch_script.bash_simple` - Simple bash script
- `stonebranch_script.bash_with_vars` - Bash script with variable substitution
- `stonebranch_script.windows_batch` - Windows batch script
- `stonebranch_script.python` - Python script
- `stonebranch_task_unix.run_script` - Task that executes a script

## Usage

```bash
# Set your API token
export STONEBRANCH_API_TOKEN="your-token"

# Initialize and apply
terraform init
terraform plan
terraform apply
```

## Variables

| Name | Description | Default |
|------|-------------|---------|
| `agent_var` | Variable name for the agent | `agent_name` |

## Script Attributes

| Attribute | Description |
|-----------|-------------|
| `name` | Unique name for the script |
| `content` | The script content |
| `description` | Optional description |
| `resolve_variables` | Enable Stonebranch variable substitution |

## Using Scripts in Tasks

To use a script in a task, set:
- `command_or_script = "Script"`
- `script = stonebranch_script.<name>.name`

## Notes

- Scripts are stored centrally and can be reused across multiple tasks
- When `resolve_variables = true`, Stonebranch variables like `${_var}` are substituted at runtime
- The script content supports any scripting language that can run on the target agent

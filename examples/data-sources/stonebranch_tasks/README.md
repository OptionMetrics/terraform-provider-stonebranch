# Stonebranch Tasks Data Source Example

This example demonstrates how to query tasks from Stonebranch Universal Controller.

## Data Sources Used

- `stonebranch_tasks` - Retrieves a list of tasks with optional filtering

## Filter Options

| Attribute | Description |
|-----------|-------------|
| `name` | Filter by task name (supports wildcards like `*`) |
| `type` | Filter by task type |
| `agent_name` | Filter by assigned agent |
| `business_services` | Filter by business service name |
| `workflow_name` | Filter by workflow membership |

## Task Attributes

Each task in the results includes:

| Attribute | Description |
|-----------|-------------|
| `sys_id` | System ID of the task |
| `name` | Name of the task |
| `type` | Type of the task (taskUnix, taskWindows, etc.) |
| `summary` | Summary/description |
| `version` | Version number |
| `agent` | Assigned agent |
| `agent_cluster` | Assigned agent cluster |
| `credentials` | Credentials used |

## Usage

```bash
# Set your API token
export STONEBRANCH_API_TOKEN="your-token"

# Initialize and apply
terraform init
terraform plan
```

## Example Usage

```hcl
# Find all backup-related tasks
data "stonebranch_tasks" "backup" {
  name = "*backup*"
}

# Find tasks in a workflow
data "stonebranch_tasks" "workflow_tasks" {
  workflow_name = "nightly-batch"
}

# Output task names
output "backup_tasks" {
  value = [for t in data.stonebranch_tasks.backup.tasks : t.name]
}
```

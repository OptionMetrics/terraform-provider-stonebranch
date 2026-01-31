# Stonebranch Task Instances Data Source Example

This example demonstrates how to query task execution history from Stonebranch Universal Controller.

## Data Sources Used

- `stonebranch_task_instances` - Retrieves a list of task instances with filtering

## Filter Options

| Attribute | Description |
|-----------|-------------|
| `task_name` | **Required.** Filter by task name (use `*` for wildcard) |
| `status` | Filter by status (Running, Success, Failed, Waiting, etc.) |
| `type` | Filter by task type |
| `agent_name` | Filter by agent |
| `updated_time_type` | Time filter type: `Today`, `Offset`, `Since`, `Older Than` |
| `updated_time` | Time value (e.g., `1h`, `30mn`, `2d`) |
| `workflow_instance_name` | Filter by parent workflow instance |
| `business_services` | Filter by business service |

## Task Instance Attributes

Each instance in the results includes:

| Attribute | Description |
|-----------|-------------|
| `sys_id` | System ID of the instance |
| `name` | Instance name |
| `type` | Task type |
| `status` | Current status |
| `status_description` | Status description |
| `trigger_time` | When the task was triggered |
| `start_time` | When execution started |
| `end_time` | When execution ended |
| `exit_code` | Exit code |
| `agent` | Agent that executed the task |
| `task_name` | Task definition name |
| `instance_number` | Instance number |
| `triggered_by` | What triggered the execution |
| `workflow_instance_name` | Parent workflow instance |

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
# Get failed tasks from the last hour
data "stonebranch_task_instances" "failed" {
  task_name         = "*"
  status            = "Failed"
  updated_time_type = "Offset"
  updated_time      = "1h"
}

# Output failed task details
output "failed_tasks" {
  value = [
    for inst in data.stonebranch_task_instances.failed.task_instances : {
      name      = inst.name
      exit_code = inst.exit_code
      end_time  = inst.end_time
    }
  ]
}
```

## Time Filter Examples

| updated_time_type | updated_time | Description |
|-------------------|--------------|-------------|
| `Today` | (not needed) | All instances from today |
| `Offset` | `30mn` | Last 30 minutes |
| `Offset` | `1h` | Last 1 hour |
| `Offset` | `2d` | Last 2 days |
| `Since` | (datetime) | Since a specific time |
| `Older Than` | `7d` | Older than 7 days |

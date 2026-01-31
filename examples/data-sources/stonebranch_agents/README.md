# Stonebranch Agents Data Source Example

This example demonstrates how to query agents from Stonebranch Universal Controller.

## Data Sources Used

- `stonebranch_agents` - Retrieves a list of agents with optional filtering

## Filter Options

| Attribute | Description |
|-----------|-------------|
| `name` | Filter by agent name (supports wildcards like `*`) |
| `type` | Filter by type: `Windows`, `Linux/Unix`, `z/OS` |
| `business_services` | Filter by business service name |

## Agent Attributes

Each agent in the results includes:

| Attribute | Description |
|-----------|-------------|
| `sys_id` | System ID of the agent |
| `name` | Name of the agent |
| `type` | Type of the agent |
| `host_name` | Hostname of the agent machine |
| `ip_address` | IP address of the agent |
| `status` | Current status |
| `version` | Agent software version |
| `os` | Operating system |
| `suspended` | Whether the agent is suspended |
| `decommissioned` | Whether the agent is decommissioned |

## Usage

```bash
# Set your API token
export STONEBRANCH_API_TOKEN="your-token"

# Initialize and apply
terraform init
terraform plan
```

## Example Usage in Tasks

```hcl
# Get available Unix agents
data "stonebranch_agents" "unix" {
  type = "Linux/Unix"
}

# Use the first available agent in a task
resource "stonebranch_task_unix" "example" {
  name    = "my-task"
  agent   = data.stonebranch_agents.unix.agents[0].name
  command = "echo hello"
}
```

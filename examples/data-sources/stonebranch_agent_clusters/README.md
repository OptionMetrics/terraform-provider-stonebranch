# Stonebranch Agent Clusters Data Source Example

This example demonstrates how to query agent clusters from Stonebranch Universal Controller.

## Data Sources Used

- `stonebranch_agent_clusters` - Retrieves a list of agent clusters with optional filtering

## Filter Options

| Attribute | Description |
|-----------|-------------|
| `name` | Filter by cluster name |
| `type` | Filter by type: `Windows`, `Linux/Unix` |
| `business_services` | Filter by business service name |

## Agent Cluster Attributes

Each cluster in the results includes:

| Attribute | Description |
|-----------|-------------|
| `sys_id` | System ID of the cluster |
| `name` | Name of the cluster |
| `type` | Type of the cluster |
| `description` | Description |
| `version` | Version number |
| `distribution` | Distribution method for task assignment |
| `suspended` | Whether the cluster is suspended |
| `limit_type` | Type of limit applied |
| `limit_amount` | Limit amount |

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
# Get available Unix agent clusters
data "stonebranch_agent_clusters" "unix" {
  type = "Linux/Unix"
}

# Use the first available cluster in a task
resource "stonebranch_task_unix" "example" {
  name          = "my-task"
  agent_cluster = data.stonebranch_agent_clusters.unix.agent_clusters[0].name
  command       = "echo hello"
}
```

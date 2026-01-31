# Stonebranch Tasks Data Source Example
#
# This example demonstrates how to query tasks from Stonebranch.
#
# Usage:
#   export STONEBRANCH_API_TOKEN="your-token"
#   terraform init
#   terraform plan

terraform {
  required_providers {
    stonebranch = {
      source = "stonebranch/stonebranch"
    }
  }
}

provider "stonebranch" {
  # Uses STONEBRANCH_API_TOKEN environment variable
}

# List all tasks
data "stonebranch_tasks" "all" {
}

# List tasks by name pattern (wildcard)
data "stonebranch_tasks" "backup_tasks" {
  name = "*backup*"
}

# List tasks assigned to a specific agent
data "stonebranch_tasks" "agent_tasks" {
  agent_name = "my-agent"
}

# List tasks in a specific workflow
data "stonebranch_tasks" "workflow_tasks" {
  workflow_name = "my-workflow"
}

# List tasks in a specific business service
data "stonebranch_tasks" "production_tasks" {
  business_services = "Production"
}

# Output examples
output "all_task_count" {
  description = "Total number of tasks"
  value       = length(data.stonebranch_tasks.all.tasks)
}

output "backup_task_names" {
  description = "List of backup task names"
  value       = [for task in data.stonebranch_tasks.backup_tasks.tasks : task.name]
}

output "task_types_summary" {
  description = "Summary of task types"
  value = {
    for task in data.stonebranch_tasks.all.tasks :
    task.type => task.name...
  }
}

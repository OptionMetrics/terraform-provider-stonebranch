# Stonebranch Task Data Source Example
#
# This example demonstrates how to look up an existing task by name.
# The data source returns common fields available across all task types.

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

# Look up an existing task by name
data "stonebranch_task" "example" {
  name = "my-existing-task"
}

# Output the task details
output "task_type" {
  description = "Type of the task"
  value       = data.stonebranch_task.example.type
}

output "task_sys_id" {
  description = "System ID of the task"
  value       = data.stonebranch_task.example.sys_id
}

output "task_agent" {
  description = "Agent assigned to the task"
  value       = data.stonebranch_task.example.agent
}

# Example: Use the data source to reference an existing task in a trigger
# resource "stonebranch_trigger_time" "scheduled" {
#   name  = "my-scheduled-trigger"
#   tasks = [data.stonebranch_task.example.name]
#   time  = "09:00"
# }

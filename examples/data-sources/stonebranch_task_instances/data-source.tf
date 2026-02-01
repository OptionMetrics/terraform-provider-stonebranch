# Stonebranch Task Instances Data Source Example
#
# This example demonstrates how to query task execution history from Stonebranch.
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

# List all task instances from today (task_name is required, use * for all)
data "stonebranch_task_instances" "today" {
  task_name         = "*"
  updated_time_type = "Today"
}

# List failed task instances from the last hour
data "stonebranch_task_instances" "failed_recent" {
  task_name         = "*"
  status            = "Failed"
  updated_time_type = "Offset"
  updated_time      = "1h"
}

# List task instances for a specific task
data "stonebranch_task_instances" "specific_task" {
  task_name         = "my-important-task"
  updated_time_type = "Today"
}

# List running task instances
data "stonebranch_task_instances" "running" {
  task_name         = "*"
  status            = "Running"
  updated_time_type = "Today"
}

# List task instances in a specific workflow
data "stonebranch_task_instances" "workflow" {
  task_name              = "*"
  workflow_instance_name = "my-workflow-instance"
  updated_time_type      = "Today"
}

# Output examples
output "today_instance_count" {
  description = "Number of task instances today"
  value       = length(data.stonebranch_task_instances.today.task_instances)
}

output "failed_tasks" {
  description = "List of failed task names from the last hour"
  value       = [for inst in data.stonebranch_task_instances.failed_recent.task_instances : inst.name]
}

output "running_tasks" {
  description = "Currently running tasks"
  value = [
    for inst in data.stonebranch_task_instances.running.task_instances : {
      name       = inst.name
      start_time = inst.start_time
      agent      = inst.agent
    }
  ]
}

output "task_status_summary" {
  description = "Summary of today's task statuses"
  value = {
    for inst in data.stonebranch_task_instances.today.task_instances :
    inst.status => inst.name...
  }
}

# Stonebranch Task Monitor Task Example
#
# This example demonstrates how to create Task Monitor tasks in Stonebranch.
# A Task Monitor task watches for specific status changes in other tasks
# (e.g., when a task completes, succeeds, or fails). Task Monitor tasks
# are used as sources for Task Monitor Triggers.
#
# Usage:
#   export STONEBRANCH_API_TOKEN="your-token"
#   terraform init
#   terraform plan
#   terraform apply

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

variable "agent_var" {
  description = "Variable containing the agent name"
  type        = string
  default     = "agent_var"
}

# Task to be monitored
resource "stonebranch_task_unix" "etl_job" {
  name        = "tf-example-etl-job"
  summary     = "ETL job to be monitored"
  agent_var   = var.agent_var
  command     = "echo 'Running ETL...'"
  exit_codes  = "0"
}

# Basic Task Monitor - watches for success status
resource "stonebranch_task_monitor" "basic" {
  name          = "tf-example-basic-monitor"
  summary       = "Basic task monitor"
  task_mon_name = stonebranch_task_unix.etl_job.name
  status_text   = "Success"
}

# Task Monitor with specific status - watch for success
resource "stonebranch_task_monitor" "success_monitor" {
  name          = "tf-example-success-monitor"
  summary       = "Monitor for task success"
  task_mon_name = stonebranch_task_unix.etl_job.name
  status_text   = "Success"
}

# Task Monitor with status and type
resource "stonebranch_task_monitor" "detailed_monitor" {
  name          = "tf-example-detailed-monitor"
  summary       = "Detailed task monitor with multiple settings"
  task_mon_name = stonebranch_task_unix.etl_job.name
  status_text   = "Finished"
  mon_type      = "Task Instance"
}

# Task to be triggered - used with Task Monitor Trigger
resource "stonebranch_task_unix" "notify_job" {
  name        = "tf-example-notify-job"
  summary     = "Notification task"
  agent_var   = var.agent_var
  command     = "echo 'ETL completed, sending notification...'"
  exit_codes  = "0"
}

# Task Monitor Trigger using the Task Monitor
resource "stonebranch_trigger_task_monitor" "on_etl_complete" {
  name         = "tf-example-on-etl-complete"
  description  = "Trigger notification when ETL job completes"
  task_monitor = stonebranch_task_monitor.success_monitor.name
  tasks        = [stonebranch_task_unix.notify_job.name]
  enabled      = false
}

# Output examples
output "monitor_id" {
  description = "System ID of the task monitor"
  value       = stonebranch_task_monitor.basic.sys_id
}

output "monitor_name" {
  description = "Name of the task monitor"
  value       = stonebranch_task_monitor.basic.name
}

output "watched_task" {
  description = "Task being monitored"
  value       = stonebranch_task_monitor.basic.task_mon_name
}

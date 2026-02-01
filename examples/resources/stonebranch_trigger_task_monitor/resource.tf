# Stonebranch Task Monitor Trigger Example
#
# This example demonstrates how to create task monitor triggers in Stonebranch.
# A task monitor trigger fires when a Task Monitor task detects that a watched
# task has reached a specific status (completed, succeeded, failed, etc.).
#
# The structure is:
# 1. Unix tasks (or other task types) that do actual work
# 2. A Task Monitor task that watches for specific task status changes
# 3. A Task Monitor Trigger that fires when the Task Monitor detects the status
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

# Task to be watched - when this completes, the monitor detects it
resource "stonebranch_task_unix" "data_extraction" {
  name        = "tf-example-data-extraction"
  summary     = "Extract data from source system"
  agent_var   = var.agent_var
  command     = "echo 'Extracting data...'"
  exit_codes  = "0"
}

# Task to be triggered when data extraction completes
resource "stonebranch_task_unix" "data_transformation" {
  name        = "tf-example-data-transformation"
  summary     = "Transform extracted data"
  agent_var   = var.agent_var
  command     = "echo 'Transforming data...'"
  exit_codes  = "0"
}

# Task Monitor task - watches for when the data extraction task finishes
resource "stonebranch_task_monitor" "extraction_monitor" {
  name          = "tf-example-extraction-monitor"
  summary       = "Monitor data extraction completion"
  task_mon_name = stonebranch_task_unix.data_extraction.name
}

# Task Monitor Trigger - fires when the Task Monitor detects completion
resource "stonebranch_trigger_task_monitor" "after_extraction" {
  name         = "tf-example-after-extraction-trigger"
  description  = "Trigger transformation after data extraction completes"
  task_monitor = stonebranch_task_monitor.extraction_monitor.name
  tasks        = [stonebranch_task_unix.data_transformation.name]

  # Triggers are created disabled by default
  enabled = false
}

# Example: Multiple tasks triggered from one monitored task
resource "stonebranch_task_unix" "send_notification" {
  name        = "tf-example-send-notification"
  summary     = "Send completion notification"
  agent_var   = var.agent_var
  command     = "echo 'Sending notification...'"
  exit_codes  = "0"
}

resource "stonebranch_task_unix" "update_dashboard" {
  name        = "tf-example-update-dashboard"
  summary     = "Update monitoring dashboard"
  agent_var   = var.agent_var
  command     = "echo 'Updating dashboard...'"
  exit_codes  = "0"
}

# Task Monitor for transformation task
resource "stonebranch_task_monitor" "transformation_monitor" {
  name          = "tf-example-transformation-monitor"
  summary       = "Monitor data transformation completion"
  task_mon_name = stonebranch_task_unix.data_transformation.name
}

resource "stonebranch_trigger_task_monitor" "multi_task_trigger" {
  name         = "tf-example-multi-task-trigger"
  description  = "Trigger multiple tasks when transformation completes"
  task_monitor = stonebranch_task_monitor.transformation_monitor.name

  # Multiple tasks can be triggered
  tasks = [
    stonebranch_task_unix.send_notification.name,
    stonebranch_task_unix.update_dashboard.name
  ]
}

# Example: Trigger with business service association
resource "stonebranch_business_service" "etl" {
  name        = "tf-example-etl-service"
  description = "ETL Pipeline Service"
}

resource "stonebranch_trigger_task_monitor" "with_business_service" {
  name         = "tf-example-bs-trigger"
  description  = "Task monitor trigger with business service"
  task_monitor = stonebranch_task_monitor.extraction_monitor.name
  tasks        = [stonebranch_task_unix.data_transformation.name]

  opswise_groups = [stonebranch_business_service.etl.name]
}

# Output examples
output "trigger_id" {
  description = "System ID of the task monitor trigger"
  value       = stonebranch_trigger_task_monitor.after_extraction.sys_id
}

output "trigger_name" {
  description = "Name of the task monitor trigger"
  value       = stonebranch_trigger_task_monitor.after_extraction.name
}

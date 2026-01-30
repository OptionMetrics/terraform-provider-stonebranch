# Stonebranch File Monitor Trigger Example
#
# This example demonstrates how to create file monitor trigger resources in Stonebranch.
# A file monitor trigger executes tasks when a file monitor task detects file events
# (file created, modified, deleted, etc.).
#
# IMPORTANT: The task_monitor field must reference an existing file monitor task
# in your StoneBranch environment. File monitor tasks are a separate resource type
# that monitors file system events on an agent.
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

# Task to execute when file is detected
resource "stonebranch_task_unix" "process_file" {
  name       = "process-incoming-file"
  summary    = "Processes incoming data files"
  command    = "/opt/scripts/process_file.sh"
  agent      = var.agent_name
  exit_codes = "0"
}

# Basic file monitor trigger
# This trigger fires when the referenced file monitor task detects a file event
resource "stonebranch_trigger_file_monitor" "incoming_files" {
  name         = "incoming-file-trigger"
  description  = "Triggers when new files arrive in the incoming directory"
  task_monitor = var.file_monitor_task_name  # Name of existing file monitor task
  tasks        = [stonebranch_task_unix.process_file.name]
}

# File monitor trigger with time restrictions
# Only fires during business hours
resource "stonebranch_trigger_file_monitor" "business_hours_only" {
  name             = "business-hours-file-trigger"
  description      = "File trigger active only during business hours"
  task_monitor     = var.file_monitor_task_name
  tasks            = [stonebranch_task_unix.process_file.name]
  time_zone        = "America/New_York"
  restricted_times = true
  enabled_start    = "08:00"
  enabled_end      = "18:00"
}

# File monitor trigger with business service assignment
resource "stonebranch_trigger_file_monitor" "with_business_service" {
  name           = "production-file-trigger"
  description    = "Production file monitoring"
  task_monitor   = var.file_monitor_task_name
  tasks          = [stonebranch_task_unix.process_file.name]
  opswise_groups = ["Production Services"]
}

# Stonebranch File Monitor Task Example
#
# This example demonstrates how to create file monitor task resources in Stonebranch.
# File monitor tasks watch for file system events such as file creation, modification,
# or deletion on an agent.
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

# Basic file monitor - watches for CSV files
resource "stonebranch_task_file_monitor" "incoming_csv" {
  name      = "monitor-incoming-csv"
  summary   = "Monitors for incoming CSV files"
  file_name = "/data/incoming/*.csv"
  agent     = var.agent_name
}

# File monitor with stability check
# Waits for file to be stable (unchanged) for 30 seconds before triggering
resource "stonebranch_task_file_monitor" "stable_files" {
  name           = "monitor-stable-files"
  summary        = "Waits for files to be stable before triggering"
  file_name      = "/data/uploads/*"
  agent          = var.agent_name
  stable_seconds = 30
}

# Recursive file monitor
resource "stonebranch_task_file_monitor" "recursive_monitor" {
  name      = "monitor-recursive"
  summary   = "Recursively monitors directory tree"
  file_name = "/data/input/*"
  agent     = var.agent_name
  recursive = true
}

# File exists monitor - checks if a specific file exists
# fm_type values: Created (default), Deleted, Changed, Exist, Missing
resource "stonebranch_task_file_monitor" "exists_check" {
  name      = "check-ready-flag"
  summary   = "Checks if ready.flag file exists"
  file_name = "/data/ready.flag"
  agent     = var.agent_name
  fm_type   = "Exist"
}

# File monitor with regex pattern
resource "stonebranch_task_file_monitor" "regex_pattern" {
  name      = "monitor-dated-files"
  summary   = "Monitors files matching date pattern"
  file_name = "/data/reports/report_[0-9]{8}\\.csv"
  agent     = var.agent_name
  use_regex = true
}

# Trigger that fires when file monitor detects files
resource "stonebranch_trigger_file_monitor" "process_incoming" {
  name         = "trigger-on-csv-arrival"
  description  = "Triggers processing when CSV files arrive"
  task_monitor = stonebranch_task_file_monitor.incoming_csv.name
  tasks        = ["process-csv-task"]  # Reference existing task
}

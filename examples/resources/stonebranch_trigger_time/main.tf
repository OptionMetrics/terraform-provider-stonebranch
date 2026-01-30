# Stonebranch Time Trigger Example
#
# This example demonstrates how to create time-based triggers in Stonebranch.
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

# Daily trigger at 9:00 AM
resource "stonebranch_trigger_time" "daily" {
  name        = "tf-example-daily-trigger"
  description = "Triggers daily at 9:00 AM"

  tasks = [stonebranch_task_unix.daily_job.name]

  time      = "09:00"
  time_zone = var.time_zone

  # Triggers are disabled by default for safety
  enabled = false
}

# Hourly trigger
resource "stonebranch_trigger_time" "hourly" {
  name        = "tf-example-hourly-trigger"
  description = "Triggers every hour"

  tasks = [stonebranch_task_unix.hourly_job.name]

  time       = "00:00"
  time_zone  = "UTC"
  time_style = "Interval"

  time_interval       = 1
  time_interval_units = "Hours"

  enabled = false
}

# Weekday-only trigger (Monday-Friday at 8:00 AM)
resource "stonebranch_trigger_time" "weekdays" {
  name        = "tf-example-weekday-trigger"
  description = "Triggers at 8:00 AM on weekdays only"

  tasks = [stonebranch_task_unix.business_job.name]

  time      = "08:00"
  time_zone = var.time_zone

  monday    = true
  tuesday   = true
  wednesday = true
  thursday  = true
  friday    = true
  saturday  = false
  sunday    = false

  enabled = false
}

# Supporting task resources
resource "stonebranch_task_unix" "daily_job" {
  name       = "tf-example-daily-job"
  summary    = "Daily batch processing job"
  agent_var  = var.agent_var
  command    = "echo 'Running daily job at $(date)'"
  exit_codes = "0"
}

resource "stonebranch_task_unix" "hourly_job" {
  name       = "tf-example-hourly-job"
  summary    = "Hourly monitoring job"
  agent_var  = var.agent_var
  command    = "echo 'Running hourly check at $(date)'"
  exit_codes = "0"
}

resource "stonebranch_task_unix" "business_job" {
  name       = "tf-example-business-job"
  summary    = "Business hours job"
  agent_var  = var.agent_var
  command    = "echo 'Running business job at $(date)'"
  exit_codes = "0"
}

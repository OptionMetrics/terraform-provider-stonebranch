# Stonebranch Cron Trigger Example
#
# This example demonstrates how to create cron trigger resources in Stonebranch.
# Cron triggers execute tasks based on cron expressions for flexible scheduling.
#
# Cron Expression Format:
#   minutes (0-59)
#   hours (0-23)
#   day_of_month (1-31)
#   month (1-12 or JAN-DEC)
#   day_of_week (0-6 or SUN-SAT, 0=Sunday)
#
# Special characters:
#   * - any value
#   , - value list separator
#   - - range of values
#   / - step values
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

# Example task to trigger
resource "stonebranch_task_unix" "example" {
  name       = "tf-example-cron-task"
  summary    = "Task triggered by cron"
  command    = "echo 'Running scheduled job at $(date)'"
  agent_var  = "agent_name"
  exit_codes = "0"
}

# Every day at noon
resource "stonebranch_trigger_cron" "daily_noon" {
  name         = "tf-example-cron-daily-noon"
  description  = "Runs every day at 12:00 PM"
  minutes      = "0"
  hours        = "12"
  day_of_month = "*"
  month        = "*"
  day_of_week  = "*"
  tasks        = [stonebranch_task_unix.example.name]
}

# Every weekday at 9 AM (Monday-Friday)
resource "stonebranch_trigger_cron" "weekdays_9am" {
  name         = "tf-example-cron-weekdays"
  description  = "Runs at 9:00 AM Monday through Friday"
  minutes      = "0"
  hours        = "9"
  day_of_month = "*"
  month        = "*"
  day_of_week  = "1-5"
  tasks        = [stonebranch_task_unix.example.name]
}

# Every 15 minutes
resource "stonebranch_trigger_cron" "every_15_minutes" {
  name         = "tf-example-cron-15min"
  description  = "Runs every 15 minutes"
  minutes      = "0,15,30,45"
  hours        = "*"
  day_of_month = "*"
  month        = "*"
  day_of_week  = "*"
  tasks        = [stonebranch_task_unix.example.name]
}

# First day of every month at midnight
resource "stonebranch_trigger_cron" "monthly" {
  name         = "tf-example-cron-monthly"
  description  = "Runs at midnight on the 1st of each month"
  minutes      = "0"
  hours        = "0"
  day_of_month = "1"
  month        = "*"
  day_of_week  = "*"
  tasks        = [stonebranch_task_unix.example.name]
}

# Every hour during business hours (8 AM - 6 PM)
resource "stonebranch_trigger_cron" "business_hours" {
  name         = "tf-example-cron-business"
  description  = "Runs hourly during business hours"
  minutes      = "0"
  hours        = "8-18"
  day_of_month = "*"
  month        = "*"
  day_of_week  = "1-5"
  time_zone    = "America/New_York"
  tasks        = [stonebranch_task_unix.example.name]
}

# Sunday maintenance window at 2 AM
resource "stonebranch_trigger_cron" "sunday_maintenance" {
  name         = "tf-example-cron-maintenance"
  description  = "Sunday maintenance window at 2 AM"
  minutes      = "0"
  hours        = "2"
  day_of_month = "*"
  month        = "*"
  day_of_week  = "0"
  time_zone    = "UTC"
  tasks        = [stonebranch_task_unix.example.name]
  enabled      = false  # Disabled by default for safety
}

# Quarter-end processing (last day of March, June, September, December)
resource "stonebranch_trigger_cron" "quarter_end" {
  name         = "tf-example-cron-quarter"
  description  = "Runs on the last day of each quarter"
  minutes      = "0"
  hours        = "23"
  day_of_month = "L"
  month        = "3,6,9,12"
  day_of_week  = "*"
  tasks        = [stonebranch_task_unix.example.name]
}

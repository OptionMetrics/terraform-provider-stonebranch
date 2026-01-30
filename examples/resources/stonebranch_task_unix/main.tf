# Stonebranch Unix Task Example
#
# This example demonstrates how to create Unix/Linux tasks in Stonebranch.
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
  # Optionally set base_url if not using the default
}

# Simple Unix task that runs a command
resource "stonebranch_task_unix" "hello" {
  name    = "tf-example-unix-hello"
  summary = "Simple Unix task that echoes a message"

  agent_var = var.agent_var

  command    = "echo 'Hello from Terraform!'"
  exit_codes = "0"
}

# Unix task that runs a script resource
resource "stonebranch_task_unix" "with_script" {
  name    = "tf-example-unix-script"
  summary = "Unix task that executes a script resource"

  agent_var = var.agent_var

  command_or_script = "Script"
  script            = stonebranch_script.backup.name

  exit_codes = "0"
}

# Unix task with retry configuration
resource "stonebranch_task_unix" "with_retry" {
  name    = "tf-example-unix-retry"
  summary = "Unix task with automatic retry on failure"

  agent_var = var.agent_var

  command = "/bin/false"  # This will fail, demonstrating retry

  retry_maximum  = 3
  retry_interval = 60

  exit_codes = "0"
}

# Supporting script resource
resource "stonebranch_script" "backup" {
  name    = "tf-example-backup-script"
  content = <<-EOT
    #!/bin/bash
    echo "Starting backup at $(date)"
    echo "Backup completed"
  EOT
}

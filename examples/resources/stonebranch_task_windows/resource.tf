# Stonebranch Windows Task Example
#
# This example demonstrates how to create Windows tasks in Stonebranch.
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

# Simple Windows task that runs a command
resource "stonebranch_task_windows" "hello" {
  name    = "tf-example-windows-hello"
  summary = "Simple Windows task that echoes a message"

  agent_var = var.agent_var

  command    = "echo Hello from Terraform!"
  exit_codes = "0"
}

# Windows task that runs a batch script
resource "stonebranch_task_windows" "with_script" {
  name    = "tf-example-windows-script"
  summary = "Windows task that executes a batch script"

  agent_var = var.agent_var

  command_or_script = "Script"
  script            = stonebranch_script.windows_batch.name

  exit_codes = "0"
}

# Windows PowerShell task
resource "stonebranch_task_windows" "powershell" {
  name    = "tf-example-windows-powershell"
  summary = "Windows task running PowerShell"

  agent_var = var.agent_var

  command    = "powershell.exe -Command \"Get-Date; Write-Host 'Hello from PowerShell'\""
  exit_codes = "0"
}

# Supporting batch script resource
resource "stonebranch_script" "windows_batch" {
  name    = "tf-example-windows-batch"
  content = <<-EOT
    @echo off
    echo Starting batch script...
    echo Current date: %date%
    echo Current time: %time%
    echo Batch script completed
  EOT
}

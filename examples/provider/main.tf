# When using dev overrides, skip "terraform init" and run "terraform plan" directly.
#
# Usage:
#   cd examples/provider
#   export STONEBRANCH_API_TOKEN="your-token"
#   terraform plan

terraform {
  required_providers {
    stonebranch = {
      source = "stonebranch/stonebranch"
    }
  }
}

provider "stonebranch" {
  # api_token = var.stonebranch_token  # or set STONEBRANCH_API_TOKEN env var
  # base_url  = "https://optionmetricsdev.stonebranch.cloud"  # optional, this is the default
}

# Example: Simple Unix task that runs a command
resource "stonebranch_task_unix" "hello_world" {
  name    = "terraform-hello-world"
  summary = "A simple task managed by Terraform"

  # Agent configuration - ONE of these is REQUIRED:
  agent = var.agent_name  # or use agent_cluster for cluster

  # Command to execute
  command = "echo 'Hello from Terraform!'"

  # Exit code handling (optional)
  exit_codes           = "0"
  exit_code_processing = "Success Exitcode Range"
}

variable "agent_name" {
  description = "Name of the StoneBranch agent to run tasks on"
  type        = string
  default     = "DEV_UA_CLOUD_LINUX_UE1_02"
}

# Example: Script resource
resource "stonebranch_script" "backup_script" {
  name        = "terraform-backup-script"
  description = "A backup script managed by Terraform"
  content     = <<-EOT
    #!/bin/bash
    echo "Starting backup..."
    date
    echo "Backup completed"
  EOT
}

# Example: Unix task that references a script resource
resource "stonebranch_task_unix" "script_task" {
  name              = "terraform-script-task"
  summary           = "Task that runs a script resource"
  command_or_script = "Script"
  script            = stonebranch_script.backup_script.name
  agent             = var.agent_name
  exit_codes        = "0"
}

# Example: Task with retry configuration
# resource "stonebranch_task_unix" "with_retry" {
#   name          = "terraform-retry-example"
#   summary       = "Task with retry configuration"
#   command       = "/opt/scripts/important-job.sh"
#   agent         = "your-agent-name"
#
#   retry_maximum  = 3
#   retry_interval = 300  # 5 minutes
# }

# Example: Time trigger that runs a task daily at 9:00 AM
resource "stonebranch_trigger_time" "daily_backup" {
  name        = "terraform-daily-backup-trigger"
  description = "Triggers the backup task every day at 9:00 AM"
  enabled     = false  # Set to true to activate

  # Reference the task(s) to run
  tasks = [stonebranch_task_unix.script_task.name]

  # Schedule configuration
  time      = "09:00"
  time_zone = "America/New_York"
}

# Example: Credential resource for SFTP authentication
resource "stonebranch_credential" "sftp_creds" {
  name             = "terraform-sftp-credentials"
  description      = "SFTP credentials managed by Terraform"
  runtime_user     = "sftpuser"
  runtime_password = var.sftp_password  # Use a variable for sensitive data
}

variable "sftp_password" {
  description = "Password for SFTP connection"
  type        = string
  sensitive   = true
  default     = "placeholder"  # Set via environment or tfvars
}

# Example: File Transfer task for SFTP download
resource "stonebranch_task_file_transfer" "download_report" {
  name    = "terraform-download-report"
  summary = "Download daily report from SFTP server"

  agent = var.agent_name

  # Transfer settings (GET = download from remote, PUT = upload to remote)
  transfer_direction = "GET"
  server_type        = "SFTP"

  # Remote server configuration
  remote_server   = "sftp.example.com"
  remote_filename = "/reports/daily_report.csv"

  # Local destination
  local_filename = "/data/reports/daily_report.csv"

  # Reference the credential resource
  remote_credentials = stonebranch_credential.sftp_creds.name
}

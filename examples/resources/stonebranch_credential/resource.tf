# Stonebranch Credential Example
#
# This example demonstrates how to create credential resources in Stonebranch.
#
# Usage:
#   export STONEBRANCH_API_TOKEN="your-token"
#   export TF_VAR_app_password="your-app-password"
#   export TF_VAR_sftp_password="your-sftp-password"
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

# Basic application credential
resource "stonebranch_credential" "app_user" {
  name             = "tf-example-app-creds"
  description      = "Application service account credentials"
  runtime_user     = var.app_user
  runtime_password = var.app_password
}

# SFTP credential
resource "stonebranch_credential" "sftp" {
  name             = "tf-example-sftp-creds"
  description      = "SFTP server credentials"
  runtime_user     = var.sftp_user
  runtime_password = var.sftp_password
}

# Example: Using credential in a Unix task
resource "stonebranch_task_unix" "with_creds" {
  name    = "tf-example-task-with-creds"
  summary = "Task that runs with specific credentials"

  agent_var = var.agent_var

  command     = "/opt/app/secure_process.sh"
  credentials = stonebranch_credential.app_user.name

  exit_codes = "0"
}

# Example: Using credential in a file transfer task
resource "stonebranch_task_file_transfer" "with_creds" {
  name    = "tf-example-transfer-with-creds"
  summary = "File transfer using SFTP credentials"

  agent_var = var.agent_var

  server_type        = "SFTP"
  remote_server      = "sftp.example.com"
  remote_filename    = "/data/file.csv"
  local_filename     = "/local/file.csv"
  remote_credentials = stonebranch_credential.sftp.name
}

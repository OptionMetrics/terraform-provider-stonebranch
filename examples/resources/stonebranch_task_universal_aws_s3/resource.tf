# Stonebranch Universal Task - AWS S3 Example
#
# This example demonstrates how to create AWS S3 tasks using the
# CS AWS S3 Universal Template in Stonebranch.
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

# List all S3 buckets
resource "stonebranch_task_universal_aws_s3" "list_buckets" {
  name    = "tf-example-s3-list-buckets"
  summary = "List all S3 buckets"

  agent_var = var.agent_var

  action = "list-buckets"

  aws_access_key_id     = var.aws_credentials_name
  aws_secret_access_key = var.aws_credentials_name
  aws_default_region    = "us-east-1"
}

# List objects in a bucket
resource "stonebranch_task_universal_aws_s3" "list_objects" {
  name    = "tf-example-s3-list-objects"
  summary = "List objects in an S3 bucket"

  agent_var = var.agent_var

  action = "list-objects"
  bucket = "my-bucket-name"
  prefix = "logs/"

  aws_access_key_id     = var.aws_credentials_name
  aws_secret_access_key = var.aws_credentials_name
  aws_default_region    = "us-east-1"

  show_details = true
}

# Upload a file to S3
# Note: For upload-file action, use 'prefix' for the S3 path, not 's3_key'
resource "stonebranch_task_universal_aws_s3" "upload_file" {
  name    = "tf-example-s3-upload"
  summary = "Upload a file to S3"

  agent_var = var.agent_var

  action     = "upload-file"
  bucket     = "my-bucket-name"
  sourcefile = "/path/to/local/file.txt"
  prefix     = "uploads/"  # S3 key prefix for the uploaded file

  upload_write_options = "True"  # Overwrite if exists
  acl                  = "private"

  aws_access_key_id     = var.aws_credentials_name
  aws_secret_access_key = var.aws_credentials_name
  aws_default_region    = "us-east-1"
}

# Download a file from S3
resource "stonebranch_task_universal_aws_s3" "download_file" {
  name    = "tf-example-s3-download"
  summary = "Download a file from S3"

  agent_var = var.agent_var

  action           = "download-file"
  bucket           = "my-bucket-name"
  s3_key           = "data/report.csv"
  target_directory = "/tmp/downloads"

  download_write_options = "False"  # Skip if exists

  aws_access_key_id     = var.aws_credentials_name
  aws_secret_access_key = var.aws_credentials_name
  aws_default_region    = "us-east-1"
}

# Copy object between buckets
resource "stonebranch_task_universal_aws_s3" "copy_object" {
  name    = "tf-example-s3-copy"
  summary = "Copy object to another bucket"

  agent_var = var.agent_var

  action        = "copy-object-to-bucket"
  bucket        = "source-bucket"
  s3_key        = "data/file.txt"
  target_bucket = "destination-bucket"
  target_s3_key = "archive/file.txt"
  operation     = "copy"

  aws_access_key_id     = var.aws_credentials_name
  aws_secret_access_key = var.aws_credentials_name
  aws_default_region    = "us-east-1"
}

# Using IAM role-based access (no credentials needed if running on EC2 with role)
resource "stonebranch_task_universal_aws_s3" "with_role" {
  name    = "tf-example-s3-with-role"
  summary = "S3 task using IAM role"

  agent_var = var.agent_var

  action = "list-buckets"

  role_based_access  = "yes"
  role_arn           = "arn:aws:iam::123456789012:role/S3AccessRole"
  service_name       = "sts"
  aws_default_region = "us-east-1"
}

# Monitor for object existence
resource "stonebranch_task_universal_aws_s3" "monitor_object" {
  name    = "tf-example-s3-monitor"
  summary = "Monitor for object arrival in S3"

  agent_var = var.agent_var

  action   = "monitor-object"
  bucket   = "my-bucket-name"
  prefix   = "incoming/"
  interval = "60"  # Check every 60 seconds

  aws_access_key_id     = var.aws_credentials_name
  aws_secret_access_key = var.aws_credentials_name
  aws_default_region    = "us-east-1"

  log_level = "DEBUG"
}

variable "agent_var" {
  description = "Variable containing the agent name"
  type        = string
  default     = "ops_agent_name"
}

variable "aws_credentials_name" {
  description = "Name of the Stonebranch credential containing AWS access keys"
  type        = string
  default     = "aws-credentials"
}

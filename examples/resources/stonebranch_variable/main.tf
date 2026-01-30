# Stonebranch Variable Example
#
# This example demonstrates how to create global variable resources in Stonebranch.
#
# Variable naming rules:
# - Must begin with a letter
# - Alphanumerics (upper or lower case) and underscore only
# - No hyphens, spaces, or special characters
# - Names are not case-sensitive
# - Do not use the prefix "ops_" (reserved for built-in variables)
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

# Simple string variable
resource "stonebranch_variable" "environment" {
  name        = "APP_ENVIRONMENT"
  value       = "production"
  description = "Current application environment"
}

# Database host variable
resource "stonebranch_variable" "db_host" {
  name        = "DATABASE_HOST"
  value       = var.database_host
  description = "Primary database hostname"
}

# Feature flag variable
resource "stonebranch_variable" "feature_flag" {
  name        = "ENABLE_NEW_FEATURE"
  value       = "true"
  description = "Feature flag for new functionality"
}

# Variable with business service assignment
resource "stonebranch_variable" "api_endpoint" {
  name          = "API_ENDPOINT"
  value         = "https://api.example.com/v1"
  description   = "External API endpoint URL"
  opswise_groups = ["Production Services"]
}

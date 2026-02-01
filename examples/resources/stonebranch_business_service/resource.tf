# Stonebranch Business Service Example
#
# This example demonstrates how to create business service resources in Stonebranch.
# Business services are used to group and organize resources such as tasks, triggers, and variables.
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

# Simple business service
resource "stonebranch_business_service" "production" {
  name        = "Production Services"
  description = "Business service for production workloads"
}

# Business service for development
resource "stonebranch_business_service" "development" {
  name        = "Development Services"
  description = "Business service for development and testing"
}

# Business service for data pipelines
resource "stonebranch_business_service" "data_pipelines" {
  name        = "Data Pipelines"
  description = "Business service for ETL and data processing tasks"
}

# Example: Assigning a variable to a business service
resource "stonebranch_variable" "app_env" {
  name           = "APP_ENVIRONMENT"
  value          = "production"
  description    = "Current application environment"
  opswise_groups = [stonebranch_business_service.production.name]
}

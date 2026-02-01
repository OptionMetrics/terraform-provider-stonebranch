# Stonebranch Agent Cluster Example
#
# This example demonstrates how to create agent clusters in Stonebranch.
# Agent clusters group multiple agents for load distribution and high availability.
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

# Basic Linux/Unix agent cluster with default settings
resource "stonebranch_agent_cluster" "linux_basic" {
  name = "tf-example-linux-cluster"
  type = "Linux/Unix"

  description = "Basic Linux agent cluster managed by Terraform"
}

# Windows agent cluster
resource "stonebranch_agent_cluster" "windows_basic" {
  name = "tf-example-windows-cluster"
  type = "Windows"

  description = "Basic Windows agent cluster managed by Terraform"
}

# Agent cluster with round-robin distribution
resource "stonebranch_agent_cluster" "round_robin" {
  name         = "tf-example-round-robin-cluster"
  type         = "Linux/Unix"
  distribution = "Round Robin"

  description = "Cluster that distributes tasks evenly across agents"
}

# Agent cluster with CPU-based distribution
resource "stonebranch_agent_cluster" "cpu_based" {
  name         = "tf-example-cpu-cluster"
  type         = "Linux/Unix"
  distribution = "Lowest CPU Utilization"

  description = "Cluster that sends tasks to least loaded agent"
}

# Agent cluster with task execution limits
resource "stonebranch_agent_cluster" "limited" {
  name        = "tf-example-limited-cluster"
  type        = "Linux/Unix"
  description = "Cluster with task execution limits"

  # Limit total concurrent tasks across all agents in cluster
  limit_type   = "Limited"
  limit_amount = 10

  # Limit concurrent tasks per individual agent
  agent_limit_type   = "Limited"
  agent_limit_amount = 2
}

# Agent cluster with business service association
resource "stonebranch_business_service" "production" {
  name        = "tf-example-production"
  description = "Production workloads"
}

resource "stonebranch_agent_cluster" "production" {
  name        = "tf-example-production-cluster"
  type        = "Linux/Unix"
  description = "Production agent cluster"

  distribution = "Lowest CPU Utilization"

  # Associate with business service
  opswise_groups = [stonebranch_business_service.production.name]

  # Skip inactive/suspended agents
  ignore_inactive_agents  = true
  ignore_suspended_agents = true
}

# Output examples
output "linux_cluster_id" {
  description = "System ID of the Linux cluster"
  value       = stonebranch_agent_cluster.linux_basic.sys_id
}

output "production_cluster_name" {
  description = "Name of the production cluster (for use in tasks)"
  value       = stonebranch_agent_cluster.production.name
}

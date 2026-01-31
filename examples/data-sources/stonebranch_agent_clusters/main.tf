# Stonebranch Agent Clusters Data Source Example
#
# This example demonstrates how to query agent clusters from Stonebranch.
#
# Usage:
#   export STONEBRANCH_API_TOKEN="your-token"
#   terraform init
#   terraform plan

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

# List all agent clusters
data "stonebranch_agent_clusters" "all" {
}

# List only Linux/Unix agent clusters
data "stonebranch_agent_clusters" "unix" {
  type = "Linux/Unix"
}

# List only Windows agent clusters
data "stonebranch_agent_clusters" "windows" {
  type = "Windows"
}

# List agent clusters in a specific business service
data "stonebranch_agent_clusters" "production" {
  business_services = "Production"
}

# Output examples
output "all_cluster_count" {
  description = "Total number of agent clusters"
  value       = length(data.stonebranch_agent_clusters.all.agent_clusters)
}

output "unix_clusters" {
  description = "List of Linux/Unix agent cluster names"
  value       = [for cluster in data.stonebranch_agent_clusters.unix.agent_clusters : cluster.name]
}

# Example: Use agent cluster data in a task
resource "stonebranch_task_unix" "example" {
  count = length(data.stonebranch_agent_clusters.unix.agent_clusters) > 0 ? 1 : 0

  name          = "tf-example-using-cluster-data"
  summary       = "Task using agent cluster discovered via data source"
  agent_cluster = data.stonebranch_agent_clusters.unix.agent_clusters[0].name
  command       = "echo 'Running on cluster discovered via data source'"
}

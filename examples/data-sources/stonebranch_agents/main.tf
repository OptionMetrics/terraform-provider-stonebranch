# Stonebranch Agents Data Source Example
#
# This example demonstrates how to query agents from Stonebranch.
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

# List all agents
data "stonebranch_agents" "all" {
}

# List only Linux/Unix agents
data "stonebranch_agents" "unix" {
  type = "Linux/Unix"
}

# List only Windows agents
data "stonebranch_agents" "windows" {
  type = "Windows"
}

# List agents in a specific business service
data "stonebranch_agents" "production" {
  business_services = "Production"
}

# List agents by name pattern (wildcard)
data "stonebranch_agents" "web_servers" {
  name = "web-*"
}

# Output examples
output "all_agent_count" {
  description = "Total number of agents"
  value       = length(data.stonebranch_agents.all.agents)
}

output "unix_agents" {
  description = "List of Linux/Unix agent names"
  value       = [for agent in data.stonebranch_agents.unix.agents : agent.name]
}

output "first_unix_agent" {
  description = "Name of the first Unix agent (for use in tasks)"
  value       = length(data.stonebranch_agents.unix.agents) > 0 ? data.stonebranch_agents.unix.agents[0].name : null
}

# Example: Use agent data in a task
resource "stonebranch_task_unix" "example" {
  count = length(data.stonebranch_agents.unix.agents) > 0 ? 1 : 0

  name    = "tf-example-using-agent-data"
  summary = "Task using agent discovered via data source"
  agent   = data.stonebranch_agents.unix.agents[0].name
  command = "echo 'Running on agent discovered via data source'"
}

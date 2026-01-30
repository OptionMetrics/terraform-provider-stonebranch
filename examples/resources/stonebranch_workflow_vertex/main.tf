# Stonebranch Workflow Vertex Example
#
# This example demonstrates how to add tasks to a workflow as vertices.
# A vertex represents a task instance within a workflow.
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

# Create a workflow
resource "stonebranch_task_workflow" "example" {
  name    = "tf-example-workflow-with-tasks"
  summary = "Workflow demonstrating vertex management"
}

# Create tasks to add to the workflow
resource "stonebranch_task_unix" "step1" {
  name       = "tf-example-step1"
  summary    = "First step in the workflow"
  agent_var  = "agent_name"
  command    = "echo 'Starting workflow...'"
  exit_codes = "0"
}

resource "stonebranch_task_unix" "step2" {
  name       = "tf-example-step2"
  summary    = "Second step in the workflow"
  agent_var  = "agent_name"
  command    = "echo 'Processing data...'"
  exit_codes = "0"
}

resource "stonebranch_task_unix" "step3" {
  name       = "tf-example-step3"
  summary    = "Final step in the workflow"
  agent_var  = "agent_name"
  command    = "echo 'Workflow complete!'"
  exit_codes = "0"
}

# Add tasks to the workflow as vertices
resource "stonebranch_workflow_vertex" "step1" {
  workflow_name = stonebranch_task_workflow.example.name
  task_name     = stonebranch_task_unix.step1.name
  vertex_x      = "100"
  vertex_y      = "100"
}

resource "stonebranch_workflow_vertex" "step2" {
  workflow_name = stonebranch_task_workflow.example.name
  task_name     = stonebranch_task_unix.step2.name
  vertex_x      = "100"
  vertex_y      = "200"
}

resource "stonebranch_workflow_vertex" "step3" {
  workflow_name = stonebranch_task_workflow.example.name
  task_name     = stonebranch_task_unix.step3.name
  vertex_x      = "100"
  vertex_y      = "300"
}

# Vertex with an alias (useful when same task appears multiple times)
resource "stonebranch_task_unix" "reusable" {
  name       = "tf-example-reusable-task"
  summary    = "A task that can be used multiple times"
  agent_var  = "agent_name"
  command    = "echo 'Reusable task'"
  exit_codes = "0"
}

resource "stonebranch_workflow_vertex" "reusable_first" {
  workflow_name = stonebranch_task_workflow.example.name
  task_name     = stonebranch_task_unix.reusable.name
  alias         = "FirstInstance"
  vertex_x      = "300"
  vertex_y      = "100"
}

resource "stonebranch_workflow_vertex" "reusable_second" {
  workflow_name = stonebranch_task_workflow.example.name
  task_name     = stonebranch_task_unix.reusable.name
  alias         = "SecondInstance"
  vertex_x      = "300"
  vertex_y      = "200"
}

# Output the vertex IDs (useful for creating edges)
output "step1_vertex_id" {
  value = stonebranch_workflow_vertex.step1.vertex_id
}

output "step2_vertex_id" {
  value = stonebranch_workflow_vertex.step2.vertex_id
}

output "step3_vertex_id" {
  value = stonebranch_workflow_vertex.step3.vertex_id
}

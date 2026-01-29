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

# Example: Task with script content instead of command
# resource "stonebranch_task_unix" "script_example" {
#   name              = "terraform-script-example"
#   summary           = "Task that runs a script"
#   command_or_script = "Script"
#   script            = <<-EOT
#     #!/bin/bash
#     echo "Starting script..."
#     date
#     echo "Script completed"
#   EOT
#   agent = "your-agent-name"
# }

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

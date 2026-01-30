# Terraform Provider for StoneBranch

A Terraform provider for managing resources in [StoneBranch Universal Controller](https://www.stonebranch.com/).

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.23
- A StoneBranch Universal Controller instance with API access
- A valid API token with appropriate permissions

## Building the Provider

```bash
git clone <repository-url>
cd terraform-provider-stonebranch
make build
```

This creates the `terraform-provider-stonebranch` binary in the project root.

## Local Development Setup

### 1. Configure Development Overrides

Create or edit `~/.terraformrc` to point Terraform to your local build:

```hcl
provider_installation {
  dev_overrides {
    "registry.terraform.io/stonebranch/stonebranch" = "/path/to/terraform-provider-stonebranch"
  }
  direct {}
}
```

Or use the provided dev config:

```bash
export TF_CLI_CONFIG_FILE=/path/to/terraform-provider-stonebranch/examples/dev.tfrc
```

### 2. Set Environment Variables

```bash
export STONEBRANCH_API_TOKEN="your-bearer-token"
export STONEBRANCH_BASE_URL="https://your-instance.stonebranch.cloud"  # optional
```

### 3. Build and Test

```bash
# Build the provider
make build

# Navigate to example directory
cd examples/provider

# IMPORTANT: Skip "terraform init" when using dev overrides!
# Just run plan/apply directly:

# Plan changes
terraform plan

# Apply changes
terraform apply
```

## Provider Configuration

```hcl
terraform {
  required_providers {
    stonebranch = {
      source = "registry.terraform.io/stonebranch/stonebranch"
    }
  }
}

provider "stonebranch" {
  # API token for authentication (required)
  # Can also use STONEBRANCH_API_TOKEN environment variable
  api_token = var.stonebranch_token

  # Base URL for the StoneBranch API (optional)
  # Can also use STONEBRANCH_BASE_URL environment variable
  # Defaults to: https://optionmetricsdev.stonebranch.cloud
  base_url = "https://your-instance.stonebranch.cloud"
}
```

## Authentication

The provider uses Bearer token authentication. Obtain your token from the StoneBranch Universal Controller:

1. Log into your StoneBranch instance
2. Navigate to user settings or API token management
3. Generate or copy your API token

You can provide the token via:
- The `api_token` provider attribute
- The `STONEBRANCH_API_TOKEN` environment variable (recommended for security)

## Resources

- [stonebranch_task_unix](#stonebranch_task_unix) - Unix/Linux command tasks
- [stonebranch_task_windows](#stonebranch_task_windows) - Windows command tasks
- [stonebranch_task_file_transfer](#stonebranch_task_file_transfer) - File transfer tasks
- [stonebranch_script](#stonebranch_script) - Reusable scripts
- [stonebranch_trigger_time](#stonebranch_trigger_time) - Time-based triggers
- [stonebranch_credential](#stonebranch_credential) - Authentication credentials

### stonebranch_task_unix

Manages a StoneBranch Unix/Linux task.

#### Example Usage

```hcl
# Simple task with a command
resource "stonebranch_task_unix" "hello" {
  name    = "terraform-hello-world"
  summary = "A simple task managed by Terraform"
  command = "echo 'Hello from Terraform!'"
  agent   = "my-linux-agent"
}

# Task with script content
resource "stonebranch_task_unix" "script" {
  name              = "terraform-script-task"
  summary           = "Task that runs a script"
  command_or_script = "Script"
  script            = <<-EOT
    #!/bin/bash
    echo "Starting..."
    date
    echo "Done"
  EOT
  agent = "my-linux-agent"
}

# Task with retry configuration
resource "stonebranch_task_unix" "with_retry" {
  name           = "terraform-retry-task"
  command        = "/opt/scripts/job.sh"
  agent          = "my-linux-agent"
  retry_maximum  = 3
  retry_interval = 300
}
```

#### Argument Reference

| Attribute | Type | Required | Description |
|-----------|------|----------|-------------|
| `name` | string | Yes | Unique name of the task |
| `summary` | string | No | Description of the task |
| `agent` | string | No* | Agent to run the task on |
| `agent_cluster` | string | No* | Agent cluster to run the task on |
| `command` | string | No | Command to execute |
| `command_or_script` | string | No | `Command` or `Script` |
| `script` | string | No | Script content (when `command_or_script = "Script"`) |
| `runtime_dir` | string | No | Working directory |
| `parameters` | string | No | Parameters to pass |
| `credentials` | string | No | Credentials to use |
| `exit_codes` | string | No | Success exit codes (e.g., `"0"` or `"0,1,2"`) |
| `exit_code_processing` | string | No | `Success Exitcode Range` or `Failure Exitcode Range` |
| `retry_maximum` | int | No | Max retry attempts |
| `retry_interval` | int | No | Seconds between retries |
| `retry_indefinitely` | bool | No | Retry forever |
| `run_as_sudo` | bool | No | Run with sudo |
| `opswise_groups` | list | No | Business service names |

*One of `agent` or `agent_cluster` is required by the API.

#### Attribute Reference

| Attribute | Description |
|-----------|-------------|
| `sys_id` | System ID assigned by StoneBranch |
| `version` | Version number for optimistic locking |

#### Import

Tasks can be imported using the task name:

```bash
terraform import stonebranch_task_unix.example "task-name"
```

### stonebranch_task_windows

Manages a StoneBranch Windows task.

#### Example Usage

```hcl
# Simple Windows task with a command
resource "stonebranch_task_windows" "hello" {
  name    = "terraform-windows-hello"
  summary = "A simple Windows task managed by Terraform"
  command = "echo Hello from Terraform!"
  agent   = "my-windows-agent"
}

# Windows task with elevated privileges
resource "stonebranch_task_windows" "admin_task" {
  name         = "terraform-admin-task"
  command      = "net user"
  agent        = "my-windows-agent"
  elevate_user = true
}
```

#### Argument Reference

| Attribute | Type | Required | Description |
|-----------|------|----------|-------------|
| `name` | string | Yes | Unique name of the task |
| `summary` | string | No | Description of the task |
| `agent` | string | No* | Agent to run the task on |
| `agent_cluster` | string | No* | Agent cluster to run the task on |
| `command` | string | No | Command to execute |
| `command_or_script` | string | No | `Command` or `Script` |
| `script` | string | No | Script resource name (when `command_or_script = "Script"`) |
| `runtime_dir` | string | No | Working directory |
| `parameters` | string | No | Parameters to pass |
| `credentials` | string | No | Credentials to use |
| `exit_codes` | string | No | Success exit codes (e.g., `"0"` or `"0,1,2"`) |
| `exit_code_processing` | string | No | `Success Exitcode Range` or `Failure Exitcode Range` |
| `retry_maximum` | int | No | Max retry attempts |
| `retry_interval` | int | No | Seconds between retries |
| `retry_indefinitely` | bool | No | Retry forever |
| `elevate_user` | bool | No | Run with administrator privileges |
| `desktop_interact` | bool | No | Allow desktop interaction |
| `create_console` | bool | No | Create console window |
| `opswise_groups` | list | No | Business service names |

*One of `agent` or `agent_cluster` is required by the API.

#### Attribute Reference

| Attribute | Description |
|-----------|-------------|
| `sys_id` | System ID assigned by StoneBranch |
| `version` | Version number for optimistic locking |

#### Import

Tasks can be imported using the task name:

```bash
terraform import stonebranch_task_windows.example "task-name"
```

### stonebranch_script

Manages a reusable script resource that can be referenced by tasks.

#### Example Usage

```hcl
resource "stonebranch_script" "backup" {
  name    = "backup-script"
  content = <<-EOT
    #!/bin/bash
    tar -czf /backup/data.tar.gz /data
  EOT
}

# Reference the script in a task
resource "stonebranch_task_unix" "backup_job" {
  name              = "backup-job"
  command_or_script = "Script"
  script            = stonebranch_script.backup.name
  agent             = "my-linux-agent"
}
```

### stonebranch_trigger_time

Manages a time-based trigger for scheduling task execution.

#### Example Usage

```hcl
resource "stonebranch_trigger_time" "daily" {
  name      = "daily-trigger"
  tasks     = [stonebranch_task_unix.my_task.name]
  time      = "08:00"
  time_zone = "America/New_York"
}
```

### stonebranch_credential

Manages authentication credentials for task execution.

#### Example Usage

```hcl
resource "stonebranch_credential" "service_account" {
  name             = "service-account-creds"
  runtime_user     = "svc_user"
  runtime_password = var.service_password
}
```

## Development

### Project Structure

```
terraform-provider-stonebranch/
├── main.go                          # Provider entry point
├── internal/
│   ├── provider/
│   │   ├── provider.go              # Provider configuration and schema
│   │   └── resources/               # Resource implementations
│   │       ├── helpers.go           # Shared helper functions
│   │       ├── task_unix.go         # Unix task resource
│   │       ├── taskwindows.go       # Windows task resource
│   │       ├── task_file_transfer.go # File transfer task
│   │       ├── script.go            # Script resource
│   │       ├── trigger_time.go      # Time trigger resource
│   │       └── credential.go        # Credential resource
│   ├── acctest/
│   │   └── acctest.go               # Acceptance test helpers
│   └── client/
│       └── client.go                # StoneBranch API HTTP client
├── examples/
│   ├── dev.tfrc                     # Development override config
│   └── provider/
│       └── main.tf                  # Example Terraform configuration
├── Makefile                         # Build automation
├── CLAUDE.md                        # AI assistant context
├── ROADMAP.md                       # Development roadmap
└── openapi.yaml                     # StoneBranch API specification
```

### Useful Commands

```bash
make build      # Build the provider binary
make test       # Run tests
make fmt        # Format Go code
make clean      # Remove built binary
```

### Running Tests

```bash
# Unit tests
make test

# Acceptance tests (requires valid credentials)
export STONEBRANCH_API_TOKEN="your-token"
export STONEBRANCH_BASE_URL="https://your-instance.stonebranch.cloud"
TF_ACC=1 go test -v ./...
```

## Troubleshooting

### "Provider not found" error

Ensure your `~/.terraformrc` or `TF_CLI_CONFIG_FILE` points to the correct binary location.

### Authentication errors

1. Verify your token is valid and not expired
2. Check that the token has appropriate permissions
3. Ensure the base URL is correct (no trailing slash)

### API errors

Enable debug logging:

```bash
export TF_LOG=DEBUG
terraform plan
```

## License

[Add license information]

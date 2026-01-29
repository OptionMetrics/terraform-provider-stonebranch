# CLAUDE.md - Project Context for AI Assistants

This file provides context for Claude and other AI assistants working on this project.

## Project Overview

This is a custom Terraform provider for **StoneBranch Universal Controller**, built using the **HashiCorp Terraform Plugin Framework** (not the older SDK v2).

### Key Details

- **Provider name**: `stonebranch`
- **Module path**: `terraform-provider-stonebranch`
- **Base API URL**: `https://optionmetricsdev.stonebranch.cloud`
- **Authentication**: Bearer token in `Authorization` header
- **API Spec**: See `openapi.yaml` (OpenAPI 3.0.1, version 7.9.1.0)

## What Has Been Implemented

### Step 1: Project Scaffold (COMPLETE)

1. **Go module setup** (`go.mod`)
   - Go 1.23
   - `terraform-plugin-framework` v1.13.0
   - `terraform-plugin-log` v0.9.0

2. **Main entry point** (`main.go`)
   - Provider server setup
   - Debug flag support
   - Registry address: `registry.terraform.io/stonebranch/stonebranch`

3. **Provider configuration** (`internal/provider/provider.go`)
   - Schema with `api_token` (required, sensitive) and `base_url` (optional)
   - Environment variable fallbacks: `STONEBRANCH_API_TOKEN`, `STONEBRANCH_BASE_URL`
   - Default base URL configured
   - Proper error diagnostics for missing token

4. **API client** (`internal/client/client.go`)
   - Generic HTTP client with Bearer token auth
   - Methods: `Get`, `Post`, `Put`, `Delete`
   - JSON request/response handling
   - Proper error types (`APIError`)
   - 30-second timeout

5. **Development tooling**
   - `Makefile` with build, test, install targets
   - `examples/dev.tfrc` for local development overrides
   - `examples/provider/main.tf` sample configuration
   - `.gitignore`

### Step 2: Task Unix Resource (COMPLETE)

Implemented `stonebranch_task_unix` resource in `internal/provider/resource_task_unix.go`:

1. **Full CRUD operations**
   - Create via `POST /resources/task` (type hardcoded to `taskUnix`)
   - Read via `GET /resources/task?taskname=X`
   - Update via `PUT /resources/task`
   - Delete via `DELETE /resources/task?taskid=X`

2. **Supported attributes**
   - Identity: `sys_id` (computed), `name` (required), `version` (computed)
   - Basic: `summary`
   - Agent: `agent`, `agent_cluster`, `agent_var`, `agent_cluster_var` (one of agent/agent_cluster required)
   - Command: `command`, `command_or_script`, `script`, `runtime_dir`, `parameters`
   - Credentials: `credentials`, `credentials_var`
   - Exit codes: `exit_codes`, `exit_code_processing`
   - Output: `output_type`, `wait_for_output`, `output_return_file`, etc. (computed, server defaults)
   - Retry: `retry_maximum`, `retry_interval`, `retry_indefinitely`, `retry_suppress_failure` (computed, server defaults)
   - Unix-specific: `run_as_sudo`
   - Business services: `opswise_groups`

3. **Import support** via task name

4. **Design decision**: Each task type is a separate resource (e.g., `stonebranch_task_unix`, `stonebranch_task_windows`) rather than a single generic resource with a type field.

## Game Plan - Next Steps

### Step 3: Add Additional Task Types

Each task type should be a separate resource:
- `stonebranch_task_windows` - Windows tasks
- `stonebranch_task_sql` - SQL/Database tasks
- `stonebranch_task_workflow` - Workflow tasks
- `stonebranch_task_email` - Email tasks
- etc.

### Step 4: Add Other Resources

- Triggers/schedules
- Credentials
- Business services
- Agent clusters

### Step 5: Data Sources

Implement read-only data sources for:
- Looking up existing tasks by name
- Listing agents/agent clusters
- Querying task instances

### Step 6: Testing & Documentation

- Unit tests for client
- Acceptance tests for resources
- Generated documentation

## API Patterns (from openapi.yaml)

### Common Patterns

1. **Resource identification**: Most resources can be identified by either `{resource}id` (sysId) or `{resource}name`

2. **CRUD on single endpoint**:
   - `GET /resources/{type}?{type}name=X` - Read
   - `POST /resources/{type}` - Create (body contains resource)
   - `PUT /resources/{type}` - Update (body contains resource with sysId)
   - `DELETE /resources/{type}?{type}id=X` - Delete

3. **List endpoints**:
   - `/resources/{type}/list` - List all
   - `/resources/{type}/listadv` - Advanced search with filters

4. **Response format**:
   - Success: Often returns text message like "Successfully created..."
   - GET returns JSON object with resource data
   - Errors return text/plain with error message

### Task Type Hierarchy

```
TaskWsData (base)
  └── TaskAgentWsData
        └── TaskDistributedAgentWsData
              ├── TaskUnixWsData (type = "taskUnix")
              ├── TaskWindowsWsData (type = "taskWindows")
              └── TaskIbmiWsData (type = "taskIbmi")
```

### Key Task Fields (from TaskWsData)

- `sysId` - Internal ID (used for updates/deletes)
- `name` - Required, unique identifier
- `type` - Task type discriminator
- `summary` - Description
- `version` - Read-only, for optimistic locking
- `credentials` - Reference to credentials record
- `retryMaximum`, `retryInterval` - Retry configuration
- `variables` - Array of task variables
- `actions` - Notification actions

## Build & Test Commands

```bash
# Build
make build

# Test locally
export STONEBRANCH_API_TOKEN="your-token"
export TF_CLI_CONFIG_FILE=./examples/dev.tfrc
terraform -chdir=examples/provider plan

# Run Go tests
make test
```

## File Locations

| Purpose | Path |
|---------|------|
| Provider config | `internal/provider/provider.go` |
| API client | `internal/client/client.go` |
| Task Unix resource | `internal/provider/resource_task_unix.go` |
| Data sources | `internal/provider/datasource_*.go` (to be created) |
| API spec | `openapi.yaml` |
| Examples | `examples/` |

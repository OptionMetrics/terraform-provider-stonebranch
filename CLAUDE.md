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
   - Go 1.24
   - `terraform-plugin-framework` v1.15.0
   - `terraform-plugin-log` v0.10.0
   - `terraform-plugin-testing` v1.14.0 (acceptance tests)
   - `godotenv` v1.5.1 (.env file loading)
   - `testify` v1.10.0 (assertions)

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

Implemented `stonebranch_task_unix` resource in `internal/provider/resources/task_unix.go`:

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

### Step 2a: Task Windows Resource (COMPLETE)

Implemented `stonebranch_task_windows` resource in `internal/provider/resources/taskwindows.go`:

**Note**: The file is named `taskwindows.go` (not `task_windows.go`) because Go interprets `_windows.go` suffix as a platform-specific build constraint that only compiles on Windows.

1. **Full CRUD operations**
   - Create via `POST /resources/task` (type hardcoded to `taskWindows`)
   - Read via `GET /resources/task?taskname=X`
   - Update via `PUT /resources/task`
   - Delete via `DELETE /resources/task?taskid=X`

2. **Supported attributes**
   - Identity: `sys_id` (computed), `name` (required), `version` (computed)
   - Basic: `summary`
   - Agent: `agent`, `agent_cluster`, `agent_var`, `agent_cluster_var`
   - Command: `command`, `command_or_script`, `script`, `runtime_dir`, `parameters`
   - Credentials: `credentials`, `credentials_var`
   - Exit codes: `exit_codes`, `exit_code_processing`
   - Output: `output_type`, `wait_for_output`, `output_return_file`, etc.
   - Retry: `retry_maximum`, `retry_interval`, `retry_indefinitely`, `retry_suppress_failure`
   - Windows-specific: `elevate_user`, `desktop_interact`, `create_console` (computed, server defaults)
   - Business services: `opswise_groups`

3. **Import support** via task name

### Step 2b: Script Resource (COMPLETE)

Implemented `stonebranch_script` resource in `internal/provider/resources/script.go`:

1. **Full CRUD operations**
   - Create via `POST /resources/script`
   - Read via `GET /resources/script?scriptname=X`
   - Update via `PUT /resources/script`
   - Delete via `DELETE /resources/script?scriptid=X`

2. **Supported attributes**
   - Identity: `sys_id` (computed), `name` (required), `version` (computed)
   - Content: `script_type`, `content` (required)
   - Optional: `description`, `resolve_variables`
   - Business services: `opswise_groups`

3. **Import support** via script name

4. **Integration with tasks**: Unix tasks can reference scripts using `command_or_script = "Script"` and `script = stonebranch_script.my_script.name`

### Step 2c: Time Trigger Resource (COMPLETE)

Implemented `stonebranch_trigger_time` resource in `internal/provider/resources/trigger_time.go`:

1. **Full CRUD operations**
   - Create via `POST /resources/trigger` (type = "triggerTime")
   - Read via `GET /resources/trigger?triggername=X`
   - Update via `PUT /resources/trigger`
   - Delete via `DELETE /resources/trigger?triggerid=X`

2. **Supported attributes**
   - Identity: `sys_id` (computed), `name` (required), `version` (computed)
   - Basic: `description`, `enabled` (computed, defaults to false)
   - Tasks: `tasks` (required, list of task names to trigger)
   - Time: `time` (required), `time_zone`, `time_style`, `time_interval`, `time_interval_units`
   - Day: `day_style`, `day_interval`, `sunday`-`saturday` flags
   - Calendar: `calendar`
   - Business services: `opswise_groups`

3. **Import support** via trigger name

4. **Note**: Triggers are created disabled by default. Use the `enabled` attribute to control this.

## Game Plan - Next Steps

### Step 3: Add Additional Task Types

Each task type should be a separate resource:
- `stonebranch_task_windows` - Windows tasks (COMPLETE)
- `stonebranch_task_file_transfer` - File transfer tasks (COMPLETE)
- `stonebranch_task_sql` - SQL/Database tasks
- `stonebranch_task_workflow` - Workflow tasks
- `stonebranch_task_email` - Email tasks
- etc.

### Step 4: Add Other Resources

- Additional triggers/schedules
- Business services
- Agent clusters

### Step 5: Data Sources

Implement read-only data sources for:
- Looking up existing tasks by name
- Listing agents/agent clusters
- Querying task instances

### Step 6: Testing & Documentation (COMPLETE)

- Unit tests for client (`internal/client/client_test.go`)
- Acceptance tests for resources (`internal/provider/resources/*_test.go`)
- Shared test helpers (`internal/acctest/acctest.go`)
- .env file support for credentials (auto-loaded in tests via godotenv)
- Generated documentation (to be done)

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

# Run unit tests (no API credentials needed)
make test

# Run only client unit tests
make testunit

# Run acceptance tests (requires API credentials)
# Tests auto-load .env file via godotenv
make testacc

# Run acceptance tests for Unix task only
make testacc-unix

# Generate test coverage report
make testcov

# Test locally with Terraform
export STONEBRANCH_API_TOKEN="your-token"  # Or use .env file
export TF_CLI_CONFIG_FILE=./examples/dev.tfrc
terraform -chdir=examples/provider plan
```

## File Locations

| Purpose | Path |
|---------|------|
| Provider config | `internal/provider/provider.go` |
| API client | `internal/client/client.go` |
| Client unit tests | `internal/client/client_test.go` |
| Resource helpers | `internal/provider/resources/helpers.go` |
| Task Unix resource | `internal/provider/resources/task_unix.go` |
| Task Unix tests | `internal/provider/resources/task_unix_test.go` |
| Task Windows resource | `internal/provider/resources/taskwindows.go` |
| Task Windows tests | `internal/provider/resources/taskwindows_test.go` |
| Task File Transfer resource | `internal/provider/resources/task_file_transfer.go` |
| Task File Transfer tests | `internal/provider/resources/task_file_transfer_test.go` |
| Script resource | `internal/provider/resources/script.go` |
| Script tests | `internal/provider/resources/script_test.go` |
| Time Trigger resource | `internal/provider/resources/trigger_time.go` |
| Time Trigger tests | `internal/provider/resources/trigger_time_test.go` |
| Credential resource | `internal/provider/resources/credential.go` |
| Credential tests | `internal/provider/resources/credential_test.go` |
| Variable resource | `internal/provider/resources/variable.go` |
| Variable tests | `internal/provider/resources/variable_test.go` |
| Test helpers | `internal/acctest/acctest.go` |
| Data sources | `internal/provider/data_sources/*.go` (to be created) |
| API spec | `openapi.yaml` |
| Examples | `examples/` |
| Environment template | `.env.example` |

## Important: Go File Naming Convention

Avoid using platform-specific suffixes in Go file names:
- `_windows.go`, `_linux.go`, `_darwin.go` - Go treats these as build constraints
- Use `taskwindows.go` instead of `task_windows.go`
- This ensures the file compiles on all platforms

## Project Structure

```
terraform-provider-stonebranch/
├── main.go                          # Provider entry point
├── internal/
│   ├── provider/
│   │   ├── provider.go              # Provider configuration
│   │   └── resources/               # Resource implementations
│   │       ├── helpers.go           # Shared helper functions
│   │       ├── task_unix.go
│   │       ├── task_unix_test.go
│   │       ├── taskwindows.go
│   │       ├── taskwindows_test.go
│   │       ├── task_file_transfer.go
│   │       ├── task_file_transfer_test.go
│   │       ├── script.go
│   │       ├── script_test.go
│   │       ├── trigger_time.go
│   │       ├── trigger_time_test.go
│   │       ├── credential.go
│   │       ├── credential_test.go
│   │       ├── variable.go
│   │       └── variable_test.go
│   ├── acctest/
│   │   └── acctest.go               # Acceptance test helpers
│   └── client/
│       ├── client.go                # API client
│       └── client_test.go           # Client unit tests
├── examples/                        # Example configurations
├── CLAUDE.md                        # AI assistant context
├── README.md                        # User documentation
├── ROADMAP.md                       # Development roadmap
└── openapi.yaml                     # StoneBranch API spec
```

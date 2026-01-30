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

### Step 2d: Variable Resource (COMPLETE)

Implemented `stonebranch_variable` resource in `internal/provider/resources/variable.go`:

1. **Full CRUD operations**
   - Create via `POST /resources/variable`
   - Read via `GET /resources/variable?variablename=X`
   - Update via `PUT /resources/variable`
   - Delete via `DELETE /resources/variable?variableid=X`

2. **Supported attributes**
   - Identity: `sys_id` (computed), `name` (required)
   - Content: `value` (required), `description`
   - Business services: `opswise_groups`

3. **Import support** via variable name

4. **Naming rules**: Variable names must begin with a letter. Allowable characters are alphanumerics (upper or lower case), and underscore (_). White spaces and hyphens are not permitted. Do not use the prefix `ops_` (reserved for built-in variables).

### Step 2e: Database Connection Resource (COMPLETE)

Implemented `stonebranch_database_connection` resource in `internal/provider/resources/database_connection.go`:

1. **Full CRUD operations**
   - Create via `POST /resources/databaseconnection`
   - Read via `GET /resources/databaseconnection?name=X`
   - Update via `PUT /resources/databaseconnection`
   - Delete via `DELETE /resources/databaseconnection?databaseconnectionid=X`

2. **Supported attributes**
   - Identity: `sys_id` (computed), `name` (required)
   - Connection: `db_driver` (required), `db_url` (required)
   - Authentication: `credentials`, `credentials_var`
   - Optional: `description`, `max_rows`
   - Business services: `opswise_groups`

3. **Import support** via database connection name

### Step 2f: SQL Task Resource (COMPLETE)

Implemented `stonebranch_task_sql` resource in `internal/provider/resources/task_sql.go`:

1. **Full CRUD operations**
   - Create via `POST /resources/task` (type = "taskSql")
   - Read via `GET /resources/task?taskname=X`
   - Update via `PUT /resources/task`
   - Delete via `DELETE /resources/task?taskid=X`

2. **Supported attributes**
   - Identity: `sys_id` (computed), `name` (required), `version` (computed)
   - Basic: `summary`
   - Connection: `database_connection` (required) - Note: named this way because `connection` is reserved in Terraform
   - SQL: `sql_statement`, `sql_command`, `column_type`, `column_op` (computed), `column_value`
   - Output: `output_type`, `output_return_file`, etc.
   - Retry: `retry_maximum`, `retry_interval`, etc.
   - Business services: `opswise_groups`

3. **Import support** via task name

### Step 2g: Workflow Task Resource (COMPLETE)

Implemented `stonebranch_task_workflow` resource in `internal/provider/resources/task_workflow.go`:

1. **Full CRUD operations**
   - Create via `POST /resources/task` (type = "taskWorkflow")
   - Read via `GET /resources/task?taskname=X`
   - Update via `PUT /resources/task`
   - Delete via `DELETE /resources/task?taskid=X`

2. **Supported attributes**
   - Identity: `sys_id` (computed), `name` (required), `version` (computed)
   - Basic: `summary`
   - Workflow options: `calculate_critical_path`, `skipped_option`, `instance_wait`, `instance_wait_lookup`, `layout_option`
   - Retry: `retry_maximum`, `retry_interval`, `retry_suppress_failure`
   - Business services: `opswise_groups`

3. **Import support** via task name

### Step 2h: Workflow Vertex Resource (COMPLETE)

Implemented `stonebranch_workflow_vertex` resource in `internal/provider/resources/workflow_vertex.go`:

1. **Full CRUD operations**
   - Create via `POST /resources/workflow/vertices?workflowname=X`
   - Read via `GET /resources/workflow/vertices?workflowname=X&vertexid=Y`
   - Update via `PUT /resources/workflow/vertices?workflowname=X`
   - Delete via `DELETE /resources/workflow/vertices?workflowname=X&vertexid=Y`

2. **Supported attributes**
   - Identity: `workflow_name` (required), `task_name` (required), `vertex_id` (computed)
   - Optional: `alias` (for multiple instances of same task), `vertex_x`, `vertex_y` (diagram position)

3. **Usage**: Add existing tasks to a workflow. Reference tasks by name and get back a vertex_id for creating edges.

### Step 2i: Workflow Edge Resource (COMPLETE)

Implemented `stonebranch_workflow_edge` resource in `internal/provider/resources/workflow_edge.go`:

1. **Full CRUD operations**
   - Create via `POST /resources/workflow/edges?workflowname=X`
   - Read via `GET /resources/workflow/edges?workflowname=X` (finds matching source/target)
   - Update via `PUT /resources/workflow/edges?workflowname=X`
   - Delete via `DELETE /resources/workflow/edges?workflowname=X&sourceid=Y&targetid=Z`

2. **Supported attributes**
   - Identity: `workflow_name` (required), `source_id` (required vertex_id), `target_id` (required vertex_id)
   - Optional: `straight_edge` (diagram display)

3. **Usage**: Create dependencies between tasks in a workflow. The source task must complete before the target task runs.

### Step 2j: Email Connection Resource (COMPLETE)

Implemented `stonebranch_email_connection` resource in `internal/provider/resources/email_connection.go`:

1. **Full CRUD operations**
   - Create via `POST /resources/emailconnection`
   - Read via `GET /resources/emailconnection?connectionname=X`
   - Update via `PUT /resources/emailconnection`
   - Delete via `DELETE /resources/emailconnection?connectionid=X`

2. **Supported attributes**
   - Identity: `sys_id` (computed), `name` (required), `version` (computed)
   - SMTP: `smtp` (required), `smtp_port`, `smtp_ssl`, `smtp_starttls`
   - Sender: `email_address`
   - Authentication: `authentication`, `authentication_type`, `default_user`, `default_password` (sensitive), `oauth_client`
   - IMAP (for reading): `imap`, `imap_port`, `imap_ssl`, `imap_starttls`, `trash_folder`
   - Optional: `description`
   - Business services: `opswise_groups`

3. **Import support** via email connection name

### Step 2k: Email Task Resource (COMPLETE)

Implemented `stonebranch_task_email` resource in `internal/provider/resources/task_email.go`:

1. **Full CRUD operations**
   - Create via `POST /resources/task` (type = "taskEmail")
   - Read via `GET /resources/task?taskname=X`
   - Update via `PUT /resources/task`
   - Delete via `DELETE /resources/task?taskid=X`

2. **Supported attributes**
   - Identity: `sys_id` (computed), `name` (required), `version` (computed)
   - Basic: `summary`
   - Connection: `email_connection`, `email_connection_var`
   - Template: `template`, `template_var`
   - Recipients: `to_recipients`, `cc_recipients`, `bcc_recipients`, `reply_to`
   - Content: `subject`, `body`
   - Attachments: `attach_local_file`, `local_attachments_path`, `local_attachment`
   - Report: `report_var`, `list_report_format`
   - Exit codes: `exit_codes`
   - Retry: `retry_maximum`, `retry_interval`, `retry_indefinitely`, `retry_suppress_failure`
   - Business services: `opswise_groups`

3. **Import support** via task name

### Step 2l: Cron Trigger Resource (COMPLETE)

Implemented `stonebranch_trigger_cron` resource in `internal/provider/resources/trigger_cron.go`:

1. **Full CRUD operations**
   - Create via `POST /resources/trigger` (type = "triggerCron")
   - Read via `GET /resources/trigger?triggername=X`
   - Update via `PUT /resources/trigger`
   - Delete via `DELETE /resources/trigger?triggerid=X`

2. **Supported attributes**
   - Identity: `sys_id` (computed), `name` (required), `version` (computed)
   - Basic: `description`, `enabled`
   - Tasks: `tasks` (required, list of task names)
   - Cron fields: `minutes`, `hours`, `day_of_month`, `month`, `day_of_week` (all required)
   - Day logic: `day_logic` (And/Or for combining day_of_month and day_of_week)
   - Scheduling: `time_zone`, `calendar`
   - Business services: `opswise_groups`

3. **Import support** via trigger name

4. **Note**: Triggers are created disabled by default. Use the `enabled` attribute to control this.

### Step 2m: Business Service Resource (COMPLETE)

Implemented `stonebranch_business_service` resource in `internal/provider/resources/business_service.go`:

1. **Full CRUD operations**
   - Create via `POST /resources/businessservice`
   - Read via `GET /resources/businessservice?busservicename=X`
   - Update via `PUT /resources/businessservice`
   - Delete via `DELETE /resources/businessservice?busserviceid=X`

2. **Supported attributes**
   - Identity: `sys_id` (computed), `name` (required), `version` (computed)
   - Content: `description`

3. **Import support** via business service name

4. **Usage**: Business services are used to group and organize resources. Other resources reference business services through the `opswise_groups` attribute.

## Game Plan - Next Steps

### Step 3: Add Additional Task Types

Each task type should be a separate resource:
- `stonebranch_task_windows` - Windows tasks (COMPLETE)
- `stonebranch_task_file_transfer` - File transfer tasks (COMPLETE)
- `stonebranch_task_sql` - SQL/Database tasks (COMPLETE)
- `stonebranch_task_workflow` - Workflow tasks (COMPLETE)
- `stonebranch_task_email` - Email tasks (COMPLETE)
- etc.

### Step 4: Add Other Resources

- `stonebranch_variable` - Global variables (COMPLETE)
- `stonebranch_database_connection` - Database connections (COMPLETE)
- `stonebranch_email_connection` - Email connections (COMPLETE)
- `stonebranch_workflow_vertex` - Tasks within workflows (COMPLETE)
- `stonebranch_workflow_edge` - Dependencies between workflow tasks (COMPLETE)
- `stonebranch_business_service` - Business services (COMPLETE)
- Additional triggers/schedules
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
| Cron Trigger resource | `internal/provider/resources/trigger_cron.go` |
| Cron Trigger tests | `internal/provider/resources/trigger_cron_test.go` |
| Credential resource | `internal/provider/resources/credential.go` |
| Credential tests | `internal/provider/resources/credential_test.go` |
| Variable resource | `internal/provider/resources/variable.go` |
| Variable tests | `internal/provider/resources/variable_test.go` |
| Database Connection resource | `internal/provider/resources/database_connection.go` |
| Database Connection tests | `internal/provider/resources/database_connection_test.go` |
| Email Connection resource | `internal/provider/resources/email_connection.go` |
| Email Connection tests | `internal/provider/resources/email_connection_test.go` |
| Task SQL resource | `internal/provider/resources/task_sql.go` |
| Task SQL tests | `internal/provider/resources/task_sql_test.go` |
| Task Email resource | `internal/provider/resources/task_email.go` |
| Task Email tests | `internal/provider/resources/task_email_test.go` |
| Task Workflow resource | `internal/provider/resources/task_workflow.go` |
| Task Workflow tests | `internal/provider/resources/task_workflow_test.go` |
| Workflow Vertex resource | `internal/provider/resources/workflow_vertex.go` |
| Workflow Vertex tests | `internal/provider/resources/workflow_vertex_test.go` |
| Workflow Edge resource | `internal/provider/resources/workflow_edge.go` |
| Workflow Edge tests | `internal/provider/resources/workflow_edge_test.go` |
| Business Service resource | `internal/provider/resources/business_service.go` |
| Business Service tests | `internal/provider/resources/business_service_test.go` |
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
│   │       ├── trigger_cron.go
│   │       ├── trigger_cron_test.go
│   │       ├── credential.go
│   │       ├── credential_test.go
│   │       ├── variable.go
│   │       ├── variable_test.go
│   │       ├── database_connection.go
│   │       ├── database_connection_test.go
│   │       ├── email_connection.go
│   │       ├── email_connection_test.go
│   │       ├── task_sql.go
│   │       ├── task_sql_test.go
│   │       ├── task_email.go
│   │       ├── task_email_test.go
│   │       ├── task_workflow.go
│   │       ├── task_workflow_test.go
│   │       ├── workflow_vertex.go
│   │       ├── workflow_vertex_test.go
│   │       ├── workflow_edge.go
│   │       ├── workflow_edge_test.go
│   │       ├── business_service.go
│   │       └── business_service_test.go
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

# Plan: sb2tf CLI Utility

## Overview

Add a command-line utility `sb2tf` to this Terraform provider project that reads existing resources from the StoneBranch Universal Controller API and generates Terraform configuration files (.tf).

## Use Cases

1. **Bootstrap a new Terraform project** - Export existing resources to start managing them with Terraform
2. **Test reproducibility** - Verify that Terraform configs can recreate existing infrastructure
3. **Migration** - Help migrate manually-created resources to Infrastructure as Code

## Design Decisions (Based on User Input)

| Decision | Choice |
|----------|--------|
| Import blocks | No - just resource definitions |
| File organization | Support both via `--format` flag (single/grouped) |
| Dependencies | Yes - automatically export dependent resources |
| CLI framework | Cobra |

---

## Project Structure

```
cmd/sb2tf/
├── main.go                          # CLI entry point, Cobra setup
├── cli/
│   ├── root.go                      # Root command, global flags (--token, --url, --output)
│   ├── list.go                      # List available resources by type
│   └── export.go                    # Export resources to Terraform HCL
└── generator/
    ├── generator.go                 # Core HCL generation logic
    ├── templates.go                 # HCL templates per resource type
    ├── resources.go                 # Resource type registry (API endpoint, TF resource name, fields)
    └── dependencies.go              # Dependency resolution logic
```

### Reused Code

- `internal/client/client.go` - API client (authentication, HTTP methods) - import directly

---

## CLI Interface

```bash
# Authentication (same env vars as provider)
export STONEBRANCH_API_TOKEN="your-token"
export STONEBRANCH_BASE_URL="https://your-instance.stonebranch.cloud"

# List resources
sb2tf list                                    # Show all resource types
sb2tf list tasks                              # List all tasks (shows name, type)
sb2tf list tasks --filter "prod-*"            # Filter by name pattern
sb2tf list triggers                           # List all triggers
sb2tf list variables                          # List all variables

# Export single resource
sb2tf export task_unix my_task                # Export one Unix task
sb2tf export task_workflow my_workflow        # Export workflow + dependencies

# Export multiple resources
sb2tf export tasks --filter "prod-*"          # Export tasks matching pattern
sb2tf export triggers --all                   # Export all triggers
sb2tf export --all                            # Export everything

# Output options
sb2tf export tasks --all --output ./terraform/           # Write to directory
sb2tf export tasks --all --format single                 # One file per resource (default)
sb2tf export tasks --all --format grouped                # Group by type

# Other flags
sb2tf export task_workflow my_wf --no-deps               # Skip dependency resolution
sb2tf export --dry-run                                   # Show what would be exported
```

---

## Supported Resource Types

| Category | CLI Name | TF Resource | API Endpoint | List Endpoint |
|----------|----------|-------------|--------------|---------------|
| **Tasks** | | | | |
| | task_unix | stonebranch_task_unix | /resources/task | /resources/task/listadv?type=taskUnix |
| | task_windows | stonebranch_task_windows | /resources/task | /resources/task/listadv?type=taskWindows |
| | task_sql | stonebranch_task_sql | /resources/task | /resources/task/listadv?type=taskSql |
| | task_email | stonebranch_task_email | /resources/task | /resources/task/listadv?type=taskEmail |
| | task_workflow | stonebranch_task_workflow | /resources/task | /resources/task/listadv?type=taskWorkflow |
| | task_file_monitor | stonebranch_task_file_monitor | /resources/task | /resources/task/listadv?type=taskFileMonitor |
| | task_file_transfer | stonebranch_task_file_transfer | /resources/task | /resources/task/listadv?type=taskFileTransfer |
| | task_timer | stonebranch_task_timer | /resources/task | /resources/task/listadv?type=taskTimer |
| | task_monitor | stonebranch_task_monitor | /resources/task | /resources/task/listadv?type=taskMonitor |
| | task_stored_procedure | stonebranch_task_stored_procedure | /resources/task | /resources/task/listadv?type=taskStoredProc |
| | task_web_service | stonebranch_task_web_service | /resources/task | /resources/task/listadv?type=taskWebService |
| | task_universal_aws_s3 | stonebranch_task_universal_aws_s3 | /resources/task | /resources/task/listadv?type=taskUniversal |
| **Triggers** | | | | |
| | trigger_time | stonebranch_trigger_time | /resources/trigger | /resources/trigger/listadv?type=triggerTime |
| | trigger_cron | stonebranch_trigger_cron | /resources/trigger | /resources/trigger/listadv?type=triggerCron |
| | trigger_file_monitor | stonebranch_trigger_file_monitor | /resources/trigger | /resources/trigger/listadv?type=triggerFm |
| | trigger_task_monitor | stonebranch_trigger_task_monitor | /resources/trigger | /resources/trigger/listadv?type=triggerTm |
| **Connections** | | | | |
| | database_connection | stonebranch_database_connection | /resources/databaseconnection | /resources/databaseconnection/list |
| | email_connection | stonebranch_email_connection | /resources/emailconnection | /resources/emailconnection/list |
| **Other** | | | | |
| | script | stonebranch_script | /resources/script | /resources/script/list |
| | variable | stonebranch_variable | /resources/variable | /resources/variable/list |
| | credential | stonebranch_credential | /resources/credential | /resources/credential/list |
| | business_service | stonebranch_business_service | /resources/businessservice | /resources/businessservice/list |
| | agent_cluster | stonebranch_agent_cluster | /resources/agentcluster | /resources/agentcluster/list |
| | calendar | stonebranch_calendar | /resources/calendar | /resources/calendar/list |
| **Workflow** | | | | |
| | workflow_vertex | stonebranch_workflow_vertex | /resources/workflow/vertices | (per-workflow) |
| | workflow_edge | stonebranch_workflow_edge | /resources/workflow/edges | (per-workflow) |

---

## Dependency Resolution

When exporting a workflow with `--deps` (default), sb2tf will:

1. **Export the workflow task itself**
2. **Fetch workflow vertices** via `GET /resources/workflow/vertices?workflowname=X`
3. **For each vertex**: Export the referenced task (and its dependencies)
4. **Fetch workflow edges** via `GET /resources/workflow/edges?workflowname=X`
5. **Export edges** with correct vertex references

### Task Dependencies

When exporting any task, also check and export:
- `credentials` → stonebranch_credential
- `script` → stonebranch_script (if command_or_script = "Script")
- `database_connection` → stonebranch_database_connection (for SQL/stored proc tasks)
- `email_connection` → stonebranch_email_connection (for email tasks)
- `agent_cluster` → stonebranch_agent_cluster

### Trigger Dependencies

- `tasks` → List of task names to export
- `calendar` → stonebranch_calendar
- `task_monitor` → For file_monitor and task_monitor triggers

---

## HCL Generation

### Template Approach

Use Go `text/template` to generate HCL. Example for a variable:

```go
const variableTemplate = `resource "stonebranch_variable" "{{.ResourceName}}" {
  name        = "{{.Name}}"
{{- if .Value}}
  value       = "{{.Value}}"
{{- end}}
{{- if .Description}}
  description = "{{.Description}}"
{{- end}}
{{- if .OpswiseGroups}}
  opswise_groups = [{{range $i, $v := .OpswiseGroups}}{{if $i}}, {{end}}"{{$v}}"{{end}}]
{{- end}}
}
`
```

### Resource Name Sanitization

Convert StoneBranch names to valid Terraform identifiers:
- `My Task Name` → `my_task_name`
- `prod-task-01` → `prod_task_01`
- `123task` → `task_123` (prepend if starts with digit)

### Example Output

```hcl
# Generated by sb2tf from StoneBranch Universal Controller
# Source: my_unix_task

resource "stonebranch_task_unix" "my_unix_task" {
  name              = "my_unix_task"
  summary           = "Runs the daily batch job"
  agent             = "linux-prod-01"
  command           = "/opt/scripts/daily.sh"
  credentials       = "service_account"
  exit_codes        = "0"

  opswise_groups    = ["production", "batch-jobs"]
}

resource "stonebranch_credential" "service_account" {
  name        = "service_account"
  runtime_user = "svc_batch"
  # Note: password cannot be exported for security
}
```

---

## Implementation Steps

### Phase 1: CLI Skeleton
1. Create `cmd/sb2tf/main.go` - Cobra root command
2. Implement `cli/root.go` with global flags:
   - `--token` (env: STONEBRANCH_API_TOKEN)
   - `--url` (env: STONEBRANCH_BASE_URL)
   - `--output` (default: stdout)
   - `--format` (single|grouped)
3. Implement `cli/list.go` - list command

**Files to create:**
- `cmd/sb2tf/main.go`
- `cmd/sb2tf/cli/root.go`
- `cmd/sb2tf/cli/list.go`

### Phase 2: Generator Framework
1. Create `generator/resources.go` - resource type registry
2. Create `generator/generator.go` - main generation logic
3. Create `generator/templates.go` - HCL templates

**Files to create:**
- `cmd/sb2tf/generator/resources.go`
- `cmd/sb2tf/generator/generator.go`
- `cmd/sb2tf/generator/templates.go`

### Phase 3: Export Command
1. Implement `cli/export.go` - export command
2. Wire up generator to export command

**Files to create/modify:**
- `cmd/sb2tf/cli/export.go`

### Phase 4: Simple Resources
Implement templates and API models for:
1. variable
2. script
3. credential
4. business_service

### Phase 5: Task Resources
Implement templates for all task types:
1. task_unix
2. task_windows
3. task_sql
4. task_email
5. task_workflow
6. (remaining task types)

### Phase 6: Triggers and Connections
1. trigger_time, trigger_cron, trigger_file_monitor, trigger_task_monitor
2. database_connection, email_connection
3. agent_cluster, calendar

### Phase 7: Workflow Support
1. Implement workflow dependency resolution
2. Generate workflow_vertex and workflow_edge resources
3. Handle cross-references between resources

### Phase 8: Polish
1. Add `--dry-run` flag
2. Add `--filter` wildcards
3. Update Makefile with build targets
4. Add README section for sb2tf

---

## Verification

After implementation, verify with:

```bash
# Build the utility
make build-sb2tf

# List resources from a test instance
./bin/sb2tf list tasks

# Export a simple resource
./bin/sb2tf export variable my_var > test.tf
cat test.tf

# Export a workflow with dependencies
./bin/sb2tf export task_workflow my_workflow --output ./exported/

# Verify generated Terraform is valid
cd ./exported && terraform init && terraform validate
```

---

## Files to Modify

| File | Change |
|------|--------|
| `Makefile` | Add `build-sb2tf` and `install-sb2tf` targets |
| `go.mod` | Add Cobra dependency |
| `CLAUDE.md` | Document sb2tf utility |
| `.goreleaser.yaml` | Add sb2tf build configuration |

---

## GoReleaser Configuration

The project uses GoReleaser to build and publish releases to GitHub. To include `sb2tf` in releases, add a second build entry to `.goreleaser.yaml`:

```yaml
# GoReleaser configuration for terraform-provider-stonebranch
# Documentation: https://goreleaser.com

version: 2

before:
  hooks:
    - go mod tidy

builds:
  # Existing Terraform provider build
  - id: terraform-provider-stonebranch
    binary: terraform-provider-stonebranch_v{{ .Version }}
    env:
      - CGO_ENABLED=0
    goos:
      - darwin
      - linux
      - windows
    goarch:
      - amd64
      - arm64
    ignore:
      - goos: windows
        goarch: arm64
    ldflags:
      - -s -w
      - -X main.version={{ .Version }}

  # NEW: sb2tf CLI utility build
  - id: sb2tf
    main: ./cmd/sb2tf
    binary: sb2tf
    env:
      - CGO_ENABLED=0
    goos:
      - darwin
      - linux
      - windows
    goarch:
      - amd64
      - arm64
    ignore:
      - goos: windows
        goarch: arm64
    ldflags:
      - -s -w
      - -X main.version={{ .Version }}

archives:
  # Provider archive (unchanged)
  - id: provider-zip
    builds:
      - terraform-provider-stonebranch
    formats:
      - zip
    name_template: "terraform-provider-stonebranch_{{ .Version }}_{{ .Os }}_{{ .Arch }}"

  # NEW: sb2tf archive (separate downloads)
  - id: sb2tf-zip
    builds:
      - sb2tf
    formats:
      - zip
    name_template: "sb2tf_{{ .Version }}_{{ .Os }}_{{ .Arch }}"

checksum:
  name_template: "checksums_{{ .Version }}_SHA256SUMS"
  algorithm: sha256

snapshot:
  version_template: "{{ incpatch .Version }}-dev"

changelog:
  sort: asc
  filters:
    exclude:
      - "^docs:"
      - "^test:"
      - "^chore:"
      - Merge pull request
      - Merge branch

# Publish to GitHub Releases
release:
  github:
    owner: OptionMetrics
    name: terraform-provider-stonebranch
  draft: false
  prerelease: auto
  name_template: "v{{ .Version }}"
```

### Key Changes:

1. **Second build entry** (`id: sb2tf`) - builds from `./cmd/sb2tf` with binary name `sb2tf`
2. **Separate archive** (`id: sb2tf-zip`) - creates `sb2tf_v1.0.0_darwin_amd64.zip` etc.
3. **Updated checksum** - single checksum file covers both binaries
4. **Same platforms** - builds for darwin/linux/windows on amd64/arm64

### Release Assets

After running `goreleaser release`, GitHub releases will include:

```
terraform-provider-stonebranch_1.0.0_darwin_amd64.zip
terraform-provider-stonebranch_1.0.0_darwin_arm64.zip
terraform-provider-stonebranch_1.0.0_linux_amd64.zip
terraform-provider-stonebranch_1.0.0_linux_arm64.zip
terraform-provider-stonebranch_1.0.0_windows_amd64.zip
sb2tf_1.0.0_darwin_amd64.zip
sb2tf_1.0.0_darwin_arm64.zip
sb2tf_1.0.0_linux_amd64.zip
sb2tf_1.0.0_linux_arm64.zip
sb2tf_1.0.0_windows_amd64.zip
checksums_1.0.0_SHA256SUMS
```

### User Installation

Users can download the sb2tf binary directly from GitHub releases:

```bash
# macOS (Apple Silicon)
curl -LO https://github.com/OptionMetrics/terraform-provider-stonebranch/releases/latest/download/sb2tf_1.0.0_darwin_arm64.zip
unzip sb2tf_1.0.0_darwin_arm64.zip
chmod +x sb2tf
sudo mv sb2tf /usr/local/bin/

# Linux (x86_64)
curl -LO https://github.com/OptionMetrics/terraform-provider-stonebranch/releases/latest/download/sb2tf_1.0.0_linux_amd64.zip
unzip sb2tf_1.0.0_linux_amd64.zip
chmod +x sb2tf
sudo mv sb2tf /usr/local/bin/
```

---

## Estimated Complexity

- **Phase 1-3 (Core)**: ~400-500 lines of Go
- **Phase 4-6 (Resources)**: ~100-150 lines per resource type (templates + models)
- **Phase 7 (Workflows)**: ~200-300 lines for dependency resolution
- **Phase 8 (Polish)**: ~100 lines

Total: ~2000-3000 lines of new Go code

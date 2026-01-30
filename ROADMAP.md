# Terraform Provider Stonebranch - Project Roadmap

This document outlines the development roadmap for the Stonebranch Universal Controller Terraform Provider. The provider enables infrastructure-as-code management of Stonebranch resources.

## Current Status

**Provider Version:** Pre-release (0.x)

### Implemented Resources

| Resource | Status | Tests | Docs |
|----------|--------|-------|------|
| `stonebranch_task_unix` | ✅ Complete | ✅ | ❌ |
| `stonebranch_script` | ✅ Complete | ✅ | ❌ |
| `stonebranch_trigger_time` | ✅ Complete | ✅ | ❌ |
| `stonebranch_task_file_transfer` | ✅ Complete | ✅ | ❌ |
| `stonebranch_credential` | ✅ Complete | ✅ | ❌ |

### Implemented Data Sources

None yet.

---

## Phase 1: Core Task Types (High Priority)

These are the most commonly used task types that form the foundation of most automation workflows.

### Task Resources

| Resource | Description | Priority | Status |
|----------|-------------|----------|--------|
| `stonebranch_task_unix` | Unix/Linux command execution | P0 | ✅ Complete |
| `stonebranch_task_windows` | Windows command execution | P0 | 🔲 Not Started |
| `stonebranch_task_sql` | SQL query execution | P0 | 🔲 Not Started |
| `stonebranch_task_stored_procedure` | Stored procedure execution | P1 | 🔲 Not Started |
| `stonebranch_task_file_transfer` | File transfer (FTP/SFTP) | P0 | ✅ Complete |
| `stonebranch_task_email` | Email notifications | P1 | 🔲 Not Started |
| `stonebranch_task_web_service` | REST/SOAP web service calls | P1 | 🔲 Not Started |

### Supporting Resources

| Resource | Description | Priority | Status |
|----------|-------------|----------|--------|
| `stonebranch_script` | Reusable scripts | P0 | ✅ Complete |
| `stonebranch_credential` | Authentication credentials | P0 | ✅ Complete |
| `stonebranch_variable` | Global/scoped variables | P0 | 🔲 Not Started |
| `stonebranch_database_connection` | Database connection definitions | P0 | 🔲 Not Started |
| `stonebranch_email_connection` | Email server connections | P1 | 🔲 Not Started |

### Deliverables
- [ ] `stonebranch_task_windows` resource with full CRUD
- [ ] `stonebranch_task_sql` resource with full CRUD
- [ ] `stonebranch_task_stored_procedure` resource with full CRUD
- [ ] `stonebranch_task_email` resource with full CRUD
- [ ] `stonebranch_task_web_service` resource with full CRUD
- [ ] `stonebranch_variable` resource with full CRUD
- [ ] `stonebranch_database_connection` resource with full CRUD
- [ ] `stonebranch_email_connection` resource with full CRUD
- [ ] Acceptance tests for all resources
- [ ] Example configurations for each resource

---

## Phase 2: Workflow & Orchestration

Resources for building complex workflows and orchestration patterns.

### Workflow Resources

| Resource | Description | Priority | Status |
|----------|-------------|----------|--------|
| `stonebranch_task_workflow` | Workflow/DAG definitions | P0 | 🔲 Not Started |
| `stonebranch_workflow_vertex` | Workflow task nodes | P1 | 🔲 Not Started |
| `stonebranch_workflow_edge` | Workflow task connections | P1 | 🔲 Not Started |

### Trigger Resources

| Resource | Description | Priority | Status |
|----------|-------------|----------|--------|
| `stonebranch_trigger_time` | Time-based scheduling | P0 | ✅ Complete |
| `stonebranch_trigger_cron` | CRON expression triggers | P0 | 🔲 Not Started |
| `stonebranch_trigger_file_monitor` | File arrival triggers | P0 | 🔲 Not Started |
| `stonebranch_trigger_task_monitor` | Task completion triggers | P1 | 🔲 Not Started |
| `stonebranch_trigger_manual` | Manual/on-demand triggers | P1 | 🔲 Not Started |
| `stonebranch_trigger_temporary` | One-time triggers | P2 | 🔲 Not Started |
| `stonebranch_trigger_composite` | Compound triggers | P2 | 🔲 Not Started |
| `stonebranch_trigger_universal` | Universal event triggers | P2 | 🔲 Not Started |

### Approval & Control

| Resource | Description | Priority | Status |
|----------|-------------|----------|--------|
| `stonebranch_task_manual` | Manual intervention tasks | P1 | 🔲 Not Started |
| `stonebranch_task_approval` | Approval workflow tasks | P2 | 🔲 Not Started |
| `stonebranch_virtual_resource` | Concurrency control | P1 | 🔲 Not Started |

### Deliverables
- [ ] `stonebranch_task_workflow` resource with vertex/edge management
- [ ] All trigger type resources
- [ ] `stonebranch_virtual_resource` for concurrency control
- [ ] Workflow composition examples
- [ ] Acceptance tests for all resources

---

## Phase 3: Calendar & Scheduling

Resources for advanced scheduling with business calendars.

### Calendar Resources

| Resource | Description | Priority | Status |
|----------|-------------|----------|--------|
| `stonebranch_calendar` | Business calendars | P0 | 🔲 Not Started |
| `stonebranch_custom_day` | Holiday/special day definitions | P1 | 🔲 Not Started |

### Deliverables
- [ ] `stonebranch_calendar` resource with full CRUD
- [ ] `stonebranch_custom_day` resource with full CRUD
- [ ] Calendar integration with triggers
- [ ] Acceptance tests for all resources

---

## Phase 4: Infrastructure & Agents

Resources for managing execution infrastructure.

### Agent Resources

| Resource | Description | Priority | Status |
|----------|-------------|----------|--------|
| `stonebranch_agent_cluster` | Agent cluster definitions | P0 | 🔲 Not Started |
| `stonebranch_agent_cluster_membership` | Cluster member management | P1 | 🔲 Not Started |

### Data Sources (Read-Only)

| Data Source | Description | Priority | Status |
|-------------|-------------|----------|--------|
| `stonebranch_agent` | Look up agent details | P0 | 🔲 Not Started |
| `stonebranch_agents` | List/filter agents | P1 | 🔲 Not Started |
| `stonebranch_agent_cluster` | Look up cluster details | P1 | 🔲 Not Started |

### Deliverables
- [ ] `stonebranch_agent_cluster` resource with full CRUD
- [ ] Agent data sources for lookups
- [ ] Acceptance tests for all resources

---

## Phase 5: Organizational Resources

Resources for organizing and categorizing other resources.

### Business Service Resources

| Resource | Description | Priority | Status |
|----------|-------------|----------|--------|
| `stonebranch_business_service` | Business service groups | P0 | 🔲 Not Started |

### Email Templates

| Resource | Description | Priority | Status |
|----------|-------------|----------|--------|
| `stonebranch_email_template` | Reusable email templates | P1 | 🔲 Not Started |

### Deliverables
- [ ] `stonebranch_business_service` resource with full CRUD
- [ ] `stonebranch_email_template` resource with full CRUD
- [ ] Acceptance tests for all resources

---

## Phase 6: Enterprise Integrations

Resources for enterprise system integrations.

### SAP Integration

| Resource | Description | Priority | Status |
|----------|-------------|----------|--------|
| `stonebranch_sap_connection` | SAP system connections | P2 | 🔲 Not Started |
| `stonebranch_task_sap` | SAP job execution | P2 | 🔲 Not Started |

### PeopleSoft Integration

| Resource | Description | Priority | Status |
|----------|-------------|----------|--------|
| `stonebranch_peoplesoft_connection` | PeopleSoft connections | P2 | 🔲 Not Started |
| `stonebranch_task_peoplesoft` | PeopleSoft job execution | P2 | 🔲 Not Started |

### Mainframe Integration

| Resource | Description | Priority | Status |
|----------|-------------|----------|--------|
| `stonebranch_task_zos` | z/OS job execution | P2 | 🔲 Not Started |
| `stonebranch_trigger_zos` | z/OS event triggers | P2 | 🔲 Not Started |
| `stonebranch_task_ibmi` | IBM i job execution | P2 | 🔲 Not Started |

### Deliverables
- [ ] Enterprise connection resources
- [ ] Enterprise task types
- [ ] Acceptance tests (may require mock server)

---

## Phase 7: Advanced Features

### Webhooks & Events

| Resource | Description | Priority | Status |
|----------|-------------|----------|--------|
| `stonebranch_webhook` | Webhook definitions | P1 | 🔲 Not Started |
| `stonebranch_universal_event_template` | Event templates | P2 | 🔲 Not Started |

### Universal Templates

| Resource | Description | Priority | Status |
|----------|-------------|----------|--------|
| `stonebranch_universal_template` | Custom task templates | P2 | 🔲 Not Started |
| `stonebranch_task_universal` | Universal template tasks | P2 | 🔲 Not Started |

### Monitoring

| Resource | Description | Priority | Status |
|----------|-------------|----------|--------|
| `stonebranch_task_monitor` | File/resource monitoring | P1 | 🔲 Not Started |
| `stonebranch_trigger_application_monitor` | Application triggers | P2 | 🔲 Not Started |

### Other Task Types

| Resource | Description | Priority | Status |
|----------|-------------|----------|--------|
| `stonebranch_task_sleep` | Delay/wait tasks | P2 | 🔲 Not Started |
| `stonebranch_task_recurring` | Recurring tasks | P2 | 🔲 Not Started |
| `stonebranch_task_ucmd` | Universal command tasks | P2 | 🔲 Not Started |
| `stonebranch_task_application_control` | Application control | P3 | 🔲 Not Started |
| `stonebranch_task_critical_endpoint` | Critical path tasks | P3 | 🔲 Not Started |

### Deliverables
- [ ] Webhook and event resources
- [ ] Universal template support
- [ ] Monitoring task types
- [ ] Remaining task types

---

## Phase 8: Administration & Security

### User Management

| Resource | Description | Priority | Status |
|----------|-------------|----------|--------|
| `stonebranch_user` | User accounts | P2 | 🔲 Not Started |
| `stonebranch_user_group` | User groups | P2 | 🔲 Not Started |

### LDAP & OAuth

| Resource | Description | Priority | Status |
|----------|-------------|----------|--------|
| `stonebranch_ldap` | LDAP configuration | P3 | 🔲 Not Started |
| `stonebranch_oauth_client` | OAuth client config | P2 | 🔲 Not Started |

### SNMP

| Resource | Description | Priority | Status |
|----------|-------------|----------|--------|
| `stonebranch_snmp_manager` | SNMP manager config | P3 | 🔲 Not Started |

### Deliverables
- [ ] User/group management resources
- [ ] Authentication configuration resources
- [ ] Acceptance tests (may require admin permissions)

---

## Phase 9: Data Sources

Read-only data sources for querying existing resources.

| Data Source | Description | Priority | Status |
|-------------|-------------|----------|--------|
| `stonebranch_task` | Look up any task by name | P0 | 🔲 Not Started |
| `stonebranch_tasks` | List/filter tasks | P1 | 🔲 Not Started |
| `stonebranch_trigger` | Look up trigger by name | P1 | 🔲 Not Started |
| `stonebranch_triggers` | List/filter triggers | P2 | 🔲 Not Started |
| `stonebranch_credential` | Look up credential | P1 | 🔲 Not Started |
| `stonebranch_script` | Look up script | P1 | 🔲 Not Started |
| `stonebranch_calendar` | Look up calendar | P1 | 🔲 Not Started |
| `stonebranch_business_service` | Look up business service | P2 | 🔲 Not Started |
| `stonebranch_variable` | Look up variable | P1 | 🔲 Not Started |
| `stonebranch_database_connection` | Look up DB connection | P2 | 🔲 Not Started |
| `stonebranch_task_instance` | Query task instances | P2 | 🔲 Not Started |

### Deliverables
- [ ] Single-resource data sources (lookup by name)
- [ ] List data sources with filtering
- [ ] Acceptance tests for all data sources

---

## Phase 10: Operations & Deployment

Resources and features for CI/CD and operations.

### Bundle & Promotion

| Resource | Description | Priority | Status |
|----------|-------------|----------|--------|
| `stonebranch_bundle` | Resource bundles | P2 | 🔲 Not Started |
| `stonebranch_promotion_target` | Promotion targets | P2 | 🔲 Not Started |

### Simulation

| Resource | Description | Priority | Status |
|----------|-------------|----------|--------|
| `stonebranch_simulation` | Task simulation | P3 | 🔲 Not Started |

### Data Sources

| Data Source | Description | Priority | Status |
|-------------|-------------|----------|--------|
| `stonebranch_system_status` | System health status | P2 | 🔲 Not Started |
| `stonebranch_cluster_nodes` | Cluster node info | P3 | 🔲 Not Started |
| `stonebranch_metrics` | System metrics | P3 | 🔲 Not Started |

### Deliverables
- [ ] Bundle management resources
- [ ] Promotion workflow resources
- [ ] Operational data sources

---

## Testing Strategy

### Unit Tests

Each resource/data source should have:
- [ ] Schema validation tests
- [ ] Model conversion tests (API ↔ Terraform)
- [ ] Nil/empty value handling tests

Location: `internal/provider/*_test.go` (non-acceptance tests)

### Acceptance Tests

Each resource/data source should have:
- [ ] Create and read test
- [ ] Update test (modify attributes)
- [ ] Import test (import existing resource)
- [ ] Delete test (handled by test framework cleanup)
- [ ] Edge cases (optional fields, computed fields)

Location: `internal/provider/*_test.go` (with `TestAcc` prefix)

### Test Infrastructure

- [ ] Mock server for offline testing
- [ ] Test fixtures for common scenarios
- [ ] CI/CD pipeline integration
- [ ] Test coverage reporting

### Current Test Coverage

| Component | Coverage |
|-----------|----------|
| API Client | ✅ Unit tests |
| `stonebranch_task_unix` | ✅ Acceptance tests |
| `stonebranch_script` | ✅ Acceptance tests |
| `stonebranch_trigger_time` | ✅ Acceptance tests |
| `stonebranch_task_file_transfer` | ✅ Acceptance tests |
| `stonebranch_credential` | ✅ Acceptance tests |

---

## Documentation

### Provider Documentation

- [ ] Provider configuration guide
- [ ] Authentication setup
- [ ] Environment variables
- [ ] Troubleshooting guide

### Resource Documentation

Each resource should have:
- [ ] Description and use cases
- [ ] Full argument reference
- [ ] Attribute reference (computed values)
- [ ] Import syntax
- [ ] Example configurations
- [ ] Related resources

### Guides & Tutorials

- [ ] Getting started guide
- [ ] Basic Unix task example
- [ ] Workflow creation guide
- [ ] Scheduling with calendars
- [ ] Credential management best practices
- [ ] Migration guide (from manual configuration)

### Generated Documentation

- [ ] Set up `tfplugindocs` generation
- [ ] Schema descriptions for all attributes
- [ ] Example files in `examples/` directory

### Documentation Location

```
docs/
├── index.md                    # Provider overview
├── guides/
│   ├── getting-started.md
│   ├── authentication.md
│   └── workflows.md
├── resources/
│   ├── task_unix.md
│   ├── task_windows.md
│   └── ...
└── data-sources/
    ├── agent.md
    └── ...
```

---

## Resource Summary

### Total Resources Planned

| Category | Count | Implemented |
|----------|-------|-------------|
| Task Types | 20 | 2 |
| Trigger Types | 12 | 1 |
| Connection Types | 5 | 0 |
| Supporting Resources | 15 | 2 |
| Data Sources | 11 | 0 |
| **Total** | **63** | **5** |

### Priority Breakdown

| Priority | Description | Count |
|----------|-------------|-------|
| P0 | Essential/blocking | 18 |
| P1 | Important | 20 |
| P2 | Nice to have | 18 |
| P3 | Future/low usage | 7 |

---

## Contributing

### Adding a New Resource

1. Create resource file: `internal/provider/resource_{name}.go`
2. Define schema with all attributes
3. Implement CRUD operations
4. Register in `provider.go` Resources() method
5. Create acceptance tests: `internal/provider/resource_{name}_test.go`
6. Add example: `examples/resources/{name}/main.tf`
7. Update this roadmap

### Adding a New Data Source

1. Create data source file: `internal/provider/datasource_{name}.go`
2. Define schema (read-only attributes)
3. Implement Read operation
4. Register in `provider.go` DataSources() method
5. Create acceptance tests
6. Add example
7. Update this roadmap

### Code Style

- Follow existing patterns in `resource_task_unix.go`
- Use meaningful attribute descriptions
- Handle nil values gracefully
- Include proper error diagnostics

---

## Milestones

### v0.1.0 - Foundation (Current)
- [x] Provider scaffold
- [x] API client
- [x] Unix task resource
- [x] Script resource
- [x] Time trigger resource
- [x] File transfer task resource
- [x] Credential resource
- [ ] Basic documentation

### v0.2.0 - Core Tasks
- [ ] Windows task resource
- [ ] SQL task resource
- [ ] Email task resource
- [ ] Web service task resource
- [ ] Variable resource
- [ ] Database connection resource

### v0.3.0 - Workflows
- [ ] Workflow task resource
- [ ] Additional trigger types
- [ ] Virtual resources

### v0.4.0 - Scheduling
- [ ] Calendar resource
- [ ] Custom day resource
- [ ] Cron trigger

### v0.5.0 - Infrastructure
- [ ] Agent cluster resource
- [ ] Agent data sources
- [ ] Business service resource

### v0.6.0 - Data Sources
- [ ] Task data sources
- [ ] Trigger data sources
- [ ] General lookup data sources

### v1.0.0 - Production Ready
- [ ] All P0 and P1 resources
- [ ] Complete documentation
- [ ] Full test coverage
- [ ] Terraform Registry publication

---

## API Reference

- **Base URL:** `https://optionmetricsdev.stonebranch.cloud`
- **API Version:** 7.9.1.0
- **Spec File:** `openapi.yaml`

---

## Notes

- Each task type is implemented as a separate resource for type safety
- Resources support import via name (not sysId)
- Triggers are created disabled by default
- Business services are referenced via `opswise_groups` attribute
- API uses optimistic locking via `version` field

---

*Last updated: January 2026*

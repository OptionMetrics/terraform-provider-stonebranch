# Stonebranch Time Trigger Example

This example demonstrates how to create time-based triggers in Stonebranch Universal Controller.

## Resources Created

- `stonebranch_trigger_time.daily` - Daily trigger at 9:00 AM
- `stonebranch_trigger_time.hourly` - Hourly trigger
- `stonebranch_trigger_time.weekdays` - Weekday-only trigger
- `stonebranch_task_unix.daily_job` - Task triggered daily
- `stonebranch_task_unix.hourly_job` - Task triggered hourly
- `stonebranch_task_unix.business_job` - Task triggered on weekdays

## Usage

```bash
# Set your API token
export STONEBRANCH_API_TOKEN="your-token"

# Initialize and apply
terraform init
terraform plan
terraform apply
```

## Variables

| Name | Description | Default |
|------|-------------|---------|
| `agent_var` | Variable name for the agent | `agent_name` |
| `time_zone` | Time zone for scheduling | `America/New_York` |

## Trigger Attributes

| Attribute | Description |
|-----------|-------------|
| `time` | Time of day (HH:MM format) |
| `time_zone` | IANA time zone identifier |
| `time_style` | `Interval` for recurring triggers |
| `time_interval` | Interval value (with `time_interval_units`) |
| `time_interval_units` | `Minutes`, `Hours`, etc. |
| `monday` - `sunday` | Day-of-week flags |
| `calendar` | Business calendar reference |
| `enabled` | Whether the trigger is active |

## Important Notes

- **Triggers are created disabled by default** - Set `enabled = true` to activate
- Multiple tasks can be triggered simultaneously by listing them in `tasks`
- Use business calendars for complex scheduling requirements
- Time zones use IANA identifiers (e.g., `America/New_York`, `UTC`, `Europe/London`)

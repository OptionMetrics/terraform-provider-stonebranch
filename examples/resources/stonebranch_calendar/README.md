# Stonebranch Calendar Example

This example demonstrates how to create calendars in Stonebranch Universal Controller.

## Resources Created

- `stonebranch_calendar.basic` - Simple calendar with just a name
- `stonebranch_calendar.business` - Standard business calendar (Monday-Friday)
- `stonebranch_calendar.fiscal` - Calendar with fiscal quarter definitions
- `stonebranch_calendar.weekend` - Calendar for weekend operations
- `stonebranch_trigger_time.daily_with_calendar` - Time trigger using the business calendar

## Usage

```bash
# Set your API token
export STONEBRANCH_API_TOKEN="your-token"

# Initialize and apply
terraform init
terraform plan
terraform apply
```

## Attributes

| Attribute | Description | Required |
|-----------|-------------|----------|
| `name` | Unique name of the calendar | Yes |
| `comments` | Description/comments | No |
| `business_days` | Comma-separated business days (monday,tuesday,etc.) | No |
| `first_day_of_week` | First day of week (sunday or monday) | No |
| `first_quarter_month` | Month when Q1 starts (1-12) | No |
| `first_quarter_day` | Day when Q1 starts (1-31) | No |
| `second_quarter_month` | Month when Q2 starts (1-12) | No |
| `second_quarter_day` | Day when Q2 starts (1-31) | No |
| `third_quarter_month` | Month when Q3 starts (1-12) | No |
| `third_quarter_day` | Day when Q3 starts (1-31) | No |
| `fourth_quarter_month` | Month when Q4 starts (1-12) | No |
| `fourth_quarter_day` | Day when Q4 starts (1-31) | No |
| `opswise_groups` | Business services this calendar belongs to | No |

## Notes

- Calendars are used by triggers to control when tasks run
- Business days define which days of the week are considered working days
- Quarter definitions are used for scheduling based on fiscal periods
- Reference a calendar in triggers using the `calendar` attribute

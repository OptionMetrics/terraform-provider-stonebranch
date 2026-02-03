# Example: Timer task with duration-based delay (30 seconds)
resource "stonebranch_task_timer" "wait_30_seconds" {
  name           = "Wait 30 Seconds"
  summary        = "Pauses workflow execution for 30 seconds"
  sleep_type     = "Duration"
  sleep_duration = "00:00:00:30"
}

# Example: Timer task waiting for 5 minutes
resource "stonebranch_task_timer" "wait_5_minutes" {
  name           = "Wait 5 Minutes"
  summary        = "Pauses workflow execution for 5 minutes"
  sleep_type     = "Duration"
  sleep_duration = "00:00:05:00"
}

# Example: Timer task waiting for 1 hour 30 minutes
resource "stonebranch_task_timer" "wait_1h30m" {
  name           = "Wait 1 Hour 30 Minutes"
  summary        = "Pauses workflow for 1 hour and 30 minutes"
  sleep_type     = "Duration"
  sleep_duration = "00:01:30:00"
}

# Example: Timer task waiting for a specific time (same day)
resource "stonebranch_task_timer" "wait_until_2pm" {
  name                 = "Wait Until 2PM"
  summary              = "Pauses workflow until 2:00 PM same day"
  sleep_type           = "Time"
  sleep_time           = "14:00"
  sleep_day_constraint = "Same Day"
}

# Example: Timer task waiting until next day
resource "stonebranch_task_timer" "wait_next_day_9am" {
  name                 = "Wait Until 9AM Next Day"
  summary              = "Pauses workflow until 9 AM next day"
  sleep_type           = "Time"
  sleep_time           = "09:00"
  sleep_day_constraint = "Next Day"
}

# Example: Timer task with relative time delay
resource "stonebranch_task_timer" "wait_relative_1h30m" {
  name       = "Wait Relative 1h30m"
  summary    = "Pauses workflow for 1 hour and 30 minutes from task start"
  sleep_type = "Relative Time"
  sleep_time = "01:30"
}

# Example: Timer task waiting until next business day
resource "stonebranch_task_timer" "wait_next_business_day" {
  name                 = "Wait for Next Business Day"
  summary              = "Pauses workflow until 9 AM on the next business day"
  sleep_type           = "Time"
  sleep_time           = "09:00"
  sleep_day_constraint = "Next Business Day"
}

# Example: Timer task with variables
resource "stonebranch_task_timer" "configurable_delay" {
  name           = "Configurable Delay"
  summary        = "Timer task with configurable delay via variable"
  sleep_type     = "Duration"
  sleep_duration = "00:00:10:00"

  variables = [
    {
      name        = "delay_override"
      value       = "10"
      description = "Override delay amount if needed"
    }
  ]
}

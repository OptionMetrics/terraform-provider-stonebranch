variable "agent_name" {
  description = "Name of the StoneBranch agent to run tasks on"
  type        = string
}

variable "file_monitor_task_name" {
  description = "Name of the existing file monitor task that detects file events"
  type        = string
}

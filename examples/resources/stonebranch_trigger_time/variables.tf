variable "agent_var" {
  description = "Variable name containing the agent to run tasks on (resolved at runtime)"
  type        = string
  default     = "agent_name"
}

variable "time_zone" {
  description = "Time zone for trigger scheduling"
  type        = string
  default     = "America/New_York"
}

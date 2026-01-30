variable "agent_var" {
  description = "Variable name containing the Windows agent to run tasks on (resolved at runtime)"
  type        = string
  default     = "windows_agent_name"
}

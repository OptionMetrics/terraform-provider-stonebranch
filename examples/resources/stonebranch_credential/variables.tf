variable "agent_var" {
  description = "Variable name containing the agent to run tasks on (resolved at runtime)"
  type        = string
  default     = "agent_name"
}

variable "app_user" {
  description = "Application username"
  type        = string
  default     = "appuser"
}

variable "app_password" {
  description = "Application password"
  type        = string
  sensitive   = true
}

variable "sftp_user" {
  description = "SFTP username"
  type        = string
  default     = "sftpuser"
}

variable "sftp_password" {
  description = "SFTP password"
  type        = string
  sensitive   = true
}

variable "agent_var" {
  description = "Variable name containing the agent to run tasks on (resolved at runtime)"
  type        = string
  default     = "agent_name"
}

variable "sftp_server" {
  description = "SFTP server hostname"
  type        = string
  default     = "sftp.example.com"
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

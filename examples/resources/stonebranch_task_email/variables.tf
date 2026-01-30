# Variables for Stonebranch Email Task Example

variable "smtp_host" {
  description = "SMTP server hostname"
  type        = string
  default     = "smtp.example.com"
}

variable "smtp_user" {
  description = "SMTP authentication username"
  type        = string
  default     = "smtp-user"
}

variable "smtp_password" {
  description = "SMTP authentication password"
  type        = string
  sensitive   = true
  default     = ""
}

variable "sender_email" {
  description = "Default sender email address"
  type        = string
  default     = "automation@example.com"
}

variable "notification_recipients" {
  description = "Comma-separated list of notification email recipients"
  type        = string
  default     = "team@example.com"
}

variable "critical_recipients" {
  description = "Comma-separated list of critical alert email recipients"
  type        = string
  default     = "oncall@example.com,manager@example.com"
}

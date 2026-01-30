# Variables for Stonebranch Email Connection Example

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

variable "imap_host" {
  description = "IMAP server hostname"
  type        = string
  default     = "imap.example.com"
}

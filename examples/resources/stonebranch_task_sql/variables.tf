variable "db_host" {
  description = "Database server hostname"
  type        = string
  default     = "mysql.example.com"
}

variable "db_port" {
  description = "Database server port"
  type        = number
  default     = 3306
}

variable "db_name" {
  description = "Database name"
  type        = string
  default     = "production"
}

variable "db_user" {
  description = "Database username"
  type        = string
  default     = "app_user"
}

variable "db_password" {
  description = "Database password"
  type        = string
  sensitive   = true
}

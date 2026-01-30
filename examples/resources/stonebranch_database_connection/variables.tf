# MySQL connection variables
variable "mysql_host" {
  description = "MySQL server hostname"
  type        = string
  default     = "mysql.example.com"
}

variable "mysql_port" {
  description = "MySQL server port"
  type        = number
  default     = 3306
}

variable "mysql_database" {
  description = "MySQL database name"
  type        = string
  default     = "production"
}

variable "mysql_user" {
  description = "MySQL username"
  type        = string
  default     = "app_user"
}

variable "mysql_password" {
  description = "MySQL password"
  type        = string
  sensitive   = true
}

# PostgreSQL connection variables
variable "postgres_host" {
  description = "PostgreSQL server hostname"
  type        = string
  default     = "postgres.example.com"
}

variable "postgres_port" {
  description = "PostgreSQL server port"
  type        = number
  default     = 5432
}

variable "postgres_database" {
  description = "PostgreSQL database name"
  type        = string
  default     = "analytics"
}

variable "postgres_user" {
  description = "PostgreSQL username"
  type        = string
  default     = "analytics_user"
}

variable "postgres_password" {
  description = "PostgreSQL password"
  type        = string
  sensitive   = true
}

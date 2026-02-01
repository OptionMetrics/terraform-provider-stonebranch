# Stonebranch SQL Task Example
#
# This example demonstrates how to create SQL task resources in Stonebranch.
# SQL tasks execute queries against database connections.
#
# Usage:
#   export STONEBRANCH_API_TOKEN="your-token"
#   export TF_VAR_db_password="your-db-password"
#   terraform init
#   terraform plan
#   terraform apply

terraform {
  required_providers {
    stonebranch = {
      source = "stonebranch/stonebranch"
    }
  }
}

provider "stonebranch" {
  # Uses STONEBRANCH_API_TOKEN environment variable
}

# Database credential
resource "stonebranch_credential" "db" {
  name             = "tf-example-sql-creds"
  description      = "Database credentials for SQL tasks"
  runtime_user     = var.db_user
  runtime_password = var.db_password
}

# Database connection
resource "stonebranch_database_connection" "main" {
  name        = "tf-example-sql-conn"
  description = "Main database for SQL tasks"
  db_type     = "MySQL"
  db_url      = "jdbc:mysql://${var.db_host}:${var.db_port}/${var.db_name}"
  db_driver   = "com.mysql.cj.jdbc.Driver"
  db_max_rows = 10000
  credentials = stonebranch_credential.db.name
}

# Simple SELECT query task
resource "stonebranch_task_sql" "daily_count" {
  name                = "tf-example-daily-count"
  summary             = "Count daily transactions"
  database_connection = stonebranch_database_connection.main.name
  sql_command         = "SELECT COUNT(*) as total FROM transactions WHERE DATE(created_at) = CURDATE()"
}

# Query with max rows limit
resource "stonebranch_task_sql" "recent_orders" {
  name                = "tf-example-recent-orders"
  summary             = "Fetch recent orders for processing"
  database_connection = stonebranch_database_connection.main.name
  sql_command         = "SELECT order_id, customer_id, total FROM orders WHERE status = 'pending' ORDER BY created_at"
  max_rows            = 100
}

# Data cleanup task
resource "stonebranch_task_sql" "cleanup_old_logs" {
  name                = "tf-example-cleanup-logs"
  summary             = "Remove logs older than 90 days"
  database_connection = stonebranch_database_connection.main.name
  sql_command         = "DELETE FROM audit_logs WHERE created_at < DATE_SUB(NOW(), INTERVAL 90 DAY)"
}

# Task with result processing
resource "stonebranch_task_sql" "check_threshold" {
  name                = "tf-example-check-threshold"
  summary             = "Check if error count exceeds threshold"
  database_connection = stonebranch_database_connection.main.name
  sql_command         = "SELECT COUNT(*) as error_count FROM errors WHERE created_at > DATE_SUB(NOW(), INTERVAL 1 HOUR)"
  result_processing   = "First Column"
  column_op           = ">"
  column_value        = "100"
}

# Task with retry configuration
resource "stonebranch_task_sql" "critical_report" {
  name                = "tf-example-critical-report"
  summary             = "Generate critical daily report"
  database_connection = stonebranch_database_connection.main.name
  sql_command         = <<-EOT
    SELECT
      DATE(created_at) as report_date,
      COUNT(*) as total_transactions,
      SUM(amount) as total_amount
    FROM transactions
    WHERE DATE(created_at) = DATE_SUB(CURDATE(), INTERVAL 1 DAY)
    GROUP BY DATE(created_at)
  EOT

  retry_maximum  = 3
  retry_interval = 60
}

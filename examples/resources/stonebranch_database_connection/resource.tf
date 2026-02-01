# Stonebranch Database Connection Example
#
# This example demonstrates how to create database connection resources in Stonebranch.
# Database connections define how SQL tasks connect to databases.
#
# Usage:
#   export STONEBRANCH_API_TOKEN="your-token"
#   export TF_VAR_mysql_password="your-mysql-password"
#   export TF_VAR_postgres_password="your-postgres-password"
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

# Credentials for database connections
resource "stonebranch_credential" "mysql" {
  name             = "tf-example-mysql-creds"
  description      = "MySQL database credentials"
  runtime_user     = var.mysql_user
  runtime_password = var.mysql_password
}

resource "stonebranch_credential" "postgres" {
  name             = "tf-example-postgres-creds"
  description      = "PostgreSQL database credentials"
  runtime_user     = var.postgres_user
  runtime_password = var.postgres_password
}

# MySQL database connection
resource "stonebranch_database_connection" "mysql" {
  name        = "tf-example-mysql-conn"
  description = "MySQL production database"
  db_type     = "MySQL"
  db_url      = "jdbc:mysql://${var.mysql_host}:${var.mysql_port}/${var.mysql_database}"
  db_driver   = "com.mysql.cj.jdbc.Driver"
  db_max_rows = 10000
  credentials = stonebranch_credential.mysql.name
}

# PostgreSQL database connection
resource "stonebranch_database_connection" "postgres" {
  name        = "tf-example-postgres-conn"
  description = "PostgreSQL analytics database"
  db_type     = "PostgreSQL"
  db_url      = "jdbc:postgresql://${var.postgres_host}:${var.postgres_port}/${var.postgres_database}"
  db_driver   = "org.postgresql.Driver"
  db_max_rows = 50000
  credentials = stonebranch_credential.postgres.name
}

# Oracle database connection (example without variable references)
resource "stonebranch_database_connection" "oracle" {
  name        = "tf-example-oracle-conn"
  description = "Oracle ERP database"
  db_type     = "Oracle"
  db_url      = "jdbc:oracle:thin:@//oracle.example.com:1521/ORCL"
  db_driver   = "oracle.jdbc.OracleDriver"
  db_max_rows = 0  # Unlimited
}

# SQL Server database connection
resource "stonebranch_database_connection" "sqlserver" {
  name        = "tf-example-sqlserver-conn"
  description = "SQL Server reporting database"
  db_type     = "SQL Server"
  db_url      = "jdbc:sqlserver://sqlserver.example.com:1433;databaseName=ReportingDB"
  db_driver   = "com.microsoft.sqlserver.jdbc.SQLServerDriver"
  db_max_rows = 25000
}

# Stonebranch Stored Procedure Task Example
#
# This example demonstrates how to create a stored procedure task
# that executes a database stored procedure with parameters.

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

# First, create a credential for the database connection
resource "stonebranch_credential" "db_cred" {
  name             = "my-db-credential"
  runtime_user     = "db_user"
  runtime_password = "secure_password"
}

# Create a database connection
resource "stonebranch_database_connection" "mysql" {
  name        = "my-mysql-connection"
  db_url      = "jdbc:mysql://localhost:3306/mydb"
  db_driver   = "com.mysql.cj.jdbc.Driver"
  credentials = stonebranch_credential.db_cred.name
}

# Basic stored procedure task
resource "stonebranch_task_stored_procedure" "basic" {
  name                = "my-basic-stored-proc"
  stored_proc_name    = "sp_get_user_count"
  database_connection = stonebranch_database_connection.mysql.name
}

# Stored procedure task with input and output parameters
resource "stonebranch_task_stored_procedure" "with_params" {
  name                = "my-stored-proc-with-params"
  stored_proc_name    = "sp_process_order"
  database_connection = stonebranch_database_connection.mysql.name
  summary             = "Process an order and return the order ID"

  parameters = [
    {
      param_var   = "customer_id_param"
      description = "Customer ID"
      param_mode  = "Input"
      param_type  = "INTEGER"
      input_value = "12345"
      position    = 1
    },
    {
      param_var   = "order_amount_param"
      description = "Order amount"
      param_mode  = "Input"
      param_type  = "DECIMAL"
      input_value = "99.99"
      position    = 2
    },
    {
      param_var      = "order_id_param"
      description    = "Generated Order ID"
      param_mode     = "Output"
      param_type     = "INTEGER"
      output_value   = "order_id_var"
      variable_scope = "Self"
      position       = 3
    }
  ]
}

# Stored procedure task using a connection variable
resource "stonebranch_task_stored_procedure" "with_connection_var" {
  name             = "my-dynamic-stored-proc"
  stored_proc_name = "sp_daily_report"
  connection_var   = "db_connection_name"
  summary          = "Uses a variable to determine which database connection to use"
}

# Stored procedure task with result processing
resource "stonebranch_task_stored_procedure" "with_result_processing" {
  name                = "my-result-checked-proc"
  stored_proc_name    = "sp_validate_data"
  database_connection = stonebranch_database_connection.mysql.name

  result_processing = "Row Count"
  result_op         = ">"
  result_value      = "0"
  exit_codes        = "0"
}

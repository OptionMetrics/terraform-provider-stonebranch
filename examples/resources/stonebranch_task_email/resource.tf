# Stonebranch Email Task Example
#
# This example demonstrates how to create email task resources in Stonebranch.
# Email tasks send emails via configured email connections.
#
# Usage:
#   export STONEBRANCH_API_TOKEN="your-token"
#   export TF_VAR_smtp_password="your-smtp-password"
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

# First, create an email connection for the tasks to use
resource "stonebranch_email_connection" "main" {
  name             = "tf-example-email-conn"
  description      = "Email connection for automation notifications"
  smtp             = var.smtp_host
  smtp_port        = 587
  smtp_starttls    = true
  authentication   = true
  default_user     = var.smtp_user
  default_password = var.smtp_password
  email_address    = var.sender_email
}

# Simple notification email task
resource "stonebranch_task_email" "simple_notification" {
  name          = "tf-example-simple-email"
  summary       = "Send simple notification email"
  email_connection = stonebranch_email_connection.main.name
  to_recipients = var.notification_recipients
  subject       = "Automation Notification"
  body          = "This is an automated notification from Stonebranch."
}

# Email task with multiple recipients
resource "stonebranch_task_email" "multi_recipient" {
  name           = "tf-example-multi-recipient-email"
  summary        = "Send email to multiple recipients"
  connection     = stonebranch_email_connection.main.name
  to_recipients  = "primary@example.com,secondary@example.com"
  cc_recipients  = "manager@example.com"
  bcc_recipients = "archive@example.com"
  subject        = "Report Summary"
  body           = <<-EOT
    Hello Team,

    Please find the daily report summary attached.

    This is an automated message.

    Best regards,
    Automation System
  EOT
}

# Email task with retry configuration
resource "stonebranch_task_email" "with_retry" {
  name                   = "tf-example-email-with-retry"
  summary                = "Critical notification with retry"
  connection             = stonebranch_email_connection.main.name
  to_recipients          = var.critical_recipients
  subject                = "[CRITICAL] System Alert"
  body                   = "A critical condition has been detected. Please investigate immediately."
  retry_maximum          = 5
  retry_interval         = 120  # 2 minutes between retries
  retry_suppress_failure = false
}

# Email task using variables for dynamic content
resource "stonebranch_task_email" "dynamic_content" {
  name          = "tf-example-dynamic-email"
  summary       = "Email with dynamic content using variables"
  email_connection = stonebranch_email_connection.main.name
  to_recipients = "$${ops_recipient_list}"  # Uses Stonebranch variable
  subject       = "[$(ops_environment)] Job Status Update"
  body          = <<-EOT
    Environment: $(ops_environment)
    Job Name: $(ops_job_name)
    Status: $(ops_exit_code)
    Start Time: $(ops_start_time)
    End Time: $(ops_end_time)
  EOT
}

# Email task with reply-to address
resource "stonebranch_task_email" "with_reply_to" {
  name          = "tf-example-email-reply-to"
  summary       = "Email with custom reply-to"
  email_connection = stonebranch_email_connection.main.name
  to_recipients = var.notification_recipients
  reply_to      = "support@example.com"
  subject       = "Action Required: Please Review"
  body          = "Please reply to this email if you have any questions."
}

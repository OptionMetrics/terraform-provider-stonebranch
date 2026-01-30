# Stonebranch Email Connection Example
#
# This example demonstrates how to create email connection resources in Stonebranch.
# Email connections define SMTP and IMAP server settings for sending and receiving emails.
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

# Basic SMTP connection (no authentication)
resource "stonebranch_email_connection" "basic" {
  name          = "tf-example-email-basic"
  description   = "Basic SMTP connection without authentication"
  smtp          = "smtp.example.com"
  smtp_port     = 25
  email_address = "noreply@example.com"
}

# SMTP connection with STARTTLS (recommended for port 587)
resource "stonebranch_email_connection" "starttls" {
  name               = "tf-example-email-starttls"
  description        = "SMTP connection with STARTTLS encryption"
  smtp               = var.smtp_host
  smtp_port          = 587
  smtp_starttls      = true
  authentication     = true
  default_user       = var.smtp_user
  default_password   = var.smtp_password
  email_address      = var.sender_email
}

# SMTP connection with SSL/TLS (for port 465)
resource "stonebranch_email_connection" "ssl" {
  name             = "tf-example-email-ssl"
  description      = "SMTP connection with SSL/TLS encryption"
  smtp             = var.smtp_host
  smtp_port        = 465
  smtp_ssl         = true
  authentication   = true
  default_user     = var.smtp_user
  default_password = var.smtp_password
  email_address    = var.sender_email
}

# Full-featured email connection with both SMTP and IMAP
resource "stonebranch_email_connection" "full" {
  name             = "tf-example-email-full"
  description      = "Full email connection with SMTP and IMAP"

  # SMTP settings for sending
  smtp             = var.smtp_host
  smtp_port        = 587
  smtp_starttls    = true

  # Authentication
  authentication   = true
  default_user     = var.smtp_user
  default_password = var.smtp_password

  # Sender address
  email_address    = var.sender_email

  # IMAP settings for reading (used by Email Monitor tasks)
  imap             = var.imap_host
  imap_port        = 993
  imap_ssl         = true
  trash_folder     = "Trash"
}

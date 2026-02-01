# Stonebranch Web Service Task Example
#
# This example demonstrates how to create web service tasks
# that call REST or SOAP web services.

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

# Basic GET request
resource "stonebranch_task_web_service" "get_example" {
  name      = "my-get-request"
  url       = "https://api.example.com/users"
  mime_type = "application/json"
}

# POST request with JSON payload
resource "stonebranch_task_web_service" "post_example" {
  name        = "my-post-request"
  url         = "https://api.example.com/users"
  http_method = "POST"
  mime_type   = "application/json"
  payload     = jsonencode({
    name  = "John Doe"
    email = "john@example.com"
  })
  summary = "Create a new user"
}

# Request with custom headers
resource "stonebranch_task_web_service" "with_headers" {
  name      = "my-request-with-headers"
  url       = "https://api.example.com/data"
  mime_type = "application/json"

  http_headers = [
    {
      name  = "Authorization"
      value = "Bearer ${var.api_token}"
    },
    {
      name  = "Accept"
      value = "application/json"
    },
    {
      name  = "X-Custom-Header"
      value = "custom-value"
    }
  ]
}

# Request with URL parameters
resource "stonebranch_task_web_service" "with_params" {
  name      = "my-request-with-params"
  url       = "https://api.example.com/search"
  mime_type = "application/json"

  url_parameters = [
    {
      name  = "query"
      value = "search-term"
    },
    {
      name  = "limit"
      value = "10"
    }
  ]
}

# Request with authentication
resource "stonebranch_credential" "api_creds" {
  name             = "my-api-credentials"
  runtime_user     = "api_user"
  runtime_password = "api_password"
}

resource "stonebranch_task_web_service" "with_auth" {
  name        = "my-authenticated-request"
  url         = "https://api.example.com/secure"
  mime_type   = "application/json"
  http_auth   = "Basic"
  credentials = stonebranch_credential.api_creds.name
}

# Request with response processing
resource "stonebranch_task_web_service" "with_response_check" {
  name                      = "my-request-with-response-check"
  url                       = "https://api.example.com/status"
  mime_type                 = "application/json"
  response_processing_type  = "Status Code"
  status_code_range         = "200-299"
}

# Request with retry configuration
resource "stonebranch_task_web_service" "with_retry" {
  name           = "my-request-with-retry"
  url            = "https://api.example.com/data"
  mime_type      = "application/json"
  timeout        = 60
  retry_maximum  = 3
  retry_interval = 30
}

# Variable for API token (would be set via environment or tfvars)
variable "api_token" {
  description = "API token for authentication"
  type        = string
  sensitive   = true
  default     = "example-token"
}

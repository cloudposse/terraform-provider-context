terraform {
  required_providers {
    context = {
      source = "registry.terraform.io/cloudposse/context"
    }
  }
}

provider "context" {
  delimiter = "-"
  enabled   = true

  properties = {
    namespace   = {}
    environment = {}
    stage       = {}
    name        = {}
  }

  property_order = ["namespace", "environment", "stage", "name"]

  values = {
    namespace   = "acme-corporation"
    environment = "production"
    stage       = "blue-green"
    name        = "my-application-service"
  }
}

# Without max_length, the full label would be:
#   "acme-corporation-production-blue-green-my-application-service" (62 chars)

# Truncate to 32 characters (e.g. for resource names with length limits)
data "context_label" "truncated" {
  max_length = 32
  truncate   = true
}

# Truncate to 63 characters (e.g. S3 bucket name limit)
data "context_label" "s3_bucket" {
  max_length = 63
  truncate   = true
}

output "full_label" {
  description = "Full label without truncation"
  value       = data.context_label.s3_bucket.rendered
}

output "truncated_label" {
  description = "Label truncated to 32 characters"
  value       = data.context_label.truncated.rendered
}

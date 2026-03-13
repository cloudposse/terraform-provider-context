terraform {
  required_providers {
    context = {
      source = "registry.terraform.io/cloudposse/context"
    }
  }
}

# This example shows how to use property validation to enforce naming conventions.
# The provider validates values at configuration time, catching errors early.
provider "context" {
  delimiter = "-"
  enabled   = true

  properties = {
    namespace = {
      required         = true
      min_length       = 2
      max_length       = 10
      validation_regex = "^[a-z]+$"
    }
    environment = {
      required         = true
      validation_regex = "^(dev|staging|prod)$"
    }
    region = {
      required         = true
      validation_regex = "^[a-z]{2}-[a-z]+-[0-9]+$"
    }
    name = {
      required   = true
      max_length = 32
    }
  }

  property_order = ["namespace", "environment", "region", "name"]

  values = {
    namespace   = "acme"
    environment = "prod"
    region      = "us-east-1"
    name        = "api-gateway"
  }
}

data "context_label" "this" {}

data "context_tags" "this" {}

output "label" {
  description = "Generated label from validated properties"
  value       = data.context_label.this.rendered
}

output "tags" {
  description = "Generated tags from validated properties"
  value       = data.context_tags.this.tags
}

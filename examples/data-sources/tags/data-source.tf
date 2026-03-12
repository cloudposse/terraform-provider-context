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
    namespace   = { required = true }
    environment = { required = true }
    stage       = {}
    name        = { required = true }
  }

  property_order = ["namespace", "environment", "stage", "name"]

  tags_key_case   = "title"
  tags_value_case = "lower"

  values = {
    namespace   = "Acme"
    environment = "Production"
    stage       = "Blue"
    name        = "WebApp"
  }
}

# Generate tags using the provider's default case settings
data "context_tags" "default" {}

# Override case transformations for a specific use case (e.g. CloudFormation-style uppercase keys)
data "context_tags" "uppercase_keys" {
  tags_key_case   = "upper"
  tags_value_case = "upper"
}

# Add extra values that aren't part of the provider context
data "context_tags" "with_extras" {
  values = {
    "cost_center" = "engineering"
    "team"        = "platform"
  }
}

output "tags_map" {
  description = "Tags as a map (e.g. for aws_instance.tags)"
  value       = data.context_tags.default.tags
}

output "tags_as_list" {
  description = "Tags as a list of {Key, Value} objects (e.g. for aws_autoscaling_group.tag)"
  value       = data.context_tags.default.tags_as_list
}

output "tags_uppercase" {
  description = "Same tags with uppercase keys and values"
  value       = data.context_tags.uppercase_keys.tags
}

output "tags_with_extras" {
  description = "Tags with additional custom values merged in"
  value       = data.context_tags.with_extras.tags
}

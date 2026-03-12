terraform {
  required_providers {
    context = {
      source  = "cloudposse/context"
      version = ">= 0.4.0"
    }
  }
}

# A complete example demonstrating the full power of the context provider.
#
# Scenario: A multi-tenant SaaS platform that needs consistent naming and tagging
# across all infrastructure resources. The context provider generates labels for
# resource names and tags for cost allocation and ownership tracking.

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
    tenant = {
      required = true
    }
    environment = {
      required         = true
      validation_regex = "^(dev|staging|prod)$"
    }
    name = {
      required   = true
      max_length = 32
    }
    # "attributes" is not included in tags by default
    attributes = {
      include_in_tags = false
    }
  }

  property_order = ["namespace", "tenant", "environment", "name", "attributes"]

  tags_key_case   = "title"
  tags_value_case = "none"

  values = {
    namespace   = "acme"
    tenant      = "platform"
    environment = "prod"
    name        = "web"
  }
}

# --- Labels ---

# Full delimited label: "acme-platform-prod-web"
data "context_label" "default" {}

# Short label for S3 buckets (namespace-name only): "acme-web"
data "context_label" "s3_bucket" {
  properties = ["namespace", "name"]
  values = {
    attributes = "assets"
  }
}

# Templated label for ARN-style paths: "acme/platform/prod/web"
data "context_label" "path" {
  template = "{{.namespace}}/{{.tenant}}/{{.environment}}/{{.name}}"
}

# Label with value overrides for a specific component
data "context_label" "database" {
  values = {
    name = "postgres"
  }
}

# Truncated label for resources with length limits
data "context_label" "truncated" {
  max_length = 20
  truncate   = true
}

# --- Tags ---

# Default tags with title-case keys
data "context_tags" "default" {}

# Tags with additional metadata merged in
data "context_tags" "with_metadata" {
  values = {
    managed_by = "terraform"
    team       = "infrastructure"
  }
}

# Tags with snake_case keys for systems that prefer it
data "context_tags" "snake_case" {
  tags_key_case = "snake"
}

# --- Config ---

# Read back the provider configuration (useful for passing to child modules)
data "context_config" "this" {}

# --- Outputs ---

output "default_label" {
  description = "Full delimited label"
  value       = data.context_label.default.rendered
  # => "acme-platform-prod-web"
}

output "s3_bucket_label" {
  description = "Short label for S3 bucket naming"
  value       = data.context_label.s3_bucket.rendered
  # => "acme-web"
}

output "path_label" {
  description = "Path-style label using template"
  value       = data.context_label.path.rendered
  # => "acme/platform/prod/web"
}

output "database_label" {
  description = "Label with overridden name"
  value       = data.context_label.database.rendered
  # => "acme-platform-prod-postgres"
}

output "truncated_label" {
  description = "Label truncated to 20 characters"
  value       = data.context_label.truncated.rendered
  # => "acme-platform-prod-w"
}

output "default_tags" {
  description = "Tags as a map with title-case keys"
  value       = data.context_tags.default.tags
  # => {"Environment" = "prod", "Name" = "web", "Namespace" = "acme", "Tenant" = "platform"}
}

output "tags_as_list" {
  description = "Tags as a list of {Key, Value} objects"
  value       = data.context_tags.default.tags_as_list
}

output "tags_with_metadata" {
  description = "Tags with extra metadata merged in"
  value       = data.context_tags.with_metadata.tags
}

output "snake_case_tags" {
  description = "Tags with snake_case keys"
  value       = data.context_tags.snake_case.tags
  # => {"environment" = "prod", "name" = "web", "namespace" = "acme", "tenant" = "platform"}
}

output "provider_delimiter" {
  description = "The delimiter configured on the provider"
  value       = data.context_config.this.delimiter
}

output "provider_enabled" {
  description = "Whether the provider is enabled"
  value       = data.context_config.this.enabled
}

output "provider_values" {
  description = "All values configured on the provider"
  value       = data.context_config.this.values
}

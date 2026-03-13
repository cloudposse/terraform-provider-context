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
    tenant      = {}
    environment = {}
    stage       = {}
    name        = { required = true }
  }

  property_order = ["namespace", "tenant", "environment", "stage", "name"]

  values = {
    namespace   = "acme"
    tenant      = "platform"
    environment = "us-east-1"
    stage       = "prod"
    name        = "api"
  }
}

# Full label using all properties: "acme-platform-us-east-1-prod-api"
data "context_label" "full" {}

# Short label using only namespace and name: "acme-api"
data "context_label" "short" {
  properties = ["namespace", "name"]
}

# Environment-scoped label: "acme-us-east-1-prod-api"
data "context_label" "env_scoped" {
  properties = ["namespace", "environment", "stage", "name"]
}

# Tenant-scoped label: "acme-platform-api"
data "context_label" "tenant_scoped" {
  properties = ["namespace", "tenant", "name"]
}

output "full_label" {
  description = "Label with all properties"
  value       = data.context_label.full.rendered
}

output "short_label" {
  description = "Short label with only namespace and name"
  value       = data.context_label.short.rendered
}

output "env_scoped_label" {
  description = "Label scoped to environment"
  value       = data.context_label.env_scoped.rendered
}

output "tenant_scoped_label" {
  description = "Label scoped to tenant"
  value       = data.context_label.tenant_scoped.rendered
}

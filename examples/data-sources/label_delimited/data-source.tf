terraform {
  required_providers {
    context = {
      source = "registry.terraform.io/cloudposse/context"
    }
  }
}

provider "context" {
  delimiter = "~"
  enabled   = false
  properties = {
    Namespace = {}
    Tenant    = {}
    Stage     = {}
    Name      = {}
  }

  property_order = ["Namespace", "Tenant", "Stage", "Name"]

  values = {
    "Namespace" = "cp"
    "Tenant"    = "core"
    "Stage"     = "prod"
    "Name"      = "example"
  }
}

data "context_label" "example" {
  values = {
    "Namespace" = "cp"
    "Tenant"    = "core"
  }
}

output "rendered" {
  value = data.context_label.example.rendered
}

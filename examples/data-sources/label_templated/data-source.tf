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
    namespace   = {}
    tenant      = {}
    stage       = {}
    environment = {}
    name        = {}
  }

  property_order = ["namespace", "tenant", "stage", "environment", "name"]

  values = {
    "namespace"   = "cp"
    "tenant"      = "core"
    "stage"       = "prod"
    "environment" = "ue1"
    "name"        = "example"
  }
}


data "context_label" "example" {
  template = "{{.namespace}}/{{.tenant}}/{{.stage}}/{{.name}}"
  values = {
    "tenant" = "plat"
    "stage"  = "dev"
  }
}

output "rendered" {
  value = data.context_label.example.rendered
}

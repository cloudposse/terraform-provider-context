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
    Namespace = { required = true }
    Tenant    = { required = true }
    Stage     = {}
    Name      = {}
  }

  property_order = ["Namespace", "Tenant", "Stage", "Name"]

  values = {
    //"Namespace" = "cp"
    //"Tenant" = "core"
    "Stage" = "prod"
    "Name"  = "example"
  }
}

data "context_label" "example" {
  values = {
    "Namespace" = "cp"
    "Tenant"    = "core"
  }
  #template = "{{.Namespace}}~{{.Tenant}}~{{.Stage}}~{{.Name}}"
  # values = {
  #   "Namespace" = "mdc"
  # }
  //properties = ["Tenant", "Stage", "Name", "Namespace"]
}

output "rendered" {
  value = data.context_label.example.rendered
}

terraform {
  required_providers {
    context = {
      source = "registry.terraform.io/cloudposse/context"
    }
  }
}

provider "context" {
  properties = {
    Namespace = { required = true }
    Tenant    = { required = true }
    Stage     = {}
    Name      = {}
  }

  values = {
    "Namespace" = "cp"
    "Tenant"    = "core"
    "Stage"     = "prod"
    "Name"      = "example"
  }
}

data "context_config" "example" {

}

locals {
  context = data.context_config.example
}

output "values" {
  value = local.context.values.Namespace
}

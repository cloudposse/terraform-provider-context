terraform {
  required_providers {
    context = {
      source = "registry.terraform.io/cloudposse/context"
    }
  }
}

provider "context" {
  properties = {
    Namespace = {
      required      = true
      tags_key_case = "upper"
    }
    Tenant = {
      required        = true
      tags_value_case = "lower"
    }
    Stage = {
      tags_key_case   = "snake"
      tags_value_case = "title"
    }
    Name = {}
  }

  tags_key_case   = "title"
  tags_value_case = "none"

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

output "tags" {
  value = local.context.tags
}

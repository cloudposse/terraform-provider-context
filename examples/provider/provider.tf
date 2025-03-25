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
      required = true
      # Override tag casing for this property
      tags_key_case = "upper"
    }
    Tenant = {
      required = true
      # Override tag casing for this property
      tags_value_case = "lower"
    }
    Stage = {
      # Override both key and value casing for this property
      tags_key_case   = "snake"
      tags_value_case = "title"
    }
    Name = {}
  }

  # Default casing for all properties (unless overridden at property level)
  tags_key_case   = "title"
  tags_value_case = "none"

  values = {
    "Namespace" = "cp"
    "Tenant"    = "core"
    "Stage"     = "prod"
    "Name"      = "example"
  }
}

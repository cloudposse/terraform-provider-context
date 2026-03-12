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

  # Remove any characters that aren't lowercase alphanumeric or hyphens.
  # Useful for generating DNS-safe names, S3 bucket names, etc.
  replace_chars_regex = "[^a-z0-9-]"

  properties = {
    namespace   = {}
    environment = {}
    name        = {}
  }

  property_order = ["namespace", "environment", "name"]

  values = {
    namespace   = "Acme Corp"
    environment = "US-East-1"
    name        = "My App!"
  }
}

# Characters matching the regex are stripped before joining.
# "Acme Corp" -> "cme orp" -> further stripped -> "cmeorp"
# Result: "cmeorp-s-ast-1-ypp"
data "context_label" "dns_safe" {}

# Override replace_chars_regex per label to allow underscores
data "context_label" "with_underscores" {
  replace_chars_regex = "[^a-z0-9_-]"
}

output "dns_safe_label" {
  description = "Label with non-DNS characters removed"
  value       = data.context_label.dns_safe.rendered
}

output "underscore_label" {
  description = "Label allowing underscores"
  value       = data.context_label.with_underscores.rendered
}

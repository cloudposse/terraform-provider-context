run "tags_example" {
  command = plan

  assert {
    condition     = lookup(data.context_tags.default.tags, "Namespace", "") == "acme"
    error_message = "Default tags should have Namespace = 'acme'"
  }

  assert {
    condition     = lookup(data.context_tags.default.tags, "Environment", "") == "production"
    error_message = "Default tags should have Environment = 'production'"
  }

  assert {
    condition     = lookup(data.context_tags.uppercase_keys.tags, "NAMESPACE", "") == "ACME"
    error_message = "Uppercase tags should have NAMESPACE = 'ACME'"
  }

  assert {
    condition     = lookup(data.context_tags.with_overrides.tags, "Environment", "") == "staging"
    error_message = "Overridden tags should have Environment = 'staging' (lowercased)"
  }

  assert {
    condition     = lookup(data.context_tags.with_overrides.tags, "Name", "") == "worker"
    error_message = "Overridden tags should have Name = 'worker' (lowercased)"
  }
}

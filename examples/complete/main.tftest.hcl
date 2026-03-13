run "complete_example" {
  command = plan

  assert {
    condition     = data.context_label.default.rendered == "acme-platform-prod-web"
    error_message = "Default label should be 'acme-platform-prod-web', got '${data.context_label.default.rendered}'"
  }

  assert {
    condition     = data.context_label.s3_bucket.rendered == "acme-web"
    error_message = "S3 bucket label should be 'acme-web', got '${data.context_label.s3_bucket.rendered}'"
  }

  assert {
    condition     = data.context_label.path.rendered == "acme/platform/prod/web"
    error_message = "Path label should be 'acme/platform/prod/web', got '${data.context_label.path.rendered}'"
  }

  assert {
    condition     = data.context_label.database.rendered == "acme-platform-prod-postgres"
    error_message = "Database label should be 'acme-platform-prod-postgres', got '${data.context_label.database.rendered}'"
  }

  assert {
    condition     = length(data.context_label.truncated.rendered) <= 20
    error_message = "Truncated label should be <= 20 chars, got ${length(data.context_label.truncated.rendered)}"
  }

  assert {
    condition     = lookup(data.context_tags.default.tags, "Namespace", "") == "acme"
    error_message = "Tags should include Namespace = 'acme'"
  }

  assert {
    condition     = lookup(data.context_tags.snake_case.tags, "namespace", "") == "acme"
    error_message = "Snake case tags should have lowercase key 'namespace'"
  }

  assert {
    condition     = lookup(data.context_tags.with_overrides.tags, "Environment", "") == "staging"
    error_message = "Overridden tags should have Environment = 'staging'"
  }

  assert {
    condition     = lookup(data.context_tags.with_overrides.tags, "Name", "") == "api"
    error_message = "Overridden tags should have Name = 'api'"
  }

  assert {
    condition     = data.context_config.this.delimiter == "-"
    error_message = "Config delimiter should be '-'"
  }

  assert {
    condition     = data.context_config.this.enabled == true
    error_message = "Config should be enabled"
  }
}

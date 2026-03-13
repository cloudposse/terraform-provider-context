run "label_properties_subset_example" {
  command = plan

  assert {
    condition     = data.context_label.full.rendered == "acme-platform-us-east-1-prod-api"
    error_message = "Full label should be 'acme-platform-us-east-1-prod-api', got '${data.context_label.full.rendered}'"
  }

  assert {
    condition     = data.context_label.short.rendered == "acme-api"
    error_message = "Short label should be 'acme-api', got '${data.context_label.short.rendered}'"
  }

  assert {
    condition     = data.context_label.env_scoped.rendered == "acme-us-east-1-prod-api"
    error_message = "Env-scoped label should be 'acme-us-east-1-prod-api', got '${data.context_label.env_scoped.rendered}'"
  }

  assert {
    condition     = data.context_label.tenant_scoped.rendered == "acme-platform-api"
    error_message = "Tenant-scoped label should be 'acme-platform-api', got '${data.context_label.tenant_scoped.rendered}'"
  }
}

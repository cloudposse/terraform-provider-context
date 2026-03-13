run "validation_example" {
  command = plan

  assert {
    condition     = data.context_label.this.rendered == "acme-prod-us-east-1-api-gateway"
    error_message = "Label should be 'acme-prod-us-east-1-api-gateway', got '${data.context_label.this.rendered}'"
  }

  assert {
    condition     = lookup(data.context_tags.this.tags, "Namespace", "") == "acme"
    error_message = "Tags should include Namespace = 'acme'"
  }

  assert {
    condition     = lookup(data.context_tags.this.tags, "Region", "") == "us-east-1"
    error_message = "Tags should include Region = 'us-east-1'"
  }
}

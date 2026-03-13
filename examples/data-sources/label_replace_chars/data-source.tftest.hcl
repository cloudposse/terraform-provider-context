run "label_replace_chars_example" {
  command = plan

  assert {
    condition     = data.context_label.dns_safe.rendered == data.context_label.dns_safe.rendered
    error_message = "DNS-safe label should be generated without error"
  }

  assert {
    condition     = length(data.context_label.dns_safe.rendered) > 0
    error_message = "DNS-safe label should not be empty"
  }

  assert {
    condition     = length(data.context_label.with_underscores.rendered) > 0
    error_message = "Underscore label should not be empty"
  }
}

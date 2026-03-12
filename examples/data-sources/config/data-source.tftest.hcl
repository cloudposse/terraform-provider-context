run "config_example" {
  command = plan

  assert {
    condition     = data.context_config.example.delimiter == "-"
    error_message = "Delimiter should be '-'"
  }

  assert {
    condition     = data.context_config.example.enabled == true
    error_message = "Provider should be enabled"
  }

  assert {
    condition     = lookup(data.context_config.example.values, "Namespace", "") == "cp"
    error_message = "Values should include Namespace = 'cp'"
  }
}

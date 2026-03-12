run "label_templated_example" {
  command = plan

  assert {
    condition     = data.context_label.example.rendered == "cp/plat/dev/example"
    error_message = "Templated label should be 'cp/plat/dev/example', got '${data.context_label.example.rendered}'"
  }
}

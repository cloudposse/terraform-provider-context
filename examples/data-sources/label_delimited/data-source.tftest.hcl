run "label_delimited_example" {
  command = plan

  assert {
    condition     = data.context_label.example.rendered == "cp~plat~dev~ue1~example"
    error_message = "Delimited label should be 'cp~plat~dev~ue1~example', got '${data.context_label.example.rendered}'"
  }
}

run "label_max_length_example" {
  command = plan

  assert {
    condition     = length(data.context_label.truncated.rendered) <= 32
    error_message = "Truncated label should be <= 32 chars, got ${length(data.context_label.truncated.rendered)}"
  }

  assert {
    condition     = length(data.context_label.s3_bucket.rendered) <= 63
    error_message = "S3 bucket label should be <= 63 chars, got ${length(data.context_label.s3_bucket.rendered)}"
  }

  assert {
    condition     = length(data.context_label.s3_bucket.rendered) > length(data.context_label.truncated.rendered)
    error_message = "S3 bucket label should be longer than truncated label"
  }
}

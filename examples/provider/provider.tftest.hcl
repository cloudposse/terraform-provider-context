run "provider_example" {
  command = plan

  assert {
    condition     = lookup(data.context_tags.test.tags, "NAMESPACE", "") == "TESTNAMESPACE"
    error_message = "Expected NAMESPACE=TESTNAMESPACE with upper key/value case"
  }

  assert {
    condition     = lookup(data.context_tags.test.tags, "environment", "") == "testenvironment"
    error_message = "Expected environment=testenvironment with lower key/value case"
  }

  assert {
    condition     = lookup(data.context_tags.test.tags, "stage", "") == "test_stage"
    error_message = "Expected stage=test_stage with snake key/value case"
  }

  assert {
    condition     = lookup(data.context_tags.test.tags, "Name", "") == "TestName"
    error_message = "Expected Name=TestName with title key / none value case"
  }
}

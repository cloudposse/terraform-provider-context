provider "context" {
  delimiter       = "-"
  enabled         = true
  tags_key_case   = "title"
  tags_value_case = "none"

  properties = {
    namespace = {
      required        = true
      include_in_tags = true
      tags_key_case   = "upper"
      tags_value_case = "upper"
    }
    environment = {
      required        = true
      include_in_tags = true
      tags_key_case   = "lower"
      tags_value_case = "lower"
    }
    stage = {
      required        = false
      include_in_tags = true
      tags_key_case   = "snake"
      tags_value_case = "snake"
    }
    name = {
      required        = true
      include_in_tags = true
    }
  }

  values = {
    namespace   = "TestNamespace"
    environment = "TestEnvironment"
    stage       = "TestStage"
    name        = "TestName"
  }
}

data "context_tags" "test" {}

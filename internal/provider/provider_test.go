package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during acceptance testing. The factory function
// will be invoked for every Terraform CLI command executed to create a provider server to which the CLI can reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"context": providerserver.NewProtocol6WithError(NewProvider("test")()),
}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions about the appropriate
	// environment variables being set are common to see in a pre-check function.
}

func TestAccProvider_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
provider "context" {
  delimiter = "-"
  enabled = true
  replace_chars_regex = "/[^a-zA-Z0-9-]/"
  tags_key_case = "title"
  tags_value_case = "none"

  properties = {
    namespace = {
      required = true
      include_in_tags = true
      min_length = 3
      max_length = 10
      validation_regex = "^[a-z0-9-]+$"
    }
    environment = {
      required = true
      include_in_tags = true
      validation_regex = "^(dev|staging|prod)$"
    }
  }

  property_order = ["namespace", "environment"]

  values = {
    namespace = "test"
    environment = "dev"
  }
}

data "context_config" "test" {}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.context_config.test", "enabled", "true"),
					resource.TestCheckResourceAttr("data.context_config.test", "delimiter", "-"),
					resource.TestCheckResourceAttr("data.context_config.test", "replace_chars_regex", "/[^a-zA-Z0-9-]/"),
					resource.TestCheckResourceAttr("data.context_config.test", "tags_key_case", "title"),
					resource.TestCheckResourceAttr("data.context_config.test", "tags_value_case", "none"),
				),
			},
		},
	})
}

func TestAccProvider_invalidConfig(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
provider "context" {
  delimiter = "-"
  enabled = true
  replace_chars_regex = "[" # Invalid regex

  properties = {
    namespace = {
      required = true
      validation_regex = "[" # Invalid regex
    }
  }

  values = {
    namespace = "test"
  }
}

data "context_config" "test" {}`,
				ExpectError: regexp.MustCompile(`(?s)Error: Validation Error.*with provider\["registry.terraform.io/hashicorp/context"\].*regex is invalid: \[ for property namespace`),
			},
		},
	})
}

package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccConfigDataSource(t *testing.T) {
	testAccBasicCfg := getConfigWithProvider(`
	data "context_config" "test" {
	}
	`)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccBasicCfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.context_config.test", "id", "40ca2d12ec2fd2c98beb58579af00d3941b50a87f61165dac0b1bfd08dc47dbc"),
					resource.TestCheckResourceAttr("data.context_config.test", "delimiter", "-"),
					resource.TestCheckResourceAttr("data.context_config.test", "values.Namespace", "cp"),
					resource.TestCheckResourceAttr("data.context_config.test", "property_order.0", "Namespace"),
					resource.TestCheckResourceAttr("data.context_config.test", "tags_key_case", "title"),
					resource.TestCheckResourceAttr("data.context_config.test", "tags_value_case", "none"),
				),
			},
		},
	})
}

func TestAccConfigDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
provider "context" {
  delimiter = "-"
  enabled = true
  replace_chars_regex = "/[^a-zA-Z0-9-]/"

  properties = {
    namespace = {
      required = true
      include_in_tags = true
    }
    environment = {
      required = true
      include_in_tags = true
    }
  }

  property_order = ["namespace", "environment"]

  values = {
    namespace = "test"
    environment = "prod"
  }
}

data "context_config" "test" {}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.context_config.test", "delimiter", "-"),
					resource.TestCheckResourceAttr("data.context_config.test", "enabled", "true"),
					resource.TestCheckResourceAttr("data.context_config.test", "values.namespace", "test"),
					resource.TestCheckResourceAttr("data.context_config.test", "values.environment", "prod"),
				),
			},
		},
	})
}

func TestAccConfigDataSource_validation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
provider "context" {
  delimiter = "-"
  enabled = true

  properties = {
    namespace = {
      required = true
      min_length = 3
      max_length = 10
      validation_regex = "^[a-z0-9-]+$"
    }
  }

  values = {
    namespace = "t"
  }
}

data "context_config" "test" {}`,
				ExpectError: regexp.MustCompile(`(?s)Error running pre-apply plan: exit status 1\s+Error: Validation Error\s+with provider\["registry.terraform.io/hashicorp/context"\],\s+on terraform_plugin_test.tf line 12, in provider "context":\s+12: provider "context" {\s+value is less than minimum length: value t for property namespace is less\s+than 3`),
			},
			{
				Config: `
provider "context" {
  delimiter = "-"
  enabled = true

  properties = {
    namespace = {
      required = true
      validation_regex = "^[a-z0-9-]+$"
    }
  }

  values = {
    namespace = "TEST"
  }
}

data "context_config" "test" {}`,
				ExpectError: regexp.MustCompile(`(?s)Error running pre-apply plan: exit status 1\s+Error: Validation Error\s+with provider\["registry.terraform.io/hashicorp/context"\],\s+on terraform_plugin_test.tf line 12, in provider "context":\s+12: provider "context" {\s+value does not match regex: value TEST for property namespace does not match\s+\^\[a-z0-9-\]\+\$`),
			},
		},
	})
}

func TestAccConfigDataSource_delimitedLabel(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
provider "context" {
  delimiter = "-"
  enabled = true

  properties = {
    namespace = {
      required = true
    }
    environment = {
      required = true
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
					resource.TestCheckResourceAttr("data.context_config.test", "id", "40ca2d12ec2fd2c98beb58579af00d3941b50a87f61165dac0b1bfd08dc47dbc"),
				),
			},
		},
	})
}

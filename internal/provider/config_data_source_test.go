package provider

import (
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

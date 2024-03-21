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

	// testAccBasicTruncatedCfg := getConfigWithProvider(`
	// data "context_config" "test" {
	// 	max_length = 10
	// 	truncate = true
	// }
	// `)

	// testAccLocalDelimiterCfg := getConfigWithProvider(`
	// data "context_config" "test" {
	// 	delimiter = "~"
	// }
	// `)

	// testAccLocalPropertyOrderCfg := getConfigWithProvider(`
	// data "context_config" "test" {
	// 	properties = ["Name", "Namespace", "Tenant", "Stage"]
	// }
	// `)

	// testAccLocalValuesCfg := getConfigWithProvider(`
	// data "context_config" "test" {
	// 	values = {
	// 		"Namespace" = "tst"
	// 		"Name" = "testing"
	// 	}
	// }
	// `)

	// testAccLocalTemplateCfg := getConfigWithProvider(`
	// data "context_config" "test" {
	// 	template = "{{.Namespace}}/{{.Tenant}}/{{.Stage}}/{{.Name}}"
	// }
	// `)

	// testAccLocalTemplateTruncatedCfg := getConfigWithProvider(`
	// data "context_config" "test" {
	// 	template = "{{.Namespace}}/{{.Tenant}}/{{.Stage}}/{{.Name}}"
	// 	max_length = 10
	// 	truncate = true
	// }
	// `)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccBasicCfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.context_config.test", "id", "32a5ac069253a6214ec4e22fa49ec1f634cd3014971dc7a8bb215221885a72e6"),
				),
			},
		},
	})
}

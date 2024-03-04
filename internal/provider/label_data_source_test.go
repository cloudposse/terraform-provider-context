package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccLabelDataSource(t *testing.T) {
	testAccBasicCfg := getConfigWithProvider(`
	data "context_label" "test" {
	}
	`)

	testAccBasicTruncatedCfg := getConfigWithProvider(`
	data "context_label" "test" {
		max_length = 10
		truncate = true
	}
	`)

	testAccLocalDelimiterCfg := getConfigWithProvider(`
	data "context_label" "test" {
		delimiter = "~"
	}
	`)

	testAccLocalPropertyOrderCfg := getConfigWithProvider(`
	data "context_label" "test" {
		properties = ["Name", "Namespace", "Tenant", "Stage"]
	}
	`)

	testAccLocalValuesCfg := getConfigWithProvider(`
	data "context_label" "test" {
		values = {
			"Namespace" = "tst"
			"Name" = "testing"
		}
	}
	`)

	testAccLocalTemplateCfg := getConfigWithProvider(`
	data "context_label" "test" {
		template = "{{.Namespace}}/{{.Tenant}}/{{.Stage}}/{{.Name}}"
	}
	`)

	testAccLocalTemplateTruncatedCfg := getConfigWithProvider(`
	data "context_label" "test" {
		template = "{{.Namespace}}/{{.Tenant}}/{{.Stage}}/{{.Name}}"
		max_length = 10
		truncate = true
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
					resource.TestCheckResourceAttr("data.context_label.test", "rendered", "cp-core-prod-example"),
				),
			},
			{
				Config: testAccBasicTruncatedCfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.context_label.test", "rendered", "cp-co16916"),
				),
			},
			{
				Config: testAccLocalDelimiterCfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.context_label.test", "rendered", "cp~core~prod~example"),
				),
			},
			{
				Config: testAccLocalPropertyOrderCfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.context_label.test", "rendered", "example-cp-core-prod"),
				),
			},
			{
				Config: testAccLocalValuesCfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.context_label.test", "rendered", "tst-core-prod-testing"),
				),
			},
			{
				Config: testAccLocalTemplateCfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.context_label.test", "rendered", "cp/core/prod/example"),
				),
			},
			{
				Config: testAccLocalTemplateTruncatedCfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.context_label.test", "rendered", "cp/co46903"),
				),
			},
		},
	})
}

func getConfigWithProvider(data string) string {
	return fmt.Sprintf(`
	provider "context" {
		properties = {
			Namespace = {required = true, max_length=10, validation_regex="^[a-z0-9]+$"}
			Tenant    = {}
			Stage     = {}
			Name      = {}
		}

		property_order = ["Namespace", "Tenant", "Stage", "Name"]

		values = {
			"Namespace" = "cp"
			"Tenant" = "core"
			"Stage" = "prod"
			"Name"  = "example"
		}
	}


	%s`, data)
}

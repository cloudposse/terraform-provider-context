package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccLabelDataSource(t *testing.T) {
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
		},
	})
}

const testAccBasicCfg = `
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


data "context_label" "test" {
}
`

const testAccLocalDelimiterCfg = `
provider "context" {
  properties = {
    Namespace = {}
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


data "context_label" "test" {
	delimiter = "~"
}
`

const testAccLocalPropertyOrderCfg = `
provider "context" {
  properties = {
    Namespace = {}
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


data "context_label" "test" {
	properties = ["Name", "Namespace", "Tenant", "Stage"]
}
`
const testAccLocalValuesCfg = `
provider "context" {
  properties = {
    Namespace = {}
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

data "context_label" "test" {
	values = {
		"Namespace" = "tst"
		"Name" = "testing"
	}
}
`

const testAccLocalTemplateCfg = `
provider "context" {
  properties = {
    Namespace = {}
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


data "context_label" "test" {
	template = "{{.Namespace}}/{{.Tenant}}/{{.Stage}}/{{.Name}}"
}
`

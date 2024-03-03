package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTagsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccTagsBasicCfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.context_tags.test", "tags.Namespace", "cp"),
					resource.TestCheckResourceAttr("data.context_tags.test", "tags.Tenant", "core"),
					resource.TestCheckResourceAttr("data.context_tags.test", "tags.Stage", "prod"),
					resource.TestCheckResourceAttr("data.context_tags.test", "tags.Name", "example"),
				),
			},
			{
				Config: testAccTagsBasicCfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.context_tags.test", "tags_as_list.0.Key", "Name"),
					resource.TestCheckResourceAttr("data.context_tags.test", "tags_as_list.1.Key", "Namespace"),
					resource.TestCheckResourceAttr("data.context_tags.test", "tags_as_list.2.Key", "Stage"),
					resource.TestCheckResourceAttr("data.context_tags.test", "tags_as_list.3.Key", "Tenant"),
					resource.TestCheckResourceAttr("data.context_tags.test", "tags_as_list.0.Value", "example"),
					resource.TestCheckResourceAttr("data.context_tags.test", "tags_as_list.1.Value", "cp"),
					resource.TestCheckResourceAttr("data.context_tags.test", "tags_as_list.2.Value", "prod"),
					resource.TestCheckResourceAttr("data.context_tags.test", "tags_as_list.3.Value", "core"),
				),
			},
			{
				Config: testAccTagsExcludedCfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr("data.context_tags.test", "tags.Namespace"),
				),
			},
		},
	})
}

const testAccTagsBasicCfg = `
provider "context" {
  properties = {
    Namespace = {}
    Tenant    = {}
    Stage     = {}
    Name      = {}
  }

  values = {
    "Namespace" = "cp"
    "Tenant" = "core"
    "Stage" = "prod"
    "Name"  = "example"
  }
}


data "context_tags" "test" {
}
`

const testAccTagsExcludedCfg = `
provider "context" {
  properties = {
    Namespace = {include_in_tags = false}
    Tenant    = {}
    Stage     = {}
    Name      = {}
  }

  values = {
    "Namespace" = "cp"
    "Tenant" = "core"
    "Stage" = "prod"
    "Name"  = "example"
  }
}


data "context_tags" "test" {
}
`

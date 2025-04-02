package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGcpLabelsDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTagsReplaceCfg,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.context_gcp_labels.test", "tags_as_list.0.Key", "ComponentPath"),
					resource.TestCheckResourceAttr("data.context_gcp_labels.test", "tags_as_list.1.Key", "Name"),
					resource.TestCheckResourceAttr("data.context_gcp_labels.test", "tags_as_list.2.Key", "Namespace"),
					resource.TestCheckResourceAttr("data.context_gcp_labels.test", "tags_as_list.3.Key", "Stage"),
					resource.TestCheckResourceAttr("data.context_gcp_labels.test", "tags_as_list.4.Key", "Tenant"),
					resource.TestCheckResourceAttr("data.context_gcp_labels.test", "tags_as_list.0.Value", "foo-bar-baz"),
					resource.TestCheckResourceAttr("data.context_gcp_labels.test", "tags_as_list.1.Value", "example"),
					resource.TestCheckResourceAttr("data.context_gcp_labels.test", "tags_as_list.2.Value", "cp"),
					resource.TestCheckResourceAttr("data.context_gcp_labels.test", "tags_as_list.3.Value", "prod"),
					resource.TestCheckResourceAttr("data.context_gcp_labels.test", "tags_as_list.4.Value", "core"),
				),
			},
		},
	})
}

const testAccTagsReplaceCfg = `
provider "context" {
  properties = {
    Namespace = {}
    Tenant    = {}
    Stage     = {}
    Name      = {}
    ComponentPath      = {}
  }

  values = {
    "Namespace" = "cloudposse"
    "Tenant" = "core"
    "Stage" = "prod"
    "Name"  = "example"
	"ComponentPath" = "foo/bar/baz"
  }
}


data "context_gcp_labels" "test" {
  replacement_map = {
	"cloudposse" = "cp"
	"/" = "-"
  }
}
`

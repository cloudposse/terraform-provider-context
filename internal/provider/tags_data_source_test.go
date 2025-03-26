package provider

import (
	"regexp"
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
			{
				Config: testAccTagsEmptyValueCfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr("data.context_tags.test", "tags.Name"),
				),
			},
			{
				Config: testAccTagKeyCasedCfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.context_tags.test", "tags.NAMESPACE", "cp"),
					resource.TestCheckResourceAttr("data.context_tags.test", "tags.TENANT", "core"),
					resource.TestCheckResourceAttr("data.context_tags.test", "tags.STAGE", "prod"),
					resource.TestCheckResourceAttr("data.context_tags.test", "tags.NAME", "example"),
				),
			},
		},
	})
}

func TestAccTagsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
provider "context" {
  delimiter = "-"
  enabled = true
  tags_key_case = "title"
  tags_value_case = "none"

  properties = {
    namespace = {
      required = true
      include_in_tags = true
    }
    environment = {
      required = true
      include_in_tags = true
    }
    stage = {
      required = false
      include_in_tags = false
    }
  }

  values = {
    namespace = "test"
    environment = "dev"
    stage = "beta"
  }
}

data "context_tags" "test" {}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.context_tags.test", "tags.Namespace", "test"),
					resource.TestCheckResourceAttr("data.context_tags.test", "tags.Environment", "dev"),
					resource.TestCheckNoResourceAttr("data.context_tags.test", "tags.Stage"),
					resource.TestCheckResourceAttr("data.context_tags.test", "tags_as_list.#", "2"),
				),
			},
		},
	})
}

func TestAccTagsDataSource_casing(t *testing.T) {
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
      include_in_tags = true
    }
  }

  values = {
    namespace = "TestValue"
  }
}

data "context_tags" "lower" {
  tags_key_case = "lower"
  tags_value_case = "lower"
}

data "context_tags" "upper" {
  tags_key_case = "upper"
  tags_value_case = "upper"
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.context_tags.lower", "tags.namespace", "testvalue"),
					resource.TestCheckResourceAttr("data.context_tags.upper", "tags.NAMESPACE", "TESTVALUE"),
				),
			},
		},
	})
}

func TestAccTagsDataSource_values(t *testing.T) {
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
      include_in_tags = true
    }
    environment = {
      required = true
      include_in_tags = true
    }
  }

  values = {
    namespace = "test"
    environment = "dev"
  }
}

data "context_tags" "test" {
  values = {
    environment = "prod"
  }
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.context_tags.test", "tags.Namespace", "test"),
					resource.TestCheckResourceAttr("data.context_tags.test", "tags.Environment", "prod"),
				),
			},
		},
	})
}

func TestAccTagsDataSource_validation(t *testing.T) {
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
      include_in_tags = true
      validation_regex = "^[a-z]+$"
    }
  }

  values = {
    namespace = "TEST"
  }
}

data "context_tags" "test" {}`,
				ExpectError: regexp.MustCompile(`(?s)Error running pre-apply plan: exit status 1\s+Error: Validation Error\s+with provider\["registry.terraform.io/hashicorp/context"\],\s+on terraform_plugin_test.tf line 12, in provider "context":\s+12: provider "context" {\s+value does not match regex: value TEST for property namespace does not match\s+\^\[a-z\]\+\$`),
			},
		},
	})
}

func TestAccTagsDataSource_propertySpecificCases(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
provider "context" {
  delimiter = "-"
  enabled = true
  tags_key_case = "title"
  tags_value_case = "none"

  properties = {
    namespace = {
      required = true
      include_in_tags = true
      tags_key_case = "upper"
      tags_value_case = "upper"
    }
    environment = {
      required = true
      include_in_tags = true
      tags_key_case = "lower"
      tags_value_case = "lower"
    }
    stage = {
      required = false
      include_in_tags = true
      tags_key_case = "snake"
      tags_value_case = "snake"
    }
  }

  values = {
    namespace = "TestNamespace"
    environment = "TestEnvironment"
    stage = "TestStage"
  }
}

data "context_tags" "test" {}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.context_tags.test", "tags.NAMESPACE", "TESTNAMESPACE"),
					resource.TestCheckResourceAttr("data.context_tags.test", "tags.environment", "testenvironment"),
					resource.TestCheckResourceAttr("data.context_tags.test", "tags.stage", "test_stage"),
					resource.TestCheckResourceAttr("data.context_tags.test", "tags_as_list.#", "3"),
				),
			},
		},
	})
}

func TestAccTagsDataSource_mixedCases(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
provider "context" {
  delimiter = "-"
  enabled = true
  tags_key_case = "title"
  tags_value_case = "none"

  properties = {
    namespace = {
      required = true
      include_in_tags = true
      tags_key_case = "upper"
    }
    environment = {
      required = true
      include_in_tags = true
      tags_value_case = "upper"
    }
    stage = {
      required = false
      include_in_tags = true
    }
  }

  values = {
    namespace = "TestNamespace"
    environment = "TestEnvironment"
    stage = "TestStage"
  }
}

data "context_tags" "test" {}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.context_tags.test", "tags.NAMESPACE", "TestNamespace"),
					resource.TestCheckResourceAttr("data.context_tags.test", "tags.Environment", "TESTENVIRONMENT"),
					resource.TestCheckResourceAttr("data.context_tags.test", "tags.Stage", "TestStage"),
					resource.TestCheckResourceAttr("data.context_tags.test", "tags_as_list.#", "3"),
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

const testAccTagsEmptyValueCfg = `
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
    "Name"  = ""
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

const testAccTagKeyCasedCfg = `
provider "context" {
  properties = {
    Namespace = {}
    Tenant    = {}
    Stage     = {}
    Name      = {}
  }

	tags_key_case = "upper"

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

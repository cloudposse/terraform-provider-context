package provider

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func getPropertiesSchema() schema.NestedAttributeObject {
	return schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"include_in_tags": schema.BoolAttribute{
				MarkdownDescription: "A flag to indicate if the property should be included in tags. If not set, defaults to true.",
				Optional:            true,
			},
			"max_length": schema.Int64Attribute{
				MarkdownDescription: "The maximum length of the property.",
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"min_length": schema.Int64Attribute{
				MarkdownDescription: "The minimum length of the property.",
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"required": schema.BoolAttribute{
				MarkdownDescription: "A flag to indicate if the property is required.",
				Optional:            true,
			},
			"validation_regex": schema.StringAttribute{
				MarkdownDescription: "A regular expression to validate the property.",
				Optional:            true,
			},
			"tags_key_case": schema.StringAttribute{
				MarkdownDescription: "The case to use for the key of this property in tags. If not set, uses the provider's tags_key_case.",
				Optional:            true,
				Validators:          []validator.String{stringvalidator.OneOf("none", "camel", "lower", "snake", "title", "upper")},
			},
			"tags_value_case": schema.StringAttribute{
				MarkdownDescription: "The case to use for the value of this property in tags. If not set, uses the provider's tags_value_case.",
				Optional:            true,
				Validators:          []validator.String{stringvalidator.OneOf("none", "camel", "lower", "snake", "title", "upper")},
			},
		},
	}
}

func getPropertiesDSSchema() dsschema.NestedAttributeObject {
	return dsschema.NestedAttributeObject{
		Attributes: map[string]dsschema.Attribute{
			"include_in_tags": dsschema.BoolAttribute{
				MarkdownDescription: "A flag to indicate if the property should be included in tags.",
				Optional:            true,
			},
			"max_length": dsschema.Int64Attribute{
				MarkdownDescription: "The maximum length of the property.",
				Optional:            true,
			},
			"min_length": dsschema.Int64Attribute{
				MarkdownDescription: "The minimum length of the property.",
				Optional:            true,
			},
			"required": dsschema.BoolAttribute{
				MarkdownDescription: "A flag to indicate if the property is required.",
				Optional:            true,
			},
			"validation_regex": dsschema.StringAttribute{
				MarkdownDescription: "A regular expression to validate the property.",
				Optional:            true,
			},
			"tags_key_case": dsschema.StringAttribute{
				MarkdownDescription: "The case to use for the key of this property in tags. If not set, uses the provider's tags_key_case.",
				Optional:            true,
			},
			"tags_value_case": dsschema.StringAttribute{
				MarkdownDescription: "The case to use for the value of this property in tags. If not set, uses the provider's tags_value_case.",
				Optional:            true,
			},
		},
	}
}

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
			"tags_key_case": schema.StringAttribute{
				MarkdownDescription: "The case to use for the key of this property in tags. If not set, uses the provider's tags_key_case setting. Valid values are: none, camel, lower, snake, title, upper.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(ValidCases...),
				},
			},
			"tags_value_case": schema.StringAttribute{
				MarkdownDescription: "The case to use for the value of this property in tags. If not set, uses the provider's tags_value_case setting. Valid values are: none, camel, lower, snake, title, upper.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(ValidCases...),
				},
			},
			"validation_regex": schema.StringAttribute{
				MarkdownDescription: "A regular expression to validate the property.",
				Optional:            true,
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
			"tags_key_case": dsschema.StringAttribute{
				MarkdownDescription: "The case to use for the key of this property in tags.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(ValidCases...),
				},
			},
			"tags_value_case": dsschema.StringAttribute{
				MarkdownDescription: "The case to use for the value of this property in tags.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(ValidCases...),
				},
			},
			"validation_regex": dsschema.StringAttribute{
				MarkdownDescription: "A regular expression to validate the property.",
				Optional:            true,
			},
		},
	}
}

func GetProviderSchema() map[string]dsschema.Attribute {
	return map[string]dsschema.Attribute{
		"delimiter": dsschema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The delimiter to use between label elements.",
		},
		"enabled": dsschema.BoolAttribute{
			Optional:            true,
			MarkdownDescription: "Set to false to prevent the module from creating any resources.",
		},
		"tags_key_case": dsschema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The case to use for the keys of tags created by the provider. Valid values are: none, camel, lower, snake, title, upper.",
			Validators: []validator.String{
				stringvalidator.OneOf(ValidCases...),
			},
		},
		"tags_value_case": dsschema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The case to use for the values of tags created by the provider. Valid values are: none, camel, lower, snake, title, upper.",
			Validators: []validator.String{
				stringvalidator.OneOf(ValidCases...),
			},
		},
		// ... rest of the schema ...
	}
}

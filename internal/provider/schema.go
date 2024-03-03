package provider

import (
	"context"

	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func getPropertiesSchema() schema.NestedAttributeObject {
	return schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"include_in_tags": schema.BoolAttribute{
				MarkdownDescription: "A flag to indicate if the property should be included in tags",
				Optional:            true,
			},
			"max_length": schema.Int64Attribute{
				MarkdownDescription: "The maximum length of the property",
				Optional:            true,
			},
			"min_length": schema.Int64Attribute{
				MarkdownDescription: "The minimum length of the property",
				Optional:            true,
			},
			"required": schema.BoolAttribute{
				MarkdownDescription: "A flag to indicate if the property is required",
				Optional:            true,
			},
			"validation_regex": schema.StringAttribute{
				MarkdownDescription: "A regular expression to validate the property",
				Optional:            true,
			},
		},
	}
}

func getPropertiesDSSchema() dsschema.NestedAttributeObject {
	return dsschema.NestedAttributeObject{
		Attributes: map[string]dsschema.Attribute{
			"include_in_tags": dsschema.BoolAttribute{
				MarkdownDescription: "A flag to indicate if the property should be included in tags",
				Optional:            true,
			},
			"max_length": dsschema.Int64Attribute{
				MarkdownDescription: "The maximum length of the property",
				Optional:            true,
			},
			"min_length": dsschema.Int64Attribute{
				MarkdownDescription: "The minimum length of the property",
				Optional:            true,
			},
			"required": dsschema.BoolAttribute{
				MarkdownDescription: "A flag to indicate if the property is required",
				Optional:            true,
			},
			"validation_regex": dsschema.StringAttribute{
				MarkdownDescription: "A regular expression to validate the property",
				Optional:            true,
			},
		},
	}
}

// FromFrameworkMap converts a types.Map to a map[string]T.
func FromFrameworkMap[T interface{}](ctx context.Context, m types.Map) (map[string]T, diag.Diagnostics) {
	localValues := make(map[string]T, len(m.Elements()))
	if !m.IsNull() {
		diag := m.ElementsAs(ctx, &localValues, false)

		if diag.HasError() {
			return nil, diag
		}
	}
	return localValues, nil
}

// FromFrameworkList converts a types.List to a []T.
func FromFrameworkList[T interface{}](ctx context.Context, m types.List) ([]T, diag.Diagnostics) {
	localValues := make([]T, len(m.Elements()))
	if !m.IsNull() {
		diag := m.ElementsAs(ctx, &localValues, false)

		if diag.HasError() {
			return nil, diag
		}
	}
	return localValues, nil
}

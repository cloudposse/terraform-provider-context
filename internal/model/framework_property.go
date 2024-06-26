package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type FrameworkProperty struct {
	IncludeInTags   types.Bool   `tfsdk:"include_in_tags"`
	MaxLength       types.Int64  `tfsdk:"max_length"`
	MinLength       types.Int64  `tfsdk:"min_length"`
	Required        types.Bool   `tfsdk:"required"`
	ValidationRegex types.String `tfsdk:"validation_regex"`
}

func (p *FrameworkProperty) ToModel(name string) (*Property, error) {
	options := []func(*Property){}

	if !p.Required.IsNull() && !p.Required.IsUnknown() && p.Required.ValueBool() {
		options = append(options, WithRequired())
	}

	if !p.IncludeInTags.IsNull() && !p.IncludeInTags.IsUnknown() && !p.IncludeInTags.ValueBool() {
		options = append(options, WithExcludeFromTags())
	}

	if !p.MinLength.IsNull() && !p.MinLength.IsUnknown() {
		options = append(options, WithMinLength(int(p.MinLength.ValueInt64())))
	}

	if !p.MaxLength.IsNull() && !p.MaxLength.IsUnknown() {
		options = append(options, WithMaxLength(int(p.MaxLength.ValueInt64())))
	}

	if !p.ValidationRegex.IsNull() && !p.ValidationRegex.IsUnknown() {
		options = append(options, WithValidationRegex(p.ValidationRegex.ValueString()))
	}

	return NewProperty(name, options...), nil
}

func (p FrameworkProperty) Types() map[string]attr.Type {
	return map[string]attr.Type{
		"include_in_tags":  types.BoolType,
		"max_length":       types.Int64Type,
		"min_length":       types.Int64Type,
		"required":         types.BoolType,
		"validation_regex": types.StringType,
	}
}

func (p FrameworkProperty) FromConfigProperty(cp Property) FrameworkProperty {
	return FrameworkProperty{
		IncludeInTags:   types.BoolValue(cp.IncludeInTags),
		MaxLength:       types.Int64Value(int64(cp.MaxLength)),
		MinLength:       types.Int64Value(int64(cp.MinLength)),
		Required:        types.BoolValue(cp.Required),
		ValidationRegex: types.StringValue(cp.ValidationRegex),
	}
}

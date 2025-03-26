package model

import (
	"fmt"

	"github.com/cloudposse/terraform-provider-context/pkg/cases"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type FrameworkProperty struct {
	IncludeInTags   types.Bool   `tfsdk:"include_in_tags"`
	MaxLength       types.Int64  `tfsdk:"max_length"`
	MinLength       types.Int64  `tfsdk:"min_length"`
	Required        types.Bool   `tfsdk:"required"`
	ValidationRegex types.String `tfsdk:"validation_regex"`
	TagsKeyCase     types.String `tfsdk:"tags_key_case"`
	TagsValueCase   types.String `tfsdk:"tags_value_case"`
}

func (p *FrameworkProperty) addRequiredOption(options []func(*Property)) []func(*Property) {
	if !p.Required.IsNull() && !p.Required.IsUnknown() && p.Required.ValueBool() {
		return append(options, WithRequired())
	}
	return options
}

func (p *FrameworkProperty) addIncludeInTagsOption(options []func(*Property)) []func(*Property) {
	if !p.IncludeInTags.IsNull() && !p.IncludeInTags.IsUnknown() && !p.IncludeInTags.ValueBool() {
		return append(options, WithExcludeFromTags())
	}
	return options
}

func (p *FrameworkProperty) addMinLengthOption(options []func(*Property)) []func(*Property) {
	if !p.MinLength.IsNull() && !p.MinLength.IsUnknown() {
		return append(options, WithMinLength(int(p.MinLength.ValueInt64())))
	}
	return options
}

func (p *FrameworkProperty) addMaxLengthOption(options []func(*Property)) []func(*Property) {
	if !p.MaxLength.IsNull() && !p.MaxLength.IsUnknown() {
		return append(options, WithMaxLength(int(p.MaxLength.ValueInt64())))
	}
	return options
}

func (p *FrameworkProperty) addValidationRegexOption(options []func(*Property)) []func(*Property) {
	if !p.ValidationRegex.IsNull() && !p.ValidationRegex.IsUnknown() {
		return append(options, WithValidationRegex(p.ValidationRegex.ValueString()))
	}
	return options
}

func (p *FrameworkProperty) ToModel(name string) (*Property, error) {
	options := []func(*Property){}

	options = p.addRequiredOption(options)
	options = p.addIncludeInTagsOption(options)
	options = p.addMinLengthOption(options)
	options = p.addMaxLengthOption(options)
	options = p.addValidationRegexOption(options)

	if !p.TagsKeyCase.IsNull() && !p.TagsKeyCase.IsUnknown() {
		keyCase, err := cases.FromString(p.TagsKeyCase.ValueString())
		if err != nil {
			return nil, fmt.Errorf("invalid tags_key_case for property %s: %w", name, err)
		}
		options = append(options, WithPropertyTagsKeyCase(keyCase))
	}

	if !p.TagsValueCase.IsNull() && !p.TagsValueCase.IsUnknown() {
		valueCase, err := cases.FromString(p.TagsValueCase.ValueString())
		if err != nil {
			return nil, fmt.Errorf("invalid tags_value_case for property %s: %w", name, err)
		}
		options = append(options, WithPropertyTagsValueCase(valueCase))
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
		"tags_key_case":    types.StringType,
		"tags_value_case":  types.StringType,
	}
}

func (p FrameworkProperty) FromConfigProperty(cp Property) FrameworkProperty {
	var tagsKeyCase, tagsValueCase types.String
	if cp.TagsKeyCase != nil {
		tagsKeyCase = types.StringValue(cp.TagsKeyCase.String())
	} else {
		tagsKeyCase = types.StringNull()
	}
	if cp.TagsValueCase != nil {
		tagsValueCase = types.StringValue(cp.TagsValueCase.String())
	} else {
		tagsValueCase = types.StringNull()
	}

	return FrameworkProperty{
		IncludeInTags:   types.BoolValue(cp.IncludeInTags),
		MaxLength:       types.Int64Value(int64(cp.MaxLength)),
		MinLength:       types.Int64Value(int64(cp.MinLength)),
		Required:        types.BoolValue(cp.Required),
		ValidationRegex: types.StringValue(cp.ValidationRegex),
		TagsKeyCase:     tagsKeyCase,
		TagsValueCase:   tagsValueCase,
	}
}

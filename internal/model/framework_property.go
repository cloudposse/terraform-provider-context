package model

import (
	"github.com/cloudposse/terraform-provider-context/pkg/cases"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type FrameworkProperty struct {
	IncludeInTags   types.Bool   `tfsdk:"include_in_tags"`
	MaxLength       types.Int64  `tfsdk:"max_length"`
	MinLength       types.Int64  `tfsdk:"min_length"`
	Required        types.Bool   `tfsdk:"required"`
	TagsKeyCase     types.String `tfsdk:"tags_key_case"`
	TagsValueCase   types.String `tfsdk:"tags_value_case"`
	ValidationRegex types.String `tfsdk:"validation_regex"`
}

func (p *FrameworkProperty) addRequiredOption(options []PropertyOption) []PropertyOption {
	if !p.Required.IsNull() && !p.Required.IsUnknown() && p.Required.ValueBool() {
		return append(options, WithRequired())
	}
	return options
}

func (p *FrameworkProperty) addIncludeInTagsOption(options []PropertyOption) []PropertyOption {
	if !p.IncludeInTags.IsNull() && !p.IncludeInTags.IsUnknown() && !p.IncludeInTags.ValueBool() {
		return append(options, WithExcludeFromTags())
	}
	return options
}

func (p *FrameworkProperty) addMinLengthOption(options []PropertyOption) []PropertyOption {
	if !p.MinLength.IsNull() && !p.MinLength.IsUnknown() {
		return append(options, WithMinLength(int(p.MinLength.ValueInt64())))
	}
	return options
}

func (p *FrameworkProperty) addMaxLengthOption(options []PropertyOption) []PropertyOption {
	if !p.MaxLength.IsNull() && !p.MaxLength.IsUnknown() {
		return append(options, WithMaxLength(int(p.MaxLength.ValueInt64())))
	}
	return options
}

func (p *FrameworkProperty) addValidationRegexOption(options []PropertyOption) []PropertyOption {
	if !p.ValidationRegex.IsNull() && !p.ValidationRegex.IsUnknown() {
		return append(options, WithValidationRegex(p.ValidationRegex.ValueString()))
	}
	return options
}

func (p *FrameworkProperty) addTagsKeyCaseOption(options []PropertyOption) []PropertyOption {
	if !p.TagsKeyCase.IsNull() && !p.TagsKeyCase.IsUnknown() {
		if caseType, err := cases.FromString(p.TagsKeyCase.ValueString()); err == nil {
			options = append(options, WithPropertyTagsKeyCase(caseType))
		}
	}
	return options
}

func (p *FrameworkProperty) addTagsValueCaseOption(options []PropertyOption) []PropertyOption {
	if !p.TagsValueCase.IsNull() && !p.TagsValueCase.IsUnknown() {
		if caseType, err := cases.FromString(p.TagsValueCase.ValueString()); err == nil {
			options = append(options, WithPropertyTagsValueCase(caseType))
		}
	}
	return options
}

func (p *FrameworkProperty) ToModel(name string) (*Property, error) {
	options := []PropertyOption{}

	options = p.addRequiredOption(options)
	options = p.addIncludeInTagsOption(options)
	options = p.addMinLengthOption(options)
	options = p.addMaxLengthOption(options)
	options = p.addValidationRegexOption(options)
	options = p.addTagsKeyCaseOption(options)
	options = p.addTagsValueCaseOption(options)

	return NewProperty(name, options...), nil
}

func (p *FrameworkProperty) Types() map[string]attr.Type {
	return map[string]attr.Type{
		"include_in_tags":  types.BoolType,
		"max_length":       types.Int64Type,
		"min_length":       types.Int64Type,
		"required":         types.BoolType,
		"tags_key_case":    types.StringType,
		"tags_value_case":  types.StringType,
		"validation_regex": types.StringType,
	}
}

func (p *FrameworkProperty) FromConfigProperty(cp *Property) FrameworkProperty {
	fp := FrameworkProperty{
		IncludeInTags:   types.BoolValue(cp.IncludeInTags),
		MaxLength:       types.Int64Value(int64(cp.MaxLength)),
		MinLength:       types.Int64Value(int64(cp.MinLength)),
		Required:        types.BoolValue(cp.Required),
		ValidationRegex: types.StringValue(cp.ValidationRegex),
	}
	if cp.TagsKeyCase != nil {
		fp.TagsKeyCase = types.StringValue(cp.TagsKeyCase.String())
	}
	if cp.TagsValueCase != nil {
		fp.TagsValueCase = types.StringValue(cp.TagsValueCase.String())
	}
	return fp
}

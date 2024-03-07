package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type delimitedLabelModel struct {
	Delimiter         *string
	MaxLength         int64
	PropertyNames     []string
	ReplaceCharsRegex *string
	Truncate          bool
	Values            map[string]string
}

func (m delimitedLabelModel) FromFramework(ctx context.Context, config LabelDataSourceModel) (delimitedLabelModel, diag.Diagnostics) {
	model := delimitedLabelModel{}

	var delimiter *string
	if !config.Delimiter.IsNull() {
		delimiter = config.Delimiter.ValueStringPointer()
	}
	model.Delimiter = delimiter

	if config.MaxLength.IsNull() {
		model.MaxLength = 0
	} else {
		model.MaxLength = config.MaxLength.ValueInt64()
	}

	properties, diags := FromFrameworkList[string](ctx, config.Properties)
	if diags.HasError() {
		return model, diags
	}
	model.PropertyNames = properties

	var replaceCharsRegex *string
	if !config.ReplaceCharsRegex.IsNull() {
		replaceCharsRegex = config.ReplaceCharsRegex.ValueStringPointer()
	}
	model.ReplaceCharsRegex = replaceCharsRegex

	if !config.Truncate.IsNull() {
		model.Truncate = config.Truncate.ValueBool()
	} else {
		model.Truncate = true
	}

	values, diags := FromFrameworkMap[string](ctx, config.Values)
	if diags.HasError() {
		return model, diags
	}
	model.Values = values

	return model, nil
}

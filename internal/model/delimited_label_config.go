package model

import (
	"context"

	"github.com/cloudposse/terraform-provider-context/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type DelimitedLabelModel struct {
	Delimiter         *string
	MaxLength         int64
	PropertyNames     []string
	ReplaceCharsRegex *string
	Truncate          bool
	Values            map[string]string
}

func (m DelimitedLabelModel) FromFramework(ctx context.Context, config DataSourceLabelConfig) (DelimitedLabelModel, diag.Diagnostics) {
	model := DelimitedLabelModel{}

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

	properties, diags := framework.FromFrameworkList[string](ctx, config.Properties)
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

	values, diags := framework.FromFrameworkMap[string](ctx, config.Values)
	if diags.HasError() {
		return model, diags
	}
	model.Values = values

	return model, nil
}

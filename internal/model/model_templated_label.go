package model

import (
	"context"

	"github.com/cloudposse/terraform-provider-context/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type TemplatedLabelModel struct {
	MaxLength         int64
	Template          string
	Truncate          bool
	ReplaceCharsRegex *string
	Values            map[string]string
}

func (m TemplatedLabelModel) FromFramework(ctx context.Context, config DataSourceLabelConfig) (TemplatedLabelModel, diag.Diagnostics) {
	model := TemplatedLabelModel{}

	if config.MaxLength.IsNull() {
		model.MaxLength = 0
	} else {
		model.MaxLength = config.MaxLength.ValueInt64()
	}

	if !config.Template.IsNull() {
		model.Template = config.Template.ValueString()
	}

	if !config.Truncate.IsNull() {
		model.Truncate = config.Truncate.ValueBool()
	} else {
		model.Truncate = true
	}

	var replaceCharsRegex *string
	if !config.ReplaceCharsRegex.IsNull() {
		replaceCharsRegex = config.ReplaceCharsRegex.ValueStringPointer()
	}
	model.ReplaceCharsRegex = replaceCharsRegex

	values, diags := framework.FromFrameworkMap[string](ctx, config.Values)
	if diags.HasError() {
		return model, diags
	}
	model.Values = values

	return model, nil
}

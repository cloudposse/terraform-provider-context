package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type templatedLabelModel struct {
	MaxLength int64
	Template  string
	Truncate  bool
	Values    map[string]string
}

func (m templatedLabelModel) FromFramework(ctx context.Context, config LabelDataSourceModel) (templatedLabelModel, diag.Diagnostics) {
	model := templatedLabelModel{}

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

	values, diags := FromFrameworkMap[string](ctx, config.Values)
	if diags.HasError() {
		return model, diags
	}
	model.Values = values

	return model, nil
}

package model

import "github.com/hashicorp/terraform-plugin-framework/types"

// DataSourceLabelConfig describes the label data source data model.
type DataSourceLabelConfig struct {
	Delimiter         types.String `tfsdk:"delimiter"`
	Id                types.String `tfsdk:"id"`
	MaxLength         types.Int64  `tfsdk:"max_length"`
	Properties        types.List   `tfsdk:"properties"`
	Rendered          types.String `tfsdk:"rendered"`
	ReplaceCharsRegex types.String `tfsdk:"replace_chars_regex"`
	Template          types.String `tfsdk:"template"`
	Truncate          types.Bool   `tfsdk:"truncate"`
	Values            types.Map    `tfsdk:"values"`
}

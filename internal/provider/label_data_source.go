package provider

import (
	"context"
	"fmt"

	"github.com/cloudposse/terraform-provider-context/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource                     = &LabelDataSource{}
	_ datasource.DataSourceWithConfigure        = &LabelDataSource{}
	_ datasource.DataSourceWithConfigValidators = &LabelDataSource{}
)

func NewLabelDataSource() datasource.DataSource {
	return &LabelDataSource{}
}

// LabelDataSource defines the data source implementation.
type LabelDataSource struct {
	providerData *providerData
}

// LabelDataSourceModel describes the data source data model.
type LabelDataSourceModel struct {
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

func (d *LabelDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_label"
}

func (d *LabelDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Label data source",

		Attributes: map[string]schema.Attribute{
			"delimiter": schema.StringAttribute{
				MarkdownDescription: "Delimiter to use when creating the label from properties. Conflicts with `template`.",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Label identifier",
				Computed:            true,
			},
			"max_length": schema.Int64Attribute{
				MarkdownDescription: "Maximum length of the label",
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"properties": schema.ListAttribute{
				MarkdownDescription: "List of properties to use when creating the label. Conflicts with `template`.",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"rendered": schema.StringAttribute{
				MarkdownDescription: "Rendered label",
				Computed:            true,
			},
			"replace_chars_regex": schema.StringAttribute{
				MarkdownDescription: "The regex to use for replacing characters in labels created by the provider. Any characters that match the regex will be removed from the label.",
				Optional:            true,
			},
			"template": schema.StringAttribute{
				MarkdownDescription: "Template to use when creating the label. Conflicts with `delimiter` and `properties`.",
				Optional:            true,
			},
			"truncate": schema.BoolAttribute{
				MarkdownDescription: "Truncate the label if it exceeds the maximum length. If false, an error will be returned if the label exceeds the maximum length.",
				Optional:            true,
			},
			"values": schema.MapAttribute{
				MarkdownDescription: "Map of values to override or add to the context when creating the label.",
				Optional:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (d *LabelDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	providerData, ok := req.ProviderData.(*providerData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *providerData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.providerData = providerData
}

func (d *LabelDataSource) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.Conflicting(
			path.MatchRoot("delimiter"),
			path.MatchRoot("template"),
		),
		datasourcevalidator.Conflicting(
			path.MatchRoot("properties"),
			path.MatchRoot("template"),
		),
	}
}

func (d *LabelDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config LabelDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	// Generate the label
	label, diags := readLabel(ctx, d.providerData.contextClient, &config)
	resp.Diagnostics = append(resp.Diagnostics, diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set other properties
	config.Id = types.StringValue(label)
	config.Rendered = types.StringValue(label)

	// Write to state
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)

	tflog.Trace(ctx, "create label data source")
}

// processErrors iterates through a list of errors and adds them to the diagnostics.
func processErrors(errs []error, diags *diag.Diagnostics) {
	for _, err := range errs {
		if err != nil {
			diags.AddError("Validation Error", err.Error())
		}
	}
}

// readLabel determines the type of label to create and calls the appropriate method to create it.
func readLabel(ctx context.Context, client *client.Client, config *LabelDataSourceModel) (string, diag.Diagnostics) {
	if !config.Template.IsNull() {
		return readTemplatedLabel(ctx, client, config)
	}
	return readDelimitedLabel(ctx, client, config)
}

// readTemplatedLabel creates a label using a template.
func readTemplatedLabel(ctx context.Context, client *client.Client, config *LabelDataSourceModel) (string, diag.Diagnostics) {
	model, diags := templatedLabelModel{}.FromFramework(ctx, *config)
	if diags.HasError() {
		return "", diags
	}

	label, errs := client.GetTemplatedLabel(model.Template, model.Values, model.ReplaceCharsRegex, int(model.MaxLength), model.Truncate)
	processErrors(errs, &diags)

	return label, diags
}

// readDelimitedLabel creates a label using a delimiter.
func readDelimitedLabel(ctx context.Context, client *client.Client, config *LabelDataSourceModel) (string, diag.Diagnostics) {
	model, diags := delimitedLabelModel{}.FromFramework(ctx, *config)
	if diags.HasError() {
		return "", diags
	}

	label, errs := client.GetDelimitedLabel(model.Delimiter, model.PropertyNames, model.PropertyNames, model.Values, model.ReplaceCharsRegex, int(model.MaxLength), model.Truncate)
	processErrors(errs, &diags)

	return label, diags
}

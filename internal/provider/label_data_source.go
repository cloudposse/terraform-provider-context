package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &LabelDataSource{}

func NewLabelDataSource() datasource.DataSource {
	return &LabelDataSource{}
}

// LabelDataSource defines the data source implementation.
type LabelDataSource struct {
	providerData *providerData
}

// LabelDataSourceModel describes the data source data model.
type LabelDataSourceModel struct {
	Delimiter  types.String `tfsdk:"delimiter"`
	Properties types.List   `tfsdk:"properties"`
	Rendered   types.String `tfsdk:"rendered"`
	Template   types.String `tfsdk:"template"`
	Values     types.Map    `tfsdk:"values"`
	Id         types.String `tfsdk:"id"`
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
			"properties": schema.ListAttribute{
				MarkdownDescription: "List of properties to use when creating the label. Conflicts with `template`.",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"rendered": schema.StringAttribute{
				MarkdownDescription: "Rendered label",
				Computed:            true,
			},
			"template": schema.StringAttribute{
				MarkdownDescription: "Template to use when creating the label. Conflicts with `delimiter` and `properties`.",
				Optional:            true,
			},
			"values": schema.MapAttribute{
				MarkdownDescription: "Map of values to override or add to the context when creating the label.",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Label identifier",
				Computed:            true,
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

func (d *LabelDataSource) handleValidationErrors(resp *datasource.ReadResponse, errs []error) {
	for _, err := range errs {
		if err != nil {
			resp.Diagnostics.AddError("Validation Error", err.Error())
		}
	}
}

func (d *LabelDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config LabelDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	localValues, diags := FromFrameworkMap[string](ctx, config.Values)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	localProperties, diags := FromFrameworkList[string](ctx, config.Properties)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var label string
	if !config.Template.IsNull() {
		localTemplate := config.Template.ValueString()
		templatedLabel, errs := d.providerData.contextClient.GetTemplatedLabel(localTemplate, localValues)
		d.handleValidationErrors(resp, errs)

		if resp.Diagnostics.HasError() {
			return
		}

		label = templatedLabel
	} else {
		var localDelimiter *string
		if !config.Delimiter.IsNull() {
			localDelimiter = config.Delimiter.ValueStringPointer()
		}

		delimitedLabel, errs := d.providerData.contextClient.GetDelimitedLabel(localDelimiter, localProperties, localProperties, localValues)
		d.handleValidationErrors(resp, errs)

		if resp.Diagnostics.HasError() {
			return
		}
		label = delimitedLabel
	}

	if resp.Diagnostics.HasError() {
		return
	}

	config.Id = types.StringValue("Label-id")
	config.Rendered = types.StringValue(label)

	tflog.Trace(ctx, "create label data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

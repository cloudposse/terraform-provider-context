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
var _ datasource.DataSource = &TagsDataSource{}

func NewTagsDataSource() datasource.DataSource {
	return &TagsDataSource{}
}

// TagsDataSource defines the data source implementation.
type TagsDataSource struct {
	providerData *providerData
}

// TagsDataSourceModel describes the data source data model.
type TagsDataSourceModel struct {
	Id         types.String `tfsdk:"id"`
	Values     types.Map    `tfsdk:"values"`
	Tags       types.Map    `tfsdk:"tags"`
	TagsAsList types.List   `tfsdk:"tags_as_list"`
}

func (d *TagsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tags"
}

func (d *TagsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Tags data source",

		Attributes: map[string]schema.Attribute{
			"tags": schema.MapAttribute{
				MarkdownDescription: "Map of tags.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"tags_as_list": schema.ListAttribute{
				MarkdownDescription: "List of tags in {Key='key', Value='value'} format.",
				Computed:            true,
				ElementType:         types.MapType{ElemType: types.StringType},
			},
			"values": schema.MapAttribute{
				MarkdownDescription: "Map of values to override or add to the context when creating the label.",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Tags identifier",
				Computed:            true,
			},
		},
	}
}

func (d *TagsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *TagsDataSource) handleValidationErrors(resp *datasource.ReadResponse, errs []error) {
	for _, err := range errs {
		if err != nil {
			resp.Diagnostics.AddError("Validation Error", err.Error())
		}
	}
}

func (d *TagsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config TagsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	localValues, diags := FromFrameworkMap[string](ctx, config.Values)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tags, errs := d.providerData.contextClient.GetTags(localValues)
	d.handleValidationErrors(resp, errs)
	if resp.Diagnostics.HasError() {
		return
	}

	tagsList, errs := d.providerData.contextClient.GetTagsAsList(localValues)
	d.handleValidationErrors(resp, errs)
	if resp.Diagnostics.HasError() {
		return
	}

	config.Id = types.StringValue("Tags-id")

	frameworkTags, diags := types.MapValueFrom(ctx, types.StringType, tags)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	config.Tags = frameworkTags

	frameworkTagsAsList, diags := types.ListValueFrom(ctx, types.MapType{ElemType: types.StringType}, tagsList)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	config.TagsAsList = frameworkTagsAsList

	tflog.Trace(ctx, "create tags data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

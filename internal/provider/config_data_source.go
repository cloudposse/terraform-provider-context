package provider

import (
	"context"
	"fmt"

	mapHelpers "github.com/cloudposse/terraform-provider-context/pkg/map"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &ConfigDataSource{}
	_ datasource.DataSourceWithConfigure = &ConfigDataSource{}
)

func NewConfigDataSource() datasource.DataSource {
	return &ConfigDataSource{}
}

// ConfigDataSource defines the data source implementation.
type ConfigDataSource struct {
	providerData *providerData
}

// ConfigDataSourceModel describes the data source data model.
type ConfigDataSourceModel struct {
	Delimiter         types.String `tfsdk:"delimiter"`
	Enabled           types.Bool   `tfsdk:"enabled"`
	Properties        types.Map    `tfsdk:"properties"`
	PropertyOrder     types.List   `tfsdk:"property_order"`
	ReplaceCharsRegex types.String `tfsdk:"replace_chars_regex"`
	TagsKeyCase       types.String `tfsdk:"tags_key_case"`
	TagsValueCase     types.String `tfsdk:"tags_value_case"`
	Values            types.Map    `tfsdk:"values"`
	Id                types.String `tfsdk:"id"`
}

func (d *ConfigDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_config"
}

func (d *ConfigDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Context Config data source",
		Attributes: map[string]schema.Attribute{
			"delimiter": schema.StringAttribute{
				MarkdownDescription: "Delimiter to use when creating the label from properties. Conflicts with `template`.",
				Computed:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Flag to indicate if the config is enabled.",
				Computed:            true,
			},
			"properties": schema.MapNestedAttribute{
				MarkdownDescription: "A map of properties to use for labels created by the provider.",
				Computed:            true,
				NestedObject:        getPropertiesDSSchema(),
			},
			"property_order": schema.ListAttribute{
				MarkdownDescription: "A list of properties to use for labels created by the provider.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"replace_chars_regex": schema.StringAttribute{
				MarkdownDescription: "Regex to use for replacing characters in labels created by the provider.",
				Computed:            true,
			},
			"tags_key_case": schema.StringAttribute{
				MarkdownDescription: "Case to use for keys in tags created by the provider.",
				Computed:            true,
			},
			"tags_value_case": schema.StringAttribute{
				MarkdownDescription: "Case to use for values in tags created by the provider.",
				Computed:            true,
			},
			"values": schema.MapAttribute{
				MarkdownDescription: "A map of values to use for labels created by the provider.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Config identifier",
				Computed:            true,
				Optional:            true,
			},
		},
	}
}

func (d *ConfigDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ConfigDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config ConfigDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	// delimiter
	delimiter := d.providerData.contextClient.GetDelimiter()
	config.Delimiter = types.StringValue(delimiter)

	// enabled
	enabled := d.providerData.contextClient.IsEnabled()
	config.Enabled = types.BoolValue(enabled)

	// properties
	properties := d.providerData.contextClient.GetProperties()
	propMap := make(map[string]FrameworkProperty, len(properties))
	for _, v := range properties {
		propMap[v.Name] = FrameworkProperty{}.FromClientProperty(v)
	}

	props, diag := types.MapValueFrom(ctx, types.ObjectType{AttrTypes: FrameworkProperty{}.Types()}, propMap)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}
	config.Properties = props

	// propertyOrder
	propertyOrder := d.providerData.contextClient.GetPropertyOrder()
	propOrder, diag := types.ListValueFrom(ctx, types.StringType, propertyOrder)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}
	config.PropertyOrder = propOrder

	// replaceCharsRegex
	replaceRegexChars := d.providerData.contextClient.GetReplaceCharsRegex()
	config.ReplaceCharsRegex = types.StringValue(replaceRegexChars)

	// tagsKeyCase
	tagsKeyCase := d.providerData.contextClient.GetTagsKeyCase()
	config.TagsKeyCase = types.StringValue(tagsKeyCase)

	// tagsValueCase
	tagsValueCase := d.providerData.contextClient.GetTagsValueCase()
	config.TagsValueCase = types.StringValue(tagsValueCase)

	// values
	values := d.providerData.contextClient.GetValues()
	vals, diag := types.MapValueFrom(ctx, types.StringType, values)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}
	config.Values = vals

	// id
	id := mapHelpers.HashMap(config)
	config.Id = types.StringValue(id)

	tflog.Trace(ctx, "create config data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

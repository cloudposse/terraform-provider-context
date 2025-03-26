package provider

import (
	"context"
	"fmt"

	"github.com/cloudposse/terraform-provider-context/internal/model"
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
	providerData *model.ProviderData
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

	providerData, ok := req.ProviderData.(*model.ProviderData)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *providerData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.providerData = providerData
}

func (d *ConfigDataSource) setBasicConfig(config *ConfigDataSourceModel) {
	// delimiter
	delimiter := d.providerData.ProviderConfig.GetDelimiter()
	config.Delimiter = types.StringValue(delimiter)

	// enabled
	enabled := d.providerData.ProviderConfig.IsEnabled()
	config.Enabled = types.BoolValue(enabled)
}

func (d *ConfigDataSource) setProperties(ctx context.Context, config *ConfigDataSourceModel, resp *datasource.ReadResponse) {
	properties := d.providerData.ProviderConfig.GetProperties()
	propMap := make(map[string]model.FrameworkProperty, len(properties))
	for _, v := range properties {
		propMap[v.Name] = model.FrameworkProperty{}.FromConfigProperty(v)
	}

	props, diag := types.MapValueFrom(ctx, types.ObjectType{AttrTypes: model.FrameworkProperty{}.Types()}, propMap)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}
	config.Properties = props
}

func (d *ConfigDataSource) setPropertyOrder(ctx context.Context, config *ConfigDataSourceModel, resp *datasource.ReadResponse) {
	propertyOrder := d.providerData.ProviderConfig.GetPropertyOrder()
	propOrder, diag := types.ListValueFrom(ctx, types.StringType, propertyOrder)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}
	config.PropertyOrder = propOrder
}

func (d *ConfigDataSource) setValues(ctx context.Context, config *ConfigDataSourceModel, resp *datasource.ReadResponse) {
	values := d.providerData.ProviderConfig.GetValues()
	vals, diag := types.MapValueFrom(ctx, types.StringType, values)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}
	config.Values = vals
}

//nolint:gocritic
func (d *ConfigDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config ConfigDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	d.setBasicConfig(&config)

	d.setProperties(ctx, &config, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	d.setPropertyOrder(ctx, &config, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	// replaceCharsRegex
	replaceRegexChars := d.providerData.ProviderConfig.GetReplaceCharsRegex()
	config.ReplaceCharsRegex = types.StringValue(replaceRegexChars)

	// tagsKeyCase
	tagsKeyCase := d.providerData.ProviderConfig.GetTagsKeyCase()
	config.TagsKeyCase = types.StringValue(tagsKeyCase)

	// tagsValueCase
	tagsValueCase := d.providerData.ProviderConfig.GetTagsValueCase()
	config.TagsValueCase = types.StringValue(tagsValueCase)

	d.setValues(ctx, &config, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	// id
	id := mapHelpers.HashMap(config)
	config.Id = types.StringValue(id)

	tflog.Trace(ctx, "create config data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

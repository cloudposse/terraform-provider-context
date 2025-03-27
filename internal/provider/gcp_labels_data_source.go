package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &GcpLabelsDataSource{}

func NewGcpLabelsDataSource() datasource.DataSource {
	return &GcpLabelsDataSource{}
}

// GcpLabelsDataSource extends TagsDataSource
type GcpLabelsDataSource struct {
	TagsDataSource
}

// GcpLabelsDataSourceModel describes the data source data model.
type GcpLabelsDataSourceModel struct {
	TagsDataSourceModel
	ReplacementMap types.Map `tfsdk:"replacement_map"`
}

func (d *GcpLabelsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gcp_labels"
}

func (d *GcpLabelsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "GCP Labels data source",

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
			"tags_key_case": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The case to use for the keys of tags created by the provider. Valid values are: none, camel, lower, snake, title, upper.",
				Validators: []validator.String{
					stringvalidator.OneOf(ValidCases...),
				},
			},
			"tags_value_case": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The case to use for the values of tags created by the provider. Valid values are: none, camel, lower, snake, title, upper.",
				Validators: []validator.String{
					stringvalidator.OneOf(ValidCases...),
				},
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
			"replacement_map": schema.MapAttribute{
				MarkdownDescription: "Map of strings to replace in the tag, applies to both key and value. The key is the string to replace, and the value is the string to replace it with.",
				Optional:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (d *GcpLabelsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config GcpLabelsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	localValues := d.getLocalValues(ctx, &config.TagsDataSourceModel, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	localValues = d.getLocalReplacements(&config, localValues)
	if resp.Diagnostics.HasError() {
		return
	}

	localTagsKeyCase := d.getLocalTagsKeyCase(&config.TagsDataSourceModel, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	localTagsValueCase := d.getLocalTagsValueCase(&config.TagsDataSourceModel, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	d.setTags(ctx, &config.TagsDataSourceModel, resp, localValues, localTagsKeyCase, localTagsValueCase)
	if resp.Diagnostics.HasError() {
		return
	}

	d.setTagsList(ctx, &config.TagsDataSourceModel, resp, localValues, localTagsKeyCase, localTagsValueCase)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "created GCP labels data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

func (d *GcpLabelsDataSource) getLocalReplacements(config *GcpLabelsDataSourceModel, values map[string]string) map[string]string {
	replacedValues := make(map[string]string)
	if !config.ReplacementMap.IsNull() {
		for tagKey, tagValue := range values {
			for old, newString := range config.ReplacementMap.Elements() {
				replacedValues[strings.ReplaceAll(tagKey, old, newString.String())] = strings.ReplaceAll(tagValue, old, newString.String())
			}
		}
		return replacedValues
	}

	return values
}

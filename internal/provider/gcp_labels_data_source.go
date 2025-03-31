package provider

import (
	"context"
	"github.com/cloudposse/terraform-provider-context/internal/framework"
	"github.com/cloudposse/terraform-provider-context/pkg/cases"
	mapHelpers "github.com/cloudposse/terraform-provider-context/pkg/map"
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

// GcpLabelsDataSource extends TagsDataSource.
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

//nolint:gocritic
func (d *GcpLabelsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config GcpLabelsDataSourceModel

	// Read Terraform configuration data into the model.
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	localValues := d.getLocalValues(ctx, &config.TagsDataSourceModel, resp)
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

	d.setTags(ctx, &config, resp, localValues, localTagsKeyCase, localTagsValueCase)
	if resp.Diagnostics.HasError() {
		return
	}

	d.setTagsList(ctx, &config, resp, localValues, localTagsKeyCase, localTagsValueCase)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "created GCP labels data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

func (d *GcpLabelsDataSource) setTags(ctx context.Context, config *GcpLabelsDataSourceModel, resp *datasource.ReadResponse, localValues map[string]string, localTagsKeyCase, localTagsValueCase *cases.Case) {
	tags, errs := d.providerData.ProviderConfig.GetTags(localValues, localTagsKeyCase, localTagsValueCase)
	d.handleValidationErrors(resp, errs)
	if resp.Diagnostics.HasError() {
		return
	}
	frameworkTags, diags := types.MapValueFrom(ctx, types.StringType, d.runReplaceOnTags(ctx, config, tags))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	config.Tags = frameworkTags

	tagsAsHash := mapHelpers.HashMap(tags)
	config.Id = types.StringValue(tagsAsHash)
}

func (d *GcpLabelsDataSource) setTagsList(ctx context.Context, config *GcpLabelsDataSourceModel, resp *datasource.ReadResponse, localValues map[string]string, localTagsKeyCase, localTagsValueCase *cases.Case) {
	tagsList, errs := d.providerData.ProviderConfig.GetTagsAsList(localValues, localTagsKeyCase, localTagsValueCase)
	d.handleValidationErrors(resp, errs)
	if resp.Diagnostics.HasError() {
		return
	}

	frameworkTagsAsList, diags := types.ListValueFrom(ctx, types.MapType{ElemType: types.StringType}, d.runReplaceOnTagsList(ctx, config, tagsList))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	config.TagsAsList = frameworkTagsAsList
}

func (d *GcpLabelsDataSource) runReplaceOnTags(ctx context.Context, config *GcpLabelsDataSourceModel, values map[string]string) map[string]string {
	replacedValues := make(map[string]string)
	if !config.ReplacementMap.IsNull() && !config.ReplacementMap.IsUnknown() {
		replacementMap, _ := framework.FromFrameworkMap[string](ctx, config.ReplacementMap)
		for tagKey, tagValue := range values {
			newTagKey := tagKey
			newTagValue := tagValue
			for old, newString := range replacementMap {
				newTagKey = strings.ReplaceAll(tagKey, old, newString)
				newTagValue = strings.ReplaceAll(tagValue, old, newString)
				replacedValues[newTagKey] = newTagValue
			}
		}

		return replacedValues
	}
	return values
}

func (d *GcpLabelsDataSource) runReplaceOnTagsList(ctx context.Context, config *GcpLabelsDataSourceModel, values []map[string]string) []map[string]string {
	if !config.ReplacementMap.IsNull() && !config.ReplacementMap.IsUnknown() {
		replacementMap, _ := framework.FromFrameworkMap[string](ctx, config.ReplacementMap)

		for _, TagKV := range values {
			for old, newString := range replacementMap {
				TagKV["Key"] = strings.ReplaceAll(TagKV["Key"], old, newString)
				TagKV["Value"] = strings.ReplaceAll(TagKV["Value"], old, newString)
			}
		}

		return values
	}
	return values
}

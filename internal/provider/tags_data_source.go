package provider

import (
	"context"
	"fmt"

	"github.com/cloudposse/terraform-provider-context/internal/framework"
	"github.com/cloudposse/terraform-provider-context/internal/model"
	"github.com/cloudposse/terraform-provider-context/pkg/cases"
	mapHelpers "github.com/cloudposse/terraform-provider-context/pkg/map"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
	providerData *model.ProviderData
}

// TagsDataSourceModel describes the data source data model.
type TagsDataSourceModel struct {
	Id            types.String `tfsdk:"id"`
	Values        types.Map    `tfsdk:"values"`
	Tags          types.Map    `tfsdk:"tags"`
	TagsKeyCase   types.String `tfsdk:"tags_key_case"`
	TagsValueCase types.String `tfsdk:"tags_value_case"`
	TagsAsList    types.List   `tfsdk:"tags_as_list"`
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
		},
	}
}

func (d *TagsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *TagsDataSource) handleValidationErrors(resp *datasource.ReadResponse, errs []error) {
	for _, err := range errs {
		if err != nil {
			resp.Diagnostics.AddError("Validation Error", err.Error())
		}
	}
}

func (d *TagsDataSource) getLocalValues(ctx context.Context, config *TagsDataSourceModel, resp *datasource.ReadResponse) map[string]string {
	localValues, diags := framework.FromFrameworkMap[string](ctx, config.Values)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return nil
	}
	return localValues
}

func (d *TagsDataSource) getLocalTagsKeyCase(config *TagsDataSourceModel, resp *datasource.ReadResponse) *cases.Case {
	if !config.TagsKeyCase.IsNull() {
		tagsKeyCase, err := cases.FromString(*config.TagsKeyCase.ValueStringPointer())
		if err != nil {
			resp.Diagnostics.AddError("Failed to convert tags_key_case to model", err.Error())
			return nil
		}
		return &tagsKeyCase
	}
	return nil
}

func (d *TagsDataSource) getLocalTagsValueCase(config *TagsDataSourceModel, resp *datasource.ReadResponse) *cases.Case {
	if !config.TagsValueCase.IsNull() {
		tagsValueCase, err := cases.FromString(*config.TagsValueCase.ValueStringPointer())
		if err != nil {
			resp.Diagnostics.AddError("Failed to convert tags_value_case to model", err.Error())
			return nil
		}
		return &tagsValueCase
	}
	return nil
}

//nolint:revive
func (d *TagsDataSource) setTags(ctx context.Context, config *TagsDataSourceModel, resp *datasource.ReadResponse, localValues map[string]string, localTagsKeyCase, localTagsValueCase *cases.Case) {
	tags, errs := d.providerData.ProviderConfig.GetTags(localValues, localTagsKeyCase, localTagsValueCase)
	d.handleValidationErrors(resp, errs)
	if resp.Diagnostics.HasError() {
		return
	}

	frameworkTags, diags := types.MapValueFrom(ctx, types.StringType, tags)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	config.Tags = frameworkTags

	tagsAsHash := mapHelpers.HashMap(tags)
	config.Id = types.StringValue(tagsAsHash)
}

//nolint:revive
func (d *TagsDataSource) setTagsList(ctx context.Context, config *TagsDataSourceModel, resp *datasource.ReadResponse, localValues map[string]string, localTagsKeyCase, localTagsValueCase *cases.Case) {
	tagsList, errs := d.providerData.ProviderConfig.GetTagsAsList(localValues, localTagsKeyCase, localTagsValueCase)
	d.handleValidationErrors(resp, errs)
	if resp.Diagnostics.HasError() {
		return
	}

	frameworkTagsAsList, diags := types.ListValueFrom(ctx, types.MapType{ElemType: types.StringType}, tagsList)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	config.TagsAsList = frameworkTagsAsList
}

//nolint:gocritic
func (d *TagsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config TagsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	localValues := d.getLocalValues(ctx, &config, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	localTagsKeyCase := d.getLocalTagsKeyCase(&config, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	localTagsValueCase := d.getLocalTagsValueCase(&config, resp)
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

	tflog.Trace(ctx, "created tags data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

package provider

import (
	"context"

	"github.com/cloudposse/terraform-provider-context/internal/model"
	"github.com/cloudposse/terraform-provider-context/pkg/cases"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure ContextProvider satisfies various provider interfaces.
var _ provider.Provider = &ContextProvider{}

// ContextProvider defines the provider implementation.
type ContextProvider struct {
	providerData *model.ProviderData
	// version is set to the provider version on release, "dev" when the provider is built and ran locally, and "test"
	// when running acceptance testing.
	version string
}

// ContextProviderModel describes the provider data model.
type providerConfigModel struct {
	Delimiter         types.String `tfsdk:"delimiter"`
	Enabled           types.Bool   `tfsdk:"enabled"`
	Properties        types.Map    `tfsdk:"properties"`
	PropertyOrder     types.List   `tfsdk:"property_order"`
	ReplaceCharsRegex types.String `tfsdk:"replace_chars_regex"`
	TagsKeyCase       types.String `tfsdk:"tags_key_case"`
	TagsValueCase     types.String `tfsdk:"tags_value_case"`
	Values            types.Map    `tfsdk:"values"`
}

func (p *ContextProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "context"
	resp.Version = p.version
}

func (p *ContextProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"delimiter": schema.StringAttribute{
				MarkdownDescription: "The default delimiter to use for labels created by the provider.",
				Optional:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "A boolean value to enable or disable the provider.",
				Optional:            true,
			},
			"properties": schema.MapNestedAttribute{
				MarkdownDescription: "A map of properties to use for labels created by the provider.",
				Optional:            true,
				NestedObject:        getPropertiesSchema(),
			},
			"property_order": schema.ListAttribute{
				MarkdownDescription: "The default order of properties to use for labels created by the provider.",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"replace_chars_regex": schema.StringAttribute{
				MarkdownDescription: "The regex to use for replacing characters in labels created by the provider. Any characters that match the regex will be removed from the label.",
				Optional:            true,
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
				MarkdownDescription: "A map of values to use for labels created by the provider.",
				Optional:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (p *ContextProvider) getConfigProperties(ctx context.Context, providerConfigModel *providerConfigModel, resp *provider.ConfigureResponse) []model.Property {
	properties := map[string]model.FrameworkProperty{}
	resp.Diagnostics.Append(providerConfigModel.Properties.ElementsAs(ctx, &properties, false)...)
	if resp.Diagnostics.HasError() {
		return nil
	}

	configProperties := []model.Property{}
	for k, prop := range properties {
		property, err := prop.ToModel(k)
		if err != nil {
			resp.Diagnostics.AddError("Failed to convert property to model", err.Error())
			return nil
		}
		configProperties = append(configProperties, *property)
	}
	return configProperties
}

func (p *ContextProvider) getPropertyOrder(ctx context.Context, providerConfigModel *providerConfigModel, resp *provider.ConfigureResponse) []string {
	propertyOrder := []string{}
	resp.Diagnostics.Append(providerConfigModel.PropertyOrder.ElementsAs(ctx, &propertyOrder, false)...)
	if resp.Diagnostics.HasError() {
		return nil
	}
	return propertyOrder
}

func (p *ContextProvider) getValues(ctx context.Context, providerConfigModel *providerConfigModel, resp *provider.ConfigureResponse) map[string]string {
	values := map[string]string{}
	resp.Diagnostics.Append(providerConfigModel.Values.ElementsAs(ctx, &values, false)...)
	if resp.Diagnostics.HasError() {
		return nil
	}
	return values
}

func (p *ContextProvider) getOptions(providerConfigModel *providerConfigModel, resp *provider.ConfigureResponse) []func(*model.ProviderConfig) {
	options := []func(*model.ProviderConfig){}

	if !providerConfigModel.Enabled.IsNull() {
		options = append(options, model.WithEnabled(providerConfigModel.Enabled.ValueBool()))
	}

	if !providerConfigModel.Delimiter.IsNull() {
		options = append(options, model.WithDelimiter(providerConfigModel.Delimiter.ValueString()))
	}

	if !providerConfigModel.ReplaceCharsRegex.IsNull() {
		options = append(options, model.WithReplaceCharsRegex(providerConfigModel.ReplaceCharsRegex.ValueString()))
	}

	if !providerConfigModel.TagsKeyCase.IsNull() {
		keyCase, err := cases.FromString(providerConfigModel.TagsKeyCase.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Failed to convert tags key case", err.Error())
			return nil
		}
		options = append(options, model.WithTagsKeyCase(keyCase))
	}

	if !providerConfigModel.TagsValueCase.IsNull() {
		valueCase, err := cases.FromString(providerConfigModel.TagsValueCase.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Failed to convert tags value case", err.Error())
			return nil
		}
		options = append(options, model.WithTagsValueCase(valueCase))
	}

	return options
}

func (p *ContextProvider) createAndValidateProviderConfig(configProperties []model.Property, propertyOrder []string, values map[string]string, options []func(*model.ProviderConfig), resp *provider.ConfigureResponse) *model.ProviderData {
	providerConfig, err := model.NewProviderConfig(configProperties, propertyOrder, values, options...)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create provider config", err.Error())
		return nil
	}

	if errs := providerConfig.ValidateProperties(values); len(errs) > 0 {
		for _, err := range errs {
			resp.Diagnostics.AddError("Validation Error", err.Error())
		}
		return nil
	}

	return &model.ProviderData{
		ProviderConfig: providerConfig,
	}
}

func (p *ContextProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var providerConfigModel providerConfigModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &providerConfigModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	configProperties := p.getConfigProperties(ctx, &providerConfigModel, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	propertyOrder := p.getPropertyOrder(ctx, &providerConfigModel, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	values := p.getValues(ctx, &providerConfigModel, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	options := p.getOptions(&providerConfigModel, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Data received from the configuration", map[string]any{
		"delimiter":           providerConfigModel.Delimiter.ValueString(),
		"enabled":             providerConfigModel.Enabled.ValueBool(),
		"properties":          configProperties,
		"property_order":      propertyOrder,
		"replace_chars_regex": providerConfigModel.ReplaceCharsRegex.ValueString(),
		"tags_key_case":       providerConfigModel.TagsKeyCase.ValueString(),
		"tags_value_case":     providerConfigModel.TagsValueCase.ValueString(),
		"values":              values,
	})

	providerData := p.createAndValidateProviderConfig(configProperties, propertyOrder, values, options, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	p.providerData = providerData
	resp.DataSourceData = providerData
	resp.ResourceData = providerData

	tflog.Info(ctx, "Configured provider config", map[string]any{"success": true})
}

func (p *ContextProvider) Resources(ctx context.Context) []func() resource.Resource {
	return nil
}

func (p *ContextProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewConfigDataSource,
		NewLabelDataSource,
		NewTagsDataSource,
	}
}

func NewProvider(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ContextProvider{
			version: version,
		}
	}
}

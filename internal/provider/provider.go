package provider

import (
	"context"

	"github.com/cloudposse/terraform-provider-context/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure ContextProvider satisfies various provider interfaces.
var _ provider.Provider = &ContextProvider{}

// ContextProvider defines the provider implementation.
type ContextProvider struct {
	providerData *providerData
	// version is set to the provider version on release, "dev" when the provider is built and ran locally, and "test"
	// when running acceptance testing.
	version string
}

// ContextProviderModel describes the provider data model.
type config struct {
	Delimiter     types.String `tfsdk:"delimiter"`
	Enabled       types.Bool   `tfsdk:"enabled"`
	Properties    types.Map    `tfsdk:"properties"`
	PropertyOrder types.List   `tfsdk:"property_order"`
	Values        types.Map    `tfsdk:"values"`
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
			"values": schema.MapAttribute{
				MarkdownDescription: "A map of values to use for labels created by the provider.",
				Optional:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (p *ContextProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config config

	// Get the configuration from the request
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert config to native go types
	properties := map[string]FrameworkProperty{}
	resp.Diagnostics.Append(config.Properties.ElementsAs(ctx, &properties, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clientProperties := []client.Property{}
	for k, prop := range properties {

		property, err := prop.ToModel(k)
		if err != nil {
			resp.Diagnostics.AddError("Failed to convert property to model", err.Error())
			return
		}
		clientProperties = append(clientProperties, *property)
	}

	propertyOrder := []string{}
	resp.Diagnostics.Append(config.PropertyOrder.ElementsAs(ctx, &propertyOrder, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	values := map[string]string{}
	resp.Diagnostics.Append(config.Values.ElementsAs(ctx, &values, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Data received from the configuration", map[string]any{
		"delimiter":      config.Delimiter.ValueString(),
		"enabled":        config.Enabled.ValueBool(),
		"properties":     clientProperties,
		"property_order": propertyOrder,
		"values":         values,
	})

	// Create the context client
	client, err := client.NewClient(clientProperties, propertyOrder, values)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create context client", err.Error())
		return
	}

	providerData := &providerData{
		contextClient: client,
	}

	// Set the provider data in the response
	p.providerData = providerData
	resp.DataSourceData = providerData
	resp.ResourceData = providerData

	tflog.Info(ctx, "Configured Context client", map[string]any{"success": true})
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

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ContextProvider{
			version: version,
		}
	}
}

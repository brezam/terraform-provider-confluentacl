package internal

import (
	"context"
	"os"
	"terraform-provider-confluentacl/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	envVarCloudApiKey    = "CONFLUENT_CLOUD_API_KEY"
	envVarCloudApiSecret = "CONFLUENT_CLOUD_API_SECRET"
)

var (
	_ provider.Provider = &confluentaclProvider{}
)

func NewProvider() provider.Provider {
	return &confluentaclProvider{}
}

type confluentaclProvider struct{}

type confluentaclProviderModel struct {
	ConfluentCloudApiKey    types.String `tfsdk:"confluent_cloud_api_key"`
	ConfluentCloudApiSecret types.String `tfsdk:"confluent_cloud_api_secret"`
}

// Metadata returns the provider type name.
func (p *confluentaclProvider) Metadata(ctx context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "confluentacl"
}

// Schema defines the provider-level schema for configuration data.
func (p *confluentaclProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"confluent_cloud_api_key": schema.StringAttribute{
				Optional: true,
			},
			"confluent_cloud_api_secret": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func (p *confluentaclProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring ConfluentAcl client")

	var config confluentaclProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cloudApiKey := os.Getenv(envVarCloudApiKey)
	cloudApiSecret := os.Getenv(envVarCloudApiSecret)

	if !config.ConfluentCloudApiKey.IsNull() {
		cloudApiKey = config.ConfluentCloudApiKey.ValueString()
	}
	if !config.ConfluentCloudApiSecret.IsNull() {
		cloudApiSecret = config.ConfluentCloudApiSecret.ValueString()
	}

	if cloudApiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("confluentCloudApiKey"),
			"Missing Confluent Cloud API Key", "Provider requires Confluent Cloud Cloud api key to function",
		)
	}
	if cloudApiSecret == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("confluentCloudApiSecret"),
			"Missing Confluent Cloud API Secret", "Provider requires Confluent Cloud Cloud api secret to function",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	client_ := client.New(cloudApiKey, cloudApiSecret)
	resp.DataSourceData = client_
	resp.ResourceData = client_
}

// DataSources defines the data sources implemented in the provider.
func (p *confluentaclProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewSchemaRegistryDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *confluentaclProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewAclResource,
		NewApiKeyResource,
	}
}

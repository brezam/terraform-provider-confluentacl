package internal

import (
	"context"
	"terraform-provider-confluentacl/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &SchemaRegistryDataSource{}
	_ datasource.DataSourceWithConfigure = &SchemaRegistryDataSource{}
)

type SchemaRegistryDataSource struct {
	client *client.Client
}

type SchemaRegistryDataSourceModel struct {
	ID            types.String `tfsdk:"id"`
	EnvironmentId types.String `tfsdk:"environment_id"`
	SchemaId      types.String `tfsdk:"schema_id"`
	Endpoint      types.String `tfsdk:"endpoint"`
}

func NewSchemaRegistryDataSource() datasource.DataSource {
	return &SchemaRegistryDataSource{}
}

func (r *SchemaRegistryDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_schema_registry"
}

func (r *SchemaRegistryDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*client.Client)
}

func (r *SchemaRegistryDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"environment_id": schema.StringAttribute{
				Required: true,
			},
			"schema_id": schema.StringAttribute{
				Computed: true,
			},
			"endpoint": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *SchemaRegistryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state SchemaRegistryDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	schema, err := r.client.GetFirstSchemaRegistry(state.EnvironmentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to get first schema registry in environment", err.Error())
	}
	if resp.Diagnostics.HasError() {
		return
	}
	state.ID = types.StringValue(schema.Id)
	state.Endpoint = types.StringValue(schema.Endpoint)
	state.SchemaId = types.StringValue(schema.Id)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

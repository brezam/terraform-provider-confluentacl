package internal

import (
	"context"
	"strconv"
	"terraform-provider-confluentacl/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &ApiKeyResource{}
	_ resource.ResourceWithConfigure = &ApiKeyResource{}
)

type ApiKeyResource struct {
	client *client.Client
}

type ApiKeyResourceModel struct {
	ID                 types.String `tfsdk:"id"`
	ServiceAccountName types.String `tfsdk:"service_account_name"`
	EnvironmentId      types.String `tfsdk:"environment_id"`
	ResourceId         types.String `tfsdk:"resource_id"`
	Description        types.String `tfsdk:"description"`
	ApiKey             types.String `tfsdk:"api_key"`
	ApiSecret          types.String `tfsdk:"api_secret"`
}

func NewApiKeyResource() resource.Resource {
	return &ApiKeyResource{}
}

func (r *ApiKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_api_key"
}

func (r *ApiKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*client.Client)
}

func (r *ApiKeyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"service_account_name": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"environment_id": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"resource_id": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"description": schema.StringAttribute{
				Optional: true,
			},
			"api_key": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"api_secret": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *ApiKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ApiKeyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	userId, err := r.client.GetSaNumericId(plan.ServiceAccountName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to list service accounts", err.Error())
	}
	if resp.Diagnostics.HasError() {
		return
	}
	apiKey, err := r.client.CreateApiKey(userId, plan.EnvironmentId.ValueString(), plan.ResourceId.ValueString(), plan.Description.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to create Api Key", err.Error())
	}
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = types.StringValue(strconv.Itoa(apiKey.ID))
	plan.ApiKey = types.StringValue(apiKey.Key)
	plan.ApiSecret = types.StringValue(apiKey.Secret)
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ApiKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ApiKeyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiKey, err := r.client.ReadApiKey(state.ApiKey.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read Api Key", err.Error())
	}
	if resp.Diagnostics.HasError() {
		return
	}
	if apiKey == nil {
		state.ID = types.StringNull()
		return
	}

	// Shenanigans for "empty" description
	currentDescription := state.Description.ValueString()
	readDescription := apiKey.Spec.Description
	if readDescription == "--" && (currentDescription == "--" || currentDescription == "") {
		readDescription = currentDescription
	}

	state.Description = types.StringValue(readDescription)
	state.ResourceId = types.StringValue(apiKey.Spec.Resource.ID)
}

func (r *ApiKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ApiKeyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	description := plan.Description.ValueString()
	if description == "" {
		description = "--" // Description cannot be set to empty, the request doesn't work even in the UI
	}
	err := r.client.UpdateApiKey(plan.ID.ValueString(), description, plan.EnvironmentId.ValueString(), plan.ResourceId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to update api key", err.Error())
		if resp.Diagnostics.HasError() {
			return
		}
	}
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *ApiKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ApiKeyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteApiKey(state.ID.ValueString(), state.EnvironmentId.ValueString(), state.ResourceId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete api key", err.Error())
		if resp.Diagnostics.HasError() {
			return
		}
	}
}

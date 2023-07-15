package internal

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"terraform-provider-confluentacl/internal/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource              = &AclResource{}
	_ resource.ResourceWithConfigure = &AclResource{}
)

type AclResource struct {
	client *client.Client
}

type AclResourceModel struct {
	ID                 types.String `tfsdk:"id"`
	RestEndpoint       types.String `tfsdk:"rest_endpoint"`
	ServiceAccountName types.String `tfsdk:"service_account_name"`
	ClusterId          types.String `tfsdk:"cluster_id"`
	ResourceType       types.String `tfsdk:"resource_type"`
	ResourceName       types.String `tfsdk:"resource_name"`
	PatternType        types.String `tfsdk:"pattern_type"`
	Host               types.String `tfsdk:"host"`
	Operation          types.String `tfsdk:"operation"`
	Permission         types.String `tfsdk:"permission"`
}

func NewAclResource() resource.Resource {
	return &AclResource{}
}

func (r *AclResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_acl"
}

func (r *AclResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*client.Client)
}

func (r *AclResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"cluster_id": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^lkc-.+`), "Value must start with lkc-"),
				},
			},
			"rest_endpoint": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"resource_type": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Validators: []validator.String{
					stringvalidator.OneOf("TOPIC", "GROUP", "CLUSTER", "TRANSACTIONAL_ID", "DELEGATION_TOKEN"),
				},
			},
			"resource_name": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"pattern_type": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Validators: []validator.String{
					stringvalidator.OneOf("MATCH", "LITERAL", "PREFIXED"),
				},
			},
			"host": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"operation": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Validators: []validator.String{
					stringvalidator.OneOf("READ", "WRITE", "CREATE", "DELETE", "ALTER", "DESCRIBE", "CLUSTER_ACTION", "DESCRIBE_CONFIGS", "ALTER_CONFIGS", "IDEMPOTENT_WRITE"),
				},
			},
			"permission": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Validators:    []validator.String{stringvalidator.OneOf("ALLOW", "DENY")},
			},
		},
	}
}

func (r *AclResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Get plan
	var plan AclResourceModel
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
	tflog.Info(ctx, fmt.Sprintf("UserId %d", userId))
	if userId == 0 {
		resp.Diagnostics.AddError("Could not find service account with name "+plan.ServiceAccountName.ValueString(), "")
	}
	requestBody := &client.ACLRequest{
		Principal:    fmt.Sprintf("User:%d", userId),
		ResourceName: plan.ResourceName.ValueString(),
		ResourceType: plan.ResourceType.ValueString(),
		PatternType:  plan.PatternType.ValueString(),
		Host:         plan.Host.ValueString(),
		Operation:    plan.Operation.ValueString(),
		Permission:   plan.Permission.ValueString(),
	}
	err = r.client.CreateACL(plan.RestEndpoint.ValueString(), plan.ClusterId.ValueString(), requestBody)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create ACL", err.Error())
	}
	plan.ID = types.StringValue(makeIdForAclModel(&plan))
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *AclResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AclResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	userId, err := r.client.GetSaNumericId(state.ServiceAccountName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to list service accounts", err.Error())
	}
	if resp.Diagnostics.HasError() {
		return
	}
	queryParams := &client.ACLRequest{
		Principal:    fmt.Sprintf("User:%d", userId),
		ResourceName: state.ResourceName.ValueString(),
		ResourceType: state.ResourceType.ValueString(),
		PatternType:  state.PatternType.ValueString(),
		Host:         state.Host.ValueString(),
		Operation:    state.Operation.ValueString(),
		Permission:   state.Permission.ValueString(),
	}
	aclsFound, err := r.client.ListSpecificACLs(
		state.RestEndpoint.ValueString(),
		state.ClusterId.ValueString(),
		queryParams,
	)
	if err != nil {
		resp.Diagnostics.AddError("Failure to read all acls in cluster", err.Error())
	}
	if len(aclsFound) > 1 {
		resp.Diagnostics.AddError("Expected to find 1 ACL matching spec", "")
		return
	}
	if len(aclsFound) == 1 {
		acl := aclsFound[0]
		state.ClusterId = types.StringValue(acl.ClusterId)
		state.ResourceName = types.StringValue(acl.ResourceName)
		state.ResourceType = types.StringValue(acl.ResourceType)
		state.PatternType = types.StringValue(acl.PatternType)
		state.Host = types.StringValue(acl.Host)
		state.Operation = types.StringValue(acl.Operation)
		state.Permission = types.StringValue(acl.Permission)
	} else {
		state.ID = types.StringNull()
	}
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *AclResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// There's no update.
}

func (r *AclResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AclResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	userId, err := r.client.GetSaNumericId(state.ServiceAccountName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to list service accounts", err.Error())
	}
	if resp.Diagnostics.HasError() {
		return
	}
	queryParams := &client.ACLRequest{
		Principal:    fmt.Sprintf("User:%d", userId),
		ResourceName: state.ResourceName.ValueString(),
		ResourceType: state.ResourceType.ValueString(),
		PatternType:  state.PatternType.ValueString(),
		Host:         state.Host.ValueString(),
		Operation:    state.Operation.ValueString(),
		Permission:   state.Permission.ValueString(),
	}
	r.client.DeleteAcl(state.RestEndpoint.ValueString(), state.ClusterId.ValueString(), queryParams)
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete ACL", err.Error())
	}

}

func makeIdForAclModel(model *AclResourceModel) string {
	return fmt.Sprintf("%s/%s/%s",
		model.ClusterId.ValueString(),
		model.ServiceAccountName.ValueString(),
		strings.Join([]string{
			model.ResourceType.ValueString(),
			model.ResourceName.ValueString(),
			model.PatternType.ValueString(),
			model.Host.ValueString(),
			model.Operation.ValueString(),
			model.Permission.ValueString(),
		}, "#"))
}

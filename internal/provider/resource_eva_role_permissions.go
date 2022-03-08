package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/mad-it/terraform-provider-eva/internal/eva"
)

type rolePermissionsType struct{}

func (t rolePermissionsType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Eva role permissions configuration.",

		Attributes: map[string]tfsdk.Attribute{
			"role_id": {
				MarkdownDescription: "ID of the role which permissions will be attached to",
				Computed:            true,
				Type:                types.Int64Type,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"scoped_functionalities": {
				MarkdownDescription: "list of scoped functionalities to be attached",
				Required:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}

func (t rolePermissionsType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return roleResource{
		provider: provider,
	}, diags
}

type rolePermissionsResource struct {
	provider provider
}

type inputRolePermissionsData struct {
	RoleID                types.Int64  `tfsdk:"role_id"`
	ScopedFunctionalities types.String `tfsdk:"scoped_functionalities"`
}

type inputScopedFunctionality struct {
	Functionality     types.String
	Scope             types.Int64
	RequiresElevation types.Bool
}

func (r rolePermissionsResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data inputRolePermissionsData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	var scopedFunctionalities []eva.ScopedFunctionality

	if err := json.Unmarshal([]byte(data.ScopedFunctionalities.Value), &scopedFunctionalities); err != nil {
		resp.Diagnostics.AddError("scoped_functionalities is not a valid json.", fmt.Sprintf("Unable to parse scoped_functionalities field, got error: %s", err))
		return
	}

	client_resp, err := r.provider.evaClient.AttachFunctionalitiesToRole(ctx, eva.AttachFunctionalitiesToRoleRequest{
		RoleID:                data.RoleID.Value,
		ScopedFunctionalities: scopedFunctionalities,
	})

	if err != nil {
		resp.Diagnostics.AddError("Creating role permissions failed.", fmt.Sprintf("Unable to create role permisions, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "Created role permissions.", client_resp)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r rolePermissionsResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data inputRolePermissionsData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	client_resp, err := r.provider.evaClient.GetRole(ctx, eva.GetRoleRequest{
		ID: data.RoleID.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError("Getting role unit failed.", fmt.Sprintf("Unable to get role, got error: %s", err))
		return
	}

	jsonresp, err := json.Marshal(client_resp.Result.ScopedFunctionalities)

	if err != nil {
		resp.Diagnostics.AddError("Getting role unit failed.", fmt.Sprintf("Unable to get role, got error: %s", err))
		return
	}

	data.ScopedFunctionalities = types.String{Value: string(jsonresp)}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r rolePermissionsResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var data inputRolePermissionsData

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	var scopedFunctionalities []eva.ScopedFunctionality

	if jsonParseErr := json.Unmarshal([]byte(data.ScopedFunctionalities.Value), &scopedFunctionalities); jsonParseErr != nil {
		resp.Diagnostics.AddError("Updating role permissions failed.", fmt.Sprintf("Unable to parse scoped functionalities, got error: %s", jsonParseErr))
		return
	}

	roleData, getRoleErr := r.provider.evaClient.GetRole(ctx, eva.GetRoleRequest{
		ID: data.RoleID.Value,
	})

	if getRoleErr != nil {
		resp.Diagnostics.AddError("Updating role permissions failed.", fmt.Sprintf("Unable to get current role permissions, got error: %s", getRoleErr))
		return
	}

	var currentScopedFunctionalities []eva.ScopedFunctionality

	for _, scopedFunctionality := range roleData.Result.ScopedFunctionalities {
		currentScopedFunctionalities = append(currentScopedFunctionalities, eva.ScopedFunctionality(scopedFunctionality))
	}

	_, detachErr := r.provider.evaClient.DetachFunctionalitiesFromRole(ctx, eva.DetachFunctionalitiesFromRoleRequest{
		RoleID:                data.RoleID.Value,
		ScopedFunctionalities: currentScopedFunctionalities,
	})

	if detachErr != nil {
		resp.Diagnostics.AddError("Updating role permissions failed.", fmt.Sprintf("Unable to detach current role permissions, got error: %s", detachErr))
		return
	}

	_, attachErr := r.provider.evaClient.AttachFunctionalitiesToRole(ctx, eva.AttachFunctionalitiesToRoleRequest{
		RoleID:                data.RoleID.Value,
		ScopedFunctionalities: scopedFunctionalities,
	})

	if attachErr != nil {
		resp.Diagnostics.AddError("Updating role permissions failed.", fmt.Sprintf("Unable to attach new role permissions, got error: %s", detachErr))
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r rolePermissionsResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data inputRoleData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.provider.evaClient.DeleteRole(ctx, eva.DeleteRoleRequest{
		ID: data.ID.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError("Deleting role unit failed.", fmt.Sprintf("Unable to delete role, got error: %s", err))
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r rolePermissionsResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}

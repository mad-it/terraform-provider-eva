package provider

import (
	"context"
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
			"id": {
				Type:                types.Int64Type,
				Computed:            true,
				MarkdownDescription: "Role permissions does not have a unique ID. This attribute is populated within the role_id",
			},
			"role_id": {
				MarkdownDescription: "ID of the role which permissions will be attached to",
				Required:            true,
				Type:                types.Int64Type,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"scoped_functionalities": {
				MarkdownDescription: "list of scoped functionalities to be attached",
				Required:            true,
				Attributes: tfsdk.ListNestedAttributes(
					map[string]tfsdk.Attribute{
						"functionality": {
							MarkdownDescription: "functionality identifier",
							Required:            true,
							Type:                types.StringType,
						},
						"scope": {
							MarkdownDescription: "functionality scope",
							Required:            true,
							Type:                types.Int64Type,
						},
						"requires_elevation": {
							MarkdownDescription: "whether functionality requires elevation or not",
							Required:            true,
							Type:                types.BoolType,
						},
					},
					tfsdk.ListNestedAttributesOptions{},
				),
			},
		},
	}, nil
}

func (t rolePermissionsType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return rolePermissionsResource{
		provider: provider,
	}, diags
}

type rolePermissionsResource struct {
	provider provider
}

type inputRolePermissionsData struct {
	ID                    types.Int64                `tfsdk:"id"`
	RoleID                types.Int64                `tfsdk:"role_id"`
	ScopedFunctionalities []inputScopedFunctionality `tfsdk:"scoped_functionalities"`
}

type inputScopedFunctionality struct {
	Functionality     types.String `tfsdk:"functionality"`
	Scope             types.Int64  `tfsdk:"scope"`
	RequiresElevation types.Bool   `tfsdk:"requires_elevation"`
}

func (r rolePermissionsResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data inputRolePermissionsData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	var scopedFunctionalities []eva.ScopedFunctionality

	for _, scopedFunctionality := range data.ScopedFunctionalities {
		scopedFunctionalities = append(scopedFunctionalities, eva.ScopedFunctionality{
			Functionality:     scopedFunctionality.Functionality.Value,
			Scope:             scopedFunctionality.Scope.Value,
			RequiresElevation: scopedFunctionality.RequiresElevation.Value,
		})
	}

	client_resp, err := r.provider.evaClient.AttachFunctionalitiesToRole(ctx, eva.AttachFunctionalitiesToRoleRequest{
		RoleID:                data.RoleID.Value,
		ScopedFunctionalities: scopedFunctionalities,
	})

	if err != nil {
		resp.Diagnostics.AddError("Creating role permissions failed.", fmt.Sprintf("Unable to create role permisions, got error: %s", err))
		return
	}

	data.ID = types.Int64{Value: data.RoleID.Value}

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

	var updatedScopedFunctionalities []inputScopedFunctionality

	for _, scopedFunctionality := range client_resp.Result.ScopedFunctionalities {
		updatedScopedFunctionalities = append(updatedScopedFunctionalities, inputScopedFunctionality{
			Functionality:     types.String{Value: scopedFunctionality.Functionality},
			Scope:             types.Int64{Value: scopedFunctionality.Scope},
			RequiresElevation: types.Bool{Value: scopedFunctionality.RequiresElevation},
		})
	}

	data.ScopedFunctionalities = updatedScopedFunctionalities

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

	roleData, getRoleErr := r.provider.evaClient.GetRole(ctx, eva.GetRoleRequest{
		ID: data.RoleID.Value,
	})

	if getRoleErr != nil {
		resp.Diagnostics.AddError("Updating role permissions failed.", fmt.Sprintf("Unable to get current role permissions, got error: %s", getRoleErr))
		return
	}

	var currentScopedFunctionalities []eva.ScopedFunctionality

	for _, scopedFunctionality := range roleData.Result.ScopedFunctionalities {
		currentScopedFunctionalities = append(currentScopedFunctionalities, eva.ScopedFunctionality{
			Functionality:     scopedFunctionality.Functionality,
			Scope:             scopedFunctionality.Scope,
			RequiresElevation: scopedFunctionality.RequiresElevation,
		})
	}

	_, detachErr := r.provider.evaClient.DetachFunctionalitiesFromRole(ctx, eva.DetachFunctionalitiesFromRoleRequest{
		RoleID:                data.RoleID.Value,
		ScopedFunctionalities: currentScopedFunctionalities,
	})

	if detachErr != nil {
		resp.Diagnostics.AddError("Updating role permissions failed.", fmt.Sprintf("Unable to detach current role permissions, got error: %s", detachErr))
		return
	}

	var scopedFunctionalitiesToUpdate []eva.ScopedFunctionality

	for _, scopedFunctionality := range data.ScopedFunctionalities {
		scopedFunctionalitiesToUpdate = append(scopedFunctionalitiesToUpdate, eva.ScopedFunctionality{
			Functionality:     scopedFunctionality.Functionality.Value,
			Scope:             scopedFunctionality.Scope.Value,
			RequiresElevation: scopedFunctionality.RequiresElevation.Value,
		})
	}

	_, attachErr := r.provider.evaClient.AttachFunctionalitiesToRole(ctx, eva.AttachFunctionalitiesToRoleRequest{
		RoleID:                data.RoleID.Value,
		ScopedFunctionalities: scopedFunctionalitiesToUpdate,
	})

	if attachErr != nil {
		resp.Diagnostics.AddError("Updating role permissions failed.", fmt.Sprintf("Unable to attach new role permissions, got error: %s", detachErr))
		return
	}

	data.ID = types.Int64{Value: data.RoleID.Value}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r rolePermissionsResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data inputRolePermissionsData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	roleData, getRoleErr := r.provider.evaClient.GetRole(ctx, eva.GetRoleRequest{
		ID: data.RoleID.Value,
	})

	if getRoleErr != nil {
		resp.Diagnostics.AddError("Deleting role permissions failed.", fmt.Sprintf("Unable to get current role permissions, got error: %s", getRoleErr))
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
		resp.Diagnostics.AddError("Deleting role unit failed.", fmt.Sprintf("Unable to detach current role permissions, got error: %s", detachErr))
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r rolePermissionsResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}

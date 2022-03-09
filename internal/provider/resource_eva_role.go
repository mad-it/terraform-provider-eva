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

type roleType struct{}

func (t roleType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Eva role configuration.",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "ID of the role.",
				Computed:            true,
				Type:                types.Int64Type,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"name": {
				MarkdownDescription: "Name of the role",
				Optional:            true,
				Type:                types.StringType,
			},
			"user_type": {
				MarkdownDescription: "User type this role applies to",
				Optional:            true,
				Type:                types.Int64Type,
			},
			"code": {
				MarkdownDescription: "A unique code to represent the role",
				Optional:            true,
				Type:                types.StringType,
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

func (t roleType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return roleResource{
		provider: provider,
	}, diags
}

type roleResource struct {
	provider provider
}

type inputRole struct {
	ID                    types.Int64                    `tfsdk:"id"`
	Name                  types.String                   `tfsdk:"name"`
	UserType              types.Int64                    `tfsdk:"user_type"`
	Code                  types.String                   `tfsdk:"code"`
	ScopedFunctionalities []inputRoleScopedFunctionality `tfsdk:"scoped_functionalities"`
}

type inputRoleScopedFunctionality struct {
	Functionality     types.String `tfsdk:"functionality"`
	Scope             types.Int64  `tfsdk:"scope"`
	RequiresElevation types.Bool   `tfsdk:"requires_elevation"`
}

func (s inputRole) getListOfScopedFunctionalities() []eva.RoleScopedFunctionality {
	var scopedFunctionalities []eva.RoleScopedFunctionality

	for _, scopedFunctionality := range s.ScopedFunctionalities {
		scopedFunctionalities = append(scopedFunctionalities, eva.RoleScopedFunctionality{
			Functionality:     scopedFunctionality.Functionality.Value,
			Scope:             scopedFunctionality.Scope.Value,
			RequiresElevation: scopedFunctionality.RequiresElevation.Value,
		})
	}

	return scopedFunctionalities
}

func (s inputRole) setListOfScopedFunctionalities(scopedFunctionalities []eva.RoleScopedFunctionality) {
	for _, scopedFunctionality := range scopedFunctionalities {
		s.ScopedFunctionalities = append(s.ScopedFunctionalities, inputRoleScopedFunctionality{
			Functionality:     types.String{Value: scopedFunctionality.Functionality},
			Scope:             types.Int64{Value: scopedFunctionality.Scope},
			RequiresElevation: types.Bool{Value: scopedFunctionality.RequiresElevation},
		})
	}
}

func (r roleResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data inputRole

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	createdRole, createRoleErr := r.provider.evaClient.CreateRole(ctx, eva.CreateRoleRequest{
		Name:     data.Name.Value,
		UserType: data.UserType.Value,
		Code:     data.Code.Value,
	})

	if createRoleErr != nil {
		resp.Diagnostics.AddError("Creating role unit failed.", fmt.Sprintf("Unable to create role, got error: %s", createRoleErr))
		return
	}

	data.ID = types.Int64{Value: createdRole.ID}

	_, attachPermissionsToRoleErr := r.provider.evaClient.AttachFunctionalitiesToRole(ctx, eva.AttachFunctionalitiesToRoleRequest{
		RoleID:                data.ID.Value,
		ScopedFunctionalities: data.getListOfScopedFunctionalities(),
	})

	if attachPermissionsToRoleErr != nil {
		resp.Diagnostics.AddError("Creating role permissions failed.", fmt.Sprintf("Unable to create role permisions, got error: %s", attachPermissionsToRoleErr))
		return
	}

	tflog.Trace(ctx, "Created a new role.")

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r roleResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data inputRole

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	roleData, err := r.provider.evaClient.GetRole(ctx, eva.GetRoleRequest{
		ID: data.ID.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError("Getting role unit failed.", fmt.Sprintf("Unable to get role, got error: %s", err))
		return
	}

	data.Name = types.String{Value: roleData.Result.Name}
	data.UserType = types.Int64{Value: roleData.Result.UserType}
	data.Code = types.String{Value: roleData.Result.Code}
	data.setListOfScopedFunctionalities(roleData.Result.ScopedFunctionalities)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r roleResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var data inputRole

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.provider.evaClient.UpdateRole(ctx, eva.UpdateRoleRequest{
		ID:       data.ID.Value,
		Name:     data.Name.Value,
		UserType: data.UserType.Value,
		Code:     data.Code.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError("Updating role unit failed.", fmt.Sprintf("Unable to update role, got error: %s", err))
		return
	}

	roleData, getRoleErr := r.provider.evaClient.GetRole(ctx, eva.GetRoleRequest{
		ID: data.ID.Value,
	})

	if getRoleErr != nil {
		resp.Diagnostics.AddError("Updating role permissions failed.", fmt.Sprintf("Unable to get role, got error: %s", getRoleErr))
		return
	}

	_, detachErr := r.provider.evaClient.DetachFunctionalitiesFromRole(ctx, eva.DetachFunctionalitiesFromRoleRequest{
		RoleID:                data.ID.Value,
		ScopedFunctionalities: roleData.Result.ScopedFunctionalities,
	})

	if detachErr != nil {
		resp.Diagnostics.AddError("Updating role permissions failed.", fmt.Sprintf("Unable to detach current role permissions, got error: %s", detachErr))
		return
	}

	_, attachErr := r.provider.evaClient.AttachFunctionalitiesToRole(ctx, eva.AttachFunctionalitiesToRoleRequest{
		RoleID:                data.ID.Value,
		ScopedFunctionalities: data.getListOfScopedFunctionalities(),
	})

	if attachErr != nil {
		resp.Diagnostics.AddError("Updating role permissions failed.", fmt.Sprintf("Unable to attach new role permissions, got error: %s", detachErr))
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r roleResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data inputRole

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

func (r roleResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}

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

type customFieldType struct{}

func (t customFieldType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Eva role configuration.",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Computed: true,
				Type:     types.Int64Type,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"name": {
				Optional: true,
				Type:     types.StringType,
			},
			"type_id": {
				Optional: true,
				Type:     types.Int64Type,
			},
			"type_key": {
				Optional: true,
				Type:     types.StringType,
			},
			"order": {
				Optional: true,
				Type:     types.Int64Type,
			},
			"backend_id": {
				Optional: true,
				Type:     types.StringType,
			},
			"visible_by_user_types": {
				Optional: true,
				Type:     types.Int64Type,
			},
			"editable_by_user_types": {
				Optional: true,
				Type:     types.Int64Type,
			},
			"options": {
				Optional: true,
				Attributes: tfsdk.ListNestedAttributes(
					map[string]tfsdk.Attribute{
						"is_array": {
							Optional: true,
							Type:     types.BoolType,
						},
						"is_required": {
							Optional: true,
							Type:     types.BoolType,
						},
						"minimum_value": {
							Optional: true,
							Type:     types.Int64Type,
						},
						"maximum_value": {
							Optional: true,
							Type:     types.Int64Type,
						},
						"minimum_length": {
							Optional: true,
							Type:     types.Int64Type,
						},
						"maximum_length": {
							Optional: true,
							Type:     types.Int64Type,
						},
						"minimum_date": {
							Optional: true,
							Type:     types.StringType,
						},
						"maximum_date": {
							Optional: true,
							Type:     types.StringType,
						},
						"default_value": {
							Optional: true,
							Type:     types.StringType,
						},
					},
					tfsdk.ListNestedAttributesOptions{},
				),
			},
		},
	}, nil
}

func (t customFieldType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return role{
		provider: provider,
	}, diags
}

type customField struct {
	provider provider
}

type customFieldTypeData struct {
	ID                    types.Int64                 `tfsdk:"id"`
	Name                  types.String                `tfsdk:"name"`
	UserType              types.Int64                 `tfsdk:"user_type"`
	Code                  types.String                `tfsdk:"code"`
	ScopedFunctionalities []roleFunctionalityTypeData `tfsdk:"scoped_functionalities"`
}

type customFieldasdasdasdTypeData struct {
	Functionality     types.String `tfsdk:"functionality"`
	Scope             types.Int64  `tfsdk:"scope"`
	RequiresElevation types.Bool   `tfsdk:"requires_elevation"`
}

func (r customField) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data roleProviderTypeData

	diags := req.Plan.Get(ctx, &data)
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

	diags = resp.State.Set(ctx, &roleProviderTypeData{
		ID:       types.Int64{Value: data.ID.Value},
		Name:     types.String{Value: data.Name.Value},
		UserType: types.Int64{Value: data.UserType.Value},
		Code:     types.String{Value: data.Code.Value},
	})

	tflog.Trace(ctx, "Created a new role.")

	_, attachPermissionsToRoleErr := r.provider.evaClient.AttachFunctionalitiesToRole(ctx, eva.AttachFunctionalitiesToRoleRequest{
		RoleID:                data.ID.Value,
		ScopedFunctionalities: data.getListOfFunctionalities(),
	})

	if attachPermissionsToRoleErr != nil {
		resp.Diagnostics.AddError("Creating role permissions failed. Please try to apply changes again.", fmt.Sprintf("Unable to create role permisions, got error: %s", attachPermissionsToRoleErr))
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r customField) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data roleProviderTypeData

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
	data.setListOfFunctionalities(roleData.Result.ScopedFunctionalities)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// func (r customField) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
// 	var data roleProviderTypeData

// 	diags := req.Plan.Get(ctx, &data)
// 	resp.Diagnostics.Append(diags...)

// 	if resp.Diagnostics.HasError() {
// 		return
// 	}

// 	_, err := r.provider.evaClient.UpdateCustomField(ctx, eva.UpdateRoleRequest{
// 		ID:       data.ID.Value,
// 		Name:     data.Name.Value,
// 		UserType: data.UserType.Value,
// 		Code:     data.Code.Value,
// 	})

// 	if err != nil {
// 		resp.Diagnostics.AddError("Updating role permissions failed.", fmt.Sprintf("Unable to attach new role permissions, got error: %s", detachErr))
// 		return
// 	}

// 	diags = resp.State.Set(ctx, &data)
// 	resp.Diagnostics.Append(diags...)
// }

func (r customField) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data roleProviderTypeData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.provider.evaClient.DeleteCustomField(ctx, eva.DeleteCustomFieldRequest{
		ID: data.ID.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError("Deleting custom field failed.", fmt.Sprintf("Unable to delete custom field, got error: %s", err))
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r customField) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}

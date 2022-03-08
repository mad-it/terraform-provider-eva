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

type inputRoleData struct {
	ID       types.Int64  `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	UserType types.Int64  `tfsdk:"user_type"`
	Code     types.String `tfsdk:"code"`
}

func (r roleResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data inputRoleData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	client_resp, err := r.provider.evaClient.CreateRole(ctx, eva.CreateRoleRequest{
		Name:     data.Name.Value,
		UserType: data.UserType.Value,
		Code:     data.Code.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError("Creating role unit failed.", fmt.Sprintf("Unable to create role, got error: %s", err))
		return
	}

	data.ID = types.Int64{Value: client_resp.ID}

	tflog.Trace(ctx, "Created a new role.")

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r roleResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data inputRoleData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	client_resp, err := r.provider.evaClient.GetRole(ctx, eva.GetRoleRequest{
		ID: data.ID.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError("Getting role unit failed.", fmt.Sprintf("Unable to get role, got error: %s", err))
		return
	}

	data.Name = types.String{Value: client_resp.Result.Name}
	data.UserType = types.Int64{Value: client_resp.Result.UserType}
	data.Code = types.String{Value: client_resp.Result.Code}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r roleResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var data inputRoleData

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

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r roleResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
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

func (r roleResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}

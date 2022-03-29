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

type employeeType struct{}

func (t employeeType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Eva employee configuration.",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "ID of the employee.",
				Computed:            true,
				Type:                types.Int64Type,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"first_name": {
				Optional: true,
				Type:     types.StringType,
			},
			"last_name": {
				Optional: true,
				Type:     types.StringType,
			},
			"email_address": {
				Required: true,
				Type:     types.StringType,
			},
			"password": {
				Required:  true,
				Type:      types.StringType,
				Sensitive: true,
			},
			"roles": {
				MarkdownDescription: "list of scoped functionalities to be attached",
				Optional:            true,
				Attributes: tfsdk.ListNestedAttributes(
					map[string]tfsdk.Attribute{
						"role_id": {
							MarkdownDescription: "id of the role",
							Required:            true,
							Type:                types.Int64Type,
						},
						"user_type": {
							MarkdownDescription: "functionality scope",
							Required:            true,
							Type:                types.Int64Type,
						},
						"organization_unit_id": {
							MarkdownDescription: "id of the organization unit the role applies too.",
							Optional:            true,
							Type:                types.Int64Type,
						},
					},
					tfsdk.ListNestedAttributesOptions{
						MinItems: 1,
					},
				),
			},
		},
	}, nil
}

func (t employeeType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return employee{
		provider: provider,
	}, diags
}

type roleTypeData struct {
	RoleID             types.Int64 `tfsdk:"role_id"`
	UserType           types.Int64 `tfsdk:"user_type"`
	OrganizationUnitID types.Int64 `tfsdk:"organization_unit_id"`
}

type employeeTypeData struct {
	ID           types.Int64    `tfsdk:"id"`
	FirstName    types.String   `tfsdk:"first_name"`
	LastName     types.String   `tfsdk:"last_name"`
	EmailAddress types.String   `tfsdk:"email_address"`
	Password     types.String   `tfsdk:"password"`
	Roles        []roleTypeData `tfsdk:"roles"`
}

type employee struct {
	provider provider
}

func (d employeeTypeData) getEvaUserRoles() []eva.RoleOrganizationUnitSet {
	var roles []eva.RoleOrganizationUnitSet

	for _, userRole := range d.Roles {
		roles = append(roles, eva.RoleOrganizationUnitSet{
			RoleID:             userRole.RoleID.Value,
			OrganizationUnitID: userRole.OrganizationUnitID.Value,
			UserType:           userRole.UserType.Value,
		})
	}

	return roles
}

func (d employeeTypeData) setUserRoles(userRoles []eva.UserRole) {
	for _, userRole := range userRoles {
		d.Roles = append(d.Roles, roleTypeData{
			RoleID:             types.Int64{Value: userRole.RoleID},
			OrganizationUnitID: types.Int64{Value: userRole.OrganizationUnitID},
			UserType:           types.Int64{Value: userRole.UserType},
		})
	}

}

func (r employee) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data employeeTypeData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	client_resp, err := r.provider.evaClient.CreateEmployee(ctx, eva.CreateEmployeeUserRequest{
		FirstName:    data.FirstName.Value,
		LastName:     data.LastName.Value,
		EmailAddress: data.EmailAddress.Value,
		Password:     data.Password.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError("Creating employee failed.", fmt.Sprintf("Unable to create employee, got error: %s", err))
		return
	}

	data.ID = types.Int64{Value: client_resp.ID}
	employee := employeeTypeData{
		ID:           types.Int64{Value: client_resp.ID},
		FirstName:    types.String{Value: data.FirstName.Value},
		LastName:     types.String{Value: data.LastName.Value},
		EmailAddress: types.String{Value: data.EmailAddress.Value},
		Password:     types.String{Value: data.Password.Value},
	}

	diags = resp.State.Set(ctx, &employee)
	resp.Diagnostics.Append(diags...)

	_, err = r.provider.evaClient.SetUserRole(ctx, eva.SetUserRoleRequest{
		UserId: client_resp.ID,
		Roles:  data.getEvaUserRoles(),
	})

	if err != nil {
		resp.Diagnostics.AddError("Assign the roles to the user failed.", fmt.Sprintf("Unable to assign roles, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "Created an employee.")

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r employee) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data employeeTypeData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	client_resp, err := r.provider.evaClient.GetUser(ctx, eva.GetUserRequest{
		ID: data.ID.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError("Creating employee unit failed.", fmt.Sprintf("Unable to create example, got error: %s", err))
		return
	}

	data.FirstName = types.String{Value: client_resp.FirstName}
	data.LastName = types.String{Value: client_resp.LastName}
	data.EmailAddress = types.String{Value: client_resp.EmailAddress}
	// Password cannot be read, so this value is not updated in the state.

	roles_client_resp, err := r.provider.evaClient.GetUserRole(ctx, eva.GetUserRoleRequest{
		UserId: data.ID.Value,
	})

	data.setUserRoles(roles_client_resp.Roles)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r employee) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var data employeeTypeData

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.provider.evaClient.UpdateUser(ctx, eva.UpdateUserRequest{
		ID:           data.ID.Value,
		FirstName:    data.FirstName.Value,
		LastName:     data.LastName.Value,
		EmailAddress: data.EmailAddress.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError("Updating employee unit failed.", fmt.Sprintf("Unable to update employee, got error: %s", err))
		return
	}

	// TODO: Implement UpdatePassword functionality.

	_, err = r.provider.evaClient.SetUserRole(ctx, eva.SetUserRoleRequest{
		UserId: data.ID.Value,
		Roles:  data.getEvaUserRoles(),
	})

	if err != nil {
		resp.Diagnostics.AddError("Assign the roles to the user failed.", fmt.Sprintf("Unable to assign roles, got error: %s", err))
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r employee) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data employeeTypeData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.provider.evaClient.DeleteUser(ctx, eva.DeleteUserRequest{
		ID: data.ID.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError("Deleting employee failed.", fmt.Sprintf("Unable to delete employee, got error: %s", err))
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r employee) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}

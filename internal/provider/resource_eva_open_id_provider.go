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

type openIdProviderType struct{}

func (t openIdProviderType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Eva organization unit openIdProviders configuration.",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "ID of the OpenIdProvider.",
				Computed:            true,
				Type:                types.Int64Type,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"enabled": {
				Required: true,
				Type:     types.BoolType,
			},
			"primary": {
				Optional: true,
				Type:     types.BoolType,
			},
			"base_url": {
				Required: true,
				Type:     types.StringType,
			},
			"name": {
				Optional: true,
				Type:     types.StringType,
			},
			"client_id": {
				Required: true,
				Type:     types.StringType,
			},
			"create_users": {
				Required: true,
				Type:     types.BoolType,
			},
			"first_name_claim": {
				Optional: true,
				Type:     types.StringType,
			},
			"last_name_claim": {
				Optional: true,
				Type:     types.StringType,
			},
			"email_address_claim": {
				Optional: true,
				Type:     types.StringType,
			},
			"nickname_claim": {
				Optional: true,
				Type:     types.StringType,
			},
			"user_type": {
				Required: true,
				Type:     types.Int64Type,
			},
		},
	}, nil
}

func (t openIdProviderType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return openIdProvider{
		provider: provider,
	}, diags
}

type openIdProviderTypeData struct {
	ID                types.Int64  `tfsdk:"id"`
	BaseUrl           types.String `tfsdk:"base_url"`
	ClientID          types.String `tfsdk:"client_id"`
	CreateUsers       types.Bool   `tfsdk:"create_users"`
	EmailAddressClaim types.String `tfsdk:"email_address_claim"`
	Enabled           types.Bool   `tfsdk:"enabled"`
	FirstNameClaim    types.String `tfsdk:"first_name_claim"`
	LastNameClaim     types.String `tfsdk:"last_name_claim"`
	Name              types.String `tfsdk:"name"`
	NicknameClaim     types.String `tfsdk:"nickname_claim"`
	Primary           types.Bool   `tfsdk:"primary"`
	UserType          types.Int64  `tfsdk:"user_type"`
}

type openIdProvider struct {
	provider provider
}

func (r openIdProvider) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data openIdProviderTypeData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	client_resp, err := r.provider.evaClient.CreateOpenIDProvider(ctx, eva.CreateOpenIDProviderRequest{
		BaseUrl:           data.BaseUrl.Value,
		ClientID:          data.ClientID.Value,
		CreateUsers:       data.CreateUsers.Value,
		EmailAddressClaim: data.EmailAddressClaim.Value,
		Enabled:           data.Enabled.Value,
		FirstNameClaim:    data.FirstNameClaim.Value,
		LastNameClaim:     data.LastNameClaim.Value,
		Name:              data.Name.Value,
		NicknameClaim:     data.NicknameClaim.Value,
		UserType:          data.UserType.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError("Creating openIdProvider unit failed.", fmt.Sprintf("Unable to create openIdProvider, got error: %s", err))
		return
	}

	data.ID = types.Int64{Value: client_resp.ID}

	tflog.Trace(ctx, "Created an openIdProvider.")

	if data.Primary.Value {
		_, err := r.provider.evaClient.SetPrimaryOpenIDProvider(ctx, eva.SetPrimaryOpenIDProviderRequest{
			ID: client_resp.ID,
		})

		if err != nil {
			resp.Diagnostics.AddWarning("Setting primary open Id provider failed.", fmt.Sprintf("Unable to set primary open Id provider, got error: %s", err))
		}
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r openIdProvider) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data openIdProviderTypeData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	client_resp, err := r.provider.evaClient.GetOpenIDProvider(ctx, eva.GetOpenIDProviderRequest{
		ID: data.ID.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError("Getting openIdProvider data failed.", fmt.Sprintf("Unable to get openIdProvider, got error: %s", err))
		return
	}

	data.BaseUrl = types.String{Value: client_resp.BaseUrl}
	data.ClientID = types.String{Value: client_resp.ClientID}
	data.ID = types.Int64{Value: client_resp.ID}
	data.CreateUsers = types.Bool{Value: client_resp.CreateUsers}
	data.Enabled = types.Bool{Value: client_resp.Enabled}
	data.EmailAddressClaim = types.String{Value: client_resp.EmailAddressClaim}
	data.FirstNameClaim = types.String{Value: client_resp.FirstNameClaim}
	data.LastNameClaim = types.String{Value: client_resp.LastNameClaim}
	data.NicknameClaim = types.String{Value: client_resp.NicknameClaim}
	data.UserType = types.Int64{Value: client_resp.UserType}
	data.Name = types.String{Value: client_resp.Name}
	data.Primary = types.Bool{Value: client_resp.Primary}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r openIdProvider) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var plan openIdProviderTypeData
	var state openIdProviderTypeData

	stateDiags := req.State.Get(ctx, &state)
	diags := req.Plan.Get(ctx, &plan)

	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(stateDiags...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.provider.evaClient.UpdateOpenIDProvider(ctx, eva.UpdateOpenIDProviderRequest{
		ID:                plan.ID.Value,
		BaseUrl:           plan.BaseUrl.Value,
		ClientID:          plan.ClientID.Value,
		CreateUsers:       plan.CreateUsers.Value,
		EmailAddressClaim: plan.EmailAddressClaim.Value,
		Enabled:           plan.Enabled.Value,
		FirstNameClaim:    plan.FirstNameClaim.Value,
		LastNameClaim:     plan.LastNameClaim.Value,
		Name:              plan.Name.Value,
		NicknameClaim:     plan.NicknameClaim.Value,
		UserType:          plan.UserType.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError("Updating openIdProvider unit failed.", fmt.Sprintf("Unable to update openIdProvider, got error: %s", err))
		return
	}

	if state.Primary.Value != plan.Primary.Value && plan.Primary.Value {
		_, err := r.provider.evaClient.SetPrimaryOpenIDProvider(ctx, eva.SetPrimaryOpenIDProviderRequest{
			ID: plan.ID.Value,
		})

		if err != nil {
			resp.Diagnostics.AddWarning("Setting primary open Id provider failed.", fmt.Sprintf("Unable to set primary open Id provider, got error: %s", err))
		}
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r openIdProvider) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data openIdProviderTypeData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.provider.evaClient.DeleteOpenIDProvider(ctx, eva.DeleteOpenIDProviderRequest{
		ID: data.ID.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError("Deleting openIdProvider unit failed.", fmt.Sprintf("Unable to delete openIdProvider, got error: %s", err))
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r openIdProvider) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}

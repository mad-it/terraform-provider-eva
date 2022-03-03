package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/mad-it/terraform-provider-eva/internal/eva"
)

type settingType struct{}

func (t settingType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Eva organization unit configration.",

		Attributes: map[string]tfsdk.Attribute{
			"key": {
				MarkdownDescription: "Key of the setting.",
				Required:            true,
				Type:                types.StringType,
			},
			"value": {
				MarkdownDescription: "Value of the setting",
				Required:            true,
				Type:                types.StringType,
			},
			"organization_unit_id": {
				MarkdownDescription: "ID of the organization unit to apply the settings for.",
				Optional:            true,
				Type:                types.Int64Type,
			},
		},
	}, nil
}

func (t settingType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return setting{
		provider: provider,
	}, diags
}

type settingTypeData struct {
	Key                types.String `tfsdk:"key"`
	Value              types.String `tfsdk:"value"`
	OrganizationUnitID types.Int64  `tfsdk:"organization_unit_id"`
}

type setting struct {
	provider provider
}

func (r setting) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data settingTypeData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: do we need an ID or some unique identifier for this resource?
	_, err := r.provider.evaClient.SetSettings(ctx, eva.SetSettingsRequest{
		Key:                data.Key.Value,
		Value:              data.Value.Value,
		OrganizationUnitID: data.OrganizationUnitID.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError("Creating setting unit failed.", fmt.Sprintf("Unable to create example, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "Created an setting.")

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r setting) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data settingTypeData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	client_resp, err := r.provider.evaClient.GetSetting(ctx, eva.GetSettingRequest{
		Key:                data.Key.Value,
		OrganizationUnitID: data.OrganizationUnitID.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError("Creating setting unit failed.", fmt.Sprintf("Unable to create example, got error: %s", err))
		return
	}

	if strings.HasPrefix(client_resp.Value, "********") &&
		strings.HasSuffix(client_resp.Value, data.Value.Value[len(data.Value.Value)-4:]) {
		// For sensitive data like passwords, EVA replaces everything except the last 4 characters with *.
		// So when the last 4 characters match in the existing state, do nothing.
	} else {
		data.Value = types.String{Value: client_resp.Value}
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r setting) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var data settingTypeData

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: do we need an ID or some unique identifier for this resource?
	_, err := r.provider.evaClient.SetSettings(ctx, eva.SetSettingsRequest{
		Key:                data.Key.Value,
		Value:              data.Value.Value,
		OrganizationUnitID: data.OrganizationUnitID.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError("Updating setting unit failed.", fmt.Sprintf("Unable to update OU, got error: %s", err))
		return
	}

	data.Value = types.String{Value: data.Value.Value}
	data.OrganizationUnitID = types.Int64{Value: data.OrganizationUnitID.Value}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r setting) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data settingTypeData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.provider.evaClient.UnsetSettings(ctx, eva.UnsetSettingsRequest{
		Key:                data.Key.Value,
		OrganizationUnitID: data.OrganizationUnitID.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError("Deleting setting unit failed.", fmt.Sprintf("Unable to delete OU, got error: %s", err))
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r setting) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("key"), req, resp)
}

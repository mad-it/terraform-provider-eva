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

type cookbookType struct{}

func (t cookbookType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Eva cookbook configuration.",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "ID of the cookbook.",
				Computed:            true,
				Type:                types.Int64Type,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"name": {
				MarkdownDescription: "Name of the cookbook",
				Optional:            true,
				Type:                types.StringType,
			},
			"recipe": {
				MarkdownDescription: "Value of the cookbook",
				Required:            true,
				Type:                types.StringType,
			},
			"is_active": {
				MarkdownDescription: "boolean of whether the cookbook is active.",
				Optional:            true,
				Type:                types.BoolType,
			},
		},
	}, nil
}

func (t cookbookType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return cookbook{
		provider: provider,
	}, diags
}

type cookbookTypeData struct {
	ID       types.Int64  `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Recipe   types.String `tfsdk:"recipe"`
	IsActive types.Bool   `tfsdk:"is_active"`
}

type cookbook struct {
	provider provider
}

func (r cookbook) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data cookbookTypeData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	client_resp, err := r.provider.evaClient.CreateAccountingRecipe(ctx, eva.CreateAccountingRecipeRequest{
		Name:     data.Name.Value,
		Recipe:   data.Recipe.Value,
		IsActive: data.IsActive.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError("Creating cookbook unit failed.", fmt.Sprintf("Unable to create cookbook, got error: %s", err))
		return
	}

	data.ID = types.Int64{Value: client_resp.ID}

	tflog.Trace(ctx, "Created an cookbook.")

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r cookbook) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data cookbookTypeData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	client_resp, err := r.provider.evaClient.GetAccountingRecipe(ctx, eva.GetAccountingRecipeRequest{
		ID: data.ID.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError("Creating cookbook unit failed.", fmt.Sprintf("Unable to create example, got error: %s", err))
		return
	}

	data.IsActive = types.Bool{Value: client_resp.Recipe.IsActive}
	data.Name = types.String{Value: client_resp.Recipe.Name}
	data.Recipe = types.String{Value: client_resp.Recipe.Recipe}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r cookbook) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var data cookbookTypeData

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.provider.evaClient.UpdateAccountingRecipe(ctx, eva.UpdateAccountingRecipeRequest{
		ID:       data.ID.Value,
		Name:     data.Name.Value,
		Recipe:   data.Recipe.Value,
		IsActive: data.IsActive.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError("Updating cookbook unit failed.", fmt.Sprintf("Unable to update cookbook, got error: %s", err))
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r cookbook) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data cookbookTypeData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.provider.evaClient.DeleteAccountingRecipe(ctx, eva.DeleteAccountingRecipeRequest{
		ID: data.ID.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError("Deleting cookbook unit failed.", fmt.Sprintf("Unable to delete cookbook, got error: %s", err))
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r cookbook) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}

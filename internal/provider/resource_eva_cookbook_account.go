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

type cookbookAccountType struct{}

func (t cookbookAccountType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Eva cookbook account configuration.",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "ID of the cookbook account",
				Computed:            true,
				Type:                types.Int64Type,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"name": {
				MarkdownDescription: "Name of the cookbook account",
				Optional:            true,
				Type:                types.StringType,
			},
			"object_account": {
				MarkdownDescription: "?????",
				Required:            true,
				Type:                types.StringType,
			},
			"booking_flags": {
				MarkdownDescription: `Booking flags for the cookbook account:
				- None = 0,
				- WithTaxInformation = 1,
				- WithoutOffsets = 2,
				- WithOrderNumber = 4,
				- WithReference = 8,
				- WithInvoiceNumber = 16,
				- WithCurrencyInformation = 32
				`,
				Required: true,
				Type:     types.Int64Type,
			},
			"type": {
				MarkdownDescription: `Type of the cookbook account:
				- GeneralLedger = 1,
				- Debtor = 2,
				- Creditor = 3
				`,
				Required: true,
				Type:     types.Int64Type,
			},
		},
	}, nil
}

func (t cookbookAccountType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return cookbookAccount{
		provider: provider,
	}, diags
}

type cookbookAccountTypeData struct {
	ID            types.Int64 `tfsdk:"id"`
	Name          string      `tfsdk:"name"`
	ObjectAccount string      `tfsdk:"object_account"`
	BookingFlags  int64       `tfsdk:"booking_flags"`
	Type          int64       `tfsdk:"type"`
}

type cookbookAccount struct {
	provider provider
}

func (r cookbookAccount) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data cookbookAccountTypeData

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	clientResponse, err := r.provider.evaClient.CreateCookbookAccount(ctx, eva.CreateCookbookAccountRequest{
		Name:          data.Name,
		ObjectAccount: data.ObjectAccount,
		BookingFlags:  data.BookingFlags,
		Type:          data.Type,
	})

	if err != nil {
		resp.Diagnostics.AddError("Creating cookbook account failed.", fmt.Sprintf("Unable to create cookbook account, got error: %s", err))
		return
	}

	data.ID = types.Int64{Value: clientResponse.ID}

	tflog.Trace(ctx, "Created a cookbook account.")

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r cookbookAccount) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data cookbookAccountTypeData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	clientResponse, err := r.provider.evaClient.GetCookbookAccount(ctx, eva.GetCookbookAccountRequest{
		ID: data.ID.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError("Getting cookbook account data failed.", fmt.Sprintf("Unable to get cookbook account, got error: %s", err))
		return
	}

	data.Name = clientResponse.Name
	data.ObjectAccount = clientResponse.ObjectAccount
	data.BookingFlags = clientResponse.BookingFlags
	data.Type = clientResponse.Type

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r cookbookAccount) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var data cookbookAccountTypeData

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.provider.evaClient.UpdateCookbookAccount(ctx, eva.UpdateCookbookAccountRequest{
		ID:            data.ID.Value,
		Name:          data.Name,
		ObjectAccount: data.ObjectAccount,
		BookingFlags:  data.BookingFlags,
		Type:          data.Type,
	})

	if err != nil {
		resp.Diagnostics.AddError("Updating cookbook account failed.", fmt.Sprintf("Unable to update cookbook account, got error: %s", err))
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r cookbookAccount) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data cookbookAccountTypeData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.provider.evaClient.DeleteCookbookAccount(ctx, eva.DeleteCookbookAccountRequest{
		ID: data.ID.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError("Deleting cookbook account failed.", fmt.Sprintf("Unable to delete cookbook account, got error: %s", err))
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r cookbookAccount) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}

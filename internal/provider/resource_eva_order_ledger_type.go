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

type orderLedgerTypeSchema struct{}

func (t orderLedgerTypeSchema) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Eva order ledger type configuration.",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "ID of the order ledger type.",
				Computed:            true,
				Type:                types.Int64Type,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"name": {
				MarkdownDescription: "name of the type",
				Required:            true,
				Type:                types.StringType,
			},
			"description": {
				MarkdownDescription: "description of the type.",
				Required:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}
func (t orderLedgerTypeSchema) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return orderLedgerType{
		provider: provider,
	}, diags
}

type orderLedgerTypeData struct {
	ID          types.Int64 `tfsdk:"id"`
	Name        string      `tfsdk:"name"`
	Description string      `tfsdk:"description"`
}

type orderLedgerType struct {
	provider provider
}

func (r orderLedgerType) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data orderLedgerTypeData

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	clientResponse, err := r.provider.evaClient.CreateOrderLedgerType(ctx, eva.CreateOrderLedgerTypeRequest{
		Name:        data.Name,
		Description: data.Description,
	})

	if err != nil {
		resp.Diagnostics.AddError("Creating order ledger type failed.", fmt.Sprintf("Unable to create order ledger type, got error: %s", err))
		return
	}

	data.ID = types.Int64{Value: clientResponse.ID}

	tflog.Trace(ctx, "Created a order ledger type.")

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r orderLedgerType) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data orderLedgerTypeData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	clientResponse, err := r.provider.evaClient.ListOrderLedgerTypes(ctx)

	if err != nil {
		resp.Diagnostics.AddError("Getting order ledger type data failed.", fmt.Sprintf("Unable to get order ledger type, got error: %s", err))
		return
	}

	var orderLedgerTypeFound bool = false

	for _, orderLedgerType := range clientResponse.Result {
		if orderLedgerType.ID == data.ID.Value {
			data.Name = orderLedgerType.Name
			data.Description = orderLedgerType.Description
			orderLedgerTypeFound = true
			break
		}
	}

	if !orderLedgerTypeFound {
		resp.Diagnostics.AddError("Could not find order ledger type data.", fmt.Sprintf("Unable to get order ledger type with ID: %d", data.ID.Value))
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r orderLedgerType) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var data orderLedgerTypeData

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.provider.evaClient.UpdateOrderLedgerType(ctx, eva.UpdateOrderLedgerTypeRequest{
		ID:          data.ID.Value,
		Name:        data.Name,
		Description: data.Description,
	})

	if err != nil {
		resp.Diagnostics.AddError("Updating order ledger type failed.", fmt.Sprintf("Unable to update order ledger type, got error: %s", err))
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r orderLedgerType) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data orderLedgerTypeData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.provider.evaClient.DeleteOrderLedgerType(ctx, eva.DeleteOrderLedgerTypeRequest{
		ID: data.ID.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError("Deleting order ledger type failed.", fmt.Sprintf("Unable to delete order ledger type, got error: %s", err))
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r orderLedgerType) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}

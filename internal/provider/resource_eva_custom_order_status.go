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

type customOrderStatusType struct{}

func (t customOrderStatusType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Eva custom order status configuration.",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "ID of the custom order status.",
				Computed:            true,
				Type:                types.Int64Type,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"name": {
				MarkdownDescription: "name of the status",
				Required:            true,
				Type:                types.StringType,
			},
			"description": {
				MarkdownDescription: "description of the status.",
				Optional:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}

func (t customOrderStatusType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return customOrderStatus{
		provider: provider,
	}, diags
}

type customOrderStatusTypeData struct {
	ID          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

type customOrderStatus struct {
	provider provider
}

func (r customOrderStatus) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data customOrderStatusTypeData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	clientResponse, err := r.provider.evaClient.CreateCustomOrderStatus(ctx, eva.CreateCustomOrderStatusRequest{
		Name:        data.Name.Value,
		Description: data.Description.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError("Creating custom order status failed.", fmt.Sprintf("Unable to create custom order status, got error: %s", err))
		return
	}

	data.ID = types.Int64{Value: clientResponse.ID}

	tflog.Trace(ctx, "Created a custom order status.")

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r customOrderStatus) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data customOrderStatusTypeData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	clientResponse, err := r.provider.evaClient.ListCustomOrderStatus(ctx)

	if err != nil {
		resp.Diagnostics.AddError("Getting custom order status data failed.", fmt.Sprintf("Unable to get custom order status, got error: %s", err))
		return
	}

	var customOrderStatusFound bool = false

	for _, customOrderStatus := range clientResponse.Result {
		if customOrderStatus.ID == data.ID.Value {
			data.Name = types.String{Value: customOrderStatus.Name}
			data.Description = types.String{Value: customOrderStatus.Description}
			customOrderStatusFound = true
			break
		}
	}

	if !customOrderStatusFound {
		resp.Diagnostics.AddError("Could not find custom order status data.", fmt.Sprintf("Unable to get custom order status, got error: %s", err))
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r customOrderStatus) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var data customOrderStatusTypeData

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.provider.evaClient.UpdateCustomOrderStatus(ctx, eva.UpdateCustomOrderStatusRequest{
		ID:          data.ID.Value,
		Name:        data.Name.Value,
		Description: data.Description.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError("Updating custom order status failed.", fmt.Sprintf("Unable to update custom order status, got error: %s", err))
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r customOrderStatus) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data customOrderStatusTypeData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.provider.evaClient.DeleteCustomOrderStatus(ctx, eva.DeleteCustomOrderStatusRequest{
		ID: data.ID.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError("Deleting custom order status failed.", fmt.Sprintf("Unable to delete custom order status, got error: %s", err))
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r customOrderStatus) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}

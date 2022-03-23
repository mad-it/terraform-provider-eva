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

type organizationUnitType struct{}

func (t organizationUnitType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Eva organization unit configration.",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "Unique ID of the shop",
				Computed:            true,
				Type:                types.Int64Type,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"name": {
				MarkdownDescription: "Name of the shop",
				Required:            true,
				Type:                types.StringType,
			},
			"parent_id": {
				MarkdownDescription: "ID of the parent shop",
				Required:            true,
				Type:                types.Int64Type,
			},
			"currency_id": {
				MarkdownDescription: "Currency of the shop",
				Required:            true,
				Type:                types.StringType,
			},
			"phone_number": {
				MarkdownDescription: "Phone number of the shop",
				Optional:            true,
				Type:                types.StringType,
			},
			"email_address": {
				MarkdownDescription: "Email of the shop",
				Optional:            true,
				Type:                types.StringType,
			},
			"backend_id": {
				MarkdownDescription: "Unique reference value of the shop",
				Optional:            true,
				Type:                types.StringType,
			},
			"type": {
				MarkdownDescription: `Type of the shop. This type is a bit-wise operation.
				- None = 0
				- Shop = 1
				- WebShop = 2
				- Container = 4
				- Pickup = 8
				- Warehouse = 16
				- Country = 36 (Is always a container 32+4)
				- Franchise = 64
				- EVA = 128
				- TestOrganizationUnit = 256
				- DisableLogin = 512
				- ExternalSupplier = 1024
				- Consignment = 3072 (Is always a ExternalSupplier 2048+1024)
				- B2b = 4096
				- Region = 8196 (Is always a container 8192+4)
				- ReturnsPortal = 16384
				`,
				Optional: true,
				Type:     types.Int64Type,
			},
			"address": {
				MarkdownDescription: "Address information of the shop",
				Optional:            true,
				Attributes: tfsdk.SingleNestedAttributes(
					map[string]tfsdk.Attribute{
						"address1": {
							MarkdownDescription: "Address1 of the shop",
							Optional:            true,
							Type:                types.StringType,
						},
						"address2": {
							MarkdownDescription: "Address2 of the shop",
							Optional:            true,
							Type:                types.StringType,
						},
						"house_number": {
							MarkdownDescription: "House number of the shop",
							Optional:            true,
							Type:                types.StringType,
						},
						"zip_code": {
							MarkdownDescription: "ZipCode of the shop",
							Optional:            true,
							Type:                types.StringType,
						},
						"city": {
							MarkdownDescription: "City of the shop",
							Optional:            true,
							Type:                types.StringType,
						},
						"country_id": {
							MarkdownDescription: "Country ID of the shop",
							Optional:            true,
							Type:                types.StringType,
						},
						"latitude": {
							MarkdownDescription: "latitude of the shop",
							Optional:            true,
							Type:                types.Float64Type,
						},
						"longitude": {
							MarkdownDescription: "latitude of the shop",
							Optional:            true,
							Type:                types.Float64Type,
						},
					},
				),
			},
		},
	}, nil
}

func (t organizationUnitType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return organizationUnit{
		provider: provider,
	}, diags
}

type address struct {
	Address1    types.String  `tfsdk:"address1"`
	Address2    types.String  `tfsdk:"address2"`
	HouseNumber types.String  `tfsdk:"house_number"`
	ZipCode     types.String  `tfsdk:"zip_code"`
	City        types.String  `tfsdk:"city"`
	CountryID   types.String  `tfsdk:"country_id"`
	Latitude    types.Float64 `tfsdk:"latitude"`
	Longitude   types.Float64 `tfsdk:"longitude"`
}
type organizationUnitData struct {
	Id           types.Int64  `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	PhoneNumber  types.String `tfsdk:"phone_number"`
	EmailAddress types.String `tfsdk:"email_address"`
	CurrencyId   types.String `tfsdk:"currency_id"`
	ParentId     types.Int64  `tfsdk:"parent_id"`
	BackendId    types.String `tfsdk:"backend_id"`
	Address      *address     `tfsdk:"address"`
	Type         types.Int64  `tfsdk:"type"`
}

type organizationUnit struct {
	provider provider
}

func (r organizationUnit) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data organizationUnitData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	var organizationUnitRequest = eva.CreateOrganizationUnitRequest{
		Name:                data.Name.Value,
		PhoneNumber:         data.PhoneNumber.Value,
		BackendID:           data.BackendId.Value,
		EmailAddress:        data.EmailAddress.Value,
		ParentID:            data.ParentId.Value,
		CurrencyID:          data.CurrencyId.Value,
		CostPriceCurrencyID: data.CurrencyId.Value,
		Type:                data.Type.Value,
	}
	// check if address input is not empty
	if data.Address != nil {
		organizationUnitRequest.Latitude = data.Address.Latitude.Value
		organizationUnitRequest.Longitude = data.Address.Longitude.Value

		organizationUnitRequest.Address = &eva.Address{
			Address1:    data.Address.Address1.Value,
			Address2:    data.Address.Address2.Value,
			HouseNumber: data.Address.HouseNumber.Value,
			ZipCode:     data.Address.ZipCode.Value,
			City:        data.Address.City.Value,
			CountryID:   data.Address.CountryID.Value,
		}
	}

	client_resp, err := r.provider.evaClient.CreateOrganizationUnit(ctx, organizationUnitRequest)

	if err != nil {
		resp.Diagnostics.AddError("Creating organization unit failed.", fmt.Sprintf("Unable to create example, got error: %s", err))
		return
	}

	data.Id = types.Int64{Value: client_resp.ID}

	tflog.Trace(ctx, "Created an organization unit.")

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r organizationUnit) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data organizationUnitData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	client_resp, err := r.provider.evaClient.GetOrganizationUnitDetailed(ctx, eva.GetOrganizationUnitDetailedRequest{
		ID: data.Id.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError("Creating organization unit failed.", fmt.Sprintf("Unable to create example, got error: %s", err))
		return
	}

	data.BackendId = types.String{Value: client_resp.BackendID}
	data.CurrencyId = types.String{Value: client_resp.CurrencyID}
	data.Id = types.Int64{Value: client_resp.ID}
	data.EmailAddress = types.String{Value: client_resp.EmailAddress}
	data.PhoneNumber = types.String{Value: client_resp.PhoneNumber}
	data.Name = types.String{Value: client_resp.Name}
	data.ParentId = types.Int64{Value: client_resp.ParentID}
	data.Address = &address{
		Address1:    types.String{Value: client_resp.Address.Address1},
		Address2:    types.String{Value: client_resp.Address.Address2},
		HouseNumber: types.String{Value: client_resp.Address.HouseNumber},
		ZipCode:     types.String{Value: client_resp.Address.ZipCode},
		City:        types.String{Value: client_resp.Address.City},
		CountryID:   types.String{Value: client_resp.Address.CountryID},
		Latitude:    types.Float64{Value: client_resp.Latitude},
		Longitude:   types.Float64{Value: client_resp.Longitude},
	}
	data.Type = types.Int64{Value: client_resp.Type}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r organizationUnit) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var data organizationUnitData

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	var organizationUnitRequest = eva.UpdateOrganizationUnitRequest{
		ID:                  data.Id.Value,
		Name:                data.Name.Value,
		PhoneNumber:         data.PhoneNumber.Value,
		EmailAddress:        data.EmailAddress.Value,
		CostPriceCurrencyID: data.CurrencyId.Value,
		Type:                data.Type.Value,
	}
	// check if address input is not empty
	if data.Address != nil {
		organizationUnitRequest.Latitude = data.Address.Latitude.Value
		organizationUnitRequest.Longitude = data.Address.Longitude.Value

		organizationUnitRequest.Address = &eva.Address{
			Address1:    data.Address.Address1.Value,
			Address2:    data.Address.Address2.Value,
			HouseNumber: data.Address.HouseNumber.Value,
			ZipCode:     data.Address.ZipCode.Value,
			City:        data.Address.City.Value,
			CountryID:   data.Address.CountryID.Value,
		}
	}

	_, err := r.provider.evaClient.UpdateOrganizationUnit(ctx, organizationUnitRequest)

	if err != nil {
		resp.Diagnostics.AddError("Updating organization unit failed.", fmt.Sprintf("Unable to update OU, got error: %s", err))
		return
	}

	data.Name = types.String{Value: data.Name.Value}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r organizationUnit) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data organizationUnitData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.provider.evaClient.DeleteOrganizationUnit(ctx, eva.DeleteOrganizationUnitRequest{
		ID: data.Id.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError("Deleting organization unit failed.", fmt.Sprintf("Unable to delete OU, got error: %s", err))
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r organizationUnit) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}

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

type organizationUnitTypes int64

// This enum is copy pasted from the EVA SDK.
const (
	None organizationUnitTypes = 0
	/**
	 * A shop represents a physical store where products can be sold.
	 */
	Shop = 1
	/**
	 * A WebShop represents an online channel that allows delivery and reservation orders, but no carry out sales.
	 */
	WebShop = 2
	/**
	 * A container is an OrganizationUnit purely used to group some other OrganizationUnits to allow easier configuration.
	 */
	Container = 4
	/**
	 * Pickup can be combined with type Shop to allow reservation orders in the store.
	 */
	Pickup = 8
	/**
	 * A warehouse represents an OrganizationUnit where delivery orders can be shipped. The stock of these organizationunits can be made available for delivery orders from (web)shops.
	 */
	Warehouse = 16
	/**
	 * A Country is a special case of the Container type that represents a Country division in the OrganizationUnits structure.
	 */
	Country = 32
	/**
	 * A shop can be flagged as franchiser to allow some special flows.
	 */
	Franchise = 64
	/**
	 * The type EVA indicates that the shop is running EVA in the store. This will trigger Tasks etc that will not be generated for Shops that are not (yet) converted to running EVA. P/a non-EVA stores will receive an email for pickupordrs instead of a StockReservationTask.
	 */
	EVA = 128
	/**
	 * TestOrganizationUnit can be used to test some things in a production environment. This is not advised :warning:. These stores will be excluded from a lot of processes.
	 */
	TestOrganizationUnit = 256
	/**
	 * OrganizationUnits with DisableLogin cannot be selected in the Login process.
	 */
	DisableLogin = 512
	/**
	 * An external supplier is an organization that is not part of your internal organization structure but that you would still like to have available in EVA to for example create purchase Orders for to replenish your warehouse or stores.
	 */
	ExternalSupplier = 1024
	/**
	 * Some suppliers deliver their stock in consignment.
	 */
	Consignment = 3072
	/**
	 * For Business-to-business orders this type can be set. Orders in these organizationunits will be ex-tax.
	 */
	B2b = 4096
	/**
	 * A Region is a special case of the Container type that represents a subdivision under Country OrganizationUnits.
	 */
	Region = 8196
	/**
	 * An OrganizationUnit that is meant to be used by customers for returning Orders.
	 */
	ReturnsPortal = 16384
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
			"types": {
				MarkdownDescription: `Types of the shop. Possible values are:
				- None
				- Shop
				- WebShop
				- Container
				- Pickup
				- Warehouse
				- Country
				- Franchise
				- EVA
				- TestOrganizationUnit
				- DisableLogin
				- ExternalSupplier
				- Consignment
				- B2b
				- Region
				- ReturnsPortal
				`,
				Optional: true,
				Type:     types.SetType{},
			},
			"address": {
				MarkdownDescription: "Address information of the shop",
				Optional:            true,
				Type:                types.ObjectType{},
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

type organizationUnitData struct {
	Id           types.Int64  `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	PhoneNumber  types.String `tfsdk:"phone_number"`
	EmailAddress types.String `tfsdk:"email_address"`
	CurrencyId   types.String `tfsdk:"currency_id"`
	ParentId     types.Int64  `tfsdk:"parent_id"`
	BackendId    types.String `tfsdk:"backend_id"`
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

	client_resp, err := r.provider.evaClient.CreateOrganizationUnit(ctx, eva.CreateOrUpdateOrganizationUnitRequest{
		Name:                data.Name.Value,
		PhoneNumber:         data.PhoneNumber.Value,
		BackendID:           data.BackendId.Value,
		EmailAddress:        data.EmailAddress.Value,
		ParentID:            data.ParentId.Value,
		CurrencyID:          data.CurrencyId.Value,
		CostPriceCurrencyID: "EUR",
	})

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

	_, err := r.provider.evaClient.UpdateOrganizationUnit(ctx, eva.CreateOrUpdateOrganizationUnitRequest{
		ID:           data.Id.Value,
		Name:         data.Name.Value,
		PhoneNumber:  data.PhoneNumber.Value,
		EmailAddress: data.EmailAddress.Value,
		ParentID:     data.ParentId.Value,
		CurrencyID:   data.CurrencyId.Value,
		BackendID:    data.BackendId.Value,
	})

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

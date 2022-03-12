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

type stencilType struct{}

func (t stencilType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Eva stencil configuration.",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "ID of the stencil.",
				Computed:            true,
				Type:                types.Int64Type,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"name": {
				MarkdownDescription: "Name of the stencil.",
				Required:            true,
				Type:                types.StringType,
			},
			"organization_unit_id": {
				MarkdownDescription: "Organization that stencil belongs to",
				Optional:            true,
				Type:                types.Int64Type,
			},
			"language_id": {
				MarkdownDescription: "Language unique identifier of the stencil",
				Optional:            true,
				Type:                types.StringType,
			},
			"country_id": {
				MarkdownDescription: "Country unique identifier of the stencil",
				Optional:            true,
				Type:                types.StringType,
			},
			"header": {
				MarkdownDescription: "Header of the stencil template",
				Optional:            true,
				Type:                types.StringType,
			},
			"template": {
				MarkdownDescription: "Template of the stencil",
				Optional:            true,
				Type:                types.StringType,
			},
			"footer": {
				MarkdownDescription: "Footer of the stencil",
				Optional:            true,
				Type:                types.StringType,
			},
			"helpers": {
				MarkdownDescription: "Helper script for the stencil template",
				Optional:            true,
				Type:                types.StringType,
			},
			"type": {
				MarkdownDescription: `Type of the stencil:
				- Template = 1,
				- Partial = 2,
				- Layout = 3
				`,
				Optional: true,
				Type:     types.Int64Type,
			},
			"layout": {
				MarkdownDescription: "Type of the stencil",
				Optional:            true,
				Type:                types.StringType,
			},
			"destination": {
				MarkdownDescription: `Destination of the stencil:
				- Mail = 1,
				- Sms = 2,
				- Pdf = 4,
				- Thermal = 8,
				- Notification = 16
				`,
				Optional: true,
				Type:     types.Int64Type,
			},
			"paper_properties": {
				MarkdownDescription: "Paper's properties block configuration",
				Optional:            true,
				Attributes: tfsdk.SingleNestedAttributes(
					map[string]tfsdk.Attribute{
						"wait_for_network_idle": {
							MarkdownDescription: "Destination of the stencil",
							Required:            true,
							Type:                types.BoolType,
						},
						"wait_for_js": {
							MarkdownDescription: "Destination of the stencil",
							Required:            true,
							Type:                types.BoolType,
						},
						"format": {
							MarkdownDescription: `Paper format:
							- A3 = 1,
							- A4 = 2,
							- A5 = 3,
							- Legal = 4,
							- Letter = 5,
							- Tabloid = 6,
							- Auto = 7,
							- Ledger = 8,
							- A0 = 9,
							- A1 = 10,
							- A2 = 11,
							- A6 = 12
							`,
							Optional: true,
							Type:     types.Int64Type,
						},
						"orientation": {
							MarkdownDescription: `Paper orientation:
							- Portrait = 1,
							- Landscape = 2
							`,
							Optional: true,
							Type:     types.Int64Type,
						},
						"thermal_printer_template_type": {
							MarkdownDescription: "Paper orientation",
							Optional:            true,
							Type:                types.Int64Type,
						},
						"size": {
							MarkdownDescription: "Paper size block configuration",
							Required:            true,
							Attributes: tfsdk.SingleNestedAttributes(
								map[string]tfsdk.Attribute{
									"width": {
										MarkdownDescription: "Paper width",
										Required:            true,
										Type:                types.StringType,
									},
									"height": {
										MarkdownDescription: "Paper height",
										Required:            true,
										Type:                types.StringType,
									},
									"device_scale_factor": {
										MarkdownDescription: "Device scale factor",
										Optional:            true,
										Type:                types.Float64Type,
									},
								},
							),
						},
						"margin": {
							MarkdownDescription: "Paper's margin block configuration",
							Optional:            true,
							Attributes: tfsdk.SingleNestedAttributes(
								map[string]tfsdk.Attribute{
									"top": {
										MarkdownDescription: "Paper margin top",
										Required:            true,
										Type:                types.Int64Type,
									},
									"left": {
										MarkdownDescription: "Paper margin left",
										Required:            true,
										Type:                types.Int64Type,
									},
									"bottom": {
										MarkdownDescription: "Paper margin bottom",
										Optional:            true,
										Type:                types.Int64Type,
									},
									"right": {
										MarkdownDescription: "Paper margin right",
										Optional:            true,
										Type:                types.Int64Type,
									},
								},
							),
						},
					},
				),
			},
		},
	}, nil
}

func (t stencilType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return stencil{
		provider: provider,
	}, diags
}

type stencilTypeData struct {
	ID                 types.Int64              `tfsdk:"id"`
	Name               types.String             `tfsdk:"name"`
	OrganizationUnitID types.Int64              `tfsdk:"organization_unit_id"`
	LanguageID         types.String             `tfsdk:"language_id"`
	CountryID          types.String             `tfsdk:"country_id"`
	Header             types.String             `tfsdk:"header"`
	Template           types.String             `tfsdk:"template"`
	Footer             types.String             `tfsdk:"footer"`
	Helpers            types.String             `tfsdk:"helpers"`
	Type               types.Int64              `tfsdk:"type"`
	Layout             types.String             `tfsdk:"layout"`
	Destination        types.Int64              `tfsdk:"destination"`
	PaperProperties    *paperPropertiesTypeData `tfsdk:"paper_properties"`
}

type paperMarginTypeData struct {
	Top    types.Int64 `tfsdk:"top"`
	Left   types.Int64 `tfsdk:"left"`
	Bottom types.Int64 `tfsdk:"bottom"`
	Right  types.Int64 `tfsdk:"right"`
}

type paperSizeTypeData struct {
	Width             types.String  `tfsdk:"width"`
	Height            types.String  `tfsdk:"height"`
	DeviceScaleFactor types.Float64 `tfsdk:"device_scale_factor"`
}

type paperPropertiesTypeData struct {
	WaitForNetworkIdle         types.Bool           `tfsdk:"wait_for_network_idle"`
	WaitForJS                  types.Bool           `tfsdk:"wait_for_js"`
	Size                       *paperSizeTypeData   `tfsdk:"size"`
	Format                     types.Int64          `tfsdk:"format"`
	Orientation                types.Int64          `tfsdk:"orientation"`
	Margin                     *paperMarginTypeData `tfsdk:"margin"`
	ThermalPrinterTemplateType types.Int64          `tfsdk:"thermal_printer_template_type"`
}

type stencil struct {
	provider provider
}

func (s stencil) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data stencilTypeData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	clientResponse, err := s.provider.evaClient.CreateMessageTemplate(ctx, eva.CreateMessageTemplateRequest{
		Name:               data.Name.Value,
		OrganizationUnitID: data.OrganizationUnitID.Value,
		LanguageID:         data.LanguageID.Value,
		CountryID:          data.CountryID.Value,
		Header:             data.Header.Value,
		Template:           data.Template.Value,
		Footer:             data.Footer.Value,
		Helpers:            data.Helpers.Value,
		Type:               data.Type.Value,
		Layout:             data.Layout.Value,
		Destination:        data.Destination.Value,
		PaperProperties: &eva.PaperProperties{
			WaitForNetworkIdle:         data.PaperProperties.WaitForNetworkIdle.Value,
			WaitForJS:                  data.PaperProperties.WaitForJS.Value,
			Format:                     data.PaperProperties.Format.Value,
			Orientation:                data.PaperProperties.Orientation.Value,
			ThermalPrinterTemplateType: data.PaperProperties.ThermalPrinterTemplateType.Value,
			Size: eva.PaperSize{
				Width:             data.PaperProperties.Size.Width.Value,
				Height:            data.PaperProperties.Size.Height.Value,
				DeviceScaleFactor: data.PaperProperties.Size.DeviceScaleFactor.Value,
			},
			Margin: eva.PaperMargin{
				Top:    data.PaperProperties.Margin.Top.Value,
				Left:   data.PaperProperties.Margin.Left.Value,
				Bottom: data.PaperProperties.Margin.Bottom.Value,
				Right:  data.PaperProperties.Margin.Right.Value,
			},
		},
	})

	if err != nil {
		resp.Diagnostics.AddError("Creating stencil unit failed.", fmt.Sprintf("Unable to create stencil, got error: %s", err))
		return
	}

	data.ID = types.Int64{Value: clientResponse.ID}

	tflog.Trace(ctx, "Created a stencil")

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (s stencil) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data stencilTypeData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	clientResponse, err := s.provider.evaClient.GetMessageTemplateByID(ctx, eva.GetMessageTemplateByIDRequest{
		ID: data.ID.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError("Getting stencil data failed.", fmt.Sprintf("Unable to get stencil, got error: %s", err))
		return
	}

	data.Name = types.String{Value: clientResponse.Name}
	data.OrganizationUnitID = types.Int64{Value: clientResponse.OrganizationUnitID}
	data.LanguageID = types.String{Value: clientResponse.LanguageID}
	data.CountryID = types.String{Value: clientResponse.CountryID}
	data.Header = types.String{Value: clientResponse.Header}
	data.Template = types.String{Value: clientResponse.Template}
	data.Footer = types.String{Value: clientResponse.Footer}
	data.Helpers = types.String{Value: clientResponse.Helpers}
	data.Type = types.Int64{Value: clientResponse.Type}
	data.Layout = types.String{Value: clientResponse.Layout}
	data.Destination = types.Int64{Value: clientResponse.Destination}
	data.PaperProperties = &paperPropertiesTypeData{
		WaitForNetworkIdle:         types.Bool{Value: clientResponse.PaperProperties.WaitForNetworkIdle},
		WaitForJS:                  types.Bool{Value: clientResponse.PaperProperties.WaitForJS},
		Format:                     types.Int64{Value: clientResponse.PaperProperties.Format},
		Orientation:                types.Int64{Value: clientResponse.PaperProperties.Orientation},
		ThermalPrinterTemplateType: types.Int64{Value: clientResponse.PaperProperties.ThermalPrinterTemplateType},
		Size: &paperSizeTypeData{
			Width:             types.String{Value: clientResponse.PaperProperties.Size.Width},
			Height:            types.String{Value: clientResponse.PaperProperties.Size.Height},
			DeviceScaleFactor: types.Float64{Value: clientResponse.PaperProperties.Size.DeviceScaleFactor},
		},
		Margin: &paperMarginTypeData{
			Top:    types.Int64{Value: clientResponse.PaperProperties.Margin.Top},
			Left:   types.Int64{Value: clientResponse.PaperProperties.Margin.Left},
			Bottom: types.Int64{Value: clientResponse.PaperProperties.Margin.Bottom},
			Right:  types.Int64{Value: clientResponse.PaperProperties.Margin.Right},
		},
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (s stencil) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var plan stencilTypeData

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := s.provider.evaClient.UpdateMessageTemplate(ctx, eva.UpdateMessageTemplateRequest{
		ID:                 plan.ID.Value,
		Name:               plan.Name.Value,
		OrganizationUnitID: plan.OrganizationUnitID.Value,
		LanguageID:         plan.LanguageID.Value,
		CountryID:          plan.CountryID.Value,
		Header:             plan.Header.Value,
		Template:           plan.Template.Value,
		Footer:             plan.Footer.Value,
		Helpers:            plan.Helpers.Value,
		Layout:             plan.Layout.Value,
		Destination:        plan.Destination.Value,
		PaperProperties: &eva.PaperProperties{
			WaitForNetworkIdle:         plan.PaperProperties.WaitForNetworkIdle.Value,
			WaitForJS:                  plan.PaperProperties.WaitForJS.Value,
			Format:                     plan.PaperProperties.Format.Value,
			Orientation:                plan.PaperProperties.Orientation.Value,
			ThermalPrinterTemplateType: plan.PaperProperties.ThermalPrinterTemplateType.Value,
			Size: eva.PaperSize{
				Width:             plan.PaperProperties.Size.Width.Value,
				Height:            plan.PaperProperties.Size.Height.Value,
				DeviceScaleFactor: plan.PaperProperties.Size.DeviceScaleFactor.Value,
			},
			Margin: eva.PaperMargin{
				Top:    plan.PaperProperties.Margin.Top.Value,
				Left:   plan.PaperProperties.Margin.Left.Value,
				Bottom: plan.PaperProperties.Margin.Bottom.Value,
				Right:  plan.PaperProperties.Margin.Right.Value,
			},
		},
	})

	if err != nil {
		resp.Diagnostics.AddError("Updating stencil unit failed.", fmt.Sprintf("Unable to update stencil, got error: %s", err))
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (s stencil) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data stencilTypeData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := s.provider.evaClient.DeleteMessageTemplate(ctx, eva.DeleteMessageTemplateRequesst{
		ID: data.ID.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError("Deleting stencil unit failed.", fmt.Sprintf("Unable to delete stencil, got error: %s", err))
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r stencil) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/mad-it/terraform-provider-eva/internal/eva"
)

type provider struct {
	evaClient eva.Client

	// configured is set to true at the end of the Configure method.
	// This can be used in Resource and DataSource implementations to verify
	// that the provider was previously configured.
	configured bool

	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

type providerData struct {
	Endpoint types.String `tfsdk:"endpoint"`
	Token    types.String `tfsdk:"token"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {

	var data providerData
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Token.Null && (data.Username.Null || data.Password.Null) {
		resp.Diagnostics.AddError("No valid credentials provided.", "Both token and username/password are not filed. Please provide one of these.")
		return
	}

	p.evaClient = *eva.NewClient(data.Endpoint.Value)

	if !data.Token.Null {
		p.evaClient.SetAuthorizationHeader(data.Token.Value)
	} else {

		err := p.evaClient.Login(ctx, eva.LoginCredentials{Username: data.Username.Value, Password: data.Password.Value})

		if err != nil {

			resp.Diagnostics.AddError(
				"Login to EVA failed.",
				fmt.Sprintf("An error ocurred when logging in to EVA. Error was: %s", err),
			)
			return
		}

	}

	p.configured = true
}

func (p *provider) GetResources(ctx context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{
		"eva_organization_unit": organizationUnitType{},
		"eva_setting":           settingType{},
		"eva_cookbook":          cookbookType{},
	}, nil
}

func (p *provider) GetDataSources(ctx context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return map[string]tfsdk.DataSourceType{}, nil
}

func (p *provider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"endpoint": {
				MarkdownDescription: "The base URL of the EVA API.",
				Required:            true,
				Type:                types.StringType,
			},
			"token": {
				MarkdownDescription: "The base URL of the EVA API.",
				Optional:            true,
				Type:                types.StringType,
			},
			"username": {
				MarkdownDescription: "Username used to log in to EVA.",
				Optional:            true,
				Type:                types.StringType,
			},
			"password": {
				MarkdownDescription: "Password used to log in to EVA.",
				Optional:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}

func New(version string) func() tfsdk.Provider {
	return func() tfsdk.Provider {
		return &provider{
			version: version,
		}
	}
}

func convertProviderType(in tfsdk.Provider) (provider, diag.Diagnostics) {
	var diags diag.Diagnostics

	p, ok := in.(*provider)

	if !ok {
		diags.AddError(
			"Unexpected Provider Instance Type",
			fmt.Sprintf("While creating the data source or resource, an unexpected provider type (%T) was received. This is always a bug in the provider code and should be reported to the provider developers.", p),
		)
		return provider{}, diags
	}

	if p == nil {
		diags.AddError(
			"Unexpected Provider Instance Type",
			"While creating the data source or resource, an unexpected empty provider instance was received. This is always a bug in the provider code and should be reported to the provider developers.",
		)
		return provider{}, diags
	}

	return *p, diags
}

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
	client eva.Client

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
	url      types.String `tfsdk:"url"`
	username types.String `tfsdk:"username"`
	password types.String `tfsdk:"password"`
}

func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
	var data providerData
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	fmt.Println("provider started")

	if resp.Diagnostics.HasError() {
		return
	}

	p.client = *eva.NewClient(data.url.Value)
	fmt.Println("client created")

	err := p.client.Login(eva.LoginCredentials{Username: data.username.Value, Password: data.password.Value})

	fmt.Println("logged in")
	if err != nil {

		diags.AddError(
			"Login to EVA failed.",
			fmt.Sprintf("An error ocurred when logging in to EVA. Error was: %s", err),
		)
		return
	}

	p.configured = true
}

func (p *provider) GetResources(ctx context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{
		"scaffolding_example": exampleResourceType{},
	}, nil
}

func (p *provider) GetDataSources(ctx context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return map[string]tfsdk.DataSourceType{}, nil
}

func (p *provider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"url": {
				MarkdownDescription: "The base URL of the EVA API.",
				Validators:          []tfsdk.AttributeValidator{
					// TODO add regex to validate its https://api.<eva-url>
				},
				Required: true,
				Type:     types.StringType,
			},
			"username": {
				MarkdownDescription: "Username used to log in to EVA.",
				Required:            true,
				Type:                types.StringType,
			},
			"password": {
				MarkdownDescription: "Password used to log in to EVA.",
				Required:            true,
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

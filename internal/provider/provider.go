package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &dxProvider{}

type dxProvider struct {
	Version string
}

type dxPrefix string

type dxProviderModel struct {
	Prefix      types.String `tfsdk:"prefix"`
	Domain      types.String `tfsdk:"domain"`
	Environment types.String `tfsdk:"environment"`
	Location    types.String `tfsdk:"location"`
}

// New crea una nuova istanza del provider
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &dxProvider{
			Version: version,
		}
	}
}

func (p *dxProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "dx"
	resp.Version = p.Version
}

func (p *dxProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The dx provider is used to generate and manage naming of Azure resources.",
		Attributes: map[string]schema.Attribute{
			"prefix": schema.StringAttribute{
				Optional:    true,
				Description: "Prefix that define the repository domain",
				Validators: []validator.String{
					stringvalidator.LengthBetween(2, 2),
				},
			},
			"domain": schema.StringAttribute{
				Optional:    true,
				Description: "The team domain name",
			},
			"environment": schema.StringAttribute{
				Optional:    true,
				Description: "Environment where the resources will be deployed",
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"d", "u", "p"}...),
				},
			},
			"location": schema.StringAttribute{
				Optional:    true,
				Description: "Location where the resources will be deployed",
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"weu", "itn", "westeurope", "italynorth"}...),
				},
			},
		},
	}
}

func (p *dxProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config dxProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Resources

func (p *dxProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewAvailableSubnetCidrResource,
	}
}

// DataSources
func (p *dxProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return nil
}

// Functions

func (p *dxProvider) Functions(_ context.Context) []func() function.Function {
	return []func() function.Function{
		NewResourceNameFunction,
	}
}

func NewResourceNameFunction() function.Function {
	return &resourceNameFunction{}
}

package provider

import (
	"context"
	"github.com/davidalpert/go-contentstack/v1/management"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure ContentStackProvider satisfies various provider interfaces.
var _ provider.Provider = &ContentStackProvider{}

// ContentStackProvider defines the provider implementation.
type ContentStackProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
	commit  string
	date    string
}

// ContentStackProviderModel describes the provider data model.
type ContentStackProviderModel struct {
	Host            types.String `tfsdk:"host"`
	ApiKey          types.String `tfsdk:"api_key"`
	ManagementToken types.String `tfsdk:"management_token"`
	Debug           types.Bool   `tfsdk:"debug"`
}

func (p *ContentStackProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "contentstack"
	resp.Version = p.version
}

func (p *ContentStackProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: `Base URL
- US (North America, or NA): https://api.contentstack.io/
- Europe (EU): https://eu-api.contentstack.com/
- Azure NA: https://azure-na-api.contentstack.com/
`,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "An API Key which uniquely identifies the stack which this provider will configure.",
				Optional:            true,
				Sensitive:           true,
			},
			"management_token": schema.StringAttribute{
				MarkdownDescription: "Management Tokens are stack-level tokens, with no users attached to them. They can do everything that authtokens can do. Since they are not personal tokens, no role-specific permissions are applicable to them. It is recommended to use these tokens for automation scripts, third-party app integrations, and for Single Sign On (SSO)-enabled organizations.",
				Optional:            true,
				Sensitive:           true,
			},
			"debug": schema.BoolAttribute{
				MarkdownDescription: "enable debug logs for the ContentStack API client",
				Optional:            true,
			},
		},
	}
}

func (p *ContentStackProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data ContentStackProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	// if data.Endpoint.IsNull() { /* ... */ }

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("CONTENTSTACK_HOST")
	apiKey := os.Getenv("CONTENTSTACK_API_KEY")
	managementToken := os.Getenv("CONTENTSTACK_MANAGEMENT_TOKEN")

	if !data.Host.IsNull() {
		host = data.Host.ValueString()
	}

	if !data.ApiKey.IsNull() {
		apiKey = data.ApiKey.ValueString()
	}

	if !data.ManagementToken.IsNull() {
		managementToken = data.ManagementToken.ValueString()
	}

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing ContentStack API Host",
			"The provider cannot create the ContentStack API client as there is a missing or empty value for the ContentStack API host. "+
				"Set the host value in the configuration or use the CONTENTSTACK_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing ContentStack API Key",
			"The provider cannot create the ContentStack API client as there is a missing or empty value for the ContentStack API Key. "+
				"Set the API Key value in the configuration or use the CONTENTSTACK_API_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if managementToken == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("management_token"),
			"Missing ContentStack Management Token",
			"The provider cannot create the ContentStack API client as there is a missing or empty value for the ContentStack Management Token. "+
				"Set the Management Token value in the configuration or use the CONTENTSTACK_MANAGEMENT_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new ContentStack API client using the configuration values
	debugClient := false
	if !data.Debug.IsUnknown() {
		debugClient = data.Debug.ValueBool()
	}
	client, err := management.NewClient(&management.Configuration{
		Host:      host,
		Key:       apiKey,
		Token:     managementToken,
		Debug:     debugClient,
		UserAgent: "terraform-provider-contentstacktypes",
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create ContentStack API Client",
			"An unexpected error occurred when creating the ContentStack API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"ContentStack Client Error: "+err.Error(),
		)
		return
	}

	// Make the ContentStack API client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *ContentStackProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewGlobalFieldResource,
	}
}

func (p *ContentStackProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewGlobalFieldDataSource,
	}
}

func New(version, commit, date string) func() provider.Provider {
	return func() provider.Provider {
		return &ContentStackProvider{
			version: version,
			commit:  commit,
			date:    date,
		}
	}
}

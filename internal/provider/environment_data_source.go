package provider

import (
	"context"
	"fmt"
	"github.com/davidalpert/go-contentstack/v1/management"
	cschema "github.com/davidalpert/go-contentstack/v1/schema"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &EnvironmentDataSource{}

func NewEnvironmentDataSource() datasource.DataSource {
	return &EnvironmentDataSource{}
}

// EnvironmentDataSource defines the data source implementation.
type EnvironmentDataSource struct {
	client *management.Client
}

// EnvironmentDataSourceModel describes the data source data model.
type EnvironmentDataSourceModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	URLs          types.Map    `tfsdk:"urls"`
	UID           types.String `tfsdk:"uid"`
	CreatedAt     types.String `tfsdk:"created_at"`
	UpdatedAt     types.String `tfsdk:"updated_at"`
	Version       types.Int64  `tfsdk:"version"`
	DeployContent types.Bool   `tfsdk:"deploy_content"`
}

func (d *EnvironmentDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environment"
}

func (d *EnvironmentDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Environment data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "internal terraform resource id (matches the name when the Environment has been created/imported)",
				Computed:            true,
			},
			"uid": schema.StringAttribute{
				MarkdownDescription: "internal contentstack identifier",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "name of the Environment",
				Required:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "created_at of the Global Field",
				Computed:            true,
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "updated_at of the Global Field",
				Computed:            true,
			},
			"version": schema.Int64Attribute{
				MarkdownDescription: "version number of the Environment",
				Computed:            true,
			},
			"deploy_content": schema.BoolAttribute{
				MarkdownDescription: "deploy_content",
				Computed:            true,
			},
			"urls": schema.MapAttribute{ // urls by locale
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "urls by locale",
			},
		},
	}
}

func (d *EnvironmentDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*management.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *management.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *EnvironmentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data EnvironmentDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	g, err := d.client.GetOnePublishingEnvironment(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read Environment %#v, got error: %s", data.UID.ValueString(), err))
		return
	}

	data.Name = types.StringValue(g.Name)
	data.UID = types.StringValue(g.UID)
	data.Version = types.Int64Value(int64(g.Version))
	data.CreatedAt = types.StringValue("TBD")
	data.UpdatedAt = types.StringValue("TBD")
	data.DeployContent = types.BoolValue(g.DeployContent)

	urls, dg := flattenUrlsByLocale(g.Urls)
	resp.Diagnostics.Append(dg...)

	uu, dg := types.MapValue(types.StringType, urls)
	resp.Diagnostics.Append(dg...)

	data.URLs = uu

	data.ID = data.Name

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read an environment", map[string]interface{}{
		"uid":  g.UID,
		"name": g.Name,
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenUrlsByLocale(urls []cschema.LocaleURLPair) (map[string]attr.Value, diag.Diagnostics) {
	result := map[string]attr.Value{}
	innerDiagnostics := diag.Diagnostics{}
	for _, u := range urls {
		result[u.Locale] = types.StringValue(u.Url)
	}
	return result, innerDiagnostics
}

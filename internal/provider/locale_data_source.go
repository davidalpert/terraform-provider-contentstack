package provider

import (
	"context"
	"fmt"
	"github.com/davidalpert/go-contentstack/v1/management"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &LocaleDataSource{}

func NewLocaleDataSource() datasource.DataSource {
	return &LocaleDataSource{}
}

// LocaleDataSource defines the data source implementation.
type LocaleDataSource struct {
	client *management.Client
}

// LocaleDataSourceModel describes the data source data model.
type LocaleDataSourceModel struct {
	ID                 types.String `tfsdk:"id"`
	UID                types.String `tfsdk:"uid"`
	Code               types.String `tfsdk:"code"`
	Name               types.String `tfsdk:"name"`
	FallbackLocaleCode types.String `tfsdk:"fallback_locale"`

	CreatedBy types.String `tfsdk:"created_by"`
	CreatedAt types.String `tfsdk:"created_at"`
	Version   types.Int64  `tfsdk:"_version"`
	UpdatedBy types.String `tfsdk:"updated_by"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

func (d *LocaleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_locale"
}

func (d *LocaleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Locale data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "internal terraform resource id (matches the code when the Locale has been created/imported)",
				Computed:            true,
			},
			"uid": schema.StringAttribute{
				MarkdownDescription: "internal contentstack identifier",
				Computed:            true,
			},
			"code": schema.StringAttribute{
				MarkdownDescription: "code for the Locale",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "name of the Locale",
				Optional:            true,
			},
			"fallback_locale": schema.StringAttribute{
				MarkdownDescription: "fallback code of the Locale",
				Optional:            true,
			},
			"created_by": schema.StringAttribute{
				MarkdownDescription: "created_by of the Locale",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "created_at of the Locale",
				Computed:            true,
			},
			"version": schema.Int64Attribute{
				MarkdownDescription: "version number of the Locale",
				Computed:            true,
			},
			"updated_by": schema.StringAttribute{
				MarkdownDescription: "created_by of the Locale",
				Computed:            true,
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "updated_at of the Locale",
				Computed:            true,
			},
		},
	}
}

func (d *LocaleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *LocaleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data LocaleDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	g, err := d.client.GetOneLocale(data.Code.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read Locale %#v, got error: %s", data.UID.ValueString(), err))
		return
	}

	data.Name = types.StringValue(g.Name)
	data.UID = types.StringValue(g.Uid)
	data.Version = types.Int64Value(int64(g.Version))
	data.CreatedAt = types.StringValue("TBD")
	data.UpdatedAt = types.StringValue("TBD")
	data.FallbackLocaleCode = types.StringValue(g.FallbackLocaleCode)

	data.ID = data.Name

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

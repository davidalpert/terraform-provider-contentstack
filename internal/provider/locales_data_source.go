package provider

import (
	"context"
	"fmt"
	"github.com/davidalpert/go-contentstack/v1/management"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &LocalesDataSource{}

func NewLocalesDataSource() datasource.DataSource {
	return &LocalesDataSource{}
}

// LocalesDataSource defines the data source implementation.
type LocalesDataSource struct {
	client *management.Client
}

// LocalesDataSourceModel describes the data source data model.
type LocalesDataSourceModel struct {
	Locales types.Map `tfsdk:"locales_by_code"`
}

func (d *LocalesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_locales"
}

func localesMapElementType() attr.Type {
	return types.ObjectType{
		AttrTypes: localesMapElementTypeSchema(),
	}
}

func localesMapElementTypeSchema() map[string]attr.Type {
	return map[string]attr.Type{
		"uid":                  types.StringType,
		"code":                 types.StringType,
		"name":                 types.StringType,
		"fallback_locale_code": types.StringType,
		"created_by":           types.StringType,
		"created_at":           types.StringType,
		"updated_by":           types.StringType,
		"updated_at":           types.StringType,
	}
}

func (d *LocalesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Data source to fetch all locales",

		Attributes: map[string]schema.Attribute{
			"locales_by_code": schema.MapAttribute{
				Computed:            true,
				ElementType:         localesMapElementType(),
				MarkdownDescription: "exposes a list of all locales defined in the provider's stack",
			},
		},
	}
}

func (d *LocalesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *LocalesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data LocalesDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ll, err := d.client.GetAllLocales()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read Locales from the configued stack: %s", err))
		return
	}

	values := map[string]attr.Value{}
	for _, l := range ll {
		ovv := map[string]attr.Value{
			"uid":                  types.StringValue(l.Uid),
			"code":                 types.StringValue(l.Code),
			"name":                 types.StringValue(l.Name),
			"fallback_locale_code": types.StringValue(l.FallbackLocaleCode),
			"created_by":           types.StringValue(l.CreatedBy),
			"created_at":           types.StringValue(l.CreatedAt),
			"updated_by":           types.StringValue(l.UpdatedBy),
			"updated_at":           types.StringValue(l.UpdatedAt),
		}

		ov, diag := types.ObjectValue(localesMapElementTypeSchema(), ovv)
		resp.Diagnostics.Append(diag...)
		values[l.Code] = ov
	}

	mv, diag := types.MapValue(localesMapElementType(), values)
	resp.Diagnostics.Append(diag...)
	data.Locales = mv

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

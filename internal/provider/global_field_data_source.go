package provider

import (
	"context"
	"fmt"
	"github.com/davidalpert/go-contentstack/v1/management"
	cschema "github.com/davidalpert/go-contentstack/v1/schema"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &GlobalFieldDataSource{}

func NewGlobalFieldDataSource() datasource.DataSource {
	return &GlobalFieldDataSource{}
}

// GlobalFieldDataSource defines the data source implementation.
type GlobalFieldDataSource struct {
	client *management.Client
}

// GlobalFieldDataSourceModel describes the data source data model.
type GlobalFieldDataSourceModel struct {
	CreatedAt   types.String                 `tfsdk:"created_at"`
	Description types.String                 `tfsdk:"description"`
	Fields      []SchemaFieldDataSourceModel `tfsdk:"field"`
	ID          types.String                 `tfsdk:"id"`
	Title       types.String                 `tfsdk:"title"`
	UID         types.String                 `tfsdk:"uid"`
	UpdatedAt   types.String                 `tfsdk:"updated_at"`
}

type SchemaFieldDataSourceModel struct {
	//Blocks         []BlockSet     `tfsdk:"blocks,omitempty"`
	DataType    string             `tfsdk:"data_type"`
	DisplayName string             `tfsdk:"display_name"`
	DisplayType *string            `tfsdk:"display_type,omitempty"`
	Enum        *cschema.EnumField `tfsdk:"enum,omitempty"`
	//ErrorMessages  *ErrorMessages `tfsdk:"error_messages,omitempty"`
	//FieldMetadata  FieldMetadata  `tfsdk:"field_metadata"`
	Format         *string `tfsdk:"format,omitempty"`
	InbuiltModel   *bool   `tfsdk:"inbuilt_model,omitempty"`
	Indexed        *bool   `tfsdk:"indexed,omitempty"`
	Mandatory      bool    `tfsdk:"mandatory"`
	Multiple       bool    `tfsdk:"multiple"`
	NonLocalizable *bool   `tfsdk:"non_localizable,omitempty"`
	//ReferenceTo    StrArray       `tfsdk:"reference_to,omitempty"`
	Uid    string `tfsdk:"uid"`
	Unique *bool  `tfsdk:"unique,omitempty"`
}

func BuildComputedFieldsSchema() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"data_type": schema.StringAttribute{
					MarkdownDescription: "data type of the field",
					Computed:            true,
				},
				"display_name": schema.StringAttribute{
					MarkdownDescription: "display name of the field",
					Computed:            true,
				},
				"display_type": schema.StringAttribute{
					MarkdownDescription: "display type of the field",
					Computed:            true,
				},
				"uid": schema.StringAttribute{
					MarkdownDescription: "uid of the field",
					Computed:            true,
				},
			},
		},
		Computed:            true,
		MarkdownDescription: "field schema of the Global Field",
	}
}

func (d *GlobalFieldDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_global_field"
}

func (d *GlobalFieldDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Example data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "internal terraform resource id (matches the uid when the global field has been created/imported)",
				Computed:            true,
			},
			"uid": schema.StringAttribute{
				MarkdownDescription: "uid of the Global Field",
				Required:            true,
			},
			"title": schema.StringAttribute{
				MarkdownDescription: "title of the Global Field",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "description of the Global Field",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "created_at of the Global Field",
				Computed:            true,
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "updated_at of the Global Field",
				Computed:            true,
			},
			"field": BuildComputedFieldsSchema(),
		},
	}
}

func (d *GlobalFieldDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *GlobalFieldDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data GlobalFieldDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	g, err := d.client.GetOneGlobalField(data.UID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read GlobalField %#v, got error: %s", data.UID.ValueString(), err))
		return
	}

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	data.ID = types.StringValue(g.UID)
	data.Title = types.StringValue(g.Title)
	data.Description = types.StringValue(g.Description)
	data.CreatedAt = types.StringValue(g.CreatedAt)
	data.UpdatedAt = types.StringValue(g.UpdatedAt)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source", map[string]interface{}{
		"uid": g.UID,
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

package provider

import (
	"context"
	"fmt"
	"github.com/davidalpert/go-contentstack/v1/management"
	cschema "github.com/davidalpert/go-contentstack/v1/schema"
	myboolplanmodifiers "github.com/davidalpert/terraform-provider-contentstack-admin/internal/tfutils/boolplanmodifier"
	mystringplanmodifiers "github.com/davidalpert/terraform-provider-contentstack-admin/internal/tfutils/stringplanmodifiers"
	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &GlobalFieldResource{}
var _ resource.ResourceWithImportState = &GlobalFieldResource{}

func NewGlobalFieldResource() resource.Resource {
	return &GlobalFieldResource{}
}

// GlobalFieldResource defines the resource implementation.
type GlobalFieldResource struct {
	client *management.Client
}

// GlobalFieldResourceModel describes the resource data model.
type GlobalFieldResourceModel struct {
	Description types.String                          `tfsdk:"description"`
	Fields      []GlobalFieldSchemaFieldResourceModel `tfsdk:"fields"`
	ID          types.String                          `tfsdk:"id"`
	Title       types.String                          `tfsdk:"title"`
	UID         types.String                          `tfsdk:"uid"`
}

func (data *GlobalFieldResourceModel) Update(g *cschema.GlobalField) {
	tflog.Warn(context.Background(), "Update Description", map[string]interface{}{
		"api":   g.Description,
		"model": data.Description.String(),
	})
	data.Description = types.StringValue(g.Description)
	data.Fields = make([]GlobalFieldSchemaFieldResourceModel, len(g.Schema))
	data.ID = types.StringValue(g.UID)
	data.Title = types.StringValue(g.Title)
	data.UID = types.StringValue(g.UID)

	for i, f := range g.Schema {
		data.Fields[i] = GlobalFieldSchemaFieldResourceModel{}
		data.Fields[i].Update(f)
	}
}

func (data *GlobalFieldResourceModel) Export() *cschema.GlobalField {
	g := &cschema.GlobalField{
		Description: data.Description.ValueString(),
		Schema:      make([]cschema.Field, len(data.Fields)),
		Title:       data.Title.ValueString(),
		UID:         data.UID.ValueString(),
	}

	for i, fieldData := range data.Fields {
		g.Schema[i] = fieldData.Export()
	}

	return g
}

type GlobalFieldSchemaFieldResourceModel struct {
	DataType    types.String `tfsdk:"data_type"`
	Description types.String `tfsdk:"description"`
	DisplayName types.String `tfsdk:"display_name"`
	DefaultText types.String `tfsdk:"default_text"`
	DefaultBool types.Bool   `tfsdk:"default_bool"`
	Format      types.String `tfsdk:"format"`
	Mandatory   types.Bool   `tfsdk:"mandatory"`
	Multiple    types.Bool   `tfsdk:"multiple"`
	Placeholder types.String `tfsdk:"placeholder"`
	Instruction types.String `tfsdk:"instruction"`
	Uid         types.String `tfsdk:"uid"`
	Unique      types.Bool   `tfsdk:"unique"`
}

func (data *GlobalFieldSchemaFieldResourceModel) Update(f cschema.Field) {
	data.DataType = types.StringValue(f.DataType)
	data.Description = types.StringValue(f.FieldMetadata.Description)
	data.DisplayName = types.StringValue(f.DisplayName)
	if f.Format != nil {
		data.Format = types.StringValue(*f.Format)
	} else {
		data.Format = types.StringNull()
	}
	if f.FieldMetadata.Placeholder != nil {
		data.Placeholder = types.StringValue(*f.FieldMetadata.Placeholder)
	} else {
		data.Placeholder = types.StringNull()
	}
	if f.FieldMetadata.Instruction != nil {
		data.Instruction = types.StringValue(*f.FieldMetadata.Instruction)
	} else {
		data.Instruction = types.StringNull()
	}
	if f.FieldMetadata.DefaultValue != nil {
		switch v := f.FieldMetadata.DefaultValue.(type) {
		case string:
			data.DefaultText = types.StringValue(v)
		case bool:
			data.DefaultBool = types.BoolValue(v)
		default:
			// TODO: show warning: unsupported
		}
	}
	data.Mandatory = types.BoolValue(f.Mandatory)
	data.Multiple = types.BoolValue(f.Multiple)
	data.Uid = types.StringValue(f.Uid)
	if f.Unique != nil {
		data.Unique = types.BoolValue(*f.Unique)
	} else {
		data.Unique = types.BoolNull()
	}
}

func (data *GlobalFieldSchemaFieldResourceModel) Export() cschema.Field {
	field := cschema.Field{
		DataType:      data.DataType.ValueString(),
		Uid:           data.Uid.ValueString(),
		FieldMetadata: cschema.FieldMetadata{},
	}

	if !data.DisplayName.IsNull() {
		field.DisplayName = data.DisplayName.ValueString()
	} else {
		field.DisplayName = field.Uid
	}

	if !data.DefaultBool.IsNull() {
		field.FieldMetadata.DefaultValue = data.DefaultBool.ValueBool()
	} else if !data.DefaultText.IsNull() {
		field.FieldMetadata.DefaultValue = data.DefaultText.ValueString()
	}

	if !data.Description.IsNull() {
		field.FieldMetadata.Description = data.Description.ValueString()
	}

	if !data.Placeholder.IsNull() {
		field.FieldMetadata.Placeholder = cschema.StrPtr(data.Placeholder.ValueString())
	}

	if !data.Instruction.IsNull() {
		field.FieldMetadata.Instruction = cschema.StrPtr(data.Instruction.ValueString())
	}

	if !data.Mandatory.IsNull() {
		field.Mandatory = data.Mandatory.ValueBool()
	}

	if !data.Multiple.IsNull() {
		field.Multiple = data.Multiple.ValueBool()
	}

	if !data.Unique.IsNull() {
		field.Unique = cschema.BoolPtr(data.Unique.ValueBool())
	}

	return field
}

func (r *GlobalFieldResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_global_field"
}

func (r *GlobalFieldResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "GlobalField resource",

		Attributes: map[string]schema.Attribute{
			"description": schema.StringAttribute{
				MarkdownDescription: "description of the GlobalField",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					mystringplanmodifiers.DefaultValue(""),
				},
			},
			"fields": BuildFieldsSchema(),
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "GlobalField identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"title": schema.StringAttribute{
				MarkdownDescription: "title of the GlobalField",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					mystringplanmodifiers.DefaultValueCopiedFromAnotherField("uid"),
				},
			},
			"uid": schema.StringAttribute{
				MarkdownDescription: "uid of the GlobalField",
				Required:            true,
			},
		},
	}
}

func BuildFieldsSchema() schema.SetNestedAttribute {
	return schema.SetNestedAttribute{
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"data_type": schema.StringAttribute{
					MarkdownDescription: `data type of the field:
  - text
  - boolean
  - number
  - file
  - link
  - json
  - isodate
`,
					Required: true,
					Validators: []validator.String{
						stringvalidator.OneOf("text", "boolean", "number", "file", "link", "json", "isodate"),
					},
				},
				"description": schema.StringAttribute{
					MarkdownDescription: "description of the field",
					Optional:            true,
					Computed:            true,
					Default:             stringdefault.StaticString(""),
				},
				"default_bool": schema.BoolAttribute{
					MarkdownDescription: "default boolean value for the field",
					Optional:            true,
					Computed:            true,
					//Default: booldefault.StaticBool(false),
					Validators: []validator.Bool{
						boolvalidator.ConflictsWith(
							path.MatchRelative().AtParent().AtName("default_text"),
						),
					},
				},
				"default_text": schema.StringAttribute{
					MarkdownDescription: "default boolean value for the field",
					Optional:            true,
					Computed:            true,
					//Default:             stringdefault.StaticString(""),
					Validators: []validator.String{
						stringvalidator.ConflictsWith(
							path.MatchRelative().AtParent().AtName("default_bool"),
						),
					},
				},
				"display_name": schema.StringAttribute{
					MarkdownDescription: "display name of the field",
					Optional:            true,
					Computed:            true,
					PlanModifiers: []planmodifier.String{
						mystringplanmodifiers.DefaultValueCopiedFromAnotherField("uid"),
					},
				},
				"format": schema.StringAttribute{
					MarkdownDescription: "format of the field",
					Optional:            true,
				},
				"placeholder": schema.StringAttribute{
					MarkdownDescription: "placeholder text for the field",
					Optional:            true,
					Computed:            true,
					PlanModifiers: []planmodifier.String{
						mystringplanmodifiers.DefaultValue(""),
					},
				},
				"instruction": schema.StringAttribute{
					MarkdownDescription: "instruction text for the field",
					Optional:            true,
					Computed:            true,
					PlanModifiers: []planmodifier.String{
						mystringplanmodifiers.DefaultValue(""),
					},
				},
				//"display_type": schema.StringAttribute{
				//	MarkdownDescription: "display type of the field",
				//	Optional:            true,
				//},
				"uid": schema.StringAttribute{
					MarkdownDescription: "uid of the field",
					Required:            true,
				},
				"mandatory": schema.BoolAttribute{
					MarkdownDescription: "is this field mandatory",
					Optional:            true,
					Computed:            true,
					PlanModifiers: []planmodifier.Bool{
						myboolplanmodifiers.DefaultValue(false),
					},
				},
				"multiple": schema.BoolAttribute{
					MarkdownDescription: "can this field be used multiple times",
					Optional:            true,
					Computed:            true,
					PlanModifiers: []planmodifier.Bool{
						myboolplanmodifiers.DefaultValue(false),
					},
				},
				"unique": schema.BoolAttribute{
					MarkdownDescription: "must this field be unique",
					Optional:            true,
					Computed:            true,
					PlanModifiers: []planmodifier.Bool{
						myboolplanmodifiers.DefaultValue(false),
					},
				},
			},
		},
		Required:            true,
		MarkdownDescription: "field schema of the Global Field",
	}
}

func (r *GlobalFieldResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*management.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *management.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *GlobalFieldResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *GlobalFieldResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	g := data.Export()

	created, err := r.client.CreateGlobalField(g)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create GlobalField %#v, got error: %s", data.UID.ValueString(), err))
		return
	}

	// explicitly safe the computed fields (but not the optional fields; have to leave those as "unknown")
	data.ID = types.StringValue(created.UID)

	for i := range data.Fields {
		data.Fields[i].Update(created.Schema[i])
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a GlobalField", map[string]interface{}{
		"uid": created.UID,
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GlobalFieldResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *GlobalFieldResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	g, err := r.client.GetOneGlobalField(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read GlobalField %#v, got error: %s", data.UID.ValueString(), err))
		return
	}

	data.Update(g)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GlobalFieldResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *GlobalFieldResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	g := data.Export()

	_, err := r.client.UpdateGlobalField(g)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update GlobalField %#v, got error: %s", data.ID.ValueString(), err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GlobalFieldResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *GlobalFieldResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteGlobalField(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete GlobalField %#v, got error: %s", data.ID.ValueString(), err))
		return
	}
}

func (r *GlobalFieldResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

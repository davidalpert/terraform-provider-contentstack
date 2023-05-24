package provider

import (
	"context"
	"fmt"
	"github.com/davidalpert/go-contentstack/v1/management"
	cschema "github.com/davidalpert/go-contentstack/v1/schema"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &EnvironmentResource{}
var _ resource.ResourceWithImportState = &EnvironmentResource{}

func NewEnvironmentResource() resource.Resource {
	return &EnvironmentResource{}
}

// EnvironmentResource defines the resource implementation.
type EnvironmentResource struct {
	client *management.Client
}

// EnvironmentResourceModel describes the resource data model.
type EnvironmentResourceModel struct {
	Name          types.String `tfsdk:"name"`
	URLsByLocale  types.Map    `tfsdk:"urls"`
	ID            types.String `tfsdk:"id"`
	UID           types.String `tfsdk:"uid"`
	CreatedAt     types.String `tfsdk:"created_at"`
	UpdatedAt     types.String `tfsdk:"updated_at"`
	Version       types.Int64  `tfsdk:"version"`
	DeployContent types.Bool   `tfsdk:"deploy_content"`
}

func (data *EnvironmentResourceModel) Update(g *cschema.Environment) diag.Diagnostics {
	data.Name = types.StringValue(g.Name)
	data.DeployContent = types.BoolValue(g.DeployContent)
	data.Version = types.Int64Value(int64(g.Version))
	data.UpdatedAt = types.StringValue("TBD")
	uu := make(map[string]attr.Value)
	for _, u := range g.Urls {
		uu[u.Locale] = types.StringValue(u.Url)
	}
	um, dg := types.MapValue(types.StringType, uu)
	data.URLsByLocale = um

	return dg
}

func (data *EnvironmentResourceModel) Export() (*cschema.Environment, diag.Diagnostics) {
	g := &cschema.Environment{
		Name:          data.Name.ValueString(),
		UID:           data.ID.ValueString(),
		DeployContent: data.DeployContent.ValueBool(),
		//Version: int(data.Version.ValueInt64()),
		Urls: make([]cschema.LocaleURLPair, 0),
	}

	var urlsByLocal map[string]string
	dg := data.URLsByLocale.ElementsAs(context.Background(), &urlsByLocal, false)
	for l, u := range urlsByLocal {
		g.Urls = append(g.Urls, cschema.LocaleURLPair{
			Locale: l,
			Url:    u,
		})
	}

	return g, dg
}

func (r *EnvironmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environment"
}

func (r *EnvironmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Environment resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "internal terraform resource id (matches the uid when the Environment has been created/imported)",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"uid": schema.StringAttribute{
				MarkdownDescription: "internal contentstack identifier",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "name of the Environment",
				Required:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "created_at of the Global Field",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
				Optional:            true,
			},
			"urls": schema.MapAttribute{ // urls by locale
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "urls by locale",
			},
		},
	}
}

func (r *EnvironmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *EnvironmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *EnvironmentResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	g, dg := data.Export()
	resp.Diagnostics.Append(dg...)

	wrapper := cschema.SingleEnvironmentWrapper{Environment: *g}
	created, err := r.client.CreatePublishingEnvironment(&wrapper)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Environment %#v, got error: %s", data.Name.ValueString(), err))
		return
	}

	// explicitly save the computed fields (but not the optional fields; have to leave those as "unknown")
	data.Name = types.StringValue(created.Environment.Name)
	data.UID = types.StringValue(created.Environment.UID)
	data.ID = types.StringValue(created.Environment.UID)
	data.CreatedAt = types.StringValue("TBD")
	data.UpdatedAt = types.StringValue("TBD")
	data.Version = types.Int64Value(int64(created.Environment.Version))

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a Environment", map[string]interface{}{
		"uid":  created.Environment.UID,
		"name": created.Environment.Name,
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EnvironmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *EnvironmentResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	g, err := r.client.GetOnePublishingEnvironment(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read Environment %#v, got error: %s", data.Name.ValueString(), err))
		return
	}

	data.Update(g)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EnvironmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *EnvironmentResourceModel

	// Read Terraform plan data into the model
	//resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	g, dg := data.Export()
	resp.Diagnostics.Append(dg...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("uid"), &(g.UID))...)

	tflog.Trace(ctx, "about to update an Environment", map[string]interface{}{
		"uid":  g.UID,
		"name": g.Name,
	})

	_, err := r.client.UpdatePublishingEnvironment(g)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update Environment %#v, got error: %s", data.Name.ValueString(), err))
		return
	}

	data.UpdatedAt = types.StringValue("TBD")
	data.Version = types.Int64Value(int64(g.Version))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EnvironmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *EnvironmentResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeletePublishingEnvironment(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete Environment %#v, got error: %s", data.Name.ValueString(), err))
		return
	}
}

func (r *EnvironmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

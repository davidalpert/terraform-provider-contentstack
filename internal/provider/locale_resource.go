package provider

import (
	"context"
	"fmt"
	"github.com/davidalpert/go-contentstack/v1/management"
	cschema "github.com/davidalpert/go-contentstack/v1/schema"
	mystringplanmodifiers "github.com/davidalpert/terraform-provider-contentstack/internal/tfutils/stringplanmodifiers"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &LocaleResource{}
var _ resource.ResourceWithImportState = &LocaleResource{}

func NewLocaleResource() resource.Resource {
	return &LocaleResource{}
}

// LocaleResource defines the resource implementation.
type LocaleResource struct {
	client *management.Client
}

// LocaleResourceModel describes the resource data model.
type LocaleResourceModel struct {
	//ID                 types.String `tfsdk:"id"`
	UID                types.String `tfsdk:"uid"`
	Code               types.String `tfsdk:"code"`
	Name               types.String `tfsdk:"name"`
	FallbackLocaleCode types.String `tfsdk:"fallback_code"`

	CreatedBy types.String `tfsdk:"created_by"`
	CreatedAt types.String `tfsdk:"created_at"`
	Version   types.Int64  `tfsdk:"version"`
	UpdatedAt types.String `tfsdk:"updated_at"`
	UpdatedBy types.String `tfsdk:"updated_by"`
}

func (data *LocaleResourceModel) Update(g *cschema.Locale) {
	data.UID = types.StringValue(g.Uid)
	data.Code = types.StringValue(g.Code)
	data.Name = types.StringValue(g.Name)
	if g.FallbackLocaleCode == "" {
		data.FallbackLocaleCode = types.StringNull()
	} else {
		data.FallbackLocaleCode = types.StringValue(g.FallbackLocaleCode)
	}
	data.CreatedAt = types.StringValue(g.CreatedAt)
	data.CreatedBy = types.StringValue(g.CreatedBy)
	data.UpdatedAt = types.StringValue(g.UpdatedAt)
	data.UpdatedBy = types.StringValue(g.UpdatedBy)
	data.Version = types.Int64Value(int64(g.Version))
}

func (data *LocaleResourceModel) Export() *cschema.Locale {
	g := &cschema.Locale{
		Code:               data.Code.ValueString(),
		Name:               data.Name.ValueString(),
		Uid:                data.UID.ValueString(),
		CreatedAt:          data.CreatedAt.ValueString(),
		CreatedBy:          data.CreatedBy.ValueString(),
		UpdatedAt:          data.UpdatedAt.ValueString(),
		UpdatedBy:          data.UpdatedBy.ValueString(),
		Version:            (int)(data.Version.ValueInt64()),
		FallbackLocaleCode: data.FallbackLocaleCode.ValueString(),
	}

	return g
}

func (r *LocaleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_locale"
}

func (r *LocaleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Locale resource",

		Attributes: map[string]schema.Attribute{
			//"id": schema.StringAttribute{
			//	MarkdownDescription: "internal terraform resource id (matches the uid when the Locale has been created/imported)",
			//	Computed:            true,
			//	PlanModifiers: []planmodifier.String{
			//		stringplanmodifier.UseStateForUnknown(),
			//	},
			//},
			"uid": schema.StringAttribute{
				MarkdownDescription: "internal contentstack identifier",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"code": schema.StringAttribute{
				MarkdownDescription: "code of the Locale",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "name of the Locale",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					mystringplanmodifiers.DefaultValueCopiedFromAnotherField("code"),
				},
			},
			"fallback_code": schema.StringAttribute{
				MarkdownDescription: "code of the fallback Locale",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					mystringplanmodifiers.DefaultValue("en-us"),
				},
			},
			"created_by": schema.StringAttribute{
				MarkdownDescription: "created_by of the Locale",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "created_at of the Locale",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"version": schema.Int64Attribute{
				MarkdownDescription: "version number of the Locale",
				Computed:            true,
			},
			"updated_by": schema.StringAttribute{
				MarkdownDescription: "updated_at of the Locale",
				Computed:            true,
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "updated_by of the Locale",
				Computed:            true,
			},
		},
	}
}

func (r *LocaleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *LocaleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *LocaleResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	g := data.Export()

	created, err := r.client.CreateLocale(g.Code, g.Name, g.FallbackLocaleCode)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Locale %#v, got error: %s", data.Code.ValueString(), err))
		return
	}

	// explicitly save the computed fields (but not the optional fields; have to leave those as "unknown")
	data.Code = types.StringValue(created.Code)
	data.UID = types.StringValue(created.Uid)
	//data.ID = types.StringValue(created.Code)
	data.FallbackLocaleCode = types.StringValue(created.FallbackLocaleCode)
	data.CreatedAt = types.StringValue(created.CreatedAt)
	data.CreatedBy = types.StringValue(created.CreatedBy)
	data.Version = types.Int64Value(int64(created.Version))
	data.UpdatedAt = types.StringValue(created.UpdatedAt)
	data.UpdatedBy = types.StringValue(created.UpdatedBy)

	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a Locale", map[string]interface{}{
		"uid":  created.Uid,
		"code": created.Code,
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *LocaleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *LocaleResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	g, err := r.client.GetOneLocale(data.Code.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read Locale %#v, got error: %s", data.Code.ValueString(), err))
		return
	}

	data.Update(g)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *LocaleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *LocaleResourceModel

	// Read Terraform plan data into the model
	//resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	g := data.Export()
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("uid"), &(g.Uid))...)

	tflog.Trace(ctx, "about to update a Locale", map[string]interface{}{
		"uid":  g.Uid,
		"code": g.Code,
	})

	_, err := r.client.UpdateLocale(g)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update Locale %#v, got error: %s", data.Code.ValueString(), err))
		return
	}

	data.Version = types.Int64Value(int64(g.Version))
	data.UpdatedAt = types.StringValue(g.UpdatedAt)
	data.UpdatedBy = types.StringValue(g.UpdatedBy)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *LocaleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *LocaleResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteLocale(data.Code.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete Locale %#v, got error: %s", data.Code.ValueString(), err))
		return
	}
}

func (r *LocaleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("code"), req, resp)
}

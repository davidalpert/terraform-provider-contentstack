package provider

import (
	"context"
	"fmt"
	"github.com/davidalpert/go-contentstack/v1/management"
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
var _ resource.Resource = &DefaultLocaleResource{}
var _ resource.ResourceWithImportState = &DefaultLocaleResource{}

func NewDefaultLocaleResource() resource.Resource {
	return &DefaultLocaleResource{}
}

// DefaultLocaleResource defines the resource implementation.
type DefaultLocaleResource struct {
	client *management.Client
}

func (r *DefaultLocaleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_locale"
}

func (r *DefaultLocaleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Default locale resource",

		Attributes: map[string]schema.Attribute{
			"uid": schema.StringAttribute{
				MarkdownDescription: "internal contentstack identifier",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"code": schema.StringAttribute{
				MarkdownDescription: "code of the Locale",
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
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					mystringplanmodifiers.DefaultNull(),
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

func (r *DefaultLocaleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DefaultLocaleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *LocaleResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	l, err := r.client.GetOneLocale("en-us")
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read default Locale %#v, got error: %s", data.Code.ValueString(), err))
		return
	}

	// explicitly save the computed fields (but not the optional fields; have to leave those as "unknown")
	data.Update(l)

	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "found the default Locale", map[string]interface{}{
		"uid":  l.Uid,
		"code": l.Code,
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DefaultLocaleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

func (r *DefaultLocaleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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

func (r *DefaultLocaleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *LocaleResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// don't try deleting the default Locale
	// simply exit without error and allow terraform to clean up state
}

func (r *DefaultLocaleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("code"), req, resp)
}

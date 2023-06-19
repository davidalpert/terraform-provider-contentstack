package stringplanmodifier

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

type defaultToResourceName struct {
}

// DefaultToResourceName return a string plan modifier that sets the specified value to the name of the resource
func DefaultToResourceName() planmodifier.String {
	return defaultToResourceName{}
}

func (m defaultToResourceName) Description(context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to the name of the resource")
}

func (m defaultToResourceName) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m defaultToResourceName) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// If the attribute configuration is not null, apply the current plan (terraform core has already planned to apply the configured value)
	if !req.ConfigValue.IsNull() {
		return
	}

	//// If the attribute plan is "known" and "not null", then a previous plan modifier in the sequence
	//// has already been applied, and we don't want to interfere.
	//if !req.PlanValue.IsUnknown() && !req.PlanValue.IsNull() {
	//	return
	//}

	//ctx.Value(ContextKeyResourceName)
	var cfg tfsdk.Config
	diags := req.Config.Get(ctx, &cfg)
	resp.Diagnostics.Append(diags...)

	//req.Plan.
	//cfg.GetAttribute(ctx, path.r)

	//resp.PlanValue = types.StringValue(req.)
}

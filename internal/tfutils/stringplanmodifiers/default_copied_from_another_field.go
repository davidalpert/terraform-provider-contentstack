package stringplanmodifier

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type defaultCopiedFromAnotherField struct {
	attrName string
}

// DefaultValueCopiedFromAnotherField return a string plan modifier acts
// like the DefaultValue modifier but takes it's value at plan time from
// another field.
func DefaultValueCopiedFromAnotherField(s string) planmodifier.String {
	return defaultCopiedFromAnotherField{
		attrName: s,
	}
}

func (m defaultCopiedFromAnotherField) Description(context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to the value of the %#v property", m.attrName)
}

func (m defaultCopiedFromAnotherField) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m defaultCopiedFromAnotherField) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// Do nothing if there is a known planned value.
	if !req.PlanValue.IsUnknown() {
		return
	}

	// If the attribute configuration is not null, apply the current plan (terraform core has already planned to apply the configured value)
	if !req.ConfigValue.IsNull() {
		return
	}

	var val string
	diags := req.Config.GetAttribute(ctx, req.Path.ParentPath().AtName(m.attrName), &val)
	resp.Diagnostics.Append(diags...)

	resp.PlanValue = types.StringValue(val)
}

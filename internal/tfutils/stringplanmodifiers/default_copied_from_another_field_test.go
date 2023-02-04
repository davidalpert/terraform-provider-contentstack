package stringplanmodifier

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestDefaultValueCopiedFromAnotherFieldValue(t *testing.T) {
	t.Parallel()

	type testCase struct {
		configValue   types.String
		stateValue    types.String
		defaultValue  string
		expectedValue types.String
		expectError   bool
	}
	tests := map[string]testCase{
		"non-default non-Null string": {
			stateValue:    types.StringNull(),
			configValue:   types.StringValue("beta"),
			defaultValue:  "alpha",
			expectedValue: types.StringValue("beta"),
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// simulate Terraform-Core planning strategy
			var plannedValue types.String
			if !test.configValue.IsNull() {
				plannedValue = test.configValue
			} else {
				plannedValue = test.stateValue
			}

			ctx := context.Background()
			request := planmodifier.StringRequest{
				Path:        path.Root("test"),
				ConfigValue: test.configValue,
				PlanValue:   plannedValue,
				StateValue:  test.stateValue,
			}
			response := planmodifier.StringResponse{
				PlanValue: request.PlanValue,
			}
			DefaultValue(test.defaultValue).PlanModifyString(ctx, request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}

			if diff := cmp.Diff(response.PlanValue, test.expectedValue); diff != "" {
				t.Errorf("unexpected diff (+wanted, -got): %s", diff)
			}
		})
	}
}

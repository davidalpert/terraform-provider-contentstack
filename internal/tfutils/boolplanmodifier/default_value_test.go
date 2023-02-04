package boolplanmodifier

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestDefaultValue(t *testing.T) {
	t.Parallel()

	type testCase struct {
		stateValue    types.Bool
		configValue   types.Bool
		defaultValue  bool
		expectedValue types.Bool
		expectError   bool
	}
	tests := map[string]testCase{
		"default bool on create": {
			stateValue:    types.BoolNull(),
			configValue:   types.BoolNull(),
			defaultValue:  true,
			expectedValue: types.BoolValue(true),
		},
		"default bool on removed": {
			stateValue:    types.BoolValue(false),
			configValue:   types.BoolNull(),
			defaultValue:  true,
			expectedValue: types.BoolValue(true),
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// simulate Terraform-Core planning strategy
			var plannedValue types.Bool
			if !test.configValue.IsNull() {
				plannedValue = test.configValue
			} else {
				plannedValue = test.stateValue
			}

			ctx := context.Background()
			request := planmodifier.BoolRequest{
				Path:        path.Root("test"),
				StateValue:  test.stateValue,
				ConfigValue: test.configValue,
				PlanValue:   plannedValue,
			}
			response := planmodifier.BoolResponse{
				PlanValue: request.PlanValue,
			}
			DefaultValue(test.defaultValue).PlanModifyBool(ctx, request, &response)

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

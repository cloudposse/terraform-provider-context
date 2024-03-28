package framework

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// FromFrameworkMap converts a types.Map to a map[string]T.
func FromFrameworkMap[T interface{}](ctx context.Context, m types.Map) (map[string]T, diag.Diagnostics) {
	localValues := make(map[string]T, len(m.Elements()))
	if !m.IsNull() {
		diag := m.ElementsAs(ctx, &localValues, false)

		if diag.HasError() {
			return nil, diag
		}
	}
	return localValues, nil
}

// FromFrameworkList converts a types.List to a []T.
func FromFrameworkList[T interface{}](ctx context.Context, m types.List) ([]T, diag.Diagnostics) {
	localValues := make([]T, len(m.Elements()))
	if !m.IsNull() {
		diag := m.ElementsAs(ctx, &localValues, false)

		if diag.HasError() {
			return nil, diag
		}
	}
	return localValues, nil
}

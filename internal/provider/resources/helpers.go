package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// StringValueOrNull returns a StringValue if s is non-empty, otherwise StringNull.
func StringValueOrNull(s string) types.String {
	if s == "" {
		return types.StringNull()
	}
	return types.StringValue(s)
}

// StringValueOrDefault returns the string value or a default if null/unknown/empty.
func StringValueOrDefault(s types.String, defaultValue string) string {
	if s.IsNull() || s.IsUnknown() || s.ValueString() == "" {
		return defaultValue
	}
	return s.ValueString()
}

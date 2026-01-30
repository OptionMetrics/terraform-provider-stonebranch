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

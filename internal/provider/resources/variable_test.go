package resources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sbacctest "terraform-provider-stonebranch/internal/acctest"
)

// randVarName generates a valid variable name suffix.
// StoneBranch variable naming rules:
// - Must begin with a letter
// - Alphanumerics (upper or lower case) and underscore only
// - No hyphens, spaces, or special characters
// - Not case-sensitive
// - Do not use prefix "ops_" (reserved for built-in variables)
func randVarName(n int) string {
	return strings.ToUpper(acctest.RandString(n))
}

func TestAccVariableResource_basic(t *testing.T) {
	// Variable names: letters, numbers, underscores only (no hyphens)
	rName := "TF_TEST_VAR_" + randVarName(8)
	resourceName := "stonebranch_variable.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccVariableConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "value", "test-value"),
					resource.TestCheckResourceAttrSet(resourceName, "sys_id"),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
				),
			},
			// ImportState
			{
				ResourceName:                         resourceName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        rName,
				ImportStateVerifyIdentifierAttribute: "name",
			},
			// Update
			{
				Config: testAccVariableConfig_updated(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "value", "updated-value"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated variable"),
				),
			},
		},
	})
}

func TestAccVariableResource_withDescription(t *testing.T) {
	// Variable names: letters, numbers, underscores only (no hyphens)
	rName := "TF_TEST_VAR_" + randVarName(8)
	resourceName := "stonebranch_variable.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVariableConfig_withDescription(rName, "Initial description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "Initial description"),
				),
			},
			{
				Config: testAccVariableConfig_withDescription(rName, "Changed description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "Changed description"),
				),
			},
		},
	})
}

func TestAccVariableResource_emptyValue(t *testing.T) {
	// Variable names: letters, numbers, underscores only (no hyphens)
	rName := "TF_TEST_VAR_" + randVarName(8)
	resourceName := "stonebranch_variable.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVariableConfig_emptyValue(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "sys_id"),
				),
			},
		},
	})
}

// Test configuration helpers

func testAccVariableConfig_basic(name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_variable" "test" {
  name  = %[1]q
  value = "test-value"
}
`, name)
}

func testAccVariableConfig_updated(name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_variable" "test" {
  name        = %[1]q
  value       = "updated-value"
  description = "Updated variable"
}
`, name)
}

func testAccVariableConfig_withDescription(name, description string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_variable" "test" {
  name        = %[1]q
  value       = "test-value"
  description = %[2]q
}
`, name, description)
}

func testAccVariableConfig_emptyValue(name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_variable" "test" {
  name = %[1]q
}
`, name)
}

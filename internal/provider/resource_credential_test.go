package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCredentialResource_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test-cred")
	resourceName := "stonebranch_credential.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccCredentialConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "runtime_user", "testuser"),
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
				// Password is not returned by API, so skip verification
				ImportStateVerifyIgnore: []string{"runtime_password"},
			},
			// Update
			{
				Config: testAccCredentialConfig_updated(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated credential"),
					resource.TestCheckResourceAttr(resourceName, "runtime_user", "updateduser"),
				),
			},
		},
	})
}

func TestAccCredentialResource_withDescription(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test-cred")
	resourceName := "stonebranch_credential.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCredentialConfig_withDescription(rName, "Initial description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "Initial description"),
				),
			},
			{
				Config: testAccCredentialConfig_withDescription(rName, "Changed description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "Changed description"),
				),
			},
		},
	})
}

// Test configuration helpers

func testAccCredentialConfig_basic(name string) string {
	return providerConfig() + fmt.Sprintf(`
resource "stonebranch_credential" "test" {
  name             = %[1]q
  runtime_user     = "testuser"
  runtime_password = "testpassword123"
}
`, name)
}

func testAccCredentialConfig_updated(name string) string {
	return providerConfig() + fmt.Sprintf(`
resource "stonebranch_credential" "test" {
  name             = %[1]q
  description      = "Updated credential"
  runtime_user     = "updateduser"
  runtime_password = "newpassword456"
}
`, name)
}

func testAccCredentialConfig_withDescription(name, description string) string {
	return providerConfig() + fmt.Sprintf(`
resource "stonebranch_credential" "test" {
  name             = %[1]q
  description      = %[2]q
  runtime_user     = "testuser"
  runtime_password = "testpassword123"
}
`, name, description)
}

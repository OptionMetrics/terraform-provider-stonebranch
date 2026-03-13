package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sbacctest "github.com/OptionMetrics/terraform-provider-stonebranch/internal/acctest"
)

func TestAccEmailConnectionResource_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test-email")
	resourceName := "stonebranch_email_connection.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccEmailConnectionConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "smtp", "smtp.example.com"),
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
				ImportStateVerifyIgnore:              []string{"default_password"},
			},
			// Update
			{
				Config: testAccEmailConnectionConfig_updated(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "smtp", "smtp2.example.com"),
					resource.TestCheckResourceAttr(resourceName, "smtp_port", "587"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated email connection"),
				),
			},
		},
	})
}

func TestAccEmailConnectionResource_withAuthentication(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test-email")
	resourceName := "stonebranch_email_connection.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEmailConnectionConfig_withAuth(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "smtp", "smtp.example.com"),
					resource.TestCheckResourceAttr(resourceName, "authentication", "true"),
					resource.TestCheckResourceAttr(resourceName, "default_user", "testuser"),
				),
			},
		},
	})
}

func TestAccEmailConnectionResource_withSSL(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test-email")
	resourceName := "stonebranch_email_connection.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEmailConnectionConfig_withSSL(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "smtp", "smtp.example.com"),
					resource.TestCheckResourceAttr(resourceName, "smtp_port", "465"),
					resource.TestCheckResourceAttr(resourceName, "smtp_ssl", "true"),
				),
			},
		},
	})
}

func TestAccEmailConnectionResource_withStartTLS(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test-email")
	resourceName := "stonebranch_email_connection.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEmailConnectionConfig_withStartTLS(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "smtp", "smtp.example.com"),
					resource.TestCheckResourceAttr(resourceName, "smtp_port", "587"),
					resource.TestCheckResourceAttr(resourceName, "smtp_starttls", "true"),
				),
			},
		},
	})
}

func TestAccEmailConnectionResource_withIMAP(t *testing.T) {
	t.Skip("IMAP configuration requires a different connection type that is not currently supported by the API")

	rName := acctest.RandomWithPrefix("tf-test-email")
	resourceName := "stonebranch_email_connection.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEmailConnectionConfig_withIMAP(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "smtp", "smtp.example.com"),
					resource.TestCheckResourceAttr(resourceName, "imap", "imap.example.com"),
					resource.TestCheckResourceAttr(resourceName, "imap_port", "993"),
					resource.TestCheckResourceAttr(resourceName, "imap_ssl", "true"),
				),
			},
		},
	})
}

func TestAccEmailConnectionResource_minimal(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test-email")
	resourceName := "stonebranch_email_connection.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEmailConnectionConfig_minimal(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "smtp", "smtp.example.com"),
					resource.TestCheckResourceAttrSet(resourceName, "sys_id"),
				),
			},
		},
	})
}

// Test configuration helpers

func testAccEmailConnectionConfig_basic(name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_email_connection" "test" {
  name          = %[1]q
  smtp          = "smtp.example.com"
  smtp_port     = 25
  email_address = "test@example.com"
}
`, name)
}

func testAccEmailConnectionConfig_updated(name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_email_connection" "test" {
  name          = %[1]q
  smtp          = "smtp2.example.com"
  smtp_port     = 587
  email_address = "test@example.com"
  description   = "Updated email connection"
}
`, name)
}

func testAccEmailConnectionConfig_withAuth(name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_email_connection" "test" {
  name             = %[1]q
  smtp             = "smtp.example.com"
  smtp_port        = 587
  email_address    = "test@example.com"
  authentication   = true
  default_user     = "testuser"
  default_password = "testpass"
}
`, name)
}

func testAccEmailConnectionConfig_withSSL(name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_email_connection" "test" {
  name          = %[1]q
  smtp          = "smtp.example.com"
  smtp_port     = 465
  email_address = "test@example.com"
  smtp_ssl      = true
}
`, name)
}

func testAccEmailConnectionConfig_withStartTLS(name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_email_connection" "test" {
  name          = %[1]q
  smtp          = "smtp.example.com"
  smtp_port     = 587
  email_address = "test@example.com"
  smtp_starttls = true
}
`, name)
}

func testAccEmailConnectionConfig_withIMAP(name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_email_connection" "test" {
  name          = %[1]q
  smtp          = "smtp.example.com"
  smtp_port     = 25
  email_address = "test@example.com"
  imap          = "imap.example.com"
  imap_port     = 993
  imap_ssl      = true
}
`, name)
}

func testAccEmailConnectionConfig_minimal(name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_email_connection" "test" {
  name          = %[1]q
  smtp          = "smtp.example.com"
  smtp_port     = 25
  email_address = "test@example.com"
}
`, name)
}

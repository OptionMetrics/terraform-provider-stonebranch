package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sbacctest "github.com/OptionMetrics/terraform-provider-stonebranch/internal/acctest"
)

func TestAccTaskEmailResource_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test-email-task")
	connName := acctest.RandomWithPrefix("tf-test-email-conn")
	resourceName := "stonebranch_task_email.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccTaskEmailConfig_basic(connName, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "email_connection", connName),
					resource.TestCheckResourceAttr(resourceName, "to_recipients", "test@example.com"),
					resource.TestCheckResourceAttr(resourceName, "subject", "Test Subject"),
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
				Config: testAccTaskEmailConfig_updated(connName, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "to_recipients", "updated@example.com"),
					resource.TestCheckResourceAttr(resourceName, "subject", "Updated Subject"),
					resource.TestCheckResourceAttr(resourceName, "summary", "Updated email task"),
				),
			},
		},
	})
}

func TestAccTaskEmailResource_withBody(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test-email-task")
	connName := acctest.RandomWithPrefix("tf-test-email-conn")
	resourceName := "stonebranch_task_email.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTaskEmailConfig_withBody(connName, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "to_recipients", "test@example.com"),
					resource.TestCheckResourceAttr(resourceName, "subject", "Test Subject"),
					resource.TestCheckResourceAttr(resourceName, "body", "This is the email body content."),
				),
			},
		},
	})
}

func TestAccTaskEmailResource_withCCAndBCC(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test-email-task")
	connName := acctest.RandomWithPrefix("tf-test-email-conn")
	resourceName := "stonebranch_task_email.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTaskEmailConfig_withCCAndBCC(connName, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "to_recipients", "to@example.com"),
					resource.TestCheckResourceAttr(resourceName, "cc_recipients", "cc@example.com"),
					resource.TestCheckResourceAttr(resourceName, "bcc_recipients", "bcc@example.com"),
				),
			},
		},
	})
}

func TestAccTaskEmailResource_withRetry(t *testing.T) {
	t.Skip("Retry configuration is not supported by the email task API")

	rName := acctest.RandomWithPrefix("tf-test-email-task")
	connName := acctest.RandomWithPrefix("tf-test-email-conn")
	resourceName := "stonebranch_task_email.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTaskEmailConfig_withRetry(connName, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "retry_maximum", "3"),
					resource.TestCheckResourceAttr(resourceName, "retry_interval", "60"),
				),
			},
		},
	})
}

func TestAccTaskEmailResource_minimal(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test-email-task")
	connName := acctest.RandomWithPrefix("tf-test-email-conn")
	resourceName := "stonebranch_task_email.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTaskEmailConfig_minimal(connName, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "sys_id"),
				),
			},
		},
	})
}

// Test configuration helpers

func testAccTaskEmailConfig_basic(connName, name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_email_connection" "test" {
  name          = %[1]q
  smtp          = "smtp.example.com"
  smtp_port     = 25
  email_address = "test@example.com"
}

resource "stonebranch_task_email" "test" {
  name             = %[2]q
  email_connection = stonebranch_email_connection.test.name
  to_recipients    = "test@example.com"
  subject          = "Test Subject"
}
`, connName, name)
}

func testAccTaskEmailConfig_updated(connName, name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_email_connection" "test" {
  name          = %[1]q
  smtp          = "smtp.example.com"
  smtp_port     = 25
  email_address = "test@example.com"
}

resource "stonebranch_task_email" "test" {
  name             = %[2]q
  email_connection = stonebranch_email_connection.test.name
  to_recipients    = "updated@example.com"
  subject          = "Updated Subject"
  summary          = "Updated email task"
}
`, connName, name)
}

func testAccTaskEmailConfig_withBody(connName, name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_email_connection" "test" {
  name          = %[1]q
  smtp          = "smtp.example.com"
  smtp_port     = 25
  email_address = "test@example.com"
}

resource "stonebranch_task_email" "test" {
  name             = %[2]q
  email_connection = stonebranch_email_connection.test.name
  to_recipients    = "test@example.com"
  subject          = "Test Subject"
  body             = "This is the email body content."
}
`, connName, name)
}

func testAccTaskEmailConfig_withCCAndBCC(connName, name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_email_connection" "test" {
  name          = %[1]q
  smtp          = "smtp.example.com"
  smtp_port     = 25
  email_address = "test@example.com"
}

resource "stonebranch_task_email" "test" {
  name             = %[2]q
  email_connection = stonebranch_email_connection.test.name
  to_recipients    = "to@example.com"
  cc_recipients    = "cc@example.com"
  bcc_recipients   = "bcc@example.com"
  subject          = "Test Subject"
}
`, connName, name)
}

func testAccTaskEmailConfig_withRetry(connName, name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_email_connection" "test" {
  name          = %[1]q
  smtp          = "smtp.example.com"
  smtp_port     = 25
  email_address = "test@example.com"
}

resource "stonebranch_task_email" "test" {
  name             = %[2]q
  email_connection = stonebranch_email_connection.test.name
  to_recipients    = "test@example.com"
  subject          = "Test Subject"
  retry_maximum    = 3
  retry_interval   = 60
}
`, connName, name)
}

func testAccTaskEmailConfig_minimal(connName, name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_email_connection" "test" {
  name          = %[1]q
  smtp          = "smtp.example.com"
  smtp_port     = 25
  email_address = "test@example.com"
}

resource "stonebranch_task_email" "test" {
  name             = %[2]q
  email_connection = stonebranch_email_connection.test.name
  to_recipients    = "test@example.com"
  subject          = "Test"
}
`, connName, name)
}

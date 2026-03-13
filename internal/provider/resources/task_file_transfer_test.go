package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sbacctest "github.com/OptionMetrics/terraform-provider-stonebranch/internal/acctest"
)

func TestAccTaskFileTransferResource_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test-ftp")
	resourceName := "stonebranch_task_file_transfer.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccTaskFileTransferConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
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
				Config: testAccTaskFileTransferConfig_updated(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "summary", "Updated file transfer task"),
				),
			},
		},
	})
}

func TestAccTaskFileTransferResource_withSummary(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test-ftp")
	resourceName := "stonebranch_task_file_transfer.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTaskFileTransferConfig_withSummary(rName, "Initial summary"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "summary", "Initial summary"),
				),
			},
			{
				Config: testAccTaskFileTransferConfig_withSummary(rName, "Changed summary"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "summary", "Changed summary"),
				),
			},
		},
	})
}

// Test configuration helpers

func testAccTaskFileTransferConfig_basic(name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_file_transfer" "test" {
  name                   = %[1]q
  agent_var              = "agent_name"
  remote_server          = "test.example.com"
  remote_filename        = "/remote/path/file.txt"
  local_filename         = "/local/path/file.txt"
  remote_credentials_var = "ftp_credentials"
}
`, name)
}

func testAccTaskFileTransferConfig_updated(name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_file_transfer" "test" {
  name                   = %[1]q
  summary                = "Updated file transfer task"
  agent_var              = "agent_name"
  remote_server          = "test.example.com"
  remote_filename        = "/remote/path/file.txt"
  local_filename         = "/local/path/file.txt"
  remote_credentials_var = "ftp_credentials"
}
`, name)
}

func testAccTaskFileTransferConfig_withSummary(name, summary string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_file_transfer" "test" {
  name                   = %[1]q
  summary                = %[2]q
  agent_var              = "agent_name"
  remote_server          = "test.example.com"
  remote_filename        = "/remote/path/file.txt"
  local_filename         = "/local/path/file.txt"
  remote_credentials_var = "ftp_credentials"
}
`, name, summary)
}

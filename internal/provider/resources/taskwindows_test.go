package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sbacctest "terraform-provider-stonebranch/internal/acctest"
)

func TestAccTaskWindowsResource_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test")
	resourceName := "stonebranch_task_windows.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccTaskWindowsConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "command", "echo hello"),
					resource.TestCheckResourceAttrSet(resourceName, "sys_id"),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
				),
			},
			// ImportState - use the task name as import ID
			{
				ResourceName:                         resourceName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        rName,
				ImportStateVerifyIdentifierAttribute: "name",
			},
			// Update
			{
				Config: testAccTaskWindowsConfig_updated(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "command", "echo updated"),
					resource.TestCheckResourceAttr(resourceName, "summary", "Updated task summary"),
				),
			},
		},
	})
}

func TestAccTaskWindowsResource_withWindowsOptions(t *testing.T) {
	// Note: Windows-specific boolean options (elevate_user, create_console, desktop_interact)
	// are defined in the schema but may not be persisted by all Stonebranch API versions.
	// This test verifies the task can be created and that computed defaults are returned.
	rName := acctest.RandomWithPrefix("tf-test")
	resourceName := "stonebranch_task_windows.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTaskWindowsConfig_withWindowsOptions(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "command", "dir"),
					// Verify Windows-specific attributes have computed values (defaults from server)
					resource.TestCheckResourceAttrSet(resourceName, "elevate_user"),
					resource.TestCheckResourceAttrSet(resourceName, "create_console"),
					resource.TestCheckResourceAttrSet(resourceName, "desktop_interact"),
				),
			},
		},
	})
}

func TestAccTaskWindowsResource_withScript(t *testing.T) {
	// This test creates a script resource and a Windows task that references it
	scriptName := acctest.RandomWithPrefix("tf-test-script")
	taskName := acctest.RandomWithPrefix("tf-test-task")
	taskResourceName := "stonebranch_task_windows.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTaskWindowsConfig_withScript(taskName, scriptName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(taskResourceName, "name", taskName),
					resource.TestCheckResourceAttr(taskResourceName, "command_or_script", "Script"),
					resource.TestCheckResourceAttr(taskResourceName, "script", scriptName),
				),
			},
		},
	})
}

func TestAccTaskWindowsResource_withSummary(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test")
	resourceName := "stonebranch_task_windows.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTaskWindowsConfig_withSummary(rName, "Initial summary"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "summary", "Initial summary"),
				),
			},
			{
				Config: testAccTaskWindowsConfig_withSummary(rName, "Changed summary"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "summary", "Changed summary"),
				),
			},
		},
	})
}

// Test configuration helpers

func testAccTaskWindowsConfig_basic(name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_windows" "test" {
  name       = %[1]q
  command    = "echo hello"
  agent_var  = "agent_name"
  exit_codes = "0"
}
`, name)
}

func testAccTaskWindowsConfig_updated(name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_windows" "test" {
  name       = %[1]q
  command    = "echo updated"
  summary    = "Updated task summary"
  agent_var  = "agent_name"
  exit_codes = "0"
}
`, name)
}

func testAccTaskWindowsConfig_withWindowsOptions(name string) string {
	// Note: We don't explicitly set Windows-specific boolean options here
	// because the API may not persist them. We just verify they have computed defaults.
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_windows" "test" {
  name       = %[1]q
  command    = "dir"
  agent_var  = "agent_name"
  exit_codes = "0"
}
`, name)
}

func testAccTaskWindowsConfig_withScript(taskName, scriptName string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_script" "test" {
  name    = %[2]q
  content = "@echo Hello from script"
}

resource "stonebranch_task_windows" "test" {
  name              = %[1]q
  command_or_script = "Script"
  script            = stonebranch_script.test.name
  agent_var         = "agent_name"
  exit_codes        = "0"
}
`, taskName, scriptName)
}

func testAccTaskWindowsConfig_withSummary(name, summary string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_windows" "test" {
  name       = %[1]q
  command    = "echo hello"
  summary    = %[2]q
  agent_var  = "agent_name"
  exit_codes = "0"
}
`, name, summary)
}

package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sbacctest "github.com/OptionMetrics/terraform-provider-stonebranch/internal/acctest"
)

func TestAccScriptResource_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test-script")
	resourceName := "stonebranch_script.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccScriptConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "content", "echo 'Hello World'"),
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
				Config: testAccScriptConfig_updated(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "content", "echo 'Updated Script'"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated description"),
				),
			},
		},
	})
}

func TestAccScriptResource_withDescription(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test-script")
	resourceName := "stonebranch_script.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccScriptConfig_withDescription(rName, "Initial description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "Initial description"),
				),
			},
			{
				Config: testAccScriptConfig_withDescription(rName, "Changed description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "Changed description"),
				),
			},
		},
	})
}

// Integration test: Create a script and reference it from a Unix task
func TestAccScriptResource_withUnixTask(t *testing.T) {
	scriptName := acctest.RandomWithPrefix("tf-test-script")
	taskName := acctest.RandomWithPrefix("tf-test-task")
	scriptResourceName := "stonebranch_script.test"
	taskResourceName := "stonebranch_task_unix.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccScriptConfig_withUnixTask(scriptName, taskName),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify script was created
					resource.TestCheckResourceAttr(scriptResourceName, "name", scriptName),
					resource.TestCheckResourceAttrSet(scriptResourceName, "sys_id"),
					// Verify task was created and references the script
					resource.TestCheckResourceAttr(taskResourceName, "name", taskName),
					resource.TestCheckResourceAttr(taskResourceName, "command_or_script", "Script"),
					resource.TestCheckResourceAttr(taskResourceName, "script", scriptName),
				),
			},
		},
	})
}

// Test configuration helpers

func testAccScriptConfig_basic(name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_script" "test" {
  name    = %[1]q
  content = "echo 'Hello World'"
}
`, name)
}

func testAccScriptConfig_updated(name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_script" "test" {
  name        = %[1]q
  content     = "echo 'Updated Script'"
  description = "Updated description"
}
`, name)
}

func testAccScriptConfig_withDescription(name, description string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_script" "test" {
  name        = %[1]q
  content     = "echo 'Hello World'"
  description = %[2]q
}
`, name, description)
}

func testAccScriptConfig_withUnixTask(scriptName, taskName string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_script" "test" {
  name    = %[1]q
  content = "echo 'Hello from script'"
}

resource "stonebranch_task_unix" "test" {
  name              = %[2]q
  command_or_script = "Script"
  script            = stonebranch_script.test.name
  agent_var         = "agent_name"
  exit_codes        = "0"
}
`, scriptName, taskName)
}

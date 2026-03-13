package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sbacctest "github.com/OptionMetrics/terraform-provider-stonebranch/internal/acctest"
)

func TestAccTriggerTimeResource_basic(t *testing.T) {
	taskName := acctest.RandomWithPrefix("tf-test-task")
	triggerName := acctest.RandomWithPrefix("tf-test-trigger")
	resourceName := "stonebranch_trigger_time.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read (triggers are created disabled by default)
			{
				Config: testAccTriggerTimeConfig_basic(taskName, triggerName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", triggerName),
					resource.TestCheckResourceAttrSet(resourceName, "sys_id"),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
				),
			},
			// ImportState
			{
				ResourceName:                         resourceName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        triggerName,
				ImportStateVerifyIdentifierAttribute: "name",
			},
			// Update - change description
			{
				Config: testAccTriggerTimeConfig_updated(taskName, triggerName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", triggerName),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated trigger description"),
				),
			},
		},
	})
}

func TestAccTriggerTimeResource_disabled(t *testing.T) {
	taskName := acctest.RandomWithPrefix("tf-test-task")
	triggerName := acctest.RandomWithPrefix("tf-test-trigger")
	resourceName := "stonebranch_trigger_time.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTriggerTimeConfig_disabled(taskName, triggerName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", triggerName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
				),
			},
		},
	})
}

func TestAccTriggerTimeResource_withTime(t *testing.T) {
	taskName := acctest.RandomWithPrefix("tf-test-task")
	triggerName := acctest.RandomWithPrefix("tf-test-trigger")
	resourceName := "stonebranch_trigger_time.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTriggerTimeConfig_withTime(taskName, triggerName, "09:00"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", triggerName),
					resource.TestCheckResourceAttr(resourceName, "time", "09:00"),
				),
			},
		},
	})
}

// Test configuration helpers

func testAccTriggerTimeConfig_basic(taskName, triggerName string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_unix" "test" {
  name       = %[1]q
  command    = "echo 'Triggered task'"
  agent_var  = "agent_name"
  exit_codes = "0"
}

resource "stonebranch_trigger_time" "test" {
  name  = %[2]q
  time  = "12:00"
  tasks = [stonebranch_task_unix.test.name]
}
`, taskName, triggerName)
}

func testAccTriggerTimeConfig_updated(taskName, triggerName string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_unix" "test" {
  name       = %[1]q
  command    = "echo 'Triggered task'"
  agent_var  = "agent_name"
  exit_codes = "0"
}

resource "stonebranch_trigger_time" "test" {
  name        = %[2]q
  description = "Updated trigger description"
  time        = "12:00"
  tasks       = [stonebranch_task_unix.test.name]
}
`, taskName, triggerName)
}

func testAccTriggerTimeConfig_disabled(taskName, triggerName string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_unix" "test" {
  name       = %[1]q
  command    = "echo 'Triggered task'"
  agent_var  = "agent_name"
  exit_codes = "0"
}

resource "stonebranch_trigger_time" "test" {
  name    = %[2]q
  enabled = false
  time    = "12:00"
  tasks   = [stonebranch_task_unix.test.name]
}
`, taskName, triggerName)
}

func testAccTriggerTimeConfig_withTime(taskName, triggerName, time string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_unix" "test" {
  name       = %[1]q
  command    = "echo 'Triggered task'"
  agent_var  = "agent_name"
  exit_codes = "0"
}

resource "stonebranch_trigger_time" "test" {
  name  = %[2]q
  time  = %[3]q
  tasks = [stonebranch_task_unix.test.name]
}
`, taskName, triggerName, time)
}

func TestAccTriggerTimeResource_withVariables(t *testing.T) {
	taskName := acctest.RandomWithPrefix("tf-test-task")
	triggerName := acctest.RandomWithPrefix("tf-test-trigger")
	resourceName := "stonebranch_trigger_time.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with variables
			{
				Config: testAccTriggerTimeConfig_withVariables(taskName, triggerName, "initial_value"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", triggerName),
					resource.TestCheckResourceAttr(resourceName, "variables.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "variables.0.name", "trigger_var1"),
					resource.TestCheckResourceAttr(resourceName, "variables.0.value", "initial_value"),
					resource.TestCheckResourceAttr(resourceName, "variables.0.description", "Trigger variable 1"),
					resource.TestCheckResourceAttr(resourceName, "variables.1.name", "trigger_var2"),
					resource.TestCheckResourceAttr(resourceName, "variables.1.value", "value2"),
				),
			},
			// Update variables
			{
				Config: testAccTriggerTimeConfig_withVariables(taskName, triggerName, "updated_value"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "variables.0.value", "updated_value"),
				),
			},
		},
	})
}

func testAccTriggerTimeConfig_withVariables(taskName, triggerName, varValue string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_unix" "test" {
  name       = %[1]q
  command    = "echo 'Triggered task'"
  agent_var  = "agent_name"
  exit_codes = "0"
}

resource "stonebranch_trigger_time" "test" {
  name  = %[2]q
  time  = "12:00"
  tasks = [stonebranch_task_unix.test.name]

  variables = [
    {
      name        = "trigger_var1"
      value       = %[3]q
      description = "Trigger variable 1"
    },
    {
      name  = "trigger_var2"
      value = "value2"
    }
  ]
}
`, taskName, triggerName, varValue)
}

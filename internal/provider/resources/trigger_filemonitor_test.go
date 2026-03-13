package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sbacctest "github.com/OptionMetrics/terraform-provider-stonebranch/internal/acctest"
)

func TestAccTriggerFileMonitorResource_basic(t *testing.T) {
	taskName := acctest.RandomWithPrefix("tf-test-task")
	monitorTaskName := acctest.RandomWithPrefix("tf-test-monitor")
	triggerName := acctest.RandomWithPrefix("tf-test-fm")
	resourceName := "stonebranch_trigger_file_monitor.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read (triggers are created disabled by default)
			{
				Config: testAccTriggerFileMonitorConfig_basic(taskName, monitorTaskName, triggerName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", triggerName),
					resource.TestCheckResourceAttr(resourceName, "task_monitor", monitorTaskName),
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
			// Update - add description
			{
				Config: testAccTriggerFileMonitorConfig_updated(taskName, monitorTaskName, triggerName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", triggerName),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated file monitor trigger"),
				),
			},
		},
	})
}

func TestAccTriggerFileMonitorResource_withTimeRestrictions(t *testing.T) {
	taskName := acctest.RandomWithPrefix("tf-test-task")
	monitorTaskName := acctest.RandomWithPrefix("tf-test-monitor")
	triggerName := acctest.RandomWithPrefix("tf-test-fm")
	resourceName := "stonebranch_trigger_file_monitor.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTriggerFileMonitorConfig_withTimeRestrictions(taskName, monitorTaskName, triggerName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", triggerName),
					resource.TestCheckResourceAttr(resourceName, "restricted_times", "true"),
					resource.TestCheckResourceAttr(resourceName, "enabled_start", "08:00"),
					resource.TestCheckResourceAttr(resourceName, "enabled_end", "18:00"),
				),
			},
		},
	})
}

func TestAccTriggerFileMonitorResource_withTimeZone(t *testing.T) {
	taskName := acctest.RandomWithPrefix("tf-test-task")
	monitorTaskName := acctest.RandomWithPrefix("tf-test-monitor")
	triggerName := acctest.RandomWithPrefix("tf-test-fm")
	resourceName := "stonebranch_trigger_file_monitor.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTriggerFileMonitorConfig_withTimeZone(taskName, monitorTaskName, triggerName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", triggerName),
					resource.TestCheckResourceAttr(resourceName, "time_zone", "America/New_York"),
				),
			},
		},
	})
}

// Test configuration helpers
// These tests use stonebranch_task_file_monitor for the task_monitor field

func testAccTriggerFileMonitorConfig_basic(taskName, monitorTaskName, triggerName string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_unix" "test" {
  name       = %[1]q
  command    = "echo 'Triggered task'"
  agent_var  = "agent_name"
  exit_codes = "0"
}

resource "stonebranch_task_file_monitor" "monitor" {
  name      = %[2]q
  file_name = "/tmp/incoming/*.csv"
  agent_var = "agent_name"
}

resource "stonebranch_trigger_file_monitor" "test" {
  name         = %[3]q
  task_monitor = stonebranch_task_file_monitor.monitor.name
  tasks        = [stonebranch_task_unix.test.name]
}
`, taskName, monitorTaskName, triggerName)
}

func testAccTriggerFileMonitorConfig_updated(taskName, monitorTaskName, triggerName string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_unix" "test" {
  name       = %[1]q
  command    = "echo 'Triggered task'"
  agent_var  = "agent_name"
  exit_codes = "0"
}

resource "stonebranch_task_file_monitor" "monitor" {
  name      = %[2]q
  file_name = "/tmp/incoming/*.csv"
  agent_var = "agent_name"
}

resource "stonebranch_trigger_file_monitor" "test" {
  name         = %[3]q
  description  = "Updated file monitor trigger"
  task_monitor = stonebranch_task_file_monitor.monitor.name
  tasks        = [stonebranch_task_unix.test.name]
}
`, taskName, monitorTaskName, triggerName)
}

func testAccTriggerFileMonitorConfig_withTimeRestrictions(taskName, monitorTaskName, triggerName string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_unix" "test" {
  name       = %[1]q
  command    = "echo 'Triggered task'"
  agent_var  = "agent_name"
  exit_codes = "0"
}

resource "stonebranch_task_file_monitor" "monitor" {
  name      = %[2]q
  file_name = "/tmp/incoming/*.csv"
  agent_var = "agent_name"
}

resource "stonebranch_trigger_file_monitor" "test" {
  name             = %[3]q
  description      = "File monitor with time restrictions"
  task_monitor     = stonebranch_task_file_monitor.monitor.name
  tasks            = [stonebranch_task_unix.test.name]
  restricted_times = true
  enabled_start    = "08:00"
  enabled_end      = "18:00"
}
`, taskName, monitorTaskName, triggerName)
}

func testAccTriggerFileMonitorConfig_withTimeZone(taskName, monitorTaskName, triggerName string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_unix" "test" {
  name       = %[1]q
  command    = "echo 'Triggered task'"
  agent_var  = "agent_name"
  exit_codes = "0"
}

resource "stonebranch_task_file_monitor" "monitor" {
  name      = %[2]q
  file_name = "/tmp/incoming/*.csv"
  agent_var = "agent_name"
}

resource "stonebranch_trigger_file_monitor" "test" {
  name         = %[3]q
  task_monitor = stonebranch_task_file_monitor.monitor.name
  tasks        = [stonebranch_task_unix.test.name]
  time_zone    = "America/New_York"
}
`, taskName, monitorTaskName, triggerName)
}

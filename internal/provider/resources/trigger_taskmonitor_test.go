package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sbacctest "terraform-provider-stonebranch/internal/acctest"
)

func TestAccTriggerTaskMonitorResource_basic(t *testing.T) {
	rName := "tf-test-ttm-" + acctest.RandString(8)
	taskName := "tf-test-task-" + acctest.RandString(8)
	watchedTaskName := "tf-test-watched-" + acctest.RandString(8)
	monitorTaskName := "tf-test-mon-" + acctest.RandString(8)
	resourceName := "stonebranch_trigger_task_monitor.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccTriggerTaskMonitorConfig_basic(rName, taskName, watchedTaskName, monitorTaskName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "task_monitor", monitorTaskName),
					resource.TestCheckResourceAttr(resourceName, "tasks.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tasks.0", taskName),
					resource.TestCheckResourceAttrSet(resourceName, "sys_id"),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
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
			// Update with description
			{
				Config: testAccTriggerTaskMonitorConfig_withDescription(rName, taskName, watchedTaskName, monitorTaskName, "Updated description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated description"),
				),
			},
		},
	})
}

// Note: Skipping enabled=true test because the Stonebranch API does not currently
// support enabling Task Monitor triggers via the API. They must be enabled
// through the UI. Triggers are created and remain disabled.
func TestAccTriggerTaskMonitorResource_disabled(t *testing.T) {
	rName := "tf-test-ttm-dis-" + acctest.RandString(8)
	taskName := "tf-test-task-" + acctest.RandString(8)
	watchedTaskName := "tf-test-watched-" + acctest.RandString(8)
	monitorTaskName := "tf-test-mon-" + acctest.RandString(8)
	resourceName := "stonebranch_trigger_task_monitor.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Verify triggers are created disabled by default
			{
				Config: testAccTriggerTaskMonitorConfig_basic(rName, taskName, watchedTaskName, monitorTaskName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
				),
			},
		},
	})
}

func TestAccTriggerTaskMonitorResource_multipleTasks(t *testing.T) {
	rName := "tf-test-ttm-mt-" + acctest.RandString(8)
	taskName1 := "tf-test-task1-" + acctest.RandString(8)
	taskName2 := "tf-test-task2-" + acctest.RandString(8)
	watchedTaskName := "tf-test-watched-" + acctest.RandString(8)
	monitorTaskName := "tf-test-mon-" + acctest.RandString(8)
	resourceName := "stonebranch_trigger_task_monitor.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTriggerTaskMonitorConfig_multipleTasks(rName, taskName1, taskName2, watchedTaskName, monitorTaskName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "tasks.#", "2"),
				),
			},
		},
	})
}

func TestAccTriggerTaskMonitorResource_withBusinessService(t *testing.T) {
	rName := "tf-test-ttm-bs-" + acctest.RandString(8)
	taskName := "tf-test-task-" + acctest.RandString(8)
	watchedTaskName := "tf-test-watched-" + acctest.RandString(8)
	monitorTaskName := "tf-test-mon-" + acctest.RandString(8)
	bsName := "tf-test-bs-" + acctest.RandString(8)
	resourceName := "stonebranch_trigger_task_monitor.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTriggerTaskMonitorConfig_withBusinessService(rName, taskName, watchedTaskName, monitorTaskName, bsName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "opswise_groups.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "opswise_groups.0", bsName),
				),
			},
		},
	})
}

// Test configuration helpers

func testAccTriggerTaskMonitorConfig_basic(name, taskName, watchedTaskName, monitorTaskName string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
# Task to be triggered when the monitor fires
resource "stonebranch_task_unix" "triggered" {
  name       = %[2]q
  agent_var  = "agent_var"
  command    = "echo 'Triggered task'"
  exit_codes = "0"
}

# Task to be watched (when this completes, the trigger fires)
resource "stonebranch_task_unix" "watched" {
  name       = %[3]q
  agent_var  = "agent_var"
  command    = "echo 'Watched task'"
  exit_codes = "0"
}

# Task Monitor task that watches the Unix task
resource "stonebranch_task_monitor" "monitor" {
  name          = %[4]q
  task_mon_name = stonebranch_task_unix.watched.name
  status_text   = "Success"
}

resource "stonebranch_trigger_task_monitor" "test" {
  name         = %[1]q
  task_monitor = stonebranch_task_monitor.monitor.name
  tasks        = [stonebranch_task_unix.triggered.name]
}
`, name, taskName, watchedTaskName, monitorTaskName)
}

func testAccTriggerTaskMonitorConfig_withDescription(name, taskName, watchedTaskName, monitorTaskName, description string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_unix" "triggered" {
  name       = %[2]q
  agent_var  = "agent_var"
  command    = "echo 'Triggered task'"
  exit_codes = "0"
}

resource "stonebranch_task_unix" "watched" {
  name       = %[3]q
  agent_var  = "agent_var"
  command    = "echo 'Watched task'"
  exit_codes = "0"
}

resource "stonebranch_task_monitor" "monitor" {
  name          = %[4]q
  task_mon_name = stonebranch_task_unix.watched.name
  status_text   = "Success"
}

resource "stonebranch_trigger_task_monitor" "test" {
  name         = %[1]q
  description  = %[5]q
  task_monitor = stonebranch_task_monitor.monitor.name
  tasks        = [stonebranch_task_unix.triggered.name]
}
`, name, taskName, watchedTaskName, monitorTaskName, description)
}

func testAccTriggerTaskMonitorConfig_enabled(name, taskName, watchedTaskName, monitorTaskName string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_unix" "triggered" {
  name       = %[2]q
  agent_var  = "agent_var"
  command    = "echo 'Triggered task'"
  exit_codes = "0"
}

resource "stonebranch_task_unix" "watched" {
  name       = %[3]q
  agent_var  = "agent_var"
  command    = "echo 'Watched task'"
  exit_codes = "0"
}

resource "stonebranch_task_monitor" "monitor" {
  name          = %[4]q
  task_mon_name = stonebranch_task_unix.watched.name
  status_text   = "Success"
}

resource "stonebranch_trigger_task_monitor" "test" {
  name         = %[1]q
  task_monitor = stonebranch_task_monitor.monitor.name
  tasks        = [stonebranch_task_unix.triggered.name]
  enabled      = true
}
`, name, taskName, watchedTaskName, monitorTaskName)
}

func testAccTriggerTaskMonitorConfig_multipleTasks(name, taskName1, taskName2, watchedTaskName, monitorTaskName string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_unix" "triggered1" {
  name       = %[2]q
  agent_var  = "agent_var"
  command    = "echo 'Triggered task 1'"
  exit_codes = "0"
}

resource "stonebranch_task_unix" "triggered2" {
  name       = %[3]q
  agent_var  = "agent_var"
  command    = "echo 'Triggered task 2'"
  exit_codes = "0"
}

resource "stonebranch_task_unix" "watched" {
  name       = %[4]q
  agent_var  = "agent_var"
  command    = "echo 'Watched task'"
  exit_codes = "0"
}

resource "stonebranch_task_monitor" "monitor" {
  name          = %[5]q
  task_mon_name = stonebranch_task_unix.watched.name
  status_text   = "Success"
}

resource "stonebranch_trigger_task_monitor" "test" {
  name         = %[1]q
  task_monitor = stonebranch_task_monitor.monitor.name
  tasks        = [
    stonebranch_task_unix.triggered1.name,
    stonebranch_task_unix.triggered2.name
  ]
}
`, name, taskName1, taskName2, watchedTaskName, monitorTaskName)
}

func testAccTriggerTaskMonitorConfig_withBusinessService(name, taskName, watchedTaskName, monitorTaskName, bsName string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_business_service" "test" {
  name = %[5]q
}

resource "stonebranch_task_unix" "triggered" {
  name       = %[2]q
  agent_var  = "agent_var"
  command    = "echo 'Triggered task'"
  exit_codes = "0"
}

resource "stonebranch_task_unix" "watched" {
  name       = %[3]q
  agent_var  = "agent_var"
  command    = "echo 'Watched task'"
  exit_codes = "0"
}

resource "stonebranch_task_monitor" "monitor" {
  name          = %[4]q
  task_mon_name = stonebranch_task_unix.watched.name
  status_text   = "Success"
}

resource "stonebranch_trigger_task_monitor" "test" {
  name           = %[1]q
  task_monitor   = stonebranch_task_monitor.monitor.name
  tasks          = [stonebranch_task_unix.triggered.name]
  opswise_groups = [stonebranch_business_service.test.name]
}
`, name, taskName, watchedTaskName, monitorTaskName, bsName)
}

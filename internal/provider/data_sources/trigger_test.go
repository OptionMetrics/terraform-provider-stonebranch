package data_sources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sbacctest "github.com/OptionMetrics/terraform-provider-stonebranch/internal/acctest"
)

func TestAccTriggerDataSource_basic(t *testing.T) {
	triggerName := "tf-test-trig-ds-" + acctest.RandString(8)
	taskName := "tf-test-task-" + acctest.RandString(8)
	dataSourceName := "data.stonebranch_trigger.test"
	resourceName := "stonebranch_trigger_time.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTriggerDataSourceConfig_basic(triggerName, taskName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", triggerName),
					resource.TestCheckResourceAttrPair(dataSourceName, "sys_id", resourceName, "sys_id"),
					resource.TestCheckResourceAttr(dataSourceName, "type", "triggerTime"),
					resource.TestCheckResourceAttrSet(dataSourceName, "version"),
					resource.TestCheckResourceAttr(dataSourceName, "enabled", "false"),
				),
			},
		},
	})
}

func TestAccTriggerDataSource_withDescription(t *testing.T) {
	triggerName := "tf-test-trig-ds-desc-" + acctest.RandString(8)
	taskName := "tf-test-task-" + acctest.RandString(8)
	dataSourceName := "data.stonebranch_trigger.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTriggerDataSourceConfig_withDescription(triggerName, taskName, "Test description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", triggerName),
					resource.TestCheckResourceAttr(dataSourceName, "description", "Test description"),
					resource.TestCheckResourceAttr(dataSourceName, "tasks.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "tasks.0", taskName),
				),
			},
		},
	})
}

func TestAccTriggerDataSource_cronTrigger(t *testing.T) {
	triggerName := "tf-test-trig-ds-cron-" + acctest.RandString(8)
	taskName := "tf-test-task-" + acctest.RandString(8)
	dataSourceName := "data.stonebranch_trigger.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTriggerDataSourceConfig_cronTrigger(triggerName, taskName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", triggerName),
					resource.TestCheckResourceAttr(dataSourceName, "type", "triggerCron"),
				),
			},
		},
	})
}

// Test configuration helpers

func testAccTriggerDataSourceConfig_basic(triggerName, taskName string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_unix" "test" {
  name       = %[2]q
  agent_var  = "agent_var"
  command    = "echo 'test'"
  exit_codes = "0"
}

resource "stonebranch_trigger_time" "test" {
  name  = %[1]q
  tasks = [stonebranch_task_unix.test.name]
  time  = "12:00"
}

data "stonebranch_trigger" "test" {
  name = stonebranch_trigger_time.test.name
}
`, triggerName, taskName)
}

func testAccTriggerDataSourceConfig_withDescription(triggerName, taskName, description string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_unix" "test" {
  name       = %[2]q
  agent_var  = "agent_var"
  command    = "echo 'test'"
  exit_codes = "0"
}

resource "stonebranch_trigger_time" "test" {
  name        = %[1]q
  description = %[3]q
  tasks       = [stonebranch_task_unix.test.name]
  time        = "12:00"
}

data "stonebranch_trigger" "test" {
  name = stonebranch_trigger_time.test.name
}
`, triggerName, taskName, description)
}

func testAccTriggerDataSourceConfig_cronTrigger(triggerName, taskName string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_unix" "test" {
  name       = %[2]q
  agent_var  = "agent_var"
  command    = "echo 'test'"
  exit_codes = "0"
}

resource "stonebranch_trigger_cron" "test" {
  name         = %[1]q
  tasks        = [stonebranch_task_unix.test.name]
  minutes      = "0"
  hours        = "12"
  day_of_month = "*"
  month        = "*"
  day_of_week  = "*"
}

data "stonebranch_trigger" "test" {
  name = stonebranch_trigger_cron.test.name
}
`, triggerName, taskName)
}

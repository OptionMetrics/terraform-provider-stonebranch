package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sbacctest "terraform-provider-stonebranch/internal/acctest"
)

func TestAccTriggerCronResource_basic(t *testing.T) {
	taskName := acctest.RandomWithPrefix("tf-test-task")
	triggerName := acctest.RandomWithPrefix("tf-test-cron")
	resourceName := "stonebranch_trigger_cron.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read (triggers are created disabled by default)
			{
				Config: testAccTriggerCronConfig_basic(taskName, triggerName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", triggerName),
					resource.TestCheckResourceAttr(resourceName, "minutes", "0"),
					resource.TestCheckResourceAttr(resourceName, "hours", "12"),
					resource.TestCheckResourceAttr(resourceName, "day_of_month", "*"),
					resource.TestCheckResourceAttr(resourceName, "month", "*"),
					resource.TestCheckResourceAttr(resourceName, "day_of_week", "*"),
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
			// Update - change cron expression
			{
				Config: testAccTriggerCronConfig_updated(taskName, triggerName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", triggerName),
					resource.TestCheckResourceAttr(resourceName, "minutes", "30"),
					resource.TestCheckResourceAttr(resourceName, "hours", "8"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated cron trigger"),
				),
			},
		},
	})
}

func TestAccTriggerCronResource_everyMinute(t *testing.T) {
	taskName := acctest.RandomWithPrefix("tf-test-task")
	triggerName := acctest.RandomWithPrefix("tf-test-cron")
	resourceName := "stonebranch_trigger_cron.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTriggerCronConfig_everyMinute(taskName, triggerName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", triggerName),
					resource.TestCheckResourceAttr(resourceName, "minutes", "*"),
					resource.TestCheckResourceAttr(resourceName, "hours", "*"),
					resource.TestCheckResourceAttr(resourceName, "day_of_month", "*"),
					resource.TestCheckResourceAttr(resourceName, "month", "*"),
					resource.TestCheckResourceAttr(resourceName, "day_of_week", "*"),
				),
			},
		},
	})
}

func TestAccTriggerCronResource_weekdays(t *testing.T) {
	taskName := acctest.RandomWithPrefix("tf-test-task")
	triggerName := acctest.RandomWithPrefix("tf-test-cron")
	resourceName := "stonebranch_trigger_cron.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTriggerCronConfig_weekdays(taskName, triggerName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", triggerName),
					resource.TestCheckResourceAttr(resourceName, "minutes", "0"),
					resource.TestCheckResourceAttr(resourceName, "hours", "9"),
					resource.TestCheckResourceAttr(resourceName, "day_of_month", "*"),
					resource.TestCheckResourceAttr(resourceName, "month", "*"),
					resource.TestCheckResourceAttr(resourceName, "day_of_week", "1-5"),
				),
			},
		},
	})
}

func TestAccTriggerCronResource_withTimeZone(t *testing.T) {
	taskName := acctest.RandomWithPrefix("tf-test-task")
	triggerName := acctest.RandomWithPrefix("tf-test-cron")
	resourceName := "stonebranch_trigger_cron.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTriggerCronConfig_withTimeZone(taskName, triggerName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", triggerName),
					resource.TestCheckResourceAttr(resourceName, "time_zone", "America/New_York"),
				),
			},
		},
	})
}

func TestAccTriggerCronResource_disabled(t *testing.T) {
	taskName := acctest.RandomWithPrefix("tf-test-task")
	triggerName := acctest.RandomWithPrefix("tf-test-cron")
	resourceName := "stonebranch_trigger_cron.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTriggerCronConfig_disabled(taskName, triggerName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", triggerName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
				),
			},
		},
	})
}

// Test configuration helpers

func testAccTriggerCronConfig_basic(taskName, triggerName string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_unix" "test" {
  name       = %[1]q
  command    = "echo 'Triggered task'"
  agent_var  = "agent_name"
  exit_codes = "0"
}

resource "stonebranch_trigger_cron" "test" {
  name         = %[2]q
  minutes      = "0"
  hours        = "12"
  day_of_month = "*"
  month        = "*"
  day_of_week  = "*"
  tasks        = [stonebranch_task_unix.test.name]
}
`, taskName, triggerName)
}

func testAccTriggerCronConfig_updated(taskName, triggerName string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_unix" "test" {
  name       = %[1]q
  command    = "echo 'Triggered task'"
  agent_var  = "agent_name"
  exit_codes = "0"
}

resource "stonebranch_trigger_cron" "test" {
  name         = %[2]q
  description  = "Updated cron trigger"
  minutes      = "30"
  hours        = "8"
  day_of_month = "*"
  month        = "*"
  day_of_week  = "*"
  tasks        = [stonebranch_task_unix.test.name]
}
`, taskName, triggerName)
}

func testAccTriggerCronConfig_everyMinute(taskName, triggerName string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_unix" "test" {
  name       = %[1]q
  command    = "echo 'Triggered task'"
  agent_var  = "agent_name"
  exit_codes = "0"
}

resource "stonebranch_trigger_cron" "test" {
  name         = %[2]q
  description  = "Runs every minute"
  minutes      = "*"
  hours        = "*"
  day_of_month = "*"
  month        = "*"
  day_of_week  = "*"
  tasks        = [stonebranch_task_unix.test.name]
}
`, taskName, triggerName)
}

func testAccTriggerCronConfig_weekdays(taskName, triggerName string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_unix" "test" {
  name       = %[1]q
  command    = "echo 'Triggered task'"
  agent_var  = "agent_name"
  exit_codes = "0"
}

resource "stonebranch_trigger_cron" "test" {
  name         = %[2]q
  description  = "Runs at 9am on weekdays"
  minutes      = "0"
  hours        = "9"
  day_of_month = "*"
  month        = "*"
  day_of_week  = "1-5"
  tasks        = [stonebranch_task_unix.test.name]
}
`, taskName, triggerName)
}

func testAccTriggerCronConfig_withTimeZone(taskName, triggerName string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_unix" "test" {
  name       = %[1]q
  command    = "echo 'Triggered task'"
  agent_var  = "agent_name"
  exit_codes = "0"
}

resource "stonebranch_trigger_cron" "test" {
  name         = %[2]q
  minutes      = "0"
  hours        = "12"
  day_of_month = "*"
  month        = "*"
  day_of_week  = "*"
  time_zone    = "America/New_York"
  tasks        = [stonebranch_task_unix.test.name]
}
`, taskName, triggerName)
}

func testAccTriggerCronConfig_disabled(taskName, triggerName string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_unix" "test" {
  name       = %[1]q
  command    = "echo 'Triggered task'"
  agent_var  = "agent_name"
  exit_codes = "0"
}

resource "stonebranch_trigger_cron" "test" {
  name         = %[2]q
  enabled      = false
  minutes      = "0"
  hours        = "12"
  day_of_month = "*"
  month        = "*"
  day_of_week  = "*"
  tasks        = [stonebranch_task_unix.test.name]
}
`, taskName, triggerName)
}

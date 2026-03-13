package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sbacctest "github.com/OptionMetrics/terraform-provider-stonebranch/internal/acctest"
)

func TestAccTaskTimerResource_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test-timer-task")
	resourceName := "stonebranch_task_timer.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccTaskTimerConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "sleep_type", "Duration"),
					resource.TestCheckResourceAttr(resourceName, "sleep_duration", "00:00:00:10"),
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
				Config: testAccTaskTimerConfig_updated(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "sleep_type", "Duration"),
					resource.TestCheckResourceAttr(resourceName, "sleep_duration", "00:00:05:00"),
					resource.TestCheckResourceAttr(resourceName, "summary", "Updated timer task"),
				),
			},
		},
	})
}

func TestAccTaskTimerResource_withTime(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test-timer-task")
	resourceName := "stonebranch_task_timer.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTaskTimerConfig_withTime(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "sleep_type", "Time"),
					resource.TestCheckResourceAttr(resourceName, "sleep_time", "14:30"),
					resource.TestCheckResourceAttr(resourceName, "sleep_day_constraint", "Same Day"),
				),
			},
		},
	})
}

func TestAccTaskTimerResource_withRelativeTime(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test-timer-task")
	resourceName := "stonebranch_task_timer.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTaskTimerConfig_withRelativeTime(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "sleep_type", "Relative Time"),
					resource.TestCheckResourceAttr(resourceName, "sleep_time", "01:30"),
				),
			},
		},
	})
}

func TestAccTaskTimerResource_withVariables(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test-timer-task")
	resourceName := "stonebranch_task_timer.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTaskTimerConfig_withVariables(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "variables.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "variables.0.name", "delay_seconds"),
					resource.TestCheckResourceAttr(resourceName, "variables.0.value", "30"),
				),
			},
		},
	})
}

// Test configuration helpers

func testAccTaskTimerConfig_basic(name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_timer" "test" {
  name           = %[1]q
  sleep_type     = "Duration"
  sleep_duration = "00:00:00:10"
}
`, name)
}

func testAccTaskTimerConfig_updated(name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_timer" "test" {
  name           = %[1]q
  sleep_type     = "Duration"
  sleep_duration = "00:00:05:00"
  summary        = "Updated timer task"
}
`, name)
}

func testAccTaskTimerConfig_withTime(name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_timer" "test" {
  name                 = %[1]q
  sleep_type           = "Time"
  sleep_time           = "14:30"
  sleep_day_constraint = "Same Day"
}
`, name)
}

func testAccTaskTimerConfig_withRelativeTime(name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_timer" "test" {
  name       = %[1]q
  sleep_type = "Relative Time"
  sleep_time = "01:30"
}
`, name)
}

func testAccTaskTimerConfig_withVariables(name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_timer" "test" {
  name           = %[1]q
  sleep_type     = "Duration"
  sleep_duration = "00:00:00:10"

  variables = [
    {
      name  = "delay_seconds"
      value = "30"
    }
  ]
}
`, name)
}

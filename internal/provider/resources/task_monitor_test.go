package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sbacctest "terraform-provider-stonebranch/internal/acctest"
)

func TestAccTaskMonitorResource_basic(t *testing.T) {
	rName := "tf-test-tm-" + acctest.RandString(8)
	watchedTaskName := "tf-test-watched-" + acctest.RandString(8)
	resourceName := "stonebranch_task_monitor.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccTaskMonitorConfig_basic(rName, watchedTaskName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "task_mon_name", watchedTaskName),
					resource.TestCheckResourceAttr(resourceName, "status_text", "Success"),
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
			// Update with summary
			{
				Config: testAccTaskMonitorConfig_withSummary(rName, watchedTaskName, "Updated summary"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "summary", "Updated summary"),
				),
			},
		},
	})
}

func TestAccTaskMonitorResource_lateStart(t *testing.T) {
	rName := "tf-test-tm-ls-" + acctest.RandString(8)
	watchedTaskName := "tf-test-watched-" + acctest.RandString(8)
	resourceName := "stonebranch_task_monitor.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTaskMonitorConfig_lateStart(rName, watchedTaskName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "monitor_late_start", "true"),
				),
			},
		},
	})
}

func TestAccTaskMonitorResource_withBusinessService(t *testing.T) {
	rName := "tf-test-tm-bs-" + acctest.RandString(8)
	watchedTaskName := "tf-test-watched-" + acctest.RandString(8)
	bsName := "tf-test-bs-" + acctest.RandString(8)
	resourceName := "stonebranch_task_monitor.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTaskMonitorConfig_withBusinessService(rName, watchedTaskName, bsName),
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

func testAccTaskMonitorConfig_basic(name, watchedTaskName string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_unix" "watched" {
  name       = %[2]q
  agent_var  = "agent_var"
  command    = "echo 'Watched task'"
  exit_codes = "0"
}

resource "stonebranch_task_monitor" "test" {
  name          = %[1]q
  task_mon_name = stonebranch_task_unix.watched.name
  status_text   = "Success"
}
`, name, watchedTaskName)
}

func testAccTaskMonitorConfig_withSummary(name, watchedTaskName, summary string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_unix" "watched" {
  name       = %[2]q
  agent_var  = "agent_var"
  command    = "echo 'Watched task'"
  exit_codes = "0"
}

resource "stonebranch_task_monitor" "test" {
  name          = %[1]q
  summary       = %[3]q
  task_mon_name = stonebranch_task_unix.watched.name
  status_text   = "Success"
}
`, name, watchedTaskName, summary)
}

func testAccTaskMonitorConfig_lateStart(name, watchedTaskName string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_unix" "watched" {
  name       = %[2]q
  agent_var  = "agent_var"
  command    = "echo 'Watched task'"
  exit_codes = "0"
}

resource "stonebranch_task_monitor" "test" {
  name               = %[1]q
  task_mon_name      = stonebranch_task_unix.watched.name
  monitor_late_start = true
}
`, name, watchedTaskName)
}

func testAccTaskMonitorConfig_withBusinessService(name, watchedTaskName, bsName string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_business_service" "test" {
  name = %[3]q
}

resource "stonebranch_task_unix" "watched" {
  name       = %[2]q
  agent_var  = "agent_var"
  command    = "echo 'Watched task'"
  exit_codes = "0"
}

resource "stonebranch_task_monitor" "test" {
  name           = %[1]q
  task_mon_name  = stonebranch_task_unix.watched.name
  status_text    = "Success"
  opswise_groups = [stonebranch_business_service.test.name]
}
`, name, watchedTaskName, bsName)
}

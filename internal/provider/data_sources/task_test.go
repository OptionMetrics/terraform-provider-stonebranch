package data_sources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sbacctest "github.com/OptionMetrics/terraform-provider-stonebranch/internal/acctest"
)

func TestAccTaskDataSource_basic(t *testing.T) {
	taskName := "tf-test-task-ds-" + acctest.RandString(8)
	dataSourceName := "data.stonebranch_task.test"
	resourceName := "stonebranch_task_unix.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTaskDataSourceConfig_basic(taskName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", taskName),
					resource.TestCheckResourceAttrPair(dataSourceName, "sys_id", resourceName, "sys_id"),
					resource.TestCheckResourceAttr(dataSourceName, "type", "taskUnix"),
					resource.TestCheckResourceAttrSet(dataSourceName, "version"),
				),
			},
		},
	})
}

func TestAccTaskDataSource_withSummary(t *testing.T) {
	taskName := "tf-test-task-ds-sum-" + acctest.RandString(8)
	dataSourceName := "data.stonebranch_task.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTaskDataSourceConfig_withSummary(taskName, "Test summary"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", taskName),
					resource.TestCheckResourceAttr(dataSourceName, "summary", "Test summary"),
				),
			},
		},
	})
}

// Test configuration helpers

func testAccTaskDataSourceConfig_basic(name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_unix" "test" {
  name       = %[1]q
  agent_var  = "agent_var"
  command    = "echo 'test'"
  exit_codes = "0"
}

data "stonebranch_task" "test" {
  name = stonebranch_task_unix.test.name
}
`, name)
}

func testAccTaskDataSourceConfig_withSummary(name, summary string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_unix" "test" {
  name       = %[1]q
  summary    = %[2]q
  agent_var  = "agent_var"
  command    = "echo 'test'"
  exit_codes = "0"
}

data "stonebranch_task" "test" {
  name = stonebranch_task_unix.test.name
}
`, name, summary)
}

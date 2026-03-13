package data_sources_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sbacctest "github.com/OptionMetrics/terraform-provider-stonebranch/internal/acctest"
)

func TestAccTaskInstancesDataSource_withTaskName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTaskInstancesDataSourceConfig_withTaskName(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.stonebranch_task_instances.test", "task_instances.#"),
				),
			},
		},
	})
}

// Test configuration helpers

// Note: The task instances API requires a task_name parameter.
// We use a wildcard pattern to match any task.
func testAccTaskInstancesDataSourceConfig_withTaskName() string {
	return sbacctest.ProviderConfig() + `
data "stonebranch_task_instances" "test" {
  task_name         = "*"
  updated_time_type = "Today"
}
`
}

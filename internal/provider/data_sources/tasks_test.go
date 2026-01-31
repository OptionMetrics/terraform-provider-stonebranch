package data_sources_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sbacctest "terraform-provider-stonebranch/internal/acctest"
)

func TestAccTasksDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTasksDataSourceConfig_basic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.stonebranch_tasks.test", "tasks.#"),
				),
			},
		},
	})
}

func TestAccTasksDataSource_withNameFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTasksDataSourceConfig_withNameFilter("*"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.stonebranch_tasks.test", "tasks.#"),
				),
			},
		},
	})
}

// Test configuration helpers

func testAccTasksDataSourceConfig_basic() string {
	return sbacctest.ProviderConfig() + `
data "stonebranch_tasks" "test" {
}
`
}

func testAccTasksDataSourceConfig_withNameFilter(name string) string {
	return sbacctest.ProviderConfig() + `
data "stonebranch_tasks" "test" {
  name = "` + name + `"
}
`
}

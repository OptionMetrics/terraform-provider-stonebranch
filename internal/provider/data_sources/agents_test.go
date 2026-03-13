package data_sources_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sbacctest "github.com/OptionMetrics/terraform-provider-stonebranch/internal/acctest"
)

func TestAccAgentsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAgentsDataSourceConfig_basic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.stonebranch_agents.test", "agents.#"),
				),
			},
		},
	})
}

func TestAccAgentsDataSource_withTypeFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAgentsDataSourceConfig_withTypeFilter("Linux/Unix"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.stonebranch_agents.test", "agents.#"),
				),
			},
		},
	})
}

// Test configuration helpers

func testAccAgentsDataSourceConfig_basic() string {
	return sbacctest.ProviderConfig() + `
data "stonebranch_agents" "test" {
}
`
}

func testAccAgentsDataSourceConfig_withTypeFilter(agentType string) string {
	return sbacctest.ProviderConfig() + `
data "stonebranch_agents" "test" {
  type = "` + agentType + `"
}
`
}

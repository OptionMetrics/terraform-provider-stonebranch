package data_sources_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sbacctest "github.com/OptionMetrics/terraform-provider-stonebranch/internal/acctest"
)

func TestAccAgentClustersDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAgentClustersDataSourceConfig_basic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.stonebranch_agent_clusters.test", "agent_clusters.#"),
				),
			},
		},
	})
}

func TestAccAgentClustersDataSource_withTypeFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAgentClustersDataSourceConfig_withTypeFilter("Linux/Unix"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.stonebranch_agent_clusters.test", "agent_clusters.#"),
				),
			},
		},
	})
}

// Test configuration helpers

func testAccAgentClustersDataSourceConfig_basic() string {
	return sbacctest.ProviderConfig() + `
data "stonebranch_agent_clusters" "test" {
}
`
}

func testAccAgentClustersDataSourceConfig_withTypeFilter(clusterType string) string {
	return sbacctest.ProviderConfig() + `
data "stonebranch_agent_clusters" "test" {
  type = "` + clusterType + `"
}
`
}

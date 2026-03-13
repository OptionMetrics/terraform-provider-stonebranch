package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sbacctest "github.com/OptionMetrics/terraform-provider-stonebranch/internal/acctest"
)

func TestAccAgentClusterResource_basic(t *testing.T) {
	rName := "tf-test-ac-" + acctest.RandString(8)
	resourceName := "stonebranch_agent_cluster.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccAgentClusterConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "Linux/Unix"),
					resource.TestCheckResourceAttrSet(resourceName, "sys_id"),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
					// Check defaults
					resource.TestCheckResourceAttr(resourceName, "distribution", "Any"),
					resource.TestCheckResourceAttr(resourceName, "limit_type", "Unlimited"),
					resource.TestCheckResourceAttr(resourceName, "agent_limit_type", "Unlimited"),
					resource.TestCheckResourceAttr(resourceName, "ignore_inactive_agents", "true"),
					resource.TestCheckResourceAttr(resourceName, "ignore_suspended_agents", "true"),
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
				Config: testAccAgentClusterConfig_withDescription(rName, "Updated description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated description"),
				),
			},
		},
	})
}

func TestAccAgentClusterResource_windows(t *testing.T) {
	rName := "tf-test-ac-win-" + acctest.RandString(8)
	resourceName := "stonebranch_agent_cluster.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAgentClusterConfig_windows(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "Windows"),
				),
			},
		},
	})
}

func TestAccAgentClusterResource_withDistribution(t *testing.T) {
	rName := "tf-test-ac-dist-" + acctest.RandString(8)
	resourceName := "stonebranch_agent_cluster.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAgentClusterConfig_withDistribution(rName, "Round Robin"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "distribution", "Round Robin"),
				),
			},
			{
				Config: testAccAgentClusterConfig_withDistribution(rName, "Lowest CPU Utilization"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "distribution", "Lowest CPU Utilization"),
				),
			},
		},
	})
}

func TestAccAgentClusterResource_withLimits(t *testing.T) {
	rName := "tf-test-ac-lim-" + acctest.RandString(8)
	resourceName := "stonebranch_agent_cluster.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAgentClusterConfig_withLimits(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "limit_type", "Limited"),
					resource.TestCheckResourceAttr(resourceName, "limit_amount", "10"),
					resource.TestCheckResourceAttr(resourceName, "agent_limit_type", "Limited"),
					resource.TestCheckResourceAttr(resourceName, "agent_limit_amount", "2"),
				),
			},
		},
	})
}

func TestAccAgentClusterResource_withBusinessService(t *testing.T) {
	rName := "tf-test-ac-bs-" + acctest.RandString(8)
	bsName := "tf-test-bs-" + acctest.RandString(8)
	resourceName := "stonebranch_agent_cluster.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAgentClusterConfig_withBusinessService(rName, bsName),
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

func testAccAgentClusterConfig_basic(name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_agent_cluster" "test" {
  name = %[1]q
  type = "Linux/Unix"
}
`, name)
}

func testAccAgentClusterConfig_withDescription(name, description string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_agent_cluster" "test" {
  name        = %[1]q
  type        = "Linux/Unix"
  description = %[2]q
}
`, name, description)
}

func testAccAgentClusterConfig_windows(name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_agent_cluster" "test" {
  name = %[1]q
  type = "Windows"
}
`, name)
}

func testAccAgentClusterConfig_withDistribution(name, distribution string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_agent_cluster" "test" {
  name         = %[1]q
  type         = "Linux/Unix"
  distribution = %[2]q
}
`, name, distribution)
}

func testAccAgentClusterConfig_withLimits(name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_agent_cluster" "test" {
  name               = %[1]q
  type               = "Linux/Unix"
  limit_type         = "Limited"
  limit_amount       = 10
  agent_limit_type   = "Limited"
  agent_limit_amount = 2
}
`, name)
}

func testAccAgentClusterConfig_withBusinessService(name, bsName string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_business_service" "test" {
  name = %[2]q
}

resource "stonebranch_agent_cluster" "test" {
  name          = %[1]q
  type          = "Linux/Unix"
  opswise_groups = [stonebranch_business_service.test.name]
}
`, name, bsName)
}

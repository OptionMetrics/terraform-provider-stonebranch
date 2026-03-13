package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sbacctest "github.com/OptionMetrics/terraform-provider-stonebranch/internal/acctest"
)

func TestAccTaskFileMonitorResource_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test-fm")
	resourceName := "stonebranch_task_file_monitor.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccTaskFileMonitorConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "file_name", "/tmp/incoming/*.csv"),
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
				Config: testAccTaskFileMonitorConfig_updated(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "file_name", "/tmp/incoming/*.txt"),
					resource.TestCheckResourceAttr(resourceName, "summary", "Updated file monitor"),
				),
			},
		},
	})
}

func TestAccTaskFileMonitorResource_withOptions(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test-fm")
	resourceName := "stonebranch_task_file_monitor.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTaskFileMonitorConfig_withOptions(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "file_name", "/data/incoming/*"),
					resource.TestCheckResourceAttr(resourceName, "recursive", "true"),
					resource.TestCheckResourceAttr(resourceName, "stable_seconds", "30"),
				),
			},
		},
	})
}

func TestAccTaskFileMonitorResource_existMonitor(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test-fm")
	resourceName := "stonebranch_task_file_monitor.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTaskFileMonitorConfig_existMonitor(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "fm_type", "Exist"),
				),
			},
		},
	})
}

// Test configuration helpers

func testAccTaskFileMonitorConfig_basic(name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_file_monitor" "test" {
  name      = %[1]q
  file_name = "/tmp/incoming/*.csv"
  agent_var = "agent_name"
}
`, name)
}

func testAccTaskFileMonitorConfig_updated(name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_file_monitor" "test" {
  name      = %[1]q
  summary   = "Updated file monitor"
  file_name = "/tmp/incoming/*.txt"
  agent_var = "agent_name"
}
`, name)
}

func testAccTaskFileMonitorConfig_withOptions(name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_file_monitor" "test" {
  name           = %[1]q
  summary        = "File monitor with options"
  file_name      = "/data/incoming/*"
  agent_var      = "agent_name"
  recursive      = true
  stable_seconds = 30
}
`, name)
}

func testAccTaskFileMonitorConfig_existMonitor(name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_file_monitor" "test" {
  name      = %[1]q
  summary   = "Check if file exists"
  file_name = "/data/ready.flag"
  agent_var = "agent_name"
  fm_type   = "Exist"
}
`, name)
}

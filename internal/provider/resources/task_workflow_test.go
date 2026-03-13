package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sbacctest "github.com/OptionMetrics/terraform-provider-stonebranch/internal/acctest"
)

func TestAccTaskWorkflowResource_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test-wf")
	resourceName := "stonebranch_task_workflow.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccTaskWorkflowConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
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
				Config: testAccTaskWorkflowConfig_updated(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "summary", "Updated workflow"),
				),
			},
		},
	})
}

func TestAccTaskWorkflowResource_withSummary(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test-wf")
	resourceName := "stonebranch_task_workflow.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTaskWorkflowConfig_withSummary(rName, "Initial summary"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "summary", "Initial summary"),
				),
			},
			{
				Config: testAccTaskWorkflowConfig_withSummary(rName, "Changed summary"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "summary", "Changed summary"),
				),
			},
		},
	})
}

func TestAccTaskWorkflowResource_withOptions(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test-wf")
	resourceName := "stonebranch_task_workflow.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTaskWorkflowConfig_withOptions(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "calculate_critical_path", "true"),
				),
			},
		},
	})
}

// Test configuration helpers

func testAccTaskWorkflowConfig_basic(name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_workflow" "test" {
  name = %[1]q
}
`, name)
}

func testAccTaskWorkflowConfig_updated(name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_workflow" "test" {
  name    = %[1]q
  summary = "Updated workflow"
}
`, name)
}

func testAccTaskWorkflowConfig_withSummary(name, summary string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_workflow" "test" {
  name    = %[1]q
  summary = %[2]q
}
`, name, summary)
}

func testAccTaskWorkflowConfig_withOptions(name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_workflow" "test" {
  name                    = %[1]q
  summary                 = "Workflow with options"
  calculate_critical_path = true
}
`, name)
}

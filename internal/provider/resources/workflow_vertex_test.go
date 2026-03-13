package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sbacctest "github.com/OptionMetrics/terraform-provider-stonebranch/internal/acctest"
)

func TestAccWorkflowVertexResource_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test-wfv")
	resourceName := "stonebranch_workflow_vertex.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccWorkflowVertexConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "workflow_name", rName+"-wf"),
					resource.TestCheckResourceAttr(resourceName, "task_name", rName+"-task"),
					resource.TestCheckResourceAttrSet(resourceName, "vertex_id"),
				),
			},
		},
	})
}

func TestAccWorkflowVertexResource_withAlias(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test-wfv")
	resourceName := "stonebranch_workflow_vertex.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWorkflowVertexConfig_withAlias(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "workflow_name", rName+"-wf"),
					resource.TestCheckResourceAttr(resourceName, "task_name", rName+"-task"),
					resource.TestCheckResourceAttr(resourceName, "alias", "MyTaskAlias"),
					resource.TestCheckResourceAttrSet(resourceName, "vertex_id"),
				),
			},
		},
	})
}

func TestAccWorkflowVertexResource_withPosition(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test-wfv")
	resourceName := "stonebranch_workflow_vertex.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWorkflowVertexConfig_withPosition(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "workflow_name", rName+"-wf"),
					resource.TestCheckResourceAttr(resourceName, "task_name", rName+"-task"),
					resource.TestCheckResourceAttr(resourceName, "vertex_x", "100"),
					resource.TestCheckResourceAttr(resourceName, "vertex_y", "200"),
					resource.TestCheckResourceAttrSet(resourceName, "vertex_id"),
				),
			},
		},
	})
}

// Test configuration helpers

func testAccWorkflowVertexConfig_basic(name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_workflow" "test" {
  name = "%[1]s-wf"
}

resource "stonebranch_task_unix" "test" {
  name       = "%[1]s-task"
  agent_var  = "agent_name"
  command    = "echo hello"
  exit_codes = "0"
}

resource "stonebranch_workflow_vertex" "test" {
  workflow_name = stonebranch_task_workflow.test.name
  task_name     = stonebranch_task_unix.test.name
}
`, name)
}

func testAccWorkflowVertexConfig_withAlias(name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_workflow" "test" {
  name = "%[1]s-wf"
}

resource "stonebranch_task_unix" "test" {
  name       = "%[1]s-task"
  agent_var  = "agent_name"
  command    = "echo hello"
  exit_codes = "0"
}

resource "stonebranch_workflow_vertex" "test" {
  workflow_name = stonebranch_task_workflow.test.name
  task_name     = stonebranch_task_unix.test.name
  alias         = "MyTaskAlias"
}
`, name)
}

func testAccWorkflowVertexConfig_withPosition(name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_workflow" "test" {
  name = "%[1]s-wf"
}

resource "stonebranch_task_unix" "test" {
  name       = "%[1]s-task"
  agent_var  = "agent_name"
  command    = "echo hello"
  exit_codes = "0"
}

resource "stonebranch_workflow_vertex" "test" {
  workflow_name = stonebranch_task_workflow.test.name
  task_name     = stonebranch_task_unix.test.name
  vertex_x      = "100"
  vertex_y      = "200"
}
`, name)
}

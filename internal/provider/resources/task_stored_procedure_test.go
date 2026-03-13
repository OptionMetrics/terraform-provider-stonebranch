package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sbacctest "github.com/OptionMetrics/terraform-provider-stonebranch/internal/acctest"
)

func TestAccTaskStoredProcedureResource_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test-sproc")
	connName := acctest.RandomWithPrefix("tf-test-dbconn")
	credName := acctest.RandomWithPrefix("tf-test-cred")
	resourceName := "stonebranch_task_stored_procedure.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccTaskStoredProcedureConfig_basic(rName, connName, credName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "stored_proc_name", "test_procedure"),
					resource.TestCheckResourceAttr(resourceName, "database_connection", connName),
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
				Config: testAccTaskStoredProcedureConfig_updated(rName, connName, credName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "stored_proc_name", "updated_procedure"),
					resource.TestCheckResourceAttr(resourceName, "summary", "Updated stored procedure task"),
				),
			},
		},
	})
}

func TestAccTaskStoredProcedureResource_withSummary(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test-sproc")
	connName := acctest.RandomWithPrefix("tf-test-dbconn")
	credName := acctest.RandomWithPrefix("tf-test-cred")
	resourceName := "stonebranch_task_stored_procedure.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTaskStoredProcedureConfig_withSummary(rName, connName, credName, "Initial summary"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "summary", "Initial summary"),
				),
			},
			{
				Config: testAccTaskStoredProcedureConfig_withSummary(rName, connName, credName, "Changed summary"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "summary", "Changed summary"),
				),
			},
		},
	})
}

func TestAccTaskStoredProcedureResource_withParameters(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test-sproc")
	connName := acctest.RandomWithPrefix("tf-test-dbconn")
	credName := acctest.RandomWithPrefix("tf-test-cred")
	resourceName := "stonebranch_task_stored_procedure.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTaskStoredProcedureConfig_withParameters(rName, connName, credName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "stored_proc_name", "test_with_params"),
					resource.TestCheckResourceAttr(resourceName, "parameters.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.param_mode", "Input"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.param_type", "VARCHAR"),
					resource.TestCheckResourceAttr(resourceName, "parameters.1.param_mode", "Output"),
					resource.TestCheckResourceAttr(resourceName, "parameters.1.param_type", "INTEGER"),
				),
			},
		},
	})
}

func TestAccTaskStoredProcedureResource_withRetryConfig(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test-sproc")
	connName := acctest.RandomWithPrefix("tf-test-dbconn")
	credName := acctest.RandomWithPrefix("tf-test-cred")
	resourceName := "stonebranch_task_stored_procedure.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTaskStoredProcedureConfig_withRetryConfig(rName, connName, credName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "retry_maximum", "3"),
					resource.TestCheckResourceAttr(resourceName, "retry_interval", "60"),
				),
			},
		},
	})
}

// Test configuration helpers

// testAccTaskStoredProcedureConfig_dependencies creates the required database connection and credential
func testAccTaskStoredProcedureConfig_dependencies(connName, credName string) string {
	return fmt.Sprintf(`
resource "stonebranch_credential" "test" {
  name             = %[2]q
  runtime_user     = "testuser"
  runtime_password = "testpassword123"
}

resource "stonebranch_database_connection" "test" {
  name        = %[1]q
  db_url      = "jdbc:mysql://localhost:3306/testdb"
  db_driver   = "com.mysql.cj.jdbc.Driver"
  credentials = stonebranch_credential.test.name
}
`, connName, credName)
}

func testAccTaskStoredProcedureConfig_basic(name, connName, credName string) string {
	return sbacctest.ProviderConfig() + testAccTaskStoredProcedureConfig_dependencies(connName, credName) + fmt.Sprintf(`
resource "stonebranch_task_stored_procedure" "test" {
  name                = %[1]q
  stored_proc_name    = "test_procedure"
  database_connection = stonebranch_database_connection.test.name
}
`, name)
}

func testAccTaskStoredProcedureConfig_updated(name, connName, credName string) string {
	return sbacctest.ProviderConfig() + testAccTaskStoredProcedureConfig_dependencies(connName, credName) + fmt.Sprintf(`
resource "stonebranch_task_stored_procedure" "test" {
  name                = %[1]q
  stored_proc_name    = "updated_procedure"
  database_connection = stonebranch_database_connection.test.name
  summary             = "Updated stored procedure task"
}
`, name)
}

func testAccTaskStoredProcedureConfig_withSummary(name, connName, credName, summary string) string {
	return sbacctest.ProviderConfig() + testAccTaskStoredProcedureConfig_dependencies(connName, credName) + fmt.Sprintf(`
resource "stonebranch_task_stored_procedure" "test" {
  name                = %[1]q
  stored_proc_name    = "test_procedure"
  database_connection = stonebranch_database_connection.test.name
  summary             = %[2]q
}
`, name, summary)
}

func testAccTaskStoredProcedureConfig_withParameters(name, connName, credName string) string {
	return sbacctest.ProviderConfig() + testAccTaskStoredProcedureConfig_dependencies(connName, credName) + fmt.Sprintf(`
resource "stonebranch_task_stored_procedure" "test" {
  name                = %[1]q
  stored_proc_name    = "test_with_params"
  database_connection = stonebranch_database_connection.test.name

  parameters = [
    {
      param_var   = "input_param_1"
      param_mode  = "Input"
      param_type  = "VARCHAR"
      input_value = "test_value"
      position    = 1
    },
    {
      param_var  = "output_param_1"
      param_mode = "Output"
      param_type = "INTEGER"
      position   = 2
    }
  ]
}
`, name)
}

func testAccTaskStoredProcedureConfig_withRetryConfig(name, connName, credName string) string {
	return sbacctest.ProviderConfig() + testAccTaskStoredProcedureConfig_dependencies(connName, credName) + fmt.Sprintf(`
resource "stonebranch_task_stored_procedure" "test" {
  name                = %[1]q
  stored_proc_name    = "test_procedure"
  database_connection = stonebranch_database_connection.test.name
  retry_maximum       = 3
  retry_interval      = 60
}
`, name)
}

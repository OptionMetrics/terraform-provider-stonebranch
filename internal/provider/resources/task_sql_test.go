package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sbacctest "github.com/OptionMetrics/terraform-provider-stonebranch/internal/acctest"
)

func TestAccTaskSQLResource_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test-sql")
	connName := acctest.RandomWithPrefix("tf-test-dbconn")
	credName := acctest.RandomWithPrefix("tf-test-cred")
	resourceName := "stonebranch_task_sql.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccTaskSQLConfig_basic(rName, connName, credName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "sql_command", "SELECT 1"),
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
				Config: testAccTaskSQLConfig_updated(rName, connName, credName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "sql_command", "SELECT * FROM users"),
					resource.TestCheckResourceAttr(resourceName, "summary", "Updated SQL task"),
				),
			},
		},
	})
}

func TestAccTaskSQLResource_withSummary(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test-sql")
	connName := acctest.RandomWithPrefix("tf-test-dbconn")
	credName := acctest.RandomWithPrefix("tf-test-cred")
	resourceName := "stonebranch_task_sql.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTaskSQLConfig_withSummary(rName, connName, credName, "Initial summary"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "summary", "Initial summary"),
				),
			},
			{
				Config: testAccTaskSQLConfig_withSummary(rName, connName, credName, "Changed summary"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "summary", "Changed summary"),
				),
			},
		},
	})
}

// Test configuration helpers

// testAccTaskSQLConfig_dependencies creates the required database connection and credential
func testAccTaskSQLConfig_dependencies(connName, credName string) string {
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

func testAccTaskSQLConfig_basic(name, connName, credName string) string {
	return sbacctest.ProviderConfig() + testAccTaskSQLConfig_dependencies(connName, credName) + fmt.Sprintf(`
resource "stonebranch_task_sql" "test" {
  name                = %[1]q
  database_connection = stonebranch_database_connection.test.name
  sql_command         = "SELECT 1"
}
`, name)
}

func testAccTaskSQLConfig_updated(name, connName, credName string) string {
	return sbacctest.ProviderConfig() + testAccTaskSQLConfig_dependencies(connName, credName) + fmt.Sprintf(`
resource "stonebranch_task_sql" "test" {
  name                = %[1]q
  database_connection = stonebranch_database_connection.test.name
  sql_command         = "SELECT * FROM users"
  summary             = "Updated SQL task"
}
`, name)
}

func testAccTaskSQLConfig_withSummary(name, connName, credName, summary string) string {
	return sbacctest.ProviderConfig() + testAccTaskSQLConfig_dependencies(connName, credName) + fmt.Sprintf(`
resource "stonebranch_task_sql" "test" {
  name                = %[1]q
  database_connection = stonebranch_database_connection.test.name
  sql_command         = "SELECT 1"
  summary             = %[2]q
}
`, name, summary)
}

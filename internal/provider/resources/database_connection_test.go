package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sbacctest "github.com/OptionMetrics/terraform-provider-stonebranch/internal/acctest"
)

func TestAccDatabaseConnectionResource_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test-dbconn")
	resourceName := "stonebranch_database_connection.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccDatabaseConnectionConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "db_type", "MySQL"),
					resource.TestCheckResourceAttr(resourceName, "db_url", "jdbc:mysql://localhost:3306/testdb"),
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
				Config: testAccDatabaseConnectionConfig_updated(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "db_type", "PostgreSQL"),
					resource.TestCheckResourceAttr(resourceName, "db_url", "jdbc:postgresql://localhost:5432/testdb"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated database connection"),
				),
			},
		},
	})
}

func TestAccDatabaseConnectionResource_withDescription(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test-dbconn")
	resourceName := "stonebranch_database_connection.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseConnectionConfig_withDescription(rName, "Initial description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "Initial description"),
				),
			},
			{
				Config: testAccDatabaseConnectionConfig_withDescription(rName, "Changed description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "Changed description"),
				),
			},
		},
	})
}

func TestAccDatabaseConnectionResource_withMaxRows(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test-dbconn")
	resourceName := "stonebranch_database_connection.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseConnectionConfig_withMaxRows(rName, 1000),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "db_max_rows", "1000"),
				),
			},
			{
				Config: testAccDatabaseConnectionConfig_withMaxRows(rName, 5000),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "db_max_rows", "5000"),
				),
			},
		},
	})
}

func TestAccDatabaseConnectionResource_minimal(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test-dbconn")
	resourceName := "stonebranch_database_connection.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseConnectionConfig_minimal(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "sys_id"),
				),
			},
		},
	})
}

// Test configuration helpers

func testAccDatabaseConnectionConfig_basic(name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_database_connection" "test" {
  name      = %[1]q
  db_type   = "MySQL"
  db_url    = "jdbc:mysql://localhost:3306/testdb"
  db_driver = "com.mysql.cj.jdbc.Driver"
}
`, name)
}

func testAccDatabaseConnectionConfig_updated(name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_database_connection" "test" {
  name        = %[1]q
  db_type     = "PostgreSQL"
  db_url      = "jdbc:postgresql://localhost:5432/testdb"
  db_driver   = "org.postgresql.Driver"
  description = "Updated database connection"
}
`, name)
}

func testAccDatabaseConnectionConfig_withDescription(name, description string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_database_connection" "test" {
  name        = %[1]q
  db_type     = "MySQL"
  db_url      = "jdbc:mysql://localhost:3306/testdb"
  db_driver   = "com.mysql.cj.jdbc.Driver"
  description = %[2]q
}
`, name, description)
}

func testAccDatabaseConnectionConfig_withMaxRows(name string, maxRows int) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_database_connection" "test" {
  name        = %[1]q
  db_type     = "MySQL"
  db_url      = "jdbc:mysql://localhost:3306/testdb"
  db_driver   = "com.mysql.cj.jdbc.Driver"
  db_max_rows = %[2]d
}
`, name, maxRows)
}

func testAccDatabaseConnectionConfig_minimal(name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_database_connection" "test" {
  name      = %[1]q
  db_url    = "jdbc:mysql://localhost:3306/testdb"
  db_driver = "com.mysql.cj.jdbc.Driver"
}
`, name)
}

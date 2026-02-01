package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sbacctest "terraform-provider-stonebranch/internal/acctest"
)

func TestAccTaskWebServiceResource_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test-ws")
	resourceName := "stonebranch_task_web_service.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccTaskWebServiceConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "url", "https://httpbin.org/get"),
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
				Config: testAccTaskWebServiceConfig_updated(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "url", "https://httpbin.org/post"),
					resource.TestCheckResourceAttr(resourceName, "http_method", "POST"),
					resource.TestCheckResourceAttr(resourceName, "summary", "Updated web service task"),
				),
			},
		},
	})
}

func TestAccTaskWebServiceResource_withSummary(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test-ws")
	resourceName := "stonebranch_task_web_service.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTaskWebServiceConfig_withSummary(rName, "Initial summary"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "summary", "Initial summary"),
				),
			},
			{
				Config: testAccTaskWebServiceConfig_withSummary(rName, "Changed summary"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "summary", "Changed summary"),
				),
			},
		},
	})
}

func TestAccTaskWebServiceResource_withHeaders(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test-ws")
	resourceName := "stonebranch_task_web_service.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTaskWebServiceConfig_withHeaders(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "http_headers.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "http_headers.0.name", "Content-Type"),
					resource.TestCheckResourceAttr(resourceName, "http_headers.0.value", "application/json"),
					resource.TestCheckResourceAttr(resourceName, "http_headers.1.name", "Accept"),
					resource.TestCheckResourceAttr(resourceName, "http_headers.1.value", "application/json"),
				),
			},
		},
	})
}

func TestAccTaskWebServiceResource_withPayload(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test-ws")
	resourceName := "stonebranch_task_web_service.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTaskWebServiceConfig_withPayload(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "http_method", "POST"),
					resource.TestCheckResourceAttr(resourceName, "payload", `{"key": "value"}`),
				),
			},
		},
	})
}

// Test configuration helpers

func testAccTaskWebServiceConfig_basic(name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_web_service" "test" {
  name      = %[1]q
  url       = "https://httpbin.org/get"
  mime_type = "application/json"
}
`, name)
}

func testAccTaskWebServiceConfig_updated(name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_web_service" "test" {
  name        = %[1]q
  url         = "https://httpbin.org/post"
  http_method = "POST"
  mime_type   = "application/json"
  summary     = "Updated web service task"
}
`, name)
}

func testAccTaskWebServiceConfig_withSummary(name, summary string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_web_service" "test" {
  name      = %[1]q
  url       = "https://httpbin.org/get"
  mime_type = "application/json"
  summary   = %[2]q
}
`, name, summary)
}

func testAccTaskWebServiceConfig_withHeaders(name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_web_service" "test" {
  name      = %[1]q
  url       = "https://httpbin.org/get"
  mime_type = "application/json"

  http_headers = [
    {
      name  = "Content-Type"
      value = "application/json"
    },
    {
      name  = "Accept"
      value = "application/json"
    }
  ]
}
`, name)
}

func testAccTaskWebServiceConfig_withPayload(name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_web_service" "test" {
  name        = %[1]q
  url         = "https://httpbin.org/post"
  http_method = "POST"
  mime_type   = "application/json"
  payload     = "{\"key\": \"value\"}"
}
`, name)
}

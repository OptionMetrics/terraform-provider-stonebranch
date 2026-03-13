package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sbacctest "github.com/OptionMetrics/terraform-provider-stonebranch/internal/acctest"
)

func TestAccTaskUniversalAwsS3Resource_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test")
	resourceName := "stonebranch_task_universal_aws_s3.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccTaskUniversalAwsS3Config_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "action", "list-buckets"),
					resource.TestCheckResourceAttr(resourceName, "aws_default_region", "us-east-1"),
					resource.TestCheckResourceAttrSet(resourceName, "sys_id"),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
				),
			},
			// ImportState - use the task name as import ID
			{
				ResourceName:                         resourceName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        rName,
				ImportStateVerifyIdentifierAttribute: "name",
			},
			// Update
			{
				Config: testAccTaskUniversalAwsS3Config_updated(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "action", "list-objects"),
					resource.TestCheckResourceAttr(resourceName, "bucket", "test-bucket"),
					resource.TestCheckResourceAttr(resourceName, "summary", "Updated task summary"),
				),
			},
		},
	})
}

func TestAccTaskUniversalAwsS3Resource_uploadFile(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test")
	resourceName := "stonebranch_task_universal_aws_s3.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTaskUniversalAwsS3Config_uploadFile(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "action", "upload-file"),
					resource.TestCheckResourceAttr(resourceName, "bucket", "my-bucket"),
					resource.TestCheckResourceAttr(resourceName, "sourcefile", "/tmp/test.txt"),
					resource.TestCheckResourceAttr(resourceName, "prefix", "uploads/"),
					resource.TestCheckResourceAttr(resourceName, "upload_write_options", "True"),
				),
			},
		},
	})
}

func TestAccTaskUniversalAwsS3Resource_downloadFile(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test")
	resourceName := "stonebranch_task_universal_aws_s3.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTaskUniversalAwsS3Config_downloadFile(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "action", "download-file"),
					resource.TestCheckResourceAttr(resourceName, "bucket", "my-bucket"),
					resource.TestCheckResourceAttr(resourceName, "s3_key", "data/file.txt"),
					resource.TestCheckResourceAttr(resourceName, "target_directory", "/tmp/downloads"),
					resource.TestCheckResourceAttr(resourceName, "download_write_options", "False"),
				),
			},
		},
	})
}

func TestAccTaskUniversalAwsS3Resource_copyObject(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test")
	resourceName := "stonebranch_task_universal_aws_s3.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTaskUniversalAwsS3Config_copyObject(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "action", "copy-object-to-bucket"),
					resource.TestCheckResourceAttr(resourceName, "bucket", "source-bucket"),
					resource.TestCheckResourceAttr(resourceName, "target_bucket", "dest-bucket"),
					resource.TestCheckResourceAttr(resourceName, "operation", "copy"),
				),
			},
		},
	})
}

func TestAccTaskUniversalAwsS3Resource_withRoleBasedAccess(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test")
	resourceName := "stonebranch_task_universal_aws_s3.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sbacctest.PreCheck(t) },
		ProtoV6ProviderFactories: sbacctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTaskUniversalAwsS3Config_roleBasedAccess(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "role_based_access", "yes"),
					resource.TestCheckResourceAttr(resourceName, "role_arn", "arn:aws:iam::123456789012:role/TestRole"),
					resource.TestCheckResourceAttr(resourceName, "service_name", "sts"),
				),
			},
		},
	})
}

// Test configuration helpers

func testAccTaskUniversalAwsS3Config_basic(name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_universal_aws_s3" "test" {
  name               = %[1]q
  agent_var          = "agent_name"
  action             = "list-buckets"
  aws_default_region = "us-east-1"
}
`, name)
}

func testAccTaskUniversalAwsS3Config_updated(name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_universal_aws_s3" "test" {
  name               = %[1]q
  summary            = "Updated task summary"
  agent_var          = "agent_name"
  action             = "list-objects"
  bucket             = "test-bucket"
  aws_default_region = "us-east-1"
}
`, name)
}

func testAccTaskUniversalAwsS3Config_uploadFile(name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_universal_aws_s3" "test" {
  name                 = %[1]q
  agent_var            = "agent_name"
  action               = "upload-file"
  bucket               = "my-bucket"
  sourcefile           = "/tmp/test.txt"
  prefix               = "uploads/"
  upload_write_options = "True"
  aws_default_region   = "us-east-1"
}
`, name)
}

func testAccTaskUniversalAwsS3Config_downloadFile(name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_universal_aws_s3" "test" {
  name                   = %[1]q
  agent_var              = "agent_name"
  action                 = "download-file"
  bucket                 = "my-bucket"
  s3_key                 = "data/file.txt"
  target_directory       = "/tmp/downloads"
  download_write_options = "False"
  aws_default_region     = "us-east-1"
}
`, name)
}

func testAccTaskUniversalAwsS3Config_copyObject(name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_universal_aws_s3" "test" {
  name               = %[1]q
  agent_var          = "agent_name"
  action             = "copy-object-to-bucket"
  bucket             = "source-bucket"
  s3_key             = "path/to/file.txt"
  target_bucket      = "dest-bucket"
  target_s3_key      = "archive/file.txt"
  operation          = "copy"
  aws_default_region = "us-east-1"
}
`, name)
}

func testAccTaskUniversalAwsS3Config_roleBasedAccess(name string) string {
	return sbacctest.ProviderConfig() + fmt.Sprintf(`
resource "stonebranch_task_universal_aws_s3" "test" {
  name               = %[1]q
  agent_var          = "agent_name"
  action             = "list-buckets"
  role_based_access  = "yes"
  role_arn           = "arn:aws:iam::123456789012:role/TestRole"
  service_name       = "sts"
  aws_default_region = "us-east-1"
}
`, name)
}

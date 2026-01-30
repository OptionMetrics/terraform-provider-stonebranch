package provider

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/joho/godotenv"
)

func init() {
	// Load .env file from project root for all tests in this package
	loadEnv()
}

// loadEnv loads the .env file from the project root.
func loadEnv() {
	// Get the directory of this source file
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return
	}

	// Navigate from internal/provider/ to project root
	projectRoot := filepath.Join(filepath.Dir(filename), "..", "..")
	envPath := filepath.Join(projectRoot, ".env")

	// Load if exists, don't fail if it doesn't
	if _, err := os.Stat(envPath); err == nil {
		_ = godotenv.Load(envPath)
	}
}

// testAccProtoV6ProviderFactories is used for acceptance tests.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"stonebranch": providerserver.NewProtocol6WithError(New("test")()),
}

// testAccPreCheck validates required environment variables are set.
func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("STONEBRANCH_API_TOKEN"); v == "" {
		t.Fatal("STONEBRANCH_API_TOKEN must be set for acceptance tests")
	}
}

// providerConfig returns a Terraform provider configuration block.
// Uses environment variables so no secrets in code.
func providerConfig() string {
	return `
provider "stonebranch" {
  # Configured via STONEBRANCH_API_TOKEN and STONEBRANCH_BASE_URL environment variables
}
`
}

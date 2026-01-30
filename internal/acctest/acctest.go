// Package acctest provides acceptance test helpers for the Stonebranch provider.
package acctest

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/joho/godotenv"

	"terraform-provider-stonebranch/internal/provider"
)

func init() {
	// Load .env file from project root for all tests
	LoadEnv()
}

// LoadEnv loads the .env file from the project root.
func LoadEnv() {
	// Get the directory of this source file
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return
	}

	// Navigate from internal/acctest/ to project root
	projectRoot := filepath.Join(filepath.Dir(filename), "..", "..")
	envPath := filepath.Join(projectRoot, ".env")

	// Load if exists, don't fail if it doesn't
	if _, err := os.Stat(envPath); err == nil {
		_ = godotenv.Load(envPath)
	}
}

// ProtoV6ProviderFactories is used for acceptance tests.
var ProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"stonebranch": providerserver.NewProtocol6WithError(provider.New("test")()),
}

// PreCheck validates required environment variables are set.
func PreCheck(t *testing.T) {
	if v := os.Getenv("STONEBRANCH_API_TOKEN"); v == "" {
		t.Fatal("STONEBRANCH_API_TOKEN must be set for acceptance tests")
	}
}

// ProviderConfig returns a Terraform provider configuration block.
// Uses environment variables so no secrets in code.
func ProviderConfig() string {
	return `
provider "stonebranch" {
  # Configured via STONEBRANCH_API_TOKEN and STONEBRANCH_BASE_URL environment variables
}
`
}

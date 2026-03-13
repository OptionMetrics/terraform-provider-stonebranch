// sb2tf exports StoneBranch Universal Controller resources to Terraform configuration files.
package main

import (
	"os"

	"github.com/OptionMetrics/terraform-provider-stonebranch/cmd/sb2tf/cli"
)

// version is set via ldflags during build
var version = "dev"

func main() {
	cli.SetVersion(version)
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}

//go:build tools

package tools

import (
	// Documentation generation
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"

	// Test formatting
	_ "github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt"
)

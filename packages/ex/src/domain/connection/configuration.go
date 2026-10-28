// Package connection defines the versioned console-to-project configuration.
// The runtime.md contract itself lives in the Exphere contracts; this package
// names it for the CLI and adds the supervision state machine.
package connection

import "github.com/jonatassales/tao/exphere/contracts/configuration"

// Configuration is the concrete contract inside an authored Markdown document.
type Configuration = configuration.Configuration

// Component references an existing source directory and an explicit build recipe.
type Component = configuration.Component

// Build is a direct executable invocation, never implicit shell text.
type Build = configuration.Build

// Native preserves an existing v2 descriptor and its owned resources on migration.
type Native = configuration.Native

// ErrInvalidConfiguration classifies every authored-configuration rejection.
var ErrInvalidConfiguration = configuration.ErrInvalidConfiguration

// ErrRecipeRefused marks a build recipe outside the allowlist policy.
var ErrRecipeRefused = configuration.ErrRecipeRefused

// Parse requires one explicit langdrift JSON fence; see the contract.
func Parse(markdown string) (Configuration, error) {
	return configuration.Parse(markdown)
}

// Package connection defines the versioned console-to-project configuration.
// The runtime.md contract itself lives in the Microruntimes contracts; this package
// names it for the CLI and adds the supervision state machine.
package connection

import (
	"errors"

	"github.com/jonatassales/tao/packages/contracts/configuration"
)

// Configuration is the concrete contract inside an authored Markdown document.
type Configuration = configuration.Configuration

// Component references a folder in the checkout and the Dockerfile in it.
type Component = configuration.Component

// Readiness is a component's explicit readiness kind.
type Readiness = configuration.Readiness

// Native preserves an existing v2 descriptor and its owned resources on migration.
type Native = configuration.Native

// ErrInvalidConfiguration classifies every authored-configuration rejection.
var ErrInvalidConfiguration = configuration.ErrInvalidConfiguration

// Parse requires one explicit langdrift JSON fence; see the contract.
func Parse(markdown string) (Configuration, error) {
	return configuration.Parse(markdown)
}

// ErrPassEnvMissing refuses a start whose component passes through a
// variable the agent's environment does not hold.
var ErrPassEnvMissing = errors.New("a pass-through variable is not set in the agent's environment")

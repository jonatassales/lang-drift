// Package support holds the cohesive utilities the command tree and the
// catalog share without importing each other: options and streams, global
// flags and project selection, prompting, result rendering, and usage errors.
package support

import (
	"io"
	"strings"
)

// Options configures the CLI presentation layer.
type Options struct {
	Version   string
	Directory string

	Input       io.Reader
	Output      io.Writer
	Diagnostics io.Writer
}

// WithDefaults initializes optional I/O streams.
func (options Options) WithDefaults() Options {
	if options.Input == nil {
		options.Input = strings.NewReader("")
	}

	if options.Output == nil {
		options.Output = io.Discard
	}

	if options.Diagnostics == nil {
		options.Diagnostics = io.Discard
	}

	return options
}

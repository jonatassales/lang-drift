package main

import (
	"bytes"
	"context"
	"strings"

	"github.com/jonatassales/tao/lib/cli/src/presentation/cli/support"
)

// dftOnATerminal is dftWith with standard output on a terminal, so the
// session is toned and titles the window.
func (p project) dftOnATerminal(deps dependencies, arguments ...string) (int, string, string) {
	p.t.Helper()

	var stdout, stderr bytes.Buffer

	cases, native := useCasesWith(deps)
	status := executeNativeCatalogs(context.Background(), arguments, support.Options{
		Version:     version,
		Directory:   p.root,
		Input:       strings.NewReader(""),
		Output:      &stdout,
		Diagnostics: &stderr,
		Terminal:    true,
	}, cases, native)

	return status, stdout.String(), stderr.String()
}

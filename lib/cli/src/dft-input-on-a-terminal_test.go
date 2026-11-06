package main

import (
	"bytes"
	"context"
	"strings"

	"github.com/jonatassales/tao/lib/cli/src/presentation/cli/support"
)

// dftInputOnATerminal is dftInput with standard output on a terminal, so the
// session is toned, titles the window and reads confirmations from input.
func (projectInstance project) dftInputOnATerminal(deps dependencies, input string, arguments ...string) (int, string, string) {
	projectInstance.t.Helper()

	var stdout, stderr bytes.Buffer

	cases, native := useCasesWith(deps)
	status := executeNativeCatalogs(context.Background(), arguments, support.Options{
		Version:     version,
		Directory:   projectInstance.root,
		Input:       strings.NewReader(input),
		Output:      &stdout,
		Diagnostics: &stderr,
		Terminal:    true,
	}, cases, native)

	return status, stdout.String(), stderr.String()
}

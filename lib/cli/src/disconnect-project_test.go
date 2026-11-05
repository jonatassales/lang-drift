package main

import (
	"os"
	"path/filepath"
	"testing"
)

// recordedConsole is where dft agent connect keeps the console, relative to
// the project root.
const recordedConsole = ".langdrift/agent/state/credentials/console.env"

// disconnectProject removes the recorded console, as if the project had never
// connected.
func disconnectProject(t *testing.T, root string) {
	t.Helper()

	if err := os.Remove(filepath.Join(root, filepath.FromSlash(recordedConsole))); err != nil {
		t.Fatal(err)
	}
}

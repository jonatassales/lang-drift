package main

import (
	"path/filepath"
	"strings"
	"testing"
)

// On a terminal, the run titles the window too, opens on its scope dimmed
// without echoing the typed command, and tones the paused gate as a
// warning: "console › gate       reservation $40.80 exceeds the $25.00
// budget · paused".
func TestTheRunOnATerminalIsTitledAndToned(t *testing.T) {
	t.Parallel()

	planned, deps := plannedProject(t)
	deps.Lifecycle = websiteLifecycle()

	_, stdout, stderr := planned.dftOnATerminal(deps, "agent", "run", "console-microruntime", "--model", "claude-opus-5-5", "--effort", "high")

	title := "\x1b]0;dft — " + filepath.Base(planned.root) + "\x07"
	opens := strings.HasPrefix(stdout, title+"\x1b[2m# run · ")
	paused := strings.Contains(stdout, "\x1b[33mconsole › gate       reservation $40.80 exceeds the $25.00 budget · paused\x1b[0m\n")
	if !opens || !paused || strings.Contains(stdout, "dft agent run") {
		t.Errorf("the run printed %q (stderr %s)", stdout, stderr)
	}
}

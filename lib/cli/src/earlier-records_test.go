package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jonatassales/tao/lib/cli/src/domain/finops"
)

// A ledger written by 0.4.0, when the command was still ex, decodes as
// before: its provenance names ex and nothing checks that value.
func TestLedgerOfVersion040Decodes(t *testing.T) {
	t.Parallel()

	original, err := os.ReadFile(filepath.Join("testdata", "langdrift", "v0.4.0", "finops", "ledger", "console-microruntime.ledger.json"))
	if err != nil {
		t.Fatal(err)
	}

	var ledger finops.Ledger

	decodeStrict(t, original, &ledger)

	if ledger.Mirror.SetBy != "ex" {
		t.Fatalf("setBy = %q, want the ex of 0.4.0", ledger.Mirror.SetBy)
	}
}

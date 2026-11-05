package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jonatassales/tao/lib/cli/src/domain/connection"
	"github.com/jonatassales/tao/lib/cli/src/domain/finops"
	"github.com/jonatassales/tao/lib/cli/src/domain/harness"
	"github.com/jonatassales/tao/lib/cli/src/domain/run"
	"github.com/jonatassales/tao/lib/cli/src/infrastructure/records"
)

// Every record of the sample decodes into its type, with unknown fields
// refused, and encodes back to the same bytes through the formatter. This
// pins the field order, the names, the null typing and both layouts to the
// owner's reference.
func TestSampleRecordsRoundTrip(t *testing.T) {
	t.Parallel()

	files := []struct {
		path  string
		fresh func() any
	}{
		{"finops/budgets.json", func() any { return &finops.Budget{} }},
		{"finops/forecasts.json", func() any { return &finops.Forecasts{} }},
		{"finops/gates.json", func() any { return &finops.Gates{} }},
		{"finops/ledger/0001-console-microruntime.ledger.json", func() any { return &finops.Ledger{} }},
		{"finops/prices/0001-illustrative.prices.json", func() any { return &finops.Prices{} }},
		{"harness/envelope.json", func() any { return &harness.Envelope{} }},
		{"harness/workflow.json", func() any { return &harness.Workflow{} }},
		{"microruntimes/runs/0001-console-microruntime/plan.json", func() any { return &run.Plan{} }},
		{"microruntimes/runs/0001-console-microruntime/run.json", func() any { return &run.Run{} }},
	}

	for _, file := range files {
		t.Run(file.path, func(t *testing.T) {
			t.Parallel()

			original, err := os.ReadFile(filepath.Join(sampleFixture, filepath.FromSlash(file.path)))
			if err != nil {
				t.Fatal(err)
			}

			value := file.fresh()
			decodeStrict(t, original, value)

			formatted, err := records.Format(value, records.FileProfile)
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(formatted, original) {
				t.Fatalf("round trip changed the file:\n%s", firstDifference(original, formatted))
			}
		})
	}
}

func TestSampleFencesRoundTrip(t *testing.T) {
	t.Parallel()

	fences := []struct {
		path  string
		fence string
		fresh func() any
	}{
		{"microruntimes/runs/0001-console-microruntime/0001-console-microruntime.run.md", "langdrift-run", func() any { return &run.RunFence{} }},
		{"microruntimes/runs/0001-console-microruntime/laps/0001.lap.md", "langdrift-lap", func() any { return &run.LapFence{} }},
		{"microruntimes/runs/0001-console-microruntime/laps/0002.lap.md", "langdrift-lap", func() any { return &run.LapFence{} }},
	}

	for _, fence := range fences {
		t.Run(fence.path, func(t *testing.T) {
			t.Parallel()

			body := fenceBody(t, filepath.Join(sampleFixture, filepath.FromSlash(fence.path)), fence.fence)

			value := fence.fresh()
			decodeStrict(t, []byte(body), value)

			formatted, err := records.Format(value, records.FenceProfile)
			if err != nil {
				t.Fatal(err)
			}

			if string(formatted) != body+"\n" {
				t.Fatalf("round trip changed the fence:\n%s", firstDifference([]byte(body+"\n"), formatted))
			}
		})
	}
}

// The runtime.md fence is the contracts' own rendering: the configuration
// marshaled with two spaces, every object expanded.
func TestSampleRuntimeFenceIsTheContractsRendering(t *testing.T) {
	t.Parallel()

	body := fenceBody(t, sampleRuntime, "langdrift")

	config, err := connection.Parse(string(mustRead(t, sampleRuntime)))
	if err != nil {
		t.Fatal(err)
	}

	rendered, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	if string(rendered) != body {
		t.Fatalf("the contracts render the fence differently:\n%s", firstDifference([]byte(body), rendered))
	}
}

func decodeStrict(t *testing.T, data []byte, value any) {
	t.Helper()

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(value); err != nil {
		t.Fatalf("decode: %v", err)
	}
}

func fenceBody(t *testing.T, path, language string) string {
	t.Helper()

	_, after, found := strings.Cut(string(mustRead(t, path)), "```"+language+"\n")
	if !found {
		t.Fatalf("%s has no %s fence", path, language)
	}

	body, _, closed := strings.Cut(after, "\n```")
	if !closed {
		t.Fatalf("%s: the %s fence is not closed", path, language)
	}

	return body
}

func mustRead(t *testing.T, path string) []byte {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	return data
}

// firstDifference shows the lines around the first byte that differs.
func firstDifference(want, got []byte) string {
	limit := min(len(want), len(got))
	offset := 0
	for offset < limit && want[offset] == got[offset] {
		offset++
	}

	start := max(0, offset-120)

	return "want: ..." + string(want[start:min(len(want), offset+120)]) + "\n got: ..." + string(got[start:min(len(got), offset+120)])
}

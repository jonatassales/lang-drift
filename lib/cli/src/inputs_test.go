package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/jonatassales/tao/lib/cli/src/domain/finops"
	"github.com/jonatassales/tao/lib/cli/src/domain/harness"
	"github.com/jonatassales/tao/lib/cli/src/presentation/cli/support"
)

// dftWith is dft with scripted sources: the ids and the clock of the records,
// and the driver that plays the stages.
func (p project) dftWith(deps dependencies, arguments ...string) (int, string, string) {
	p.t.Helper()

	return p.dftInput(deps, "", arguments...)
}

// dftInput is dftWith with what the person types.
func (p project) dftInput(deps dependencies, input string, arguments ...string) (int, string, string) {
	p.t.Helper()

	var stdout, stderr bytes.Buffer
	cases, native := useCasesWith(deps)
	status := executeNativeCatalogs(context.Background(), arguments, support.Options{
		Version:     version,
		Directory:   p.root,
		Input:       strings.NewReader(input),
		Output:      &stdout,
		Diagnostics: &stderr,
	}, cases, native)

	return status, stdout.String(), stderr.String()
}

// scriptedIDs hands out the given identifiers in order.
func scriptedIDs(ids ...string) func() (string, error) {
	remaining := append([]string{}, ids...)

	return func() (string, error) {
		id := remaining[0]
		remaining = remaining[1:]

		return id, nil
	}
}

// bodyFile writes an input body for a set command, taken from a sample
// mirror with its mirror block and version left out.
func bodyFile[Body any](t *testing.T, root, name, sample string, strip func(*Body)) string {
	t.Helper()

	var body Body
	decodeStrict(t, mustRead(t, filepath.Join(sampleFixture, filepath.FromSlash(sample))), &body)
	strip(&body)

	encoded, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(root, name)
	if err := os.WriteFile(path, encoded, 0o644); err != nil {
		t.Fatal(err)
	}

	return path
}

// withoutMirror zeroes a record's mirror block and version, so two records
// compare by their input alone.
func withoutMirror(value any) {
	record := reflect.ValueOf(value).Elem()
	record.FieldByName("Mirror").Set(reflect.Zero(reflect.TypeOf(harness.Mirror{})))
	record.FieldByName("Version").SetInt(0)
}

func sameInput[Body any](t *testing.T, written, sample string) {
	t.Helper()

	var got, want Body
	decodeStrict(t, mustRead(t, written), &got)
	decodeStrict(t, mustRead(t, filepath.Join(sampleFixture, filepath.FromSlash(sample))), &want)
	withoutMirror(&got)
	withoutMirror(&want)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("%s differs from the sample:\n got %+v\nwant %+v", written, got, want)
	}
}

// The eight set commands, given the sample's values, produce the sample's
// mirrors, apart from the provenance the sample attributes to a console.
func TestInputsMirrorTheSample(t *testing.T) {
	t.Parallel()

	p := newProject(t)
	deps := dependencies{
		IDs:   scriptedIDs("a759267ef752f98d6c25c4114ad65d0c", "7c2a0000000000000000000000000000", "3f810000000000000000000000000000", "c0de0000000000000000000000000000", "d1ab0000000000000000000000000000", "9d4e0000000000000000000000000000", "51b00000000000000000000000000000", "0000000000000000000000000000ffff", "a6130000000000000000000000000000"),
		Clock: func() time.Time { return time.Date(2026, 10, 27, 18, 4, 0, 0, time.UTC) },
	}

	workflow := bodyFile(t, p.root, "workflow.json", "harness/workflow.json", func(body *harness.Workflow) { withoutMirror(body) })
	prices := bodyFile(t, p.root, "prices.json", "finops/prices/0001-illustrative.prices.json", func(body *finops.Prices) { withoutMirror(body) })
	forecasts := bodyFile(t, p.root, "forecasts.json", "finops/forecasts.json", func(body *finops.Forecasts) { withoutMirror(body) })
	driver := filepath.Join(p.root, "driver.json")
	if err := os.WriteFile(driver, []byte(`{"executable":"claude","args":["-p","{{prompt}}","--output-format","json"],"passEnv":["ANTHROPIC_API_KEY"],"timeoutSeconds":1800,"usage":{"inputTokens":"usage.input_tokens","outputTokens":"usage.output_tokens","costCents":""},"unmeasured":"pause"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	commands := [][]string{
		{"init"},
		{"harness", "envelope", "set", "--per-execution-cpu-millis", "4000", "--per-execution-memory-mib", "6144", "--default-cpu-millis", "2000", "--default-memory-mib", "3072", "--not-runtime-field", "gpu", "--not-runtime-field", "disk"},
		{"harness", "workflow", "set", "--file", workflow},
		{"harness", "model", "set", "claude-sonnet-5-5", "--effort", "medium"},
		{"harness", "driver", "set", "--file", driver},
		{"finops", "budget", "set", "--period", "2026-10", "--period-cents", "50000", "--per-run-cents", "2500", "--reserve-margin-percent", "20", "--allocate", "0070-console-microruntime=2500"},
		{"finops", "gates", "set", "--default"},
		{"finops", "prices", "set", "--file", prices, "--slug", "illustrative"},
		{"finops", "forecasts", "set", "--file", forecasts},
	}

	for _, arguments := range commands {
		if status, stdout, stderr := p.dftWith(deps, arguments...); status != 0 || (arguments[0] != "init" && !strings.Contains(stdout, `"status": "set"`)) {
			t.Fatalf("dft %v = %d\n%s%s", arguments, status, stdout, stderr)
		}
	}

	folder := filepath.Join(p.root, ".langdrift")
	sameInput[harness.Envelope](t, filepath.Join(folder, "harness", "envelope.json"), "harness/envelope.json")
	sameInput[harness.Workflow](t, filepath.Join(folder, "harness", "workflow.json"), "harness/workflow.json")
	sameInput[finops.Budget](t, filepath.Join(folder, "finops", "budgets.json"), "finops/budgets.json")
	sameInput[finops.Gates](t, filepath.Join(folder, "finops", "gates.json"), "finops/gates.json")
	sameInput[finops.Prices](t, filepath.Join(folder, "finops", "prices", "0001-illustrative.prices.json"), "finops/prices/0001-illustrative.prices.json")
	sameInput[finops.Forecasts](t, filepath.Join(folder, "finops", "forecasts.json"), "finops/forecasts.json")

	var envelope harness.Envelope
	decodeStrict(t, mustRead(t, filepath.Join(folder, "harness", "envelope.json")), &envelope)
	if envelope.Mirror.Record != "envelope-7c2a" || envelope.Mirror.Via != "cli" || envelope.Mirror.Of != "cli" || envelope.Mirror.SetBy != "owner" {
		t.Errorf("envelope provenance = %+v", envelope.Mirror)
	}

	if visible := visibleFiles(t, folder); len(visible) != 8 {
		t.Errorf("the folder holds %v, want the identity, the rules and six mirrors", visible)
	}
}

func TestInputsRefuseBadValuesAsUsage(t *testing.T) {
	t.Parallel()

	p := newProject(t)
	if status, _, stderr := p.dft("", "init"); status != 0 {
		t.Fatal(stderr)
	}

	duplicate := filepath.Join(p.root, "workflow.json")
	if err := os.WriteFile(duplicate, []byte(`{"stages":[{"id":"planning","label":"Planning","skills":[],"human":"required"},{"id":"planning","label":"Again","skills":[],"human":"optional"}],"defaults":{"humanAt":[]},"efforts":["low"]}`), 0o644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		arguments []string
		stderr    string
	}{
		{arguments: []string{"harness", "workflow", "set", "--file", duplicate}, stderr: "distinct"},
		{arguments: []string{"finops", "gates", "set"}, stderr: "--default or --file"},
		{arguments: []string{"finops", "budget", "set", "--period", "2026-10", "--period-cents", "1000", "--per-run-cents", "500", "--allocate", "a=900", "--allocate", "b=900"}, stderr: "exceed"},
		{arguments: []string{"harness", "envelope", "set", "--per-execution-cpu-millis", "100", "--per-execution-memory-mib", "100", "--default-cpu-millis", "200", "--default-memory-mib", "100"}, stderr: "ceiling"},
	}

	for _, tt := range tests {
		if status, _, stderr := p.dft("", tt.arguments...); status != 2 || !strings.Contains(stderr, tt.stderr) {
			t.Errorf("dft %v = %d, %q; want exit 2 mentioning %q", tt.arguments, status, stderr, tt.stderr)
		}
	}
}

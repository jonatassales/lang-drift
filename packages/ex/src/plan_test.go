package main

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/jonatassales/tao/packages/ex/src/domain/connection"
	"github.com/jonatassales/tao/packages/ex/src/domain/finops"
	"github.com/jonatassales/tao/packages/ex/src/domain/harness"
	"github.com/jonatassales/tao/packages/ex/src/domain/run"
)

// seedProject copies the frozen fixture project, the ticket 0070 and the six
// sources its links reach, into a disposable project with the sample's
// inputs recorded, and returns it with its scripted sources.
func seedProject(t *testing.T) (project, dependencies) {
	t.Helper()

	p := newProject(t)
	copyTree(t, projectFixture, p.root)

	deps := dependencies{
		IDs:   scriptedIDs("a759267ef752f98d6c25c4114ad65d0c", "7c2a0000000000000000000000000000", "3f810000000000000000000000000000", "c0de0000000000000000000000000000", "9d4e0000000000000000000000000000", "51b00000000000000000000000000000", "0000000000000000000000000000ffff", "a6130000000000000000000000000000", "d21f0000000000000000000000000000", "701706be7e70d681d608826d3ad548b9", "e1e10000000000000000000000000000", "e2e20000000000000000000000000000", "e3e30000000000000000000000000000", "e4e40000000000000000000000000000", "e5e50000000000000000000000000000", "e6e60000000000000000000000000000", "e7e70000000000000000000000000000"),
		Clock: func() time.Time { return time.Date(2026, 10, 28, 13, 46, 50, 0, time.UTC) },
	}

	workflow := bodyFile(t, p.root, "workflow.json", "harness/workflow.json", func(body *harness.Workflow) { withoutMirror(body) })
	prices := bodyFile(t, p.root, "prices.json", "finops/prices/0001-illustrative.prices.json", func(body *finops.Prices) { withoutMirror(body) })
	forecasts := bodyFile(t, p.root, "forecasts.json", "finops/forecasts.json", func(body *finops.Forecasts) {
		withoutMirror(body)
		body.Planning = &finops.TokensThousands{Input: 35, Output: 5}
	})
	driver := driverFile(t, p.root, "pause")

	for _, arguments := range [][]string{
		{"init"},
		{"harness", "envelope", "set", "--per-execution-cpu-millis", "4000", "--per-execution-memory-mib", "6144", "--default-cpu-millis", "2000", "--default-memory-mib", "3072", "--not-runtime-field", "gpu", "--not-runtime-field", "disk"},
		{"harness", "workflow", "set", "--file", workflow},
		{"harness", "model", "set", "claude-sonnet-5-5", "--effort", "medium"},
		{"finops", "budget", "set", "--period", "2026-10", "--period-cents", "50000", "--per-run-cents", "2500", "--reserve-margin-percent", "20", "--allocate", "0070-console-microruntime=2500"},
		{"finops", "gates", "set", "--default"},
		{"finops", "prices", "set", "--file", prices, "--slug", "illustrative"},
		{"finops", "forecasts", "set", "--file", forecasts},
		{"harness", "driver", "set", "--file", driver},
	} {
		if status, stdout, stderr := p.exWith(deps, arguments...); status != 0 {
			t.Fatalf("ex %v = %d\n%s%s", arguments, status, stdout, stderr)
		}
	}

	return p, deps
}

func copyTree(t *testing.T, from, to string) {
	t.Helper()

	err := filepath.WalkDir(from, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative, err := filepath.Rel(from, path)
		if err != nil {
			return err
		}
		target := filepath.Join(to, relative)
		if entry.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(target, content, 0o644)
	})
	if err != nil {
		t.Fatal(err)
	}
}

const consoleComponent = "path=.,name=console,kind=web,image=langdrift-console:local,port=3001,readiness=/icon.svg,build=platform/microruntime/console/build-image.sh,tech=nextjs:Next.js;go:Go"

// The plan's first act reproduces the sample: the ticket, the sources its
// links reach with their hashes, the context hash, the inputs and the
// discrepancies, and the runtime file parses with the next revision.
func TestPlanReadsTheTicketAndWritesTheRuntimeFile(t *testing.T) {
	t.Parallel()

	p, deps := seedProject(t)

	status, stdout, stderr := p.exWith(deps, "agent", "plan", ".agents/workflow/tickets/0070-console-microruntime.ticket.md", "--no-judge", "--component", consoleComponent)
	if status != 0 {
		t.Fatalf("ex agent plan = %d\n%s%s", status, stdout, stderr)
	}

	for _, line := range []string{
		"ex agent plan 0070-console-microruntime.ticket.md",
		"agent › tool call: read 0070-console-microruntime.ticket.md",
		"agent › needs-triage · links 0010-exphere-and-the-langdrift-folder.spec.md · ex-microruntimes-capability.research.md · 0002-console-input-is-untrusted.adr.md · 0001-finops.gr.md · package.json · go.mod",
		"✓ read       7 documents · nothing outside the ticket's links",
		"reads      .langdrift/harness/guardrails/0001-finops.gr.md · FIN-01 · FIN-02 · FIN-03 · FIN-04 · FIN-05 · FIN-06 · FIN-07 · FIN-08",
		"bounded    cpuMillis 2000 · memoryMiB 3072 · set by owner in the CLI · record envelope-7c2a",
		"budget     $25.00 per run for model calls · set by owner in the CLI",
		"agent › resolving 1 components named in the CLI · paths, images, ports",
		"agent › tool call: inspect . · no manifest · :3001",
		"✓ .                     Next.js+Go · web 500m · 1024 MiB",
		"agent › cross-checking · 1 ports distinct · 1 images resolved · no collisions",
		"✓ size       500m of 2000m · 1024 of 3072 MiB · 1 one-replica Deployments · proposed shares",
		"limits     gpu and disk are not runtime.md fields · none requested",
		"estimate   runtime has no price source · cpu and memory stay in millicores and MiB",
		"agent › tool call: validate runtime.md draft against the schema · 0 errors",
		"✓ plan       runtime.md draft, revision 1 · proposed · authorize with ex agent run",
		"# next: once the team authorizes, ex agent run applies the recorded intent",
	} {
		if !strings.Contains(stdout, line+"\n") {
			t.Errorf("the session did not print %q:\n%s", line, stdout)
		}
	}

	runtimePath := filepath.Join(p.root, ".langdrift", "microruntimes", "runtimes", "console-microruntime", "runtime.md")
	markdown := mustRead(t, runtimePath)

	config, err := connection.Parse(string(markdown))
	if err != nil {
		t.Fatalf("the runtime file does not parse: %v", err)
	}
	if config.RuntimeID != "console-microruntime" || config.Revision != 1 || config.CPUMillis != 2000 || config.Model != "claude-sonnet-5-5" || len(config.Components) != 1 || config.ProjectID != "a759267ef752f98d6c25c4114ad65d0c" {
		t.Errorf("fence = %+v", config)
	}

	var written run.Plan
	decodeStrict(t, []byte(fenceBody(t, runtimePath, "langdrift-plan")), &written)

	var sample run.Plan
	decodeStrict(t, mustRead(t, samplePlan), &sample)

	if !reflect.DeepEqual(written.Ticket.Path, sample.Ticket.Path) || *written.Ticket.Number != "0070" || written.Ticket.Slug != "console-microruntime" || *written.Ticket.Title != *sample.Ticket.Title || *written.Ticket.Status != *sample.Ticket.Status {
		t.Errorf("ticket = %+v, want the sample's", written.Ticket)
	}

	if written.Collection.Agent != nil || written.Collection.Limits != sample.Collection.Limits || written.Collection.Reached != sample.Collection.Reached || written.Collection.Context != sample.Collection.Context {
		t.Errorf("collection = %+v\nwant %+v", written.Collection, sample.Collection)
	}

	if len(written.Sources) != len(sample.Sources) {
		t.Fatalf("%d sources, the sample has %d", len(written.Sources), len(sample.Sources))
	}
	for index, source := range written.Sources {
		want := sample.Sources[index]
		if source.Path != want.Path || *source.MediaType != *want.MediaType || *source.Bytes != *want.Bytes || *source.SHA256 != *want.SHA256 || source.Depth != want.Depth || source.State != want.State {
			t.Errorf("source %d = %+v, want %+v", index, source, want)
		}
		if (source.LinkedFrom == nil) != (want.LinkedFrom == nil) || (source.LinkedFrom != nil && *source.LinkedFrom != *want.LinkedFrom) {
			t.Errorf("source %d linked from %v, want %v", index, source.LinkedFrom, want.LinkedFrom)
		}
		if want.LinkSyntax != "named in the console" && source.LinkSyntax != want.LinkSyntax {
			t.Errorf("source %d syntax %q, want %q", index, source.LinkSyntax, want.LinkSyntax)
		}
		if source.Verdict != nil || source.Reason != nil || source.Confidence != nil {
			t.Errorf("source %d was judged without a judge", index)
		}
	}

	if written.Inputs.CPUMillis.Value != 2000 || written.Inputs.MemoryMiB.Value != 3072 || written.Inputs.ModelBudgetCents.Value != 2500 || written.Inputs.Ceiling != sample.Inputs.Ceiling && written.Inputs.Ceiling.CPUMillis != 4000 {
		t.Errorf("inputs = %+v", written.Inputs)
	}
	if !reflect.DeepEqual(written.Inputs.NotRuntimeFields, sample.Inputs.NotRuntimeFields) || len(written.Discrepancies) != 0 {
		t.Errorf("notRuntimeFields = %v, discrepancies = %v", written.Inputs.NotRuntimeFields, written.Discrepancies)
	}

	if written.State != "proposed" || written.Runtime.Revision != 1 || written.Runtime.SHA256 == "" || len(written.Facts) != 0 || written.Suggestion.By != "input" || *written.Suggestion.ModelName != "Claude Sonnet 5.5" {
		t.Errorf("plan state = %+v, suggestion = %+v", written.Runtime, written.Suggestion)
	}

	var record run.Run
	decodeStrict(t, mustRead(t, filepath.Join(filepath.Dir(runtimePath), "run.json")), &record)
	if record.Kind != "ticket" || record.State != "planned" || len(record.Stages) != 6 || !record.Stages[0].Human || record.Stages[1].Human {
		t.Errorf("run record = %+v", record)
	}

	var ledger finops.Ledger
	decodeStrict(t, mustRead(t, filepath.Join(p.root, ".langdrift", "finops", "ledger", "console-microruntime.ledger.json")), &ledger)
	if *ledger.Ticket != "0070-console-microruntime" || ledger.Bounded.TicketAllocationCents != 2500 || ledger.Planning.Cents.Observed != 0 || ledger.Planning.Model != nil || *ledger.PriceSource != ".langdrift/finops/prices/0001-illustrative.prices.json" {
		t.Errorf("ledger = %+v", ledger)
	}

	status, _, stderr = p.exWith(deps, "agent", "plan", ".agents/workflow/tickets/0070-console-microruntime.ticket.md", "--no-judge", "--component", consoleComponent)
	if status != 0 {
		t.Fatalf("replan = %d: %s", status, stderr)
	}
	again, err := connection.Parse(string(mustRead(t, runtimePath)))
	if err != nil || again.Revision != 2 || again.RuntimeUID != config.RuntimeUID {
		t.Errorf("replan gave revision %d with uid %s, want revision 2 with the same uid", again.Revision, again.RuntimeUID)
	}
}

func TestPlanRefusesWhatItCannotRead(t *testing.T) {
	t.Parallel()

	p, deps := seedProject(t)

	if status, _, stderr := p.exWith(deps, "agent", "plan", ".agents/workflow/tickets/0099-missing.ticket.md", "--no-judge", "--component", consoleComponent); status != 4 || !strings.Contains(stderr, "no ticket") {
		t.Errorf("a missing ticket gave %d: %s", status, stderr)
	}

	if status, _, stderr := p.exWith(deps, "agent", "plan", ".agents/workflow/tickets/0070-console-microruntime.ticket.md", "--no-judge"); status != 2 || !strings.Contains(stderr, "--component") {
		t.Errorf("a plan without components gave %d: %s", status, stderr)
	}

	bare := newProject(t)
	if status, _, stderr := bare.ex("", "init"); status != 0 {
		t.Fatal(stderr)
	}
	copyTree(t, projectFixture, bare.root)
	if status, _, stderr := bare.ex("", "agent", "plan", ".agents/workflow/tickets/0070-console-microruntime.ticket.md", "--no-judge", "--component", consoleComponent); status != 4 || !strings.Contains(stderr, "envelope") {
		t.Errorf("a plan without inputs gave %d: %s", status, stderr)
	}
}

// A draft that fails a cross-check is written as blocked, with the failing
// check recorded, and the command refuses.
func TestPlanBlocksADraftThatFailsACheck(t *testing.T) {
	t.Parallel()

	p, deps := seedProject(t)

	status, stdout, stderr := p.exWith(deps, "agent", "plan", ".agents/workflow/tickets/0070-console-microruntime.ticket.md", "--no-judge",
		"--component", "path=.,image=langdrift-console:local,port=3001,readiness=/icon.svg,build=rm,args=-rf /",
		"--component", "path=apps/console,image=console:local,port=3001,readiness=/,build=docker,args=build -t console:local .")
	if status != 3 || !strings.Contains(stderr, "ports-distinct") || !strings.Contains(stderr, "recipe-allowed") {
		t.Fatalf("a failing draft gave %d: %s\n%s", status, stderr, stdout)
	}

	if !strings.Contains(stdout, "✗ plan       runtime.md draft, revision 1 · blocked · schema · ports-distinct · recipe-allowed failed\n") {
		t.Errorf("the blocked line is missing:\n%s", stdout)
	}

	runtimePath := filepath.Join(p.root, ".langdrift", "microruntimes", "runtimes", "console-microruntime", "runtime.md")

	var written run.Plan
	decodeStrict(t, []byte(fenceBody(t, runtimePath, "langdrift-plan")), &written)

	if written.State != "blocked" || written.Runtime.State != "blocked" {
		t.Errorf("plan state = %s, runtime state = %s", written.State, written.Runtime.State)
	}

	results := map[string]string{}
	for _, check := range written.Checks {
		results[check.ID] = check.Result
	}
	if results["schema"] != "failed" || results["ports-distinct"] != "failed" || results["recipe-allowed"] != "failed" || results["images-resolved"] != "passed" || results["within-ceiling"] != "passed" {
		t.Errorf("checks = %v", results)
	}

	if written.Totals.Proposed.CPUMillis != 1000 || written.Totals.Headroom.CPUMillis != 1000 || written.Components[1].Executed.CPUMillis != 1000 || written.Components[1].InferredFrom[0] != "apps/console/package.json" {
		t.Errorf("components = %+v, totals = %+v", written.Components, written.Totals)
	}

	var record run.Run
	decodeStrict(t, mustRead(t, filepath.Join(filepath.Dir(runtimePath), "run.json")), &record)
	if record.State != "blocked" {
		t.Errorf("run state = %s", record.State)
	}

	// A replan still finds the revision, even though the fence failed validation.
	if status, _, stderr := p.exWith(deps, "agent", "plan", ".agents/workflow/tickets/0070-console-microruntime.ticket.md", "--no-judge", "--component", consoleComponent); status != 0 {
		t.Fatalf("the replan after a blocked draft = %d: %s", status, stderr)
	}
	again, err := connection.Parse(string(mustRead(t, runtimePath)))
	if err != nil || again.Revision != 2 {
		t.Errorf("the replan gave revision %d, %v; want 2", again.Revision, err)
	}
}

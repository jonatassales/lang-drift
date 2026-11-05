package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jonatassales/tao/lib/cli/src/domain/rag"
	"github.com/jonatassales/tao/lib/cli/src/domain/run"
)

// scriptedEmbedder places each known episode at a chosen cosine from the
// query, which sits on the first axis.
type scriptedEmbedder struct {
	similarities map[string]float64
}

func (scriptedEmbedder) ID() string { return "scripted" }

func (embedder scriptedEmbedder) Embed(text string) []float32 {
	vector := make([]float32, rag.Dimensions)
	vector[0] = 1

	for title, similarity := range embedder.similarities {
		if strings.HasPrefix(text, title) {
			vector[0] = float32(similarity)
			vector[1] = float32(math.Sqrt(1 - similarity*similarity))
		}
	}

	return vector
}

// writeMemory writes a runtime's memory with its fence, so the index has an episode.
func writeMemory(t *testing.T, root, runtime, title, stacks, paths string, measured bool, cpu, memory int) {
	t.Helper()

	folder := filepath.Join(root, ".langdrift", "microruntimes", "runtimes", runtime)
	if err := os.MkdirAll(folder, 0o755); err != nil {
		t.Fatal(err)
	}

	observed := "null"
	if measured {
		observed = fmt.Sprintf(`{ "cpuMillisPeak": %d, "memoryMiBPeak": %d, "tokensThousands": { "input": 3700, "output": 240 } }`, cpu, memory)
	}

	content := fmt.Sprintf("# %s\n\n```langdrift-run\n{\n  \"run\": %q,\n  \"ticket\": \"0001-%s\",\n  \"kind\": \"ticket\",\n  \"outcome\": \"completed\",\n  \"laps\": 1,\n  \"componentPaths\": [%s],\n  \"stacks\": [%s],\n  \"measuredResources\": %t,\n  \"measuredTokens\": %t,\n  \"observed\": %s,\n  \"facts\": []\n}\n```\n", title, runtime, runtime, paths, stacks, measured, measured, observed)
	if err := os.WriteFile(filepath.Join(folder, "run.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// With the website's six episodes indexed, the plan's second act prints the
// website's figures: six candidates from 0.81 to 0.42, three kept above the
// threshold, none with measured resources.
func TestPlanRetrievesTheWebsitesEpisodes(t *testing.T) {
	t.Parallel()

	p, deps := seedProject(t)

	episodes := []struct {
		runtime, title string
		similarity     float64
	}{
		{"checkout-preview-sizing", "Size the checkout preview", 0.81},
		{"notify-api-rollout", "Roll out the notify API", 0.74},
		{"signup-redesign-preview", "Preview the sign-up redesign", 0.66},
		{"docs-search-index", "Index the docs for search", 0.52},
		{"batch-export-worker", "Run the batch export worker", 0.47},
		{"storefront-cache-warmup", "Warm the storefront cache", 0.42},
	}
	similarities := map[string]float64{}
	for _, episode := range episodes {
		similarities[episode.title] = episode.similarity
		writeMemory(t, p.root, episode.runtime, episode.title, `"python"`, `"docs"`, false, 0, 0)
	}
	deps.Embedder = scriptedEmbedder{similarities: similarities}

	status, stdout, stderr := p.dftWith(deps, "agent", "plan", ".agents/workflow/tickets/0070-console-microruntime.ticket.md", "--no-judge", "--component", consoleComponent)
	if status != 0 {
		t.Fatalf("dft agent plan = %d\n%s%s", status, stdout, stderr)
	}

	for _, line := range []string{
		"agent › RAG over project episodes · a model and effort for this ticket",
		"agent › embedding the ticket intent · 1 query vector",
		"agent › vector search · project episodes · top 6 by cosine similarity",
		"agent › 6 candidates · similarity 0.81 to 0.42",
		"agent › judging · 3 share the stack and scale · 3 differ: below 0.55",
		"agent › dropped 3 below 0.55 · kept 3",
		"agent › 3 episodes cited · none with measured resources",
		"agent › reasoning · no cited episode measured cpu or memory · sizing follows the inputs, not history",
		"✓ suggests   claude-sonnet-5-5 · effort medium",
	} {
		if !strings.Contains(stdout, line+"\n") {
			t.Errorf("the session did not print %q:\n%s", line, stdout)
		}
	}

	runtimePath := filepath.Join(p.root, ".langdrift", "microruntimes", "runtimes", "console-microruntime", "runtime.md")
	var written run.Plan
	decodeStrict(t, []byte(fenceBody(t, runtimePath, "langdrift-plan")), &written)

	considered := written.Retrieval.Considered
	if len(considered) != 6 || considered[0].Run != "checkout-preview-sizing" || considered[0].Similarity != 0.81 || considered[0].Verdict != "kept" || considered[5].Run != "storefront-cache-warmup" || considered[5].Verdict != "dropped" || considered[3].Reason == "" {
		t.Errorf("considered = %+v", considered)
	}
	if written.Retrieval.Method != "vector+structured" || written.Sizing.Basis != "input" || len(written.Sizing.MeasuredRunsCited) != 0 {
		t.Errorf("retrieval = %+v, sizing = %+v", written.Retrieval, written.Sizing)
	}
}

// A run that measured resources is cited by the next plan, which sizes
// from its peaks.
func TestPlanCitesAMeasuredRun(t *testing.T) {
	t.Parallel()

	p, deps := seedProject(t)
	writeMemory(t, p.root, "console-preview", "Run the console with extub as a microruntime", `"nextjs", "go"`, `"."`, true, 410, 690)

	if status, stdout, stderr := p.dftWith(deps, "agent", "reindex"); status != 0 || !strings.Contains(stdout, `"episodes": 1`) {
		t.Fatalf("dft agent reindex = %d\n%s%s", status, stdout, stderr)
	}

	status, stdout, stderr := p.dftWith(deps, "agent", "plan", ".agents/workflow/tickets/0070-console-microruntime.ticket.md", "--no-judge", "--component", consoleComponent)
	if status != 0 {
		t.Fatalf("dft agent plan = %d\n%s%s", status, stdout, stderr)
	}

	for _, line := range []string{
		"agent › 1 candidates · similarity",
		"agent › 1 episodes cited · 1 with measured resources",
		"agent › reasoning · 1 cited episodes measured cpu and memory · sizing follows their peaks with headroom",
		"✓ .                     Next.js+Go · web 650m · 1152 MiB",
	} {
		if !strings.Contains(stdout, line) {
			t.Errorf("the session did not print %q:\n%s", line, stdout)
		}
	}

	runtimePath := filepath.Join(p.root, ".langdrift", "microruntimes", "runtimes", "console-microruntime", "runtime.md")
	var written run.Plan
	decodeStrict(t, []byte(fenceBody(t, runtimePath, "langdrift-plan")), &written)

	if written.Sizing.Basis != "measured-runs" || len(written.Sizing.MeasuredRunsCited) != 1 || written.Sizing.MeasuredRunsCited[0] != "console-preview" {
		t.Errorf("sizing = %+v", written.Sizing)
	}
	if len(written.Retrieval.Considered) != 1 || len(written.Retrieval.Considered[0].Criteria) != 4 || written.Retrieval.Considered[0].Score != 1 {
		t.Errorf("considered = %+v", written.Retrieval.Considered)
	}
}

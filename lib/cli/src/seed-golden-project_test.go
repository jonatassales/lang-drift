package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/jonatassales/tao/lib/cli/src/application/execution"
	"github.com/jonatassales/tao/lib/cli/src/application/plan"
	"github.com/jonatassales/tao/lib/cli/src/domain/finops"
	"github.com/jonatassales/tao/lib/cli/src/domain/harness"
	"github.com/jonatassales/tao/lib/cli/src/domain/run"
)

const goldenProjectFixture = "testdata/project"

var goldenPlanArguments = []string{
	"agent",
	"plan",
	".agents/0031-TICKET-book-a-demo-runtime.ticket.md",
	"--component", "name=signup-api,path=services/signup-api,kind=api,port=8080,tech=go:Go",
	"--component", "name=notify-api,path=services/notify-api,kind=api,port=8081,tech=nodejs:Node.js",
	"--component", "name=nats,path=platform/nats,kind=broker,port=4222,tech=nats:NATS",
	"--component", "name=redis,path=platform/redis,kind=cache,port=6379,tech=redis:Redis",
	"--component", "name=signup,path=apps/signup,kind=web,port=3000,tech=nextjs:Next.js",
}

// seedGoldenProject copies testdata/project into a folder named "project",
// initializes a Git repository with one commit and a remote, and seeds the
// inputs, episodes and fake console until specs 0050 to 0066 land.
func seedGoldenProject(t *testing.T, humanAt []string) (project, dependencies) {
	t.Helper()

	baseDir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	root := filepath.Join(baseDir, "project")
	copyTree(t, goldenProjectFixture, root)

	for _, component := range []struct {
		folder string
		port   string
	}{
		{
			folder: "services/signup-api",
			port:   "8080",
		},
		{
			folder: "services/notify-api",
			port:   "8081",
		},
		{
			folder: "platform/nats",
			port:   "4222",
		},
		{
			folder: "platform/redis",
			port:   "6379",
		},
		{
			folder: "apps/signup",
			port:   "3000",
		},
	} {
		dockerfile := filepath.Join(root, filepath.FromSlash(component.folder), "Dockerfile")
		content := "FROM scratch\nEXPOSE " + component.port + "\nHEALTHCHECK CMD true\n"
		if writeErr := os.WriteFile(dockerfile, []byte(content), 0o644); writeErr != nil {
			t.Fatal(writeErr)
		}
	}

	projectInstance := project{
		t:    t,
		root: root,
	}

	similarities := map[string]float64{}
	for _, episode := range []struct {
		runtime    string
		title      string
		similarity float64
	}{
		{
			runtime:    "checkout-preview-sizing",
			title:      "Size the checkout preview",
			similarity: 0.81,
		},
		{
			runtime:    "notify-api-rollout",
			title:      "Roll out the notify API",
			similarity: 0.74,
		},
		{
			runtime:    "signup-redesign-preview",
			title:      "Preview the sign-up redesign",
			similarity: 0.66,
		},
		{
			runtime:    "docs-search-index",
			title:      "Index the docs for search",
			similarity: 0.52,
		},
		{
			runtime:    "batch-export-worker",
			title:      "Run the batch export worker",
			similarity: 0.47,
		},
		{
			runtime:    "storefront-cache-warmup",
			title:      "Warm the storefront cache",
			similarity: 0.42,
		},
	} {
		similarities[episode.title] = episode.similarity
		writeMemory(t, root, episode.runtime, episode.title, `"python"`, `"docs"`, false, 0, 0)
	}

	sourcesJSON, err := json.Marshal(map[string]any{
		"sources": []map[string]any{
			{
				"path":       ".agents/0031-TICKET-book-a-demo-runtime.ticket.md",
				"verdict":    "kept",
				"reason":     "the ticket",
				"confidence": 0.95,
				"role":       "ticket",
				"summary":    "ticket",
			},
			{
				"path":       ".agents/0002-PRD-book-a-demo.prd.md",
				"verdict":    "kept",
				"reason":     "states the use case",
				"confidence": 0.92,
				"role":       "prd",
				"summary":    "use case: book-a-demo request",
			},
			{
				"path":       ".agents/0005-DD-book-a-demo-runtime.dd.md",
				"verdict":    "kept",
				"reason":     "sets the envelope and budget",
				"confidence": 0.94,
				"role":       "design-doc",
				"summary":    "runtime envelope and per-run budget",
			},
			{
				"path":       ".agents/0011-SPEC-book-a-demo-runtime.spec.md",
				"verdict":    "kept",
				"reason":     "names the components and acceptance criteria",
				"confidence": 0.96,
				"role":       "spec",
				"summary":    "5 components, 3 acceptance criteria",
			},
			{
				"path":       ".langdrift/harness/guardrails/0001-finops.gr.md",
				"verdict":    "kept",
				"reason":     "finops guardrail",
				"confidence": 0.95,
				"role":       "guardrail",
				"summary":    "guardrail",
			},
		},
		"facts": []any{},
		"criteria": []map[string]string{
			{
				"id":      "FIN-02",
				"applies": "estimate before a paid call",
			},
			{
				"id":      "FIN-03",
				"applies": "a budget before a paid run",
			},
			{
				"id":      "FIN-04",
				"applies": "size cpu and memory to the load",
			},
		},
		"figures": []any{},
	})
	if err != nil {
		t.Fatal(err)
	}

	episodesJSON, err := json.Marshal(map[string]any{
		"episodes": []map[string]string{
			{
				"run":    "docs-search-index",
				"reason": "other runtime",
			},
			{
				"run":    "batch-export-worker",
				"reason": "no broker",
			},
			{
				"run":    "storefront-cache-warmup",
				"reason": "docs only",
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	suggestionJSON, err := json.Marshal(map[string]any{
		"model":        "claude-sonnet-5-5",
		"effort":       "medium",
		"effortReason": "low skips the broker wiring · high is not warranted by 5 components",
		"modelReason":  "Haiku 5.5 retries the broker wiring · Opus 5.5 costs more per outcome",
		"reasons": []string{
			"low skips the broker wiring · high is not warranted by 5 components",
			"Haiku 5.5 retries the broker wiring · Opus 5.5 costs more per outcome",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	console := newFakeConsole(t, map[string]string{
		plan.SourcesQuestion:    string(sourcesJSON),
		plan.EpisodesQuestion:   string(episodesJSON),
		plan.SuggestionQuestion: string(suggestionJSON),
	})

	pendingChecks := []run.CheckResult{
		{
			Name:   "build",
			Result: "pending",
		},
		{
			Name:   "check",
			Result: "pending",
		},
		{
			Name:   "test",
			Result: "pending",
		},
		{
			Name:   "guardrails",
			Result: "pending",
		},
	}
	passedChecks := []run.CheckResult{
		{
			Name:   "build",
			Result: "passed",
		},
		{
			Name:   "check",
			Result: "passed",
		},
		{
			Name:   "test",
			Result: "passed",
		},
		{
			Name:   "guardrails",
			Result: "passed",
		},
	}

	deps := dependencies{
		Clock:    func() time.Time { return time.Date(2026, 10, 28, 13, 46, 50, 0, time.UTC) },
		Embedder: scriptedEmbedder{similarities: similarities},
		Driver:   websiteDriver(),
		Lifecycle: &scriptedLifecycle{
			execution: "exec-5c1e9a",
			cleanup: execution.Cleanup{
				Confirmed:   true,
				Removed:     5,
				Total:       5,
				ClusterKept: true,
			},
		},
		Delivery: &scriptedDelivery{
			branch: execution.Branch{
				Name:     "feat/0031-book-a-demo-runtime",
				Commit:   "3f9c2e1",
				Upstream: true,
			},
			pullRequest: &run.PullRequest{
				Number: 19,
				Title:  "feat(demo): run book a demo as a microruntime",
				State:  "open",
			},
			checks: [][]run.CheckResult{pendingChecks, passedChecks},
		},
		Poll: time.Millisecond,
	}

	optionalHumanStages := map[string]bool{
		"planning":  true,
		"spec":      true,
		"tickets":   true,
		"implement": true,
		"review":    true,
	}
	workflow := bodyFile(t, baseDir, "workflow.json", "harness/workflow.json", func(body *harness.Workflow) {
		withoutMirror(body)
		body.Defaults.HumanAt = humanAt
		for index := range body.Stages {
			if optionalHumanStages[body.Stages[index].ID] {
				body.Stages[index].Human = harness.HumanOptional
			}
		}
	})
	prices := bodyFile(t, baseDir, "prices.json", "finops/prices/0001-illustrative.prices.json", func(body *finops.Prices) {
		withoutMirror(body)
	})
	forecasts := bodyFile(t, baseDir, "forecasts.json", "finops/forecasts.json", func(body *finops.Forecasts) {
		withoutMirror(body)
		body.Planning = &finops.TokensThousands{
			Input:  35,
			Output: 5,
		}
	})
	driver := driverFile(t, baseDir, "pause")

	for _, arguments := range [][]string{
		{"init"},
		{"harness", "envelope", "set", "--per-execution-cpu-millis", "4000", "--per-execution-memory-mib", "6144", "--default-cpu-millis", "2000", "--default-memory-mib", "3072", "--not-runtime-field", "gpu", "--not-runtime-field", "disk"},
		{"harness", "workflow", "set", "--file", workflow},
		{"harness", "model", "set", "claude-sonnet-5-5", "--effort", "medium"},
		{"finops", "budget", "set", "--period", "2026-10", "--period-cents", "50000", "--per-run-cents", "2500", "--reserve-margin-percent", "20", "--allocate", "0031-TICKET-book-a-demo-runtime=2500"},
		{"finops", "gates", "set", "--default"},
		{"finops", "prices", "set", "--file", prices, "--slug", "illustrative"},
		{"finops", "forecasts", "set", "--file", forecasts},
		{"harness", "driver", "set", "--file", driver},
	} {
		if status, stdout, stderr := projectInstance.dftWith(deps, arguments...); status != 0 {
			t.Fatalf("dft %v = %d\n%s%s", arguments, status, stdout, stderr)
		}
	}

	guardrailSource := filepath.Join(goldenProjectFixture, ".langdrift", "harness", "guardrails", "0001-finops.gr.md")
	guardrailTarget := filepath.Join(root, ".langdrift", "harness", "guardrails", "0001-finops.gr.md")
	if writeErr := os.WriteFile(guardrailTarget, mustRead(t, guardrailSource), 0o644); writeErr != nil {
		t.Fatal(writeErr)
	}

	connectProject(t, root, console.server.URL)

	for _, gitArguments := range [][]string{
		{"init"},
		{"config", "user.email", "test@langdrift.md"},
		{"config", "user.name", "Tester"},
		{"remote", "add", "origin", "https://github.com/example/project.git"},
		{"add", "."},
		{"commit", "-m", "chore: init"},
	} {
		command := exec.Command("git", gitArguments...)
		command.Dir = root
		if output, runErr := command.CombinedOutput(); runErr != nil {
			t.Fatalf("git %v: %v\n%s", gitArguments, runErr, output)
		}
	}

	return projectInstance, deps
}

package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jonatassales/tao/lib/cli/src/domain/run"
)

// dft agent connect checks RUNTIME_SECRET against the console's judge route
// and records the console only when the route accepts it; the plan then asks
// that route itself, each question once, as the contract's request, with the
// recorded secret as a bearer. RUNTIME_SECRET is the process's environment,
// so this test does not run in parallel.
func TestConnectThenPlanAsksTheConsolesJudge(t *testing.T) {
	console := newFakeConsole(t, minimalAnswers())

	seeded, deps := seedProject(t)
	deps.Judge = nil

	disconnectProject(t, seeded.root)
	recorded := filepath.Join(seeded.root, filepath.FromSlash(recordedConsole))

	t.Setenv("RUNTIME_SECRET", "not-the-deployment-secret")

	status, stdout, stderr := seeded.dftWith(deps, "agent", "connect", console.server.URL)
	refusedNamingTheConsole := status == 1 && strings.Contains(stderr, console.server.URL) && strings.Contains(stderr, "refused RUNTIME_SECRET")
	if !refusedNamingTheConsole {
		t.Fatalf("a refused secret gave %d: %s\n%s", status, stderr, stdout)
	}
	if _, err := os.Stat(recorded); err == nil {
		t.Fatal("a refused connection was recorded")
	}

	t.Setenv("RUNTIME_SECRET", testRuntimeSecret)

	status, stdout, stderr = seeded.dftWith(deps, "agent", "connect", console.server.URL)
	if status != 0 || !strings.Contains(stdout, "connected to "+console.server.URL) {
		t.Fatalf("dft agent connect = %d: %s\n%s", status, stderr, stdout)
	}
	if strings.Contains(stdout+stderr, testRuntimeSecret) {
		t.Errorf("connect printed the secret:\n%s%s", stdout, stderr)
	}
	if info, err := os.Stat(recorded); err != nil || info.Mode().Perm() != 0o600 {
		t.Fatalf("the console connection = %v, %v; want a private file", info, err)
	}

	status, stdout, stderr = seeded.dftWith(deps, "agent", "plan", consoleTicket, "--component", consoleComponent)
	if status != 0 {
		t.Fatalf("dft agent plan = %d: %s\n%s", status, stderr, stdout)
	}

	asked, bearers := console.questions()
	if len(asked) != 2 || asked[0].Kind != "sources" || asked[1].Kind != "suggestion" {
		t.Fatalf("the console was asked %+v", asked)
	}
	for index, question := range asked {
		carriesTheSecret := bearers[index] == "Bearer "+testRuntimeSecret
		asksThePlanner := question.Model == "claude-sonnet-5-5" && question.MaxTokens == 8000
		carriesThePrompts := question.PromptVersion == "judge-1" && question.System != "" && question.Prompt != ""

		if !carriesTheSecret || !asksThePlanner || !carriesThePrompts {
			t.Errorf("question %d = %+v with %q", index, question, bearers[index])
		}
	}

	runtimePath := filepath.Join(seeded.root, ".langdrift", "microruntimes", "runtimes", "console-microruntime", "runtime.md")
	var written run.Plan
	decodeStrict(t, []byte(fenceBody(t, runtimePath, "langdrift-plan")), &written)
	judgedByTheConsole := written.Judge != nil && written.Judge.Version == "console-judge-1"
	if written.State != "proposed" || written.Suggestion.By != "judge" || !judgedByTheConsole {
		t.Errorf("the judged plan = state %s, suggestion %+v, judge %+v", written.State, written.Suggestion, written.Judge)
	}

	ledger := ledgerOf(t, seeded)
	pricedByTheConsole := ledger.Planning.Provenance != nil && *ledger.Planning.Provenance == "console"
	if !pricedByTheConsole || ledger.Planning.Cents.Observed != 10 {
		t.Errorf("the planning line = %+v", ledger.Planning)
	}
}

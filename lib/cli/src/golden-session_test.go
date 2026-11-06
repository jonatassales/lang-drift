package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jonatassales/tao/lib/cli/src/application/session"
	"github.com/jonatassales/tao/lib/core/contracts/go/handoff"
)

var (
	planPending = map[string]string{
		"readTicket":      "0053",
		"links":           "0053",
		"readPrd":         "0053",
		"readDesignDoc":   "0053",
		"readSpec":        "0053",
		"readDocument":    "0053",
		"reads":           "0053",
		"guardrail":       "0054",
		"bounds":          "0054",
		"budget":          "0056",
		"resolve":         "0066",
		"reasonResources": "0052",
		"inspectNext":     "0066",
		"draft":           "0050",
	}
	runPending = map[string]string{
		"scope":            "0050",
		"budget":           "0056",
		"confirmed":        "0057",
		"implement":        "0059",
		"hook":             "0060",
		"tdd":              "0059",
		"devTicketWaiting": "0061",
		"devTicketClosed":  "0061",
		"runtimeReady":     "0062",
		"safeSeatMet":      "0061",
		"runtimeStopped":   "0062",
	}
	runOverBudgetPending = map[string]string{
		"scope":  "0050",
		"budget": "0056",
	}
	runAllStagesPending = map[string]string{
		"scope":              "0050",
		"budget":             "0056",
		"confirmed":          "0057",
		"confirmedSpec":      "0057",
		"confirmedTickets":   "0057",
		"confirmedImplement": "0057",
		"implement":          "0059",
		"hook":               "0060",
		"tdd":                "0059",
		"devTicketWaiting":   "0061",
		"devTicketClosed":    "0061",
		"runtimeReady":       "0062",
		"safeSeatMet":        "0061",
		"runtimeStopped":     "0062",
	}
	handoffPending = map[string]string{
		"plan.json":   "0050",
		"intent.json": "0050",
		"run.json":    "0050",
	}
)

// TestGoldenSessionMatchesTheSimulation runs dft in process against a fresh
// copy of testdata/project and compares every printed line, its kind and the
// hand-off records with the frozen goldens.
func TestGoldenSessionMatchesTheSimulation(t *testing.T) {
	t.Parallel()

	loadGolden := func(t *testing.T, name string, wantLines int) []goldenLine {
		t.Helper()

		data := mustRead(t, filepath.Join("testdata", "goldens", name+".json"))

		var lines []goldenLine
		if err := json.Unmarshal(data, &lines); err != nil {
			t.Fatalf("decode golden %s: %v", name, err)
		}

		if len(lines) != wantLines {
			t.Fatalf("golden %s has %d lines, want %d", name, len(lines), wantLines)
		}

		return lines
	}

	assertSession := func(t *testing.T, stdout string, typedCommand string, events []session.Event, golden []goldenLine, pending map[string]string) []printedLine {
		t.Helper()

		const expectedTitle = "\x1b]0;dft — project\x07"
		if !strings.HasPrefix(stdout, expectedTitle) {
			t.Fatalf("terminal output did not open with title %q: %q", expectedTitle, stdout[:min(80, len(stdout))])
		}

		if strings.Contains(stdout, typedCommand) {
			t.Errorf("session echoed typed command %q", typedCommand)
		}

		tones := map[session.Kind]string{
			session.Comment: "\x1b[2m",
			session.Success: "\x1b[32m",
			session.Warning: "\x1b[33m",
		}
		var expectedStdout strings.Builder
		expectedStdout.WriteString(expectedTitle)

		printed := make([]printedLine, 0, len(events))
		for index, event := range events {
			line := printedLineOf(event)
			printed = append(printed, line)

			tonedLine := line.text
			if tone, colored := tones[line.kind]; colored {
				tonedLine = tone + line.text + "\x1b[0m"
			}
			expectedStdout.WriteString(tonedLine + "\n")

			if strings.HasPrefix(line.text, "agent › ") {
				t.Errorf("line %d %q uses legacy actor prefix: %q", index+1, line.id, line.text)
			}

			if strings.Contains(line.text, " › ") && !strings.HasPrefix(line.text, "console › ") {
				t.Errorf("actor line %d %q does not start with console › : %q", index+1, line.id, line.text)
			}

			if !matchesLineShape(line.text) {
				t.Errorf("line %d %q does not match any of the six shapes: %q", index+1, line.id, line.text)
			}
		}

		if stdout != expectedStdout.String() {
			t.Errorf("terminal stdout bytes do not match printed session:\ngot:\n%q\nwant:\n%q", stdout, expectedStdout.String())
		}

		if differences := goldenDiff(golden, printed, pending); len(differences) > 0 {
			t.Errorf("golden differences:\n%s", strings.Join(differences, "\n"))
		}

		return printed
	}

	assertHandoffFile := func(t *testing.T, runtimeDir string, name string, verify func(string) error) {
		t.Helper()

		spec := handoffPending[name]
		raw, readErr := os.ReadFile(filepath.Join(runtimeDir, name))
		if readErr != nil {
			if spec == "" {
				t.Errorf("missing hand-off record %s: %v", name, readErr)
			}

			return
		}

		verifyErr := verify(string(raw))
		if verifyErr != nil {
			if spec == "" {
				t.Errorf("hand-off record %s does not match the four cards: %v", name, verifyErr)
			}

			return
		}

		if spec != "" {
			t.Errorf("%s already matches the four cards, remove from handoffPending under spec %s", name, spec)
		}
	}

	for _, componentDir := range []string{
		"services/signup-api",
		"services/notify-api",
		"platform/nats",
		"platform/redis",
		"apps/signup",
	} {
		dockerfilePath := filepath.Join(goldenProjectFixture, filepath.FromSlash(componentDir), "Dockerfile")
		if _, err := os.Stat(dockerfilePath); !os.IsNotExist(err) {
			t.Fatalf("fixture %s must have no Dockerfile", componentDir)
		}
	}

	t.Run("plan and default run", func(t *testing.T) {
		t.Parallel()

		planGolden := loadGolden(t, "plan", 42)
		runGolden := loadGolden(t, "run", 25)

		seeded, deps := seedGoldenProject(t, []string{"planning", "review"})
		planRecorder := &sessionRecorder{}
		deps.Observer = planRecorder

		status, stdout, stderr := seeded.dftOnATerminal(deps, goldenPlanArguments...)
		if status != 0 {
			t.Fatalf("dft agent plan = %d: %s\n%s", status, stderr, stdout)
		}

		printedPlan := assertSession(t, stdout, "dft agent plan", planRecorder.events, planGolden, planPending)

		scopeFirst := printedPlan[0].id == "scope" && printedPlan[0].kind == session.Comment
		analysisAfterLinks := printedPlan[3].id == "analysis" && printedPlan[2].id == "links"
		planAfterSuggest := eventPosition(planRecorder.events, session.Planning) == eventPosition(planRecorder.events, session.Suggest)+1
		handoffLast := printedPlan[len(printedPlan)-1].id == "handoff" && printedPlan[len(printedPlan)-1].kind == session.Comment
		if !scopeFirst || !analysisAfterLinks || !planAfterSuggest || !handoffLast {
			t.Errorf("the four comment lines of plan are out of position: first=%s, 4th=%s, last=%s", printedPlan[0].id, printedPlan[3].id, printedPlan[len(printedPlan)-1].id)
		}

		runtimeDir := filepath.Join(seeded.root, ".langdrift", "microruntimes", "runtimes", "book-a-demo-runtime")
		assertHandoffFile(t, runtimeDir, "plan.json", func(raw string) error {
			cards, decodeErr := handoff.Decode[handoff.CardData](raw)
			if decodeErr != nil {
				return decodeErr
			}

			fourCardsPresent := cards.Ticket != nil && cards.Analysis != nil && cards.Plan != nil && cards.Workflow != nil
			if !fourCardsPresent {
				return os.ErrInvalid
			}

			analysisPresent := cards.Analysis.Documents != nil && cards.Analysis.Episodes != nil && cards.Analysis.Suggestion != nil
			planPresent := cards.Plan.Components != nil && cards.Plan.Sizing != nil && cards.Plan.Price != nil && cards.Plan.Authorization != nil && cards.Plan.Evidence != nil
			if !analysisPresent || !planPresent {
				return os.ErrInvalid
			}

			evidence := cards.Plan.Evidence
			evidenceFieldsPresent := evidence.File != nil && evidence.State != nil && evidence.Revision != nil && evidence.ProposedTo != nil
			evidenceMatches := evidenceFieldsPresent && *evidence.File == "runtime.md" && *evidence.State == "draft" && *evidence.Revision == 1 && *evidence.ProposedTo == "console"
			if !evidenceMatches {
				return os.ErrInvalid
			}

			return nil
		})

		runRecorder := &sessionRecorder{}
		deps.Observer = runRecorder

		status, stdout, stderr = seeded.dftInputOnATerminal(deps, "y\n", "agent", "run", "book-a-demo-runtime", "--confirm")
		if status != 0 {
			t.Fatalf("dft agent run = %d: %s\n%s", status, stderr, stdout)
		}

		assertSession(t, stdout, "dft agent run", runRecorder.events, runGolden, runPending)

		assertHandoffFile(t, runtimeDir, "intent.json", func(raw string) error {
			intent, decodeErr := handoff.Decode[handoff.Intent](raw)
			if decodeErr != nil {
				return decodeErr
			}

			expectedIntent := handoff.Intent{
				Ticket:    "0031-TICKET-book-a-demo-runtime",
				RuntimeID: "book-a-demo-runtime",
				Revision:  1,
				Model:     "claude-sonnet-5-5",
				Effort:    "medium",
				Human: handoff.Human{
					Planning: true,
					Review:   true,
				},
				Status: handoff.IntentAuthorized,
			}
			if intent != expectedIntent {
				return os.ErrInvalid
			}

			return nil
		})

		assertHandoffFile(t, runtimeDir, "run.json", func(raw string) error {
			record, decodeErr := handoff.Decode[handoff.RunRecord](raw)
			if decodeErr != nil {
				return decodeErr
			}

			runAndOutcomeMatch := record.Run != nil && *record.Run == "book-a-demo-runtime-0001" && record.Outcome != nil && *record.Outcome == handoff.RunCompleted
			sectionsPresent := record.Intent != nil && len(record.Stages) > 0 && record.Effort != nil && record.SafeSeat != nil && record.Runtime != nil && record.Delivery != nil && record.Delivery.PullRequest != nil
			if !runAndOutcomeMatch || !sectionsPresent {
				return os.ErrInvalid
			}

			return nil
		})
	})

	t.Run("run over budget", func(t *testing.T) {
		t.Parallel()

		overBudgetGolden := loadGolden(t, "run-over-budget", 6)

		seeded, deps := seedGoldenProject(t, []string{"planning", "review"})
		if status, stdout, stderr := seeded.dftWith(deps, goldenPlanArguments...); status != 0 {
			t.Fatalf("dft agent plan = %d: %s\n%s", status, stderr, stdout)
		}

		recorder := &sessionRecorder{}
		deps.Observer = recorder

		status, stdout, stderr := seeded.dftOnATerminal(deps, "agent", "run", "book-a-demo-runtime", "--model", "claude-opus-5-5", "--effort", "high")
		if status != 0 || stderr != "" {
			t.Fatalf("dft agent run over budget = %d: %s\n%s", status, stderr, stdout)
		}

		assertSession(t, stdout, "dft agent run", recorder.events, overBudgetGolden, runOverBudgetPending)
	})

	t.Run("run to merge", func(t *testing.T) {
		t.Parallel()

		toMergeGolden := loadGolden(t, "run-to-merge", 25)

		seeded, deps := seedGoldenProject(t, []string{"planning"})
		if status, stdout, stderr := seeded.dftWith(deps, goldenPlanArguments...); status != 0 {
			t.Fatalf("dft agent plan = %d: %s\n%s", status, stderr, stdout)
		}

		recorder := &sessionRecorder{}
		deps.Observer = recorder

		status, stdout, stderr := seeded.dftInputOnATerminal(deps, "y\n", "agent", "run", "book-a-demo-runtime", "--confirm")
		if status != 0 {
			t.Fatalf("dft agent run to merge = %d: %s\n%s", status, stderr, stdout)
		}

		assertSession(t, stdout, "dft agent run", recorder.events, toMergeGolden, runPending)
	})

	t.Run("run all stages", func(t *testing.T) {
		t.Parallel()

		allStagesGolden := loadGolden(t, "run-all-stages", 31)

		seeded, deps := seedGoldenProject(t, []string{"planning", "spec", "tickets", "implement", "review"})
		if status, stdout, stderr := seeded.dftWith(deps, goldenPlanArguments...); status != 0 {
			t.Fatalf("dft agent plan = %d: %s\n%s", status, stderr, stdout)
		}

		recorder := &sessionRecorder{}
		deps.Observer = recorder

		status, stdout, stderr := seeded.dftInputOnATerminal(deps, "y\ny\ny\ny\n", "agent", "run", "book-a-demo-runtime", "--confirm")
		if status != 0 {
			t.Fatalf("dft agent run all stages = %d: %s\n%s", status, stderr, stdout)
		}

		assertSession(t, stdout, "dft agent run", recorder.events, allStagesGolden, runAllStagesPending)
	})

	t.Run("run haiku low", func(t *testing.T) {
		t.Parallel()

		haikuLowGolden := loadGolden(t, "run-haiku-low", 25)

		seeded, deps := seedGoldenProject(t, []string{"planning", "review"})
		driver := websiteDriver()
		driver.usage = map[string][2]int64{
			"grilling":    {200_000, 10_000},
			"to-spec":     {300_000, 20_000},
			"to-tickets":  {200_000, 10_000},
			"implement":   {900_000, 50_000},
			"tdd":         {200_000, 20_000},
			"code-review": {300_000, 10_000},
		}
		deps.Driver = driver

		if status, stdout, stderr := seeded.dftWith(deps, goldenPlanArguments...); status != 0 {
			t.Fatalf("dft agent plan = %d: %s\n%s", status, stderr, stdout)
		}

		recorder := &sessionRecorder{}
		deps.Observer = recorder

		status, stdout, stderr := seeded.dftInputOnATerminal(deps, "y\n", "agent", "run", "book-a-demo-runtime", "--model", "claude-haiku-5-5", "--effort", "low", "--confirm")
		if status != 0 {
			t.Fatalf("dft agent run haiku low = %d: %s\n%s", status, stderr, stdout)
		}

		assertSession(t, stdout, "dft agent run", recorder.events, haikuLowGolden, runPending)
	})
}

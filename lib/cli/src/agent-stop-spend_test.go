package main

import (
	"strings"
	"testing"
	"time"

	"github.com/jonatassales/tao/lib/cli/src/application/execution"
	"github.com/jonatassales/tao/lib/cli/src/domain/harness"
)

// A run that paid for stages before it waited keeps that spend when it is
// abandoned: with implement waiting for a person after spec and tickets were
// paid, the ledger line's observed figure stays, only the rest of the
// reservation is released, and the totals still add up.
func TestStopKeepsTheSpendObservedBeforeTheWait(t *testing.T) {
	t.Parallel()

	seeded, deps := seedProject(t)
	workflow := bodyFile(t, seeded.root, "workflow-implement.json", "harness/workflow.json", func(body *harness.Workflow) {
		withoutMirror(body)

		for index := range body.Stages {
			if body.Stages[index].ID == "implement" {
				body.Stages[index].Human = harness.HumanRequired
			}
		}
	})

	deps.Driver, deps.Lifecycle, deps.Delivery, deps.Poll = websiteDriver(), websiteLifecycle(), websiteDelivery(), time.Millisecond

	for _, arguments := range [][]string{
		{"harness", "workflow", "set", "--file", workflow},
		{"agent", "plan", consoleTicket, "--component", consoleComponent},
		{"agent", "run", "console-microruntime"},
		{"agent", "confirm", "console-microruntime", "planning"},
		{"agent", "run", "console-microruntime"},
	} {
		if status, stdout, stderr := seeded.dftWith(deps, arguments...); status != 0 {
			t.Fatalf("dft %v = %d: %s\n%s", arguments, status, stderr, stdout)
		}
	}

	record := runRecord(t, seeded)
	waitsAtImplement := record.State == "waiting" && stageState(record, "implement").State == "waiting"
	if !waitsAtImplement {
		t.Fatalf("the run does not wait at implement: %+v", record)
	}

	if status, stdout, stderr := seeded.dftWith(deps, "agent", "stop", "console-microruntime"); status != 0 {
		t.Fatalf("dft agent stop = %d: %s\n%s", status, stderr, stdout)
	}

	ledger := ledgerOf(t, seeded)
	line := ledger.Laps[0]

	keptTheSpend := line.Cents.Observed > 0 && line.Cents.Recovered == line.Cents.Reserved-line.Cents.Observed
	if !keptTheSpend || line.Closed == nil || line.Closed.Outcome != "abandoned" {
		t.Errorf("the abandoned lap's line = %+v", line)
	}

	addsUp := ledger.Totals.Cents.Recovered == ledger.Planning.Cents.Recovered+line.Cents.Recovered
	if !addsUp || ledger.Totals.Cents.Observed != ledger.Planning.Cents.Observed+line.Cents.Observed {
		t.Errorf("the totals do not add up after the closing line: %+v", ledger.Totals)
	}
}

// A cleanup nobody confirmed leaves the abandoned lap's outcome unknown: the
// reservation is held, and the lap still says who stopped it and when.
func TestStopHoldsTheReservationWhenTheCleanupIsUnconfirmed(t *testing.T) {
	t.Parallel()

	seeded, deps := plannedProject(t)
	lifecycle := websiteLifecycle()
	lifecycle.cleanup = execution.Cleanup{
		Removed: 2,
		Total:   4,
	}
	deps.Driver, deps.Lifecycle, deps.Delivery, deps.Poll = websiteDriver(), lifecycle, websiteDelivery(), time.Millisecond

	if status, stdout, stderr := seeded.dftWith(deps, "agent", "run", "console-microruntime"); status != 0 {
		t.Fatalf("dft agent run = %d: %s\n%s", status, stderr, stdout)
	}

	status, stdout, stderr := seeded.dftWith(deps, "agent", "stop", "console-microruntime", "--by", "owner")
	if status != 0 {
		t.Fatalf("dft agent stop = %d: %s\n%s", status, stderr, stdout)
	}

	expectLines(t, stdout, "✓ stopped    lap 0001 · by owner · outcome unknown, the reservation is held")

	record := runRecord(t, seeded)
	unknown := record.Outcome != nil && *record.Outcome == "unknown" && record.Laps[0].Outcome == "unknown"
	if !unknown || ledgerOf(t, seeded).Laps[0].Cents.Recovered != 0 {
		t.Errorf("after an unconfirmed cleanup: outcome %v, laps %+v, line %+v", record.Outcome, record.Laps, ledgerOf(t, seeded).Laps[0])
	}

	lapFile := string(mustRead(t, runtimeFolder(seeded)+"/laps/0001.lap.md"))
	if !strings.Contains(lapFile, "requested by owner with dft agent stop") || !strings.Contains(lapFile, "stopped by owner; the cleanup was not confirmed") {
		t.Errorf("the lap file does not say who stopped it:\n%s", lapFile)
	}
}

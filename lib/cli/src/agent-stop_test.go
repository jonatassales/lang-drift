package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jonatassales/tao/lib/cli/src/domain/connection"
)

// dft agent stop lets a run that waits at a stage go: the runtime is taken
// down, the lap closes as abandoned with who and when, the ledger line
// releases the reservation with nothing observed, the records written so far
// stay, and the runtime can be planned again.
func TestStopAbandonsAWaitingRun(t *testing.T) {
	t.Parallel()

	seeded, deps := plannedProject(t)
	lifecycle := websiteLifecycle()
	deps.Driver, deps.Lifecycle, deps.Delivery, deps.Poll = websiteDriver(), lifecycle, websiteDelivery(), time.Millisecond

	status, stdout, stderr := seeded.dftWith(deps, "agent", "run", "console-microruntime")
	waits := status == 0 && runRecord(t, seeded).State == "waiting"
	if !waits {
		t.Fatalf("dft agent run = %d, want a run that waits: %s\n%s", status, stderr, stdout)
	}

	before := filesUnder(t, runtimeFolder(seeded))

	status, stdout, stderr = seeded.dftWith(deps, "agent", "stop", "console-microruntime", "--by", "owner")
	if status != 0 {
		t.Fatalf("dft agent stop = %d: %s\n%s", status, stderr, stdout)
	}

	expectLines(t, stdout,
		"dft agent stop console-microruntime",
		"✓ abandoned  lap 0001 · by owner · $11.40 released to the envelope",
		"# next: dft agent plan plans console-microruntime again",
	)

	record := runRecord(t, seeded)
	runAbandoned := record.State == "finished" && record.Outcome != nil && *record.Outcome == "abandoned"
	lapAbandoned := record.Laps[0].Outcome == "abandoned"
	takenDownOnce := lifecycle.downs == 1
	if !runAbandoned || !lapAbandoned || !takenDownOnce {
		t.Errorf("after the stop: state %s, outcome %v, laps %+v, %d downs", record.State, record.Outcome, record.Laps, lifecycle.downs)
	}

	ledger := ledgerOf(t, seeded)
	line := ledger.Laps[0]
	released := line.Cents.Reserved == 1140 && line.Cents.Recovered == 1140 && line.Cents.Observed == 0
	if !released || line.Closed == nil {
		t.Fatalf("the abandoned lap's line = %+v", line)
	}

	closedByTheOwner := line.Closed.Outcome == "abandoned" && line.Closed.By == "owner" && line.Closed.At == "2026-10-28T13:46:50Z"
	if !closedByTheOwner {
		t.Errorf("the closing line = %+v", *line.Closed)
	}

	if ledger.Totals.Cents.Recovered != ledger.Planning.Cents.Recovered+line.Cents.Recovered {
		t.Errorf("the totals do not add up after the closing line: %+v", ledger.Totals)
	}

	lapFile := string(mustRead(t, filepath.Join(runtimeFolder(seeded), "laps", "0001.lap.md")))
	if !strings.Contains(lapFile, "| Abandoned | ") || !strings.Contains(lapFile, "| Released | ") {
		t.Errorf("the lap file does not show the abandonment:\n%s", lapFile)
	}

	for _, kept := range before {
		if _, err := os.Stat(kept); err != nil {
			t.Errorf("the stop removed %s", kept)
		}
	}

	status, stdout, stderr = seeded.dftWith(deps, "agent", "status", "console-microruntime")
	showsTheAbandonedLap := status == 0 && strings.Contains(stdout, `"outcome": "abandoned"`) && strings.Contains(stdout, `"openLap": null`)
	if !showsTheAbandonedLap {
		t.Errorf("dft agent status after the stop = %d: %s\n%s", status, stderr, stdout)
	}

	status, stdout, stderr = seeded.dftWith(deps, "agent", "plan", consoleTicket, "--component", consoleComponent)
	if status != 0 {
		t.Fatalf("a plan after the stop = %d: %s\n%s", status, stderr, stdout)
	}

	config, err := connection.Parse(string(mustRead(t, filepath.Join(runtimeFolder(seeded), "runtime.md"))))
	if err != nil || config.Revision != 2 {
		t.Errorf("the plan after the stop = revision %d, %v; want revision 2", config.Revision, err)
	}
}

// Only a run that waits can be stopped: a runtime that is planned, finished
// or failed is refused naming its state, and an unknown runtime is not found.
func TestStopRefusesARunThatDoesNotWait(t *testing.T) {
	t.Parallel()

	t.Run("planned", func(t *testing.T) {
		t.Parallel()

		seeded, deps := plannedProject(t)

		status, _, stderr := seeded.dftWith(deps, "agent", "stop", "console-microruntime")
		if status != 1 || !strings.Contains(stderr, "is planned") {
			t.Errorf("stopping a planned runtime gave %d: %s", status, stderr)
		}
	})

	t.Run("completed", func(t *testing.T) {
		t.Parallel()

		seeded, deps := plannedProject(t)
		deps.Driver, deps.Lifecycle, deps.Delivery, deps.Poll = websiteDriver(), websiteLifecycle(), websiteDelivery(), time.Millisecond

		for _, arguments := range [][]string{{"agent", "run", "console-microruntime"}, {"agent", "confirm", "console-microruntime", "planning"}, {"agent", "run", "console-microruntime"}} {
			if status, stdout, stderr := seeded.dftWith(deps, arguments...); status != 0 {
				t.Fatalf("dft %v = %d: %s\n%s", arguments, status, stderr, stdout)
			}
		}

		status, _, stderr := seeded.dftWith(deps, "agent", "stop", "console-microruntime")
		if status != 1 || !strings.Contains(stderr, "is finished") {
			t.Errorf("stopping a completed run gave %d: %s", status, stderr)
		}
	})

	t.Run("abandoned", func(t *testing.T) {
		t.Parallel()

		seeded, deps := plannedProject(t)
		deps.Driver, deps.Lifecycle, deps.Delivery, deps.Poll = websiteDriver(), websiteLifecycle(), websiteDelivery(), time.Millisecond

		for _, arguments := range [][]string{{"agent", "run", "console-microruntime"}, {"agent", "stop", "console-microruntime"}} {
			if status, stdout, stderr := seeded.dftWith(deps, arguments...); status != 0 {
				t.Fatalf("dft %v = %d: %s\n%s", arguments, status, stderr, stdout)
			}
		}

		status, _, stderr := seeded.dftWith(deps, "agent", "stop", "console-microruntime")
		if status != 1 || !strings.Contains(stderr, "is finished") {
			t.Errorf("stopping an abandoned run gave %d: %s", status, stderr)
		}
	})

	t.Run("failed", func(t *testing.T) {
		t.Parallel()

		seeded, deps := plannedProject(t)
		driver := websiteDriver()
		driver.cutOff["to-spec"] = true
		deps.Driver, deps.Lifecycle, deps.Delivery, deps.Poll = driver, websiteLifecycle(), websiteDelivery(), time.Millisecond

		if status, stdout, stderr := seeded.dftInput(deps, "y\n", "agent", "run", "console-microruntime", "--confirm"); status != 1 {
			t.Fatalf("an interrupted run gave %d: %s\n%s", status, stderr, stdout)
		}

		status, _, stderr := seeded.dftWith(deps, "agent", "stop", "console-microruntime")
		if status != 1 || !strings.Contains(stderr, "is failed") {
			t.Errorf("stopping a failed run gave %d: %s", status, stderr)
		}
	})

	t.Run("unknown", func(t *testing.T) {
		t.Parallel()

		seeded, deps := seedProject(t)

		status, _, stderr := seeded.dftWith(deps, "agent", "stop", "no-such-runtime")
		if status != 1 || !strings.Contains(stderr, "No such runtime") {
			t.Errorf("stopping an unknown runtime gave %d: %s", status, stderr)
		}
	})
}

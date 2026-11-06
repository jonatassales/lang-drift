package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/jonatassales/tao/lib/cli/src/domain/finops"
)

// dft init asks, at a terminal or with --ask, for the month's budget in
// dollars and the default gates, and records them through the same use cases
// as dft finops budget set and dft finops gates set. Flags answer without a
// terminal, an empty budget skips it, and a second init asks only for what is
// not recorded yet.
func TestInitAsksForTheBudgetAndGates(test *testing.T) {
	test.Parallel()

	october := dependencies{
		Clock: func() time.Time { return time.Date(2026, 10, 30, 9, 0, 0, 0, time.UTC) },
	}

	budgetOf := func(test *testing.T, root string) (finops.Budget, bool) {
		test.Helper()

		data, err := os.ReadFile(filepath.Join(root, ".langdrift", "finops", "budgets.json"))
		if os.IsNotExist(err) {
			return finops.Budget{}, false
		}

		if err != nil {
			test.Fatal(err)
		}

		var budget finops.Budget
		if err := json.Unmarshal(data, &budget); err != nil {
			test.Fatal(err)
		}

		return budget, true
	}

	test.Run("records the answers", func(test *testing.T) {
		test.Parallel()

		consumer := newProject(test)

		status, stdout, stderr := consumer.dftInput(october, "1500\n25.50\n20\ny\n", "init", "--ask", "--by", "lead")
		if status != 0 {
			test.Fatalf("dft init --ask = %d, %s", status, stderr)
		}

		for _, question := range []string{"Monthly budget for 2026-10", "Per-run cap", "Reserve margin", "Apply the default gates"} {
			if !strings.Contains(stderr, question) {
				test.Errorf("dft init did not ask %q:\n%s", question, stderr)
			}
		}

		budget, found := budgetOf(test, consumer.root)
		if !found {
			test.Fatal("dft init recorded no budget")
		}

		want := finops.BudgetBounds{
			PeriodCents:          150000,
			PerRunCents:          2550,
			ReserveMarginPercent: 20,
		}

		if budget.Period != "2026-10" || budget.Bounded != want || budget.Mirror.SetBy != "lead" {
			test.Errorf("budget = %s %+v set by %s, want 2026-10 %+v set by lead", budget.Period, budget.Bounded, budget.Mirror.SetBy, want)
		}

		var gates finops.Gates
		decodeStrict(test, mustRead(test, filepath.Join(consumer.root, ".langdrift", "finops", "gates.json")), &gates)

		if !reflect.DeepEqual(gates.Gates, finops.DefaultGates()) || gates.Mirror.SetBy != "lead" {
			test.Errorf("gates = %+v set by %s, want the default gates set by lead", gates.Gates, gates.Mirror.SetBy)
		}

		var result struct {
			Status   string   `json:"status"`
			Recorded []string `json:"recorded"`
			Next     []string `json:"next"`
		}
		if err := json.Unmarshal([]byte(stdout), &result); err != nil {
			test.Fatalf("result %q: %v", stdout, err)
		}

		wantNext := []string{"dft harness envelope set", "dft harness model set", "dft harness driver set"}

		if result.Status != "initialized" || !reflect.DeepEqual(result.Recorded, []string{"budget", "gates"}) || !reflect.DeepEqual(result.Next, wantNext) {
			test.Errorf("result = %+v, want initialized, budget and gates recorded, %v next", result, wantNext)
		}

		status, stdout, stderr = consumer.dftInput(october, "", "init", "--ask")
		if status != 0 || strings.Contains(stderr, "Monthly budget") || strings.Contains(stderr, "default gates") {
			test.Fatalf("a second dft init --ask = %d, asked again:\n%s", status, stderr)
		}

		if strings.Contains(stdout, "recorded") {
			test.Errorf("a second dft init recorded again: %s", stdout)
		}

		if again, _ := budgetOf(test, consumer.root); again.Mirror.Record != budget.Mirror.Record {
			test.Error("a second dft init replaced the budget")
		}
	})

	test.Run("skips an empty budget", func(test *testing.T) {
		test.Parallel()

		consumer := newProject(test)

		status, _, stderr := consumer.dftInput(october, "\nn\n", "init", "--ask")
		if status != 0 {
			test.Fatalf("dft init --ask = %d, %s", status, stderr)
		}

		if strings.Contains(stderr, "Per-run cap") {
			test.Errorf("dft init asked for the per-run cap of a skipped budget:\n%s", stderr)
		}

		if _, found := budgetOf(test, consumer.root); found {
			test.Error("dft init recorded a skipped budget")
		}
	})

	test.Run("refuses a malformed amount", func(test *testing.T) {
		test.Parallel()

		for _, amount := range []string{"abc", "-5", "12.345", "0"} {
			consumer := newProject(test)

			status, _, stderr := consumer.dftInput(october, amount+"\n", "init", "--ask")
			if status != 1 || !strings.Contains(stderr, "dollars") {
				test.Errorf("dft init with a budget of %q = %d, %s", amount, status, stderr)
			}

			if _, found := budgetOf(test, consumer.root); found {
				test.Errorf("dft init recorded a budget of %q", amount)
			}
		}
	})

	test.Run("records the flags without asking", func(test *testing.T) {
		test.Parallel()

		consumer := newProject(test)

		status, _, stderr := consumer.dftInput(october, "", "init", "--budget", "1500", "--per-run", "25", "--reserve-margin", "10", "--default-gates")
		if status != 0 || stderr != "" {
			test.Fatalf("dft init with flags = %d, %q", status, stderr)
		}

		budget, found := budgetOf(test, consumer.root)
		if !found || budget.Bounded.PeriodCents != 150000 || budget.Bounded.PerRunCents != 2500 || budget.Bounded.ReserveMarginPercent != 10 {
			test.Errorf("budget = %+v, %v", budget.Bounded, found)
		}

		if _, err := os.Stat(filepath.Join(consumer.root, ".langdrift", "finops", "gates.json")); err != nil {
			test.Errorf("dft init --default-gates recorded no gates: %v", err)
		}
	})

	test.Run("takes the budget flags together", func(test *testing.T) {
		test.Parallel()

		consumer := newProject(test)

		status, _, stderr := consumer.dftInput(october, "", "init", "--budget", "1500")
		if status != 1 || !strings.Contains(stderr, "--per-run") {
			test.Errorf("dft init --budget alone = %d, %s", status, stderr)
		}
	})

	test.Run("asks nothing with --json", func(test *testing.T) {
		test.Parallel()

		consumer := newProject(test)

		status, stdout, stderr := consumer.dftInput(october, "1500\n25\n20\ny\n", "init", "--json", "--ask")
		if status != 0 || stderr != "" || strings.Contains(stdout, "recorded") {
			test.Errorf("dft init --json --ask = %d, stdout %q, stderr %q", status, stdout, stderr)
		}
	})
}

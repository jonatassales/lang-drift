package session_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	app "github.com/jonatassales/tao/lib/cli/src/application/session"
	"github.com/jonatassales/tao/lib/cli/src/domain/finops"
	"github.com/jonatassales/tao/lib/cli/src/presentation/cli/session"
)

func TestHumanPrintsOneLinePerEvent(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	sink := session.For(&out, false, false)

	sink.Emit(app.Event{Session: app.PlanSession, ID: app.PlanScope})
	sink.Emit(app.Event{Session: app.PlanSession, ID: app.ReadTicket, Data: app.Data{"ticket": "0031-TICKET-book-a-demo-runtime"}})

	want := "# plan only · nothing is built, started or paid for in this step\nconsole › tool call: read 0031-TICKET-book-a-demo-runtime\n"
	if out.String() != want {
		t.Fatalf("printed %q, want %q", out.String(), want)
	}
}

func TestNDJSONPrintsTheWebsiteIdsAndReveals(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	sink := session.NDJSON{Writer: &out, Now: func() time.Time { return time.Date(2026, 10, 28, 13, 48, 5, 0, time.UTC) }}

	sink.Emit(app.Event{Session: app.PlanSession, ID: app.Suggest, Data: app.Data{"model": "claude-sonnet-5-5", "effort": "medium"}})
	sink.Emit(app.Event{Session: app.RunSession, ID: app.GatePaused, Data: app.Data{"reservedCents": 4080, "perRunCents": 2500}})

	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("printed %d lines, want 2", len(lines))
	}

	var first session.Record
	if err := json.Unmarshal([]byte(lines[0]), &first); err != nil {
		t.Fatal(err)
	}

	if first.At != "2026-10-28T13:48:05Z" || first.Session != "plan" || first.ID != "suggest" || first.Kind != app.Success || first.Card != "analysis" || first.Field != "suggestion" || first.Text != "✓ suggests   claude-sonnet-5-5 · effort medium" {
		t.Fatalf("first event = %+v", first)
	}

	var second session.Record
	if err := json.Unmarshal([]byte(lines[1]), &second); err != nil {
		t.Fatal(err)
	}

	if second.Kind != app.Warning || second.Card != "workflow" || second.Field != "paused" || second.Data["reservedCents"] != 4080.0 {
		t.Fatalf("second event = %+v", second)
	}
}

func TestFormatsMoneyAndTokensAsTheWebsiteDoes(t *testing.T) {
	t.Parallel()

	money := map[int]string{950: "$9.50", 1140: "$11.40", 4080: "$40.80", 160: "$1.60", 2500: "$25.00", 123456: "$1,234.56", 5: "$0.05"}
	for cents, want := range money {
		if got := finops.Dollars(cents); got != want {
			t.Errorf("Dollars(%d) = %q, want %q", cents, got, want)
		}
	}

	tokens := map[int]string{3500: "3.5M", 250: "250K", 6000: "6M", 500: "500K", 2000: "2M", 120: "120K", 2100: "2.1M", 6300: "6.3M"}
	for thousands, want := range tokens {
		if got := finops.Tokens(thousands); got != want {
			t.Errorf("Tokens(%d) = %q, want %q", thousands, got, want)
		}
	}

	if got := finops.Grouped(2000); got != "2,000" {
		t.Errorf("Number(2000) = %q", got)
	}
}

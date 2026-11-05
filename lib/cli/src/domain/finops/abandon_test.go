package finops_test

import (
	"testing"
	"time"

	"github.com/jonatassales/tao/lib/cli/src/domain/finops"
)

// abandonedAt is when the owner let the lap go, in another zone than UTC, so
// the closing line's time shows it is written in UTC.
var abandonedAt = time.Date(2026, 10, 28, 10, 46, 50, 0, time.FixedZone("BRT", -3*60*60))

// Each money rule of an abandoned lap: nothing observed releases the whole
// reservation; spend observed before the wait is kept and the rest released;
// spend above the reservation releases nothing and stays visible as
// overspend; a call cut off holds the reservation, since what it spent is
// unknown, and an overspend from it would be counted from a floor.
func TestAbandonSettlesTheReservation(t *testing.T) {
	t.Parallel()

	for _, scenario := range []struct {
		name string
		lap  finops.Lap
		want finops.Abandonment
	}{
		{
			name: "nothing observed",
			lap: finops.Lap{
				Reserved: 1140,
			},
			want: finops.Abandonment{
				Released: 1140,
			},
		},
		{
			name: "spend observed before the wait",
			lap: finops.Lap{
				Reserved: 1140,
				Observed: 300,
			},
			want: finops.Abandonment{
				Released: 840,
				Observed: 300,
			},
		},
		{
			name: "spend above the reservation",
			lap: finops.Lap{
				Reserved: 1140,
				Observed: 1400,
			},
			want: finops.Abandonment{
				Observed:  1400,
				Overspend: 260,
			},
		},
		{
			name: "a call cut off",
			lap: finops.Lap{
				Reserved:   1140,
				Observed:   300,
				Unmeasured: true,
			},
			want: finops.Abandonment{
				Observed: 300,
			},
		},
		{
			name: "a call cut off above the reservation",
			lap: finops.Lap{
				Reserved:   1140,
				Observed:   1400,
				Unmeasured: true,
			},
			want: finops.Abandonment{
				Observed: 1400,
			},
		},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			t.Parallel()

			got := finops.Abandon(scenario.lap, abandonedAt, "owner")
			got.Closing = finops.LapClosing{}

			if got != scenario.want {
				t.Errorf("Abandon = %+v, want %+v", got, scenario.want)
			}

			if got.Released > scenario.lap.Reserved {
				t.Errorf("released %d of a %d reservation", got.Released, scenario.lap.Reserved)
			}
		})
	}
}

// The closing line names the outcome, the person and the time in RFC 3339,
// in UTC.
func TestAbandonWritesTheClosingLine(t *testing.T) {
	t.Parallel()

	lap := finops.Lap{
		Reserved: 1140,
	}

	got := finops.Abandon(lap, abandonedAt, "owner").Closing

	want := finops.LapClosing{
		Outcome: finops.LapAbandoned,
		By:      "owner",
		At:      "2026-10-28T13:46:50Z",
	}
	if got != want {
		t.Errorf("closing = %+v, want %+v", got, want)
	}
}

// After the abandoned lap's line, the ledger's totals add up: what was
// reserved, observed and released, and what the period has left.
func TestAbandonKeepsTheLedgerTotalsTrue(t *testing.T) {
	t.Parallel()

	lap := finops.Lap{
		Reserved: 1140,
		Observed: 300,
	}
	abandonment := finops.Abandon(lap, abandonedAt, "owner")

	ledger := finops.Ledger{
		Planning: finops.PlanningLine{
			Cents: finops.PlanningCents{
				Reserved:  15,
				Observed:  10,
				Recovered: 5,
			},
		},
		Laps: []finops.LapLine{
			{
				Lap: "0002",
				Cents: finops.LapCents{
					Reserved:  lap.Reserved,
					Observed:  abandonment.Observed,
					Recovered: abandonment.Released,
					Overspend: abandonment.Overspend,
				},
				Closed: &abandonment.Closing,
			},
		},
	}

	finops.SettleTotals(&ledger, 50000)

	want := finops.TotalCents{
		Reserved:  1155,
		Observed:  310,
		Recovered: 845,
	}
	if ledger.Totals.Cents != want || ledger.Totals.PeriodLeftAfterCents != 49690 {
		t.Errorf("totals = %+v, want %+v and 49690 left", ledger.Totals, want)
	}

	unspent := ledger.Totals.Cents.Reserved - ledger.Totals.Cents.Observed
	if ledger.Totals.Cents.Recovered != unspent {
		t.Errorf("recovered %d, but %d of the reservations went unspent", ledger.Totals.Cents.Recovered, unspent)
	}
}

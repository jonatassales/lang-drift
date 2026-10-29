package finops_test

import (
	"testing"

	"github.com/jonatassales/tao/packages/ex/src/domain/finops"
)

var (
	opus   = finops.ModelPrice{ID: "claude-opus-5-5", InputCents: 400, OutputCents: 2000}
	sonnet = finops.ModelPrice{ID: "claude-sonnet-5-5", InputCents: 200, OutputCents: 1000}
	haiku  = finops.ModelPrice{ID: "claude-haiku-5-5", InputCents: 10, OutputCents: 50}

	expected = map[string]finops.TokensThousands{"low": {Input: 2000, Output: 120}, "medium": {Input: 3500, Output: 250}, "high": {Input: 6000, Output: 500}}
	observed = map[string]finops.TokensThousands{"low": {Input: 2100, Output: 125}, "medium": {Input: 3700, Output: 240}, "high": {Input: 6300, Output: 520}}
)

// The website's three-by-three matrix of models and efforts, in cents.
func TestEstimatesMatchTheWebsite(t *testing.T) {
	t.Parallel()

	tests := []struct {
		price                                  finops.ModelPrice
		effort                                 string
		estimate, reserved, observed, released int
	}{
		{opus, "low", 1040, 1248, 1090, 158},
		{opus, "medium", 1900, 2280, 1960, 320},
		{opus, "high", 3400, 4080, 3560, 520},
		{sonnet, "low", 520, 624, 545, 79},
		{sonnet, "medium", 950, 1140, 980, 160},
		{sonnet, "high", 1700, 2040, 1780, 260},
		{haiku, "low", 26, 32, 27, 5},
		{haiku, "medium", 48, 58, 49, 9},
		{haiku, "high", 85, 102, 89, 13},
	}

	for _, tt := range tests {
		estimate := finops.EstimateCents(expected[tt.effort], tt.price)
		reserved := finops.RequestedReservation(estimate, 20)
		spent := finops.EstimateCents(observed[tt.effort], tt.price)
		settlement := finops.Settle(estimate, reserved, spent)

		if estimate != tt.estimate || reserved != tt.reserved || spent != tt.observed || settlement.Recovered != tt.released {
			t.Errorf("%s at %s: estimate %d, reserved %d, observed %d, released %d; want %d, %d, %d, %d", tt.price.ID, tt.effort, estimate, reserved, spent, settlement.Recovered, tt.estimate, tt.reserved, tt.observed, tt.released)
		}
	}

	if finops.ObservedCents(3_700_000, 240_000, sonnet) != 980 {
		t.Errorf("raw tokens priced to %d, want 980", finops.ObservedCents(3_700_000, 240_000, sonnet))
	}
}

// The sample's two laps: the first paused at FG-03 with nothing paid, the
// second opened, observed and settled.
func TestGatesReproduceTheSampleLaps(t *testing.T) {
	t.Parallel()

	bounds := finops.Bounds{PerRunCents: 2500, TicketAllocationCents: 2500, PeriodLeftCents: 50000}
	authorization := finops.Authorization{Priced: true, Allocated: true}

	first := finops.BeforeAuthorization(authorization, 4080, bounds)
	if len(first) != 3 || first[2].ID != "FG-03" || first[2].Result != "paused" || *first[2].Detail != "requested reservation 4080 exceeds the per-run cap 2500" || finops.Allowed(first) {
		t.Fatalf("lap 0001 gates = %+v", first)
	}

	second := finops.BeforeAuthorization(authorization, 1140, bounds)
	if !finops.Allowed(second) || *second[2].Detail != "reservation 1140 within the per-run cap 2500" {
		t.Fatalf("lap 0002 gates = %+v", second)
	}

	during := finops.DuringRun(980, true, 1140)
	if during.Result != "not-reached" || *during.Detail != "observed 980 stayed below the reservation 1140" {
		t.Errorf("FG-04 = %+v", during)
	}

	if unknown := finops.DuringRun(0, false, 1140); unknown.Result != "paused" {
		t.Errorf("an unknown spend did not pause: %+v", unknown)
	}

	if reached := finops.DuringRun(1140, true, 1140); reached.Result != "paused" {
		t.Errorf("a spend at the reservation did not pause: %+v", reached)
	}

	settlement := finops.Settle(950, 1140, 980)
	if settlement.Recovered != 160 || settlement.Overspend != 0 || settlement.Variance != 30 {
		t.Errorf("settlement = %+v", settlement)
	}

	if after := finops.AfterRun(true); after == nil || after.Result != "keep-provisional" {
		t.Errorf("FG-06 = %+v", after)
	}
	if finops.AfterRun(false) != nil {
		t.Error("a lap that paid nothing got FG-06")
	}

	unpriced := finops.BeforeAuthorization(finops.Authorization{Priced: false, Allocated: true}, 100, bounds)
	if len(unpriced) != 2 || unpriced[0].Result != "blocked" {
		t.Errorf("an unpriced model was not blocked: %+v", unpriced)
	}

	overspent := finops.BeforeAuthorization(finops.Authorization{Priced: true, Allocated: true, Overspend: true}, 100, bounds)
	if len(overspent) != 3 || overspent[2].ID != "FG-05" || overspent[2].Result != "require-owner" {
		t.Errorf("an unacknowledged overspend did not require the owner: %+v", overspent)
	}
}

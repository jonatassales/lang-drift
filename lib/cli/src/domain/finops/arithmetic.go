package finops

import (
	"fmt"
	"math"
)

// EstimateCents is what tokens cost at a model's list prices, the website's
// way: the total in thousandths of a cent, rounded.
func EstimateCents(tokens TokensThousands, price ModelPrice) int {
	thousandths := tokens.Input*price.InputCents + tokens.Output*price.OutputCents

	return int(math.Round(float64(thousandths) / 1000))
}

// ObservedCents is what raw token counts cost at a model's list prices.
func ObservedCents(inputTokens, outputTokens int64, price ModelPrice) int {
	millionths := inputTokens*int64(price.InputCents) + outputTokens*int64(price.OutputCents)

	return int(math.Round(float64(millionths) / 1_000_000))
}

// RequestedReservation is the estimate plus the margin, rounded up.
func RequestedReservation(estimateCents, reserveMarginPercent int) int {
	return int(math.Ceil(float64(estimateCents) * float64(100+reserveMarginPercent) / 100))
}

// Bounds are the caps a reservation is checked against: the ones the ledger
// records.
type Bounds = LedgerBounds

// The gate results.
const (
	Passed          = "passed"
	Blocked         = "blocked"
	GatePaused      = "paused"
	NotReached      = "not-reached"
	RequireOwner    = "require-owner"
	KeepProvisional = "keep-provisional"
)

// Authorization is what a lap needs before paid work: a price source for its
// model, an allocation for its ticket, and no overspend nobody acknowledged.
type Authorization struct {
	Priced       bool
	Allocated    bool
	Overspend    bool
	Acknowledged bool
}

// BeforeAuthorization runs the gates before a lap's paid work, in the order
// the ledger records them: FG-01, FG-02, FG-05 when an earlier lap overspent,
// then FG-03 when the first ones passed. The lap may start only when every
// result is passed.
func BeforeAuthorization(authorization Authorization, reservation int, bounds Bounds) []GateResult {
	results := []GateResult{
		{ID: "FG-01", Result: resultOf(authorization.Priced, Blocked)},
		{ID: "FG-02", Result: resultOf(authorization.Allocated, Blocked)},
	}

	if authorization.Overspend {
		results = append(results, GateResult{ID: "FG-05", Result: resultOf(authorization.Acknowledged, RequireOwner)})
	}

	for _, result := range results {
		if result.Result != Passed {
			return results
		}
	}

	return append(results, reservationGate(reservation, bounds))
}

func resultOf(passed bool, otherwise string) string {
	if passed {
		return Passed
	}

	return otherwise
}

// reservationGate is FG-03: the reservation stays inside the per-run cap,
// the ticket's allocation and what is left of the period.
func reservationGate(reservation int, bounds Bounds) GateResult {
	switch {
	case reservation > bounds.PerRunCents:
		return GateResult{ID: "FG-03", Result: GatePaused, Detail: text(fmt.Sprintf("requested reservation %d exceeds the per-run cap %d", reservation, bounds.PerRunCents))}
	case reservation > bounds.TicketAllocationCents:
		return GateResult{ID: "FG-03", Result: GatePaused, Detail: text(fmt.Sprintf("requested reservation %d exceeds the ticket allocation %d", reservation, bounds.TicketAllocationCents))}
	case reservation > bounds.PeriodLeftCents:
		return GateResult{ID: "FG-03", Result: GatePaused, Detail: text(fmt.Sprintf("requested reservation %d exceeds what is left of the period, %d", reservation, bounds.PeriodLeftCents))}
	}

	return GateResult{ID: "FG-03", Result: Passed, Detail: text(fmt.Sprintf("reservation %d within the per-run cap %d", reservation, bounds.PerRunCents))}
}

// Allowed reports whether every gate passed.
func Allowed(results []GateResult) bool {
	for _, result := range results {
		if result.Result != Passed {
			return false
		}
	}

	return true
}

// DuringRun is FG-04, before each paid call: paid calls pause when observed
// spend reaches the reservation, or when it is unknown, since an unknown
// must never count as zero.
func DuringRun(observedCents int, observedKnown bool, reservedCents int) GateResult {
	switch {
	case !observedKnown:
		return GateResult{ID: "FG-04", Result: GatePaused, Detail: text("observed spend is unknown; a call reported no usage")}
	case observedCents >= reservedCents:
		return GateResult{ID: "FG-04", Result: GatePaused, Detail: text(fmt.Sprintf("observed %d reached the reservation %d", observedCents, reservedCents))}
	}

	return GateResult{ID: "FG-04", Result: NotReached, Detail: text(fmt.Sprintf("observed %d stayed below the reservation %d", observedCents, reservedCents))}
}

// AfterRun is FG-06: a lap that paid keeps its billing provisional until the
// provider bills it.
func AfterRun(paid bool) *GateResult {
	if !paid {
		return nil
	}

	return &GateResult{ID: "FG-06", Result: KeepProvisional, Detail: text("the provider has not billed the lap; billed stays unknown")}
}

// Settlement is what a lap's money came to once it ended.
type Settlement struct {
	Recovered int
	Overspend int
	Variance  int
}

// Settle releases what the reservation did not spend and names the rest.
func Settle(estimateCents, reservedCents, observedCents int) Settlement {
	return Settlement{
		Recovered: max(0, reservedCents-observedCents),
		Overspend: max(0, observedCents-reservedCents),
		Variance:  observedCents - estimateCents,
	}
}

func text(value string) *string { return &value }

package finops

import "github.com/jonatassales/tao/lib/cli/src/domain/harness"

// Ledger is one runtime's money: a planning line, a line per lap and totals.
// A value nobody measured, such as what the provider billed, is null.
type Ledger struct {
	Mirror         harness.Mirror `json:"mirror"`
	Version        int            `json:"version"`
	Sample         bool           `json:"sample,omitempty"`
	Run            string         `json:"run"`
	Ticket         *string        `json:"ticket"`
	State          string         `json:"state"`
	Currency       string         `json:"currency"`
	PriceSource    *string        `json:"priceSource"`
	ForecastSource *string        `json:"forecastSource"`
	Bounded        LedgerBounds   `json:"bounded"`
	Planning       PlanningLine   `json:"planning"`
	Laps           []LapLine      `json:"laps"`
	Totals         LedgerTotals   `json:"totals"`
}

// LedgerBounds are the caps the run's reservations were checked against.
type LedgerBounds struct {
	PerRunCents           int `json:"perRunCents"`
	TicketAllocationCents int `json:"ticketAllocationCents"`
	PeriodLeftCents       int `json:"periodLeftCents"`
}

// PlanningLine is what the planner's own calls cost.
type PlanningLine struct {
	Model *string       `json:"model"`
	Note  string        `json:"note"`
	Cents PlanningCents `json:"cents"`
}

// PlanningCents are the planning line's amounts.
type PlanningCents struct {
	Estimate  int  `json:"estimate"`
	Reserved  int  `json:"reserved"`
	Observed  int  `json:"observed"`
	Billed    *int `json:"billed"`
	Recovered int  `json:"recovered"`
	Overspend int  `json:"overspend"`
}

// LapLine is one lap's money: the choice, the tokens, the amounts and the
// gates it passed.
type LapLine struct {
	Lap             string       `json:"lap"`
	Model           string       `json:"model"`
	Effort          string       `json:"effort"`
	TokensThousands LapTokens    `json:"tokensThousands"`
	Cents           LapCents     `json:"cents"`
	Gates           []GateResult `json:"gates"`
	NothingPaid     bool         `json:"nothingPaid"`
}

// LapTokens are the tokens a lap expected and the ones it observed.
type LapTokens struct {
	Estimate TokensThousands  `json:"estimate"`
	Observed *TokensThousands `json:"observed"`
}

// LapCents are a lap's amounts, from the estimate to the variance.
type LapCents struct {
	Estimate             int  `json:"estimate"`
	RequestedReservation int  `json:"requestedReservation"`
	Reserved             int  `json:"reserved"`
	Observed             int  `json:"observed"`
	Billed               *int `json:"billed"`
	Recovered            int  `json:"recovered"`
	Overspend            int  `json:"overspend"`
	Variance             *int `json:"variance"`
}

// GateResult is one gate's decision for a lap.
type GateResult struct {
	ID     string  `json:"id"`
	Result string  `json:"result"`
	Detail *string `json:"detail,omitempty"`
}

// LedgerTotals sum the lines and say what the period has left.
type LedgerTotals struct {
	Cents                TotalCents `json:"cents"`
	PeriodLeftAfterCents int        `json:"periodLeftAfterCents"`
}

// TotalCents are the ledger's summed amounts.
type TotalCents struct {
	Reserved  int  `json:"reserved"`
	Observed  int  `json:"observed"`
	Billed    *int `json:"billed"`
	Recovered int  `json:"recovered"`
	Overspend int  `json:"overspend"`
}

// SettleTotals sums the ledger's lines and says what the period has left
// after them. Billed stays unknown while any line's is.
func SettleTotals(ledger *Ledger, periodCents int) {
	totals := TotalCents{Reserved: ledger.Planning.Cents.Reserved, Observed: ledger.Planning.Cents.Observed, Recovered: ledger.Planning.Cents.Recovered, Overspend: ledger.Planning.Cents.Overspend}
	billed, known := 0, ledger.Planning.Cents.Billed != nil
	if known {
		billed = *ledger.Planning.Cents.Billed
	}

	for _, line := range ledger.Laps {
		totals.Reserved += line.Cents.Reserved
		totals.Observed += line.Cents.Observed
		totals.Recovered += line.Cents.Recovered
		totals.Overspend += line.Cents.Overspend
		if line.Cents.Billed == nil {
			known = false

			continue
		}
		billed += *line.Cents.Billed
	}

	if known {
		totals.Billed = &billed
	}

	ledger.Totals = LedgerTotals{Cents: totals, PeriodLeftAfterCents: periodCents - totals.Observed}
}

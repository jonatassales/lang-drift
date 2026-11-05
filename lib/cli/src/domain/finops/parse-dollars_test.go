package finops

import "testing"

// ParseDollars reads what a person types in dollars as whole cents, and
// refuses anything that is not a positive amount with at most two decimals.
func TestParseDollars(test *testing.T) {
	test.Parallel()

	accepted := map[string]int{
		"1500":     150000,
		"1500.5":   150050,
		"1500.50":  150050,
		"25.05":    2505,
		"0.01":     1,
		" 12 ":     1200,
		"00012.30": 1230,
	}

	for text, want := range accepted {
		if got, err := ParseDollars(text); err != nil || got != want {
			test.Errorf("ParseDollars(%q) = %d, %v, want %d", text, got, err, want)
		}
	}

	for _, text := range []string{"", "abc", "-5", "12.345", "0", "0.00", "1,500", "$15", "1e3", "99999999999999999999"} {
		if got, err := ParseDollars(text); err == nil {
			test.Errorf("ParseDollars(%q) = %d, want an error", text, got)
		}
	}
}

package finops

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// dollarsPattern is a whole number of dollars, at most nine digits so the
// cents never overflow, with at most two decimals.
var dollarsPattern = regexp.MustCompile(`^([0-9]{1,9})(?:\.([0-9]{1,2}))?$`)

// ParseDollars reads a positive amount a person types in dollars, such as
// 1500 or 1500.50, as whole cents.
func ParseDollars(text string) (int, error) {
	parts := dollarsPattern.FindStringSubmatch(strings.TrimSpace(text))
	if parts == nil {
		return 0, fmt.Errorf("%q is not an amount in dollars, such as 1500 or 1500.50", text)
	}

	dollars, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("%q is too large an amount in dollars", text)
	}

	cents, err := strconv.Atoi((parts[2] + "00")[:2])
	if err != nil {
		return 0, fmt.Errorf("%q is not an amount in dollars, such as 1500 or 1500.50", text)
	}

	total := dollars*100 + cents
	if total <= 0 {
		return 0, fmt.Errorf("%q is not a positive amount in dollars", text)
	}

	return total, nil
}

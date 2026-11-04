package finops

import (
	"fmt"
	"strconv"
	"strings"
)

// Dollars prints cents as dollars with cents, grouped by thousands, as the
// lines and the laps show money.
func Dollars(cents int) string {
	sign := ""
	if cents < 0 {
		sign, cents = "-", -cents
	}

	return fmt.Sprintf("%s$%s.%02d", sign, Grouped(cents/100), cents%100)
}

// Tokens prints thousands of tokens compactly: 3.5M, 250K, 6M.
func Tokens(thousands int) string {
	if thousands < 1000 {
		return strconv.Itoa(thousands) + "K"
	}

	millions := float64(thousands) / 1000

	return strings.TrimSuffix(strings.TrimSuffix(fmt.Sprintf("%.1f", millions), "0"), ".") + "M"
}

// Grouped prints an integer with thousands separators.
func Grouped(value int) string {
	if value < 0 {
		return "-" + Grouped(-value)
	}

	digits := strconv.Itoa(value)

	var parts []string
	for len(digits) > 3 {
		parts = append([]string{digits[len(digits)-3:]}, parts...)
		digits = digits[:len(digits)-3]
	}

	return strings.Join(append([]string{digits}, parts...), ",")
}

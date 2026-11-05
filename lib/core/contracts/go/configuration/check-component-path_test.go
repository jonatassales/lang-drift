package configuration

import (
	"errors"
	"testing"
)

func TestCheckComponentPath(t *testing.T) {
	for _, test := range []struct {
		text string
		ok   bool
	}{
		{".", true},
		{"apps/console", true},
		{"apps/../console", false},
		{"../console", false},
		{"/apps/console", false},
		{"apps/console/", false},
		{`apps\console`, false},
		{"", false},
	} {
		err := CheckComponentPath(test.text)
		wrongError := err != nil && !errors.Is(err, ErrInvalidPath)
		unexpected := (err == nil) != test.ok || wrongError

		if unexpected {
			t.Errorf("CheckComponentPath(%q) = %v, want ok=%t", test.text, err, test.ok)
		}
	}
}

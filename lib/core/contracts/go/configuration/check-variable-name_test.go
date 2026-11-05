package configuration

import (
	"errors"
	"testing"
)

func TestCheckVariableName(t *testing.T) {
	for _, test := range []struct {
		text string
		ok   bool
	}{
		{"PORT", true},
		{"MINIO_ACCESS_KEY", true},
		{"S3", true},
		{"port", false},
		{"_PORT", false},
		{"3PORT", false},
		{"PORT-NUMBER", false},
		{"", false},
	} {
		err := CheckVariableName(test.text)
		wrongError := err != nil && !errors.Is(err, ErrInvalidVariable)
		unexpected := (err == nil) != test.ok || wrongError

		if unexpected {
			t.Errorf("CheckVariableName(%q) = %v, want ok=%t", test.text, err, test.ok)
		}
	}
}

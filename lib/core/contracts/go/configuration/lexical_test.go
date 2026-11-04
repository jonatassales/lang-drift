package configuration

import (
	"errors"
	"strings"
	"testing"
)

func TestCheckIdentifier(t *testing.T) {
	for _, test := range []struct {
		text string
		ok   bool
	}{
		{"console-local", true},
		{"a", true},
		{strings.Repeat("a", 63), true},
		{strings.Repeat("a", 64), false},
		{"Console", false},
		{"-console", false},
		{"console-", false},
		{"console_local", false},
		{"", false},
	} {
		err := CheckIdentifier(test.text)
		if (err == nil) != test.ok || (err != nil && !errors.Is(err, ErrInvalidIdentifier)) {
			t.Errorf("CheckIdentifier(%q) = %v, want ok=%t", test.text, err, test.ok)
		}
	}
}

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
		if (err == nil) != test.ok || (err != nil && !errors.Is(err, ErrInvalidPath)) {
			t.Errorf("CheckComponentPath(%q) = %v, want ok=%t", test.text, err, test.ok)
		}
	}
}

func TestLexicalPredicates(t *testing.T) {
	if !IsPlatformPath(".langdrift/state") || IsPlatformPath(".langdriftx") {
		t.Error("IsPlatformPath does not bound .langdrift exactly")
	}
	if !ValidPort(1) || !ValidPort(65535) || ValidPort(0) || ValidPort(65536) {
		t.Error("ValidPort does not bound 1 to 65535")
	}
	if !IsImageReference("console:local") || IsImageReference("console local") || IsImageReference(" ") {
		t.Error("IsImageReference accepts whitespace")
	}
}

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
		if (err == nil) != test.ok || (err != nil && !errors.Is(err, ErrInvalidVariable)) {
			t.Errorf("CheckVariableName(%q) = %v, want ok=%t", test.text, err, test.ok)
		}
	}
}

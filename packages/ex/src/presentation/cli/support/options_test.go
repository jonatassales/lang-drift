package support_test

import (
	"strings"
	"testing"

	"github.com/exphere/ex/src/presentation/cli/support"
)

func TestOptions_WithDefaults(t *testing.T) {
	t.Parallel()

	options := support.Options{}.WithDefaults()
	if options.Input == nil || options.Output == nil || options.Diagnostics == nil {
		t.Fatalf("WithDefaults() = %+v, want every stream set", options)
	}

	var output strings.Builder
	if kept := (support.Options{Output: &output}).WithDefaults(); kept.Output != &output {
		t.Fatal("WithDefaults() replaced a provided stream")
	}
}

package support_test

import (
	"errors"
	"testing"

	"github.com/exphere/ex/src/presentation/cli/support"

	"github.com/spf13/cobra"
)

func TestExactArguments(t *testing.T) {
	t.Parallel()

	command := &cobra.Command{Use: "validate"}
	validate := support.ExactArguments(1)

	if err := validate(command, []string{"demo"}); err != nil {
		t.Fatalf("ExactArguments(1) with one argument = %v", err)
	}

	var usage support.UsageError
	for _, arguments := range [][]string{nil, {"demo", "extra"}} {
		err := validate(command, arguments)
		if !errors.As(err, &usage) || err.Error() != "validate requires 1 argument(s)" || !errors.Is(err, usage.Cause) {
			t.Errorf("ExactArguments(1) with %v = %v, want a usage error", arguments, err)
		}
	}
}

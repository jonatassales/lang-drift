package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/exphere/ex/src/presentation/cli/support"
)

// exitCode maps execution failures to process exit codes.
func exitCode(
	applicationContext context.Context,
	executionError error,
) int {
	var invalidUsage support.UsageError

	switch {
	case applicationContext.Err() != nil ||
		errors.Is(executionError, context.Canceled):
		return 130

	case errors.As(executionError, &invalidUsage):
		return 2

	default:
		return 1
	}
}

// reportError writes errors to the diagnostic stream.
func reportError(
	diagnostics io.Writer,
	executionError error,
	exitStatus int,
	structuredOutput bool,
) {
	if structuredOutput {
		_ = json.NewEncoder(diagnostics).Encode(map[string]any{
			"error":    executionError.Error(),
			"exitCode": exitStatus,
		})
		return
	}

	_, _ = fmt.Fprintf(diagnostics, "ex: %v\n", executionError)
}

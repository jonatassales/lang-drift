package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os/exec"

	app "github.com/jonatassales/tao/exphere/cli/src/application"
	connected "github.com/jonatassales/tao/exphere/cli/src/application/connected"
	"github.com/jonatassales/tao/exphere/cli/src/domain/connection"
	"github.com/jonatassales/tao/exphere/cli/src/presentation/cli/support"
)

// exitCode maps execution failures to the documented contract: 0 success,
// 1 operation failed, 2 invalid usage, 3 refused by a guard or policy,
// 4 missing prerequisite, 130 interrupted. Scripts react to the code; the
// message stays for people.
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

	case errors.Is(executionError, app.ErrAdoptionRefused) ||
		errors.Is(executionError, app.ErrTargetMismatch) ||
		errors.Is(executionError, app.ErrConfigurationChanged) ||
		errors.Is(executionError, app.ErrAlreadyExists) ||
		errors.Is(executionError, connection.ErrRecipeRefused) ||
		errors.Is(executionError, connected.ErrAccessRevoked):
		return 3

	case errors.Is(executionError, app.ErrNotFound) ||
		errors.Is(executionError, exec.ErrNotFound):
		return 4

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
	if hint := hintFor(exitStatus); hint != "" {
		_, _ = fmt.Fprintf(diagnostics, "hint: %s\n", hint)
	}
}

// hintFor turns the exit class into the one next step a person can take.
func hintFor(exitStatus int) string {
	switch exitStatus {
	case 2:
		return "run ex --help for usage"
	case 3:
		return "ex refused to change what it does not own; inspect the current state before retrying"
	case 4:
		return "a prerequisite is missing; run ex doctor"
	default:
		return ""
	}
}

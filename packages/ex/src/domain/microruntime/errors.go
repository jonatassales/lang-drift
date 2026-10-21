package microruntime

import (
	"errors"
	"fmt"
)

// Domain errors classify why EX rejected a definition or refused an action.
// Callers match them with errors.Is; the wrapped message carries the detail.
var (
	ErrInvalidIdentifier     = errors.New("invalid identifier")
	ErrInvalidPath           = errors.New("invalid relative path")
	ErrInvalidQuantity       = errors.New("invalid resource quantity")
	ErrInvalidDescriptor     = errors.New("invalid runtime descriptor")
	ErrInvalidBlueprint      = errors.New("invalid runtime blueprint")
	ErrUnsupportedManifest   = errors.New("unsupported manifest")
	ErrInconsistentExecution = errors.New("inconsistent execution record")
	ErrExistingObligations   = errors.New("runtime still has obligations")
	ErrOwnershipChanged      = errors.New("ownership changed")
	ErrWorkloadFailed        = errors.New("workload did not become ready")
	ErrCleanupBlocked        = errors.New("cleanup blocked")
)

// reject wraps a classification with a formatted reason.
func reject(classification error, format string, arguments ...any) error {
	return fmt.Errorf("%w: %s", classification, fmt.Sprintf(format, arguments...))
}

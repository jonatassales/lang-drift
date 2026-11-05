package microruntime

import (
	"errors"
	"fmt"

	"github.com/jonatassales/tao/lib/core/contracts/go/configuration"
)

// Domain errors classify why the LangDrift CLI rejected a definition or refused an action.
// Callers match them with errors.Is; the wrapped message carries the detail.
var (
	ErrInvalidIdentifier     = configuration.ErrInvalidIdentifier
	ErrInvalidPath           = configuration.ErrInvalidPath
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

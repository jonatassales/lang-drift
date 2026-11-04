package catalog

import (
	"fmt"
	"strings"

	"github.com/jonatassales/tao/lib/cli/src/presentation/cli/support"
	"github.com/jonatassales/tao/lib/core/contracts/go/configuration"
)

// parseReadiness reads a component's explicit readiness: tcp, or http with
// its path as http:/path.
func parseReadiness(field string) (*configuration.Readiness, error) {
	if field == configuration.ReadinessTCP {
		return &configuration.Readiness{Kind: configuration.ReadinessTCP}, nil
	}

	if path, found := strings.CutPrefix(field, configuration.ReadinessHTTP+":"); found {
		return &configuration.Readiness{
			Kind: configuration.ReadinessHTTP,
			Path: path,
		}, nil
	}

	return nil, support.UsageError{Cause: fmt.Errorf("--component readiness %q is neither tcp nor http:/path", field)}
}

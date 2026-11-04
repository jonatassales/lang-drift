package execution

import (
	"strconv"

	"github.com/jonatassales/tao/packages/contracts/configuration"
)

// readinessSignal names what a component answers readiness with: its
// Dockerfile's HEALTHCHECK, its port, or an HTTP path.
func readinessSignal(component configuration.Component) string {
	if component.Readiness == nil {
		return "HEALTHCHECK of " + component.DockerfilePath()
	}

	if component.Readiness.Kind == configuration.ReadinessTCP {
		return "tcp " + component.Name + ":" + strconv.Itoa(component.Port)
	}

	return component.Readiness.Path
}

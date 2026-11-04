package plan

import (
	"strings"

	"github.com/jonatassales/tao/packages/ex/src/domain/microruntime"
	"github.com/jonatassales/tao/packages/ex/src/domain/run"
)

// readinessCheck names how each component is known to be ready, and fails
// naming the Dockerfile that declares no HEALTHCHECK when the fence names
// no explicit readiness either.
func readinessCheck(components []resolved) run.Check {
	var known, unknown []string
	for _, component := range components {
		if component.probe == nil {
			unknown = append(unknown, component.dockerfile+": "+component.readinessErr.Error())

			continue
		}

		switch component.probe.Kind {
		case microruntime.ProbeExec:
			known = append(known, "HEALTHCHECK of "+component.dockerfile)
		case microruntime.ProbeHTTP:
			known = append(known, "http "+component.probe.Path)
		default:
			known = append(known, string(component.probe.Kind))
		}
	}

	if len(unknown) > 0 {
		return run.Check{
			ID:     "readiness-known",
			Result: failed,
			Detail: text(strings.Join(unknown, " · ")),
		}
	}

	return run.Check{
		ID:     "readiness-known",
		Result: passed,
		Detail: text(strings.Join(known, " · ")),
	}
}

package plan

import (
	"github.com/jonatassales/tao/packages/ex/src/domain/microruntime"
	"github.com/jonatassales/tao/packages/ex/src/domain/run"
)

// readinessRecord is the probe as the plan records it, with where it came
// from: the Dockerfile's HEALTHCHECK, or the fence. An unknown readiness is
// null, never a default.
func readinessRecord(probe *microruntime.Probe, dockerfile string) *run.Readiness {
	if probe == nil {
		return nil
	}

	from := "runtime.md"
	if probe.Kind == microruntime.ProbeExec {
		from = "HEALTHCHECK of " + dockerfile
	}

	return &run.Readiness{
		Kind:                string(probe.Kind),
		From:                from,
		Path:                optional(probe.Path),
		Command:             probe.Command,
		PeriodSeconds:       probe.Period,
		TimeoutSeconds:      probe.Timeout,
		FailureThreshold:    probe.Failures,
		InitialDelaySeconds: probe.InitialDelay,
	}
}

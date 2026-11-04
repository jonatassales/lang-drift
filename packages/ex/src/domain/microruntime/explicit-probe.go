package microruntime

import "github.com/jonatassales/tao/packages/contracts/configuration"

// ExplicitProbe is the probe a fence's explicit readiness names, with the
// probe's default timing.
func ExplicitProbe(readiness configuration.Readiness) Probe {
	probe := Probe{
		Kind:     ProbeTCP,
		Period:   ReadinessPeriodSeconds,
		Timeout:  ReadinessTimeoutSeconds,
		Failures: ReadinessFailures,
	}

	if readiness.Kind == configuration.ReadinessHTTP {
		probe.Kind = ProbeHTTP
		probe.Path = readiness.Path
	}

	return probe
}

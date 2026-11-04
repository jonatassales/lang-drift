package microruntime

// ProbeKind selects how Kubernetes decides a container is ready.
type ProbeKind string

const (
	// ProbeExec runs the component's HEALTHCHECK command in the container.
	ProbeExec ProbeKind = "exec"
	// ProbeTCP opens the component's port.
	ProbeTCP ProbeKind = "tcp"
	// ProbeHTTP gets a path on the component's port.
	ProbeHTTP ProbeKind = "http"
)

// Probe is a workload's readiness probe. Period, Timeout and InitialDelay
// are seconds; Failures is the failure threshold.
type Probe struct {
	Kind         ProbeKind
	Path         string
	Command      []string
	Period       int
	Timeout      int
	Failures     int
	InitialDelay int
}

// Validate accepts an exec probe with a command, a tcp probe, or an http
// probe with an absolute path, each with positive timing.
func (probe Probe) Validate() error {
	if probe.Period < 1 || probe.Timeout < 1 || probe.Failures < 1 || probe.InitialDelay < 0 {
		return reject(ErrInvalidBlueprint, "readiness probe requires positive period, timeout and failures")
	}

	switch probe.Kind {
	case ProbeExec:
		if len(probe.Command) == 0 {
			return reject(ErrInvalidBlueprint, "exec readiness requires a command")
		}

		return nil

	case ProbeTCP:
		return nil

	case ProbeHTTP:
		if !isReadinessPath(probe.Path) {
			return reject(ErrInvalidBlueprint, "http readiness requires an absolute path")
		}

		return nil
	}

	return reject(ErrInvalidBlueprint, "readiness probe kind %q is not exec, tcp or http", probe.Kind)
}

// native renders the probe as Kubernetes expects it: its handler with the
// probe's timing.
func (probe Probe) native() map[string]any {
	rendered := map[string]any{
		"periodSeconds":    probe.Period,
		"timeoutSeconds":   probe.Timeout,
		"failureThreshold": probe.Failures,
	}

	if probe.InitialDelay > 0 {
		rendered["initialDelaySeconds"] = probe.InitialDelay
	}

	key, handler := probe.Handler()
	rendered[key] = handler

	return rendered
}

// Handler is the probe's Kubernetes handler and the key it goes under:
// httpGet or tcpSocket on the named port, or exec with the command.
func (probe Probe) Handler() (string, map[string]any) {
	switch probe.Kind {
	case ProbeHTTP:
		return "httpGet", map[string]any{
			"path": probe.Path,
			"port": httpPortName,
		}

	case ProbeTCP:
		return "tcpSocket", map[string]any{"port": httpPortName}
	}

	command := make([]any, 0, len(probe.Command))
	for _, argument := range probe.Command {
		command = append(command, argument)
	}

	return "exec", map[string]any{"command": command}
}

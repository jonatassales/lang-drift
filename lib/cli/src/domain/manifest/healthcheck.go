package manifest

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jonatassales/tao/lib/cli/src/domain/microruntime"
)

// ErrNoHealthcheck marks a Dockerfile that declares no HEALTHCHECK of its
// own, or turns it off with HEALTHCHECK NONE. One inherited from a base image
// is not seen: the component's own Dockerfile declares its readiness.
var ErrNoHealthcheck = errors.New("no HEALTHCHECK of its own")

// Healthcheck reads the last HEALTHCHECK of the Dockerfile's final stage, the
// one the built image keeps, as an exec readiness probe: the shell form runs through /bin/sh -c, the
// exec form as its array, and the options become the probe's timing, with
// the probe's defaults where one is absent.
func Healthcheck(dockerfile []byte) (microruntime.Probe, error) {
	var declared *instruction

	for _, step := range instructions(dockerfile) {
		switch step.keyword {
		case "FROM":
			declared = nil
		case "HEALTHCHECK":
			declared = &step
		}
	}

	if declared == nil {
		return microruntime.Probe{}, ErrNoHealthcheck
	}

	probe := microruntime.Probe{
		Kind:     microruntime.ProbeExec,
		Period:   microruntime.ReadinessPeriodSeconds,
		Timeout:  microruntime.ReadinessTimeoutSeconds,
		Failures: microruntime.ReadinessFailures,
	}

	rest := declared.arguments
	for strings.HasPrefix(rest, "--") {
		option, after := cutWord(rest)
		if err := applyOption(&probe, option); err != nil {
			return microruntime.Probe{}, err
		}

		rest = after
	}

	form, command := cutWord(rest)

	switch strings.ToUpper(form) {
	case "NONE":
		return microruntime.Probe{}, fmt.Errorf("%w: HEALTHCHECK NONE", ErrNoHealthcheck)

	case "CMD":
		arguments, err := commandOf(command)
		if err != nil {
			return microruntime.Probe{}, err
		}

		probe.Command = arguments

		return probe, nil
	}

	return microruntime.Probe{}, errors.New("HEALTHCHECK names no CMD")
}

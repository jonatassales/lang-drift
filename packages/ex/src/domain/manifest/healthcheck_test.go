package manifest_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/jonatassales/tao/packages/ex/src/domain/manifest"
	"github.com/jonatassales/tao/packages/ex/src/domain/microruntime"
)

func TestHealthcheckRendersTheComponentsOwnProbe(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		dockerfile string
		want       microruntime.Probe
	}{
		{
			"shell form without options",
			"FROM redis:8\nHEALTHCHECK CMD redis-cli ping | grep PONG\n",
			microruntime.Probe{
				Kind:     microruntime.ProbeExec,
				Command:  []string{"/bin/sh", "-c", "redis-cli ping | grep PONG"},
				Period:   microruntime.ReadinessPeriodSeconds,
				Timeout:  microruntime.ReadinessTimeoutSeconds,
				Failures: microruntime.ReadinessFailures,
			},
		},
		{
			"exec form with every option",
			"FROM node:24\nHEALTHCHECK --interval=10s --timeout=3s --retries=5 --start-period=1m \\\n  CMD [\"node\", \"health.js\"]\n",
			microruntime.Probe{
				Kind:         microruntime.ProbeExec,
				Command:      []string{"node", "health.js"},
				Period:       10,
				Timeout:      3,
				Failures:     5,
				InitialDelay: 60,
			},
		},
		{
			"the final stage's own",
			"FROM golang AS build\nHEALTHCHECK CMD false\nFROM alpine\nHEALTHCHECK CMD true\n",
			microruntime.Probe{
				Kind:     microruntime.ProbeExec,
				Command:  []string{"/bin/sh", "-c", "true"},
				Period:   microruntime.ReadinessPeriodSeconds,
				Timeout:  microruntime.ReadinessTimeoutSeconds,
				Failures: microruntime.ReadinessFailures,
			},
		},
		{
			"the last one wins, in any case",
			"FROM alpine\nHEALTHCHECK CMD true\n# a comment\nhealthcheck --interval=1500ms cmd [\"wget\", \"-q\", \"--spider\", \"http://127.0.0.1:3402/\"]\n",
			microruntime.Probe{
				Kind:     microruntime.ProbeExec,
				Command:  []string{"wget", "-q", "--spider", "http://127.0.0.1:3402/"},
				Period:   2,
				Timeout:  microruntime.ReadinessTimeoutSeconds,
				Failures: microruntime.ReadinessFailures,
			},
		},
	}

	for _, test := range tests {
		probe, err := manifest.Healthcheck([]byte(test.dockerfile))
		if err != nil {
			t.Errorf("%s: Healthcheck() refused: %v", test.name, err)

			continue
		}

		if !reflect.DeepEqual(probe, test.want) {
			t.Errorf("%s: Healthcheck() = %+v, want %+v", test.name, probe, test.want)
		}
	}
}

func TestHealthcheckRefusesADockerfileWithoutOne(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		dockerfile string
		absent     bool
	}{
		{"no HEALTHCHECK", "FROM redis:8\nEXPOSE 6379\n", true},
		{"HEALTHCHECK NONE", "FROM app\nHEALTHCHECK NONE\n", true},
		{"NONE after a command", "FROM app\nHEALTHCHECK CMD true\nHEALTHCHECK NONE\n", true},
		{"only a builder stage declares one", "FROM golang AS build\nHEALTHCHECK CMD true\nFROM alpine\nCOPY --from=build /app /app\n", true},
		{"no command", "FROM app\nHEALTHCHECK --interval=5s\n", false},
		{"unknown option", "FROM app\nHEALTHCHECK --every=5s CMD true\n", false},
		{"bad duration", "FROM app\nHEALTHCHECK --interval=soon CMD true\n", false},
		{"bad exec form", "FROM app\nHEALTHCHECK CMD [\"true\"\n", false},
	}

	for _, test := range tests {
		_, err := manifest.Healthcheck([]byte(test.dockerfile))
		if err == nil {
			t.Errorf("%s: Healthcheck() accepted it", test.name)

			continue
		}

		if errors.Is(err, manifest.ErrNoHealthcheck) != test.absent {
			t.Errorf("%s: Healthcheck() = %v, want absent=%t", test.name, err, test.absent)
		}
	}
}

package manifest_test

import (
	"testing"

	"github.com/jonatassales/tao/lib/cli/src/domain/manifest"
)

func TestExposedPortIsTheFirstExpose(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		dockerfile string
		want       int
	}{
		{"one port", "FROM node\nEXPOSE 3001\n", 3001},
		{"the first of several", "FROM node\nexpose 3402/tcp 9229\nEXPOSE 8080\n", 3402},
		{"none", "FROM node\nCMD [\"node\"]\n", 0},
		{"not a port", "FROM node\nEXPOSE $PORT\n", 0},
	}

	for _, test := range tests {
		if port := manifest.ExposedPort([]byte(test.dockerfile)); port != test.want {
			t.Errorf("%s: ExposedPort() = %d, want %d", test.name, port, test.want)
		}
	}
}

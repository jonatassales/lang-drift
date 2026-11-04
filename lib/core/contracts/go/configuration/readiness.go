package configuration

import (
	"errors"
	"strings"
)

// The readiness kinds a fence may name. Without one, readiness is the
// component Dockerfile's own HEALTHCHECK.
const (
	ReadinessTCP  = "tcp"
	ReadinessHTTP = "http"
)

// Readiness is an explicit readiness signal for a component whose image
// cannot run its HEALTHCHECK command, or whose page proves readiness better.
type Readiness struct {
	Kind string `json:"kind"`
	Path string `json:"path,omitempty"`
}

// Validate accepts a tcp probe without a path, or an http probe with an
// absolute path.
func (readiness Readiness) Validate() error {
	switch readiness.Kind {
	case ReadinessTCP:
		if readiness.Path != "" {
			return errors.New("tcp readiness takes no path")
		}

		return nil

	case ReadinessHTTP:
		if !IsPlainText(readiness.Path) || !strings.HasPrefix(readiness.Path, "/") || strings.ContainsAny(readiness.Path, " \t") {
			return errors.New("http readiness requires a path starting with /")
		}

		return nil
	}

	return errors.New("readiness kind is tcp or http")
}

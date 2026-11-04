package connected

import (
	"fmt"
	"path/filepath"

	"github.com/jonatassales/tao/packages/ex/src/domain/connection"
	"github.com/jonatassales/tao/packages/ex/src/domain/manifest"
	"github.com/jonatassales/tao/packages/ex/src/domain/microruntime"
	"github.com/jonatassales/tao/packages/ex/src/domain/project"
	"github.com/jonatassales/tao/packages/ex/src/infrastructure/files"
)

// dockerfileLimit bounds the Dockerfile a component is read from.
const dockerfileLimit = 1 << 20

// componentService is the component's port and readiness: the fence's when
// it names them, otherwise the first EXPOSE and the HEALTHCHECK of its own
// Dockerfile.
func componentService(owner project.Project, component connection.Component) (int, microruntime.Probe, error) {
	if component.Port != 0 && component.Readiness != nil {
		return component.Port, microruntime.ExplicitProbe(*component.Readiness), nil
	}

	dockerfile := component.DockerfilePath()

	location, err := files.Existing(owner.Root, filepath.FromSlash(dockerfile))
	if err != nil {
		return 0, microruntime.Probe{}, fmt.Errorf("the Dockerfile %s: %w", dockerfile, err)
	}

	content, err := files.ReadBounded(location, dockerfileLimit)
	if err != nil {
		return 0, microruntime.Probe{}, fmt.Errorf("the Dockerfile %s: %w", dockerfile, err)
	}

	port := component.Port
	if port == 0 {
		port = manifest.ExposedPort(content)
	}

	if port == 0 {
		return 0, microruntime.Probe{}, fmt.Errorf("%s exposes no port and the fence names none", dockerfile)
	}

	if component.Readiness != nil {
		return port, microruntime.ExplicitProbe(*component.Readiness), nil
	}

	probe, err := manifest.Healthcheck(content)
	if err != nil {
		return 0, microruntime.Probe{}, fmt.Errorf("%s: %w; declare one or name tcp or http readiness", dockerfile, err)
	}

	return port, probe, nil
}

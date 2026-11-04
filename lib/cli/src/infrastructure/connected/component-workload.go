package connected

import (
	"fmt"
	"strconv"

	"github.com/jonatassales/tao/lib/cli/src/domain/connection"
	"github.com/jonatassales/tao/lib/cli/src/domain/harness"
	"github.com/jonatassales/tao/lib/cli/src/domain/microruntime"
	"github.com/jonatassales/tao/lib/cli/src/domain/project"
	"github.com/jonatassales/tao/lib/core/contracts/go/configuration"
)

// componentWorkload is one component as a workload: named after the
// component so others address it by name, running the image derived for this
// revision, on the port and with the readiness its Dockerfile declares where
// the fence leaves them out.
func componentWorkload(owner project.Project, config connection.Configuration, component connection.Component, share harness.Resources) (microruntime.Workload, error) {
	name, err := microruntime.ParseIdentifier(component.Name)
	if err != nil {
		return microruntime.Workload{}, err
	}

	relative, err := microruntime.ParseComponentPath(component.Path)
	if err != nil {
		return microruntime.Workload{}, err
	}

	port, probe, err := componentService(owner, component)
	if err != nil {
		return microruntime.Workload{}, fmt.Errorf("component %s: %w", component.Name, err)
	}

	secret, err := envSecretName(component)
	if err != nil {
		return microruntime.Workload{}, fmt.Errorf("component %s: %w", component.Name, err)
	}

	cpu, err := microruntime.ParseQuantity(strconv.Itoa(share.CPUMillis) + "m")
	if err != nil {
		return microruntime.Workload{}, err
	}

	memory, err := microruntime.ParseQuantity(strconv.Itoa(share.MemoryMiB) + "Mi")
	if err != nil {
		return microruntime.Workload{}, err
	}

	return microruntime.Workload{
		Name:      name,
		Component: relative,
		Image:     configuration.ImageTag(config.RuntimeID, component.Name, config.Revision),
		Kind:      microruntime.DeploymentWorkload,
		Resources: microruntime.Resources{
			CPURequest:    cpu,
			CPULimit:      cpu,
			MemoryRequest: memory,
			MemoryLimit:   memory,
		},
		Service: microruntime.ServiceSpec{
			Port:     port,
			Replicas: 1,
			Probe:    probe,
		},
		EnvSecret: secret,
	}, nil
}

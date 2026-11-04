package connected

import (
	"fmt"

	"github.com/jonatassales/tao/packages/ex/src/domain/connection"
	"github.com/jonatassales/tao/packages/ex/src/domain/microruntime"
	"github.com/jonatassales/tao/packages/ex/src/domain/project"
)

// componentEnvironments reads what each component's container receives:
// its committed env file, and the variables it passes through from the
// agent's environment. A pass-through name the agent does not hold refuses
// the start with the name, before anything is built.
func componentEnvironments(owner project.Project, descriptor microruntime.Descriptor, config connection.Configuration, lookup func(string) (string, bool)) ([]microruntime.Environment, error) {
	units := map[string]microruntime.Identifier{}
	for name, unit := range descriptor.Units {
		units[unit.Component.String()] = name
	}

	var environments []microruntime.Environment

	for _, component := range config.Components {
		secret, err := envSecretName(component)
		if err != nil {
			return nil, fmt.Errorf("component %s: %w", component.Name, err)
		}

		if secret.IsZero() {
			continue
		}

		variables, err := envFileVariables(owner, component)
		if err != nil {
			return nil, err
		}

		for _, name := range component.PassEnv {
			value, found := lookup(name)
			if !found {
				return nil, fmt.Errorf("%w: component %s passes through %s", connection.ErrPassEnvMissing, component.Name, name)
			}

			variables[name] = value
		}

		environments = append(environments, microruntime.Environment{
			Unit:      units[component.Path],
			Secret:    secret,
			Variables: variables,
		})
	}

	return environments, nil
}

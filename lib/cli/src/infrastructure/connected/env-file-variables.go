package connected

import (
	"fmt"
	"path/filepath"

	"github.com/jonatassales/tao/lib/cli/src/domain/connection"
	"github.com/jonatassales/tao/lib/cli/src/domain/project"
	"github.com/jonatassales/tao/lib/cli/src/infrastructure/files"
	"github.com/jonatassales/tao/lib/core/contracts/go/configuration"
)

// envFileLimit bounds a component's committed env file.
const envFileLimit = 1 << 20

// envFileVariables reads the component's committed env file, or nothing
// when it names none. The file must not set a name it also passes through.
func envFileVariables(owner project.Project, component connection.Component) (map[string]string, error) {
	if component.EnvFile == "" {
		return map[string]string{}, nil
	}

	location, err := files.Existing(owner.Root, filepath.FromSlash(component.EnvFile))
	if err != nil {
		return nil, fmt.Errorf("component %s env file: %w", component.Name, err)
	}

	content, err := files.ReadBounded(location, envFileLimit)
	if err != nil {
		return nil, fmt.Errorf("component %s env file: %w", component.Name, err)
	}

	variables, err := configuration.ParseEnvFile(string(content))
	if err != nil {
		return nil, fmt.Errorf("component %s env file %s: %w", component.Name, component.EnvFile, err)
	}

	if err := component.CheckEnv(variables); err != nil {
		return nil, err
	}

	return variables, nil
}

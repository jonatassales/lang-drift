package connected

import (
	"github.com/jonatassales/tao/lib/cli/src/domain/connection"
	"github.com/jonatassales/tao/lib/cli/src/domain/microruntime"
)

// envSecretName is the Secret a component's container reads its
// environment from, or none when the component names no env file and no
// pass-through variable.
func envSecretName(component connection.Component) (microruntime.Identifier, error) {
	if component.EnvFile == "" && len(component.PassEnv) == 0 {
		return microruntime.Identifier{}, nil
	}

	return microruntime.ParseIdentifier(component.Name + "-env")
}

package plan

import "github.com/jonatassales/tao/lib/core/contracts/go/configuration"

// reference is the component as the fence names it: the files it points at
// and what refines them.
func (input ComponentInput) reference() configuration.Component {
	return configuration.Component{
		Name:       input.Name,
		Path:       input.Path,
		Dockerfile: input.Dockerfile,
		Port:       input.Port,
		Readiness:  input.Readiness,
		EnvFile:    input.EnvFile,
		PassEnv:    input.PassEnv,
	}
}

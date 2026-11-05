package configuration

import (
	"fmt"
	"path"
)

// DefaultDockerfile is the file a component builds from when it names none.
const DefaultDockerfile = "Dockerfile"

// maximumPassEnv bounds the variables one component reads from the agent.
const maximumPassEnv = 64

// Component is a reference to files: a folder in the checkout and the
// Dockerfile inside it. The image, the build, the port and the readiness are
// derived from that Dockerfile; the fields here only refine what it says.
type Component struct {
	ID         string     `json:"id,omitempty"`
	Name       string     `json:"name"`
	Kind       string     `json:"kind,omitempty"`
	Path       string     `json:"path"`
	Dockerfile string     `json:"dockerfile,omitempty"`
	Port       int        `json:"port,omitempty"`
	Readiness  *Readiness `json:"readiness,omitempty"`
	EnvFile    string     `json:"envFile,omitempty"`
	PassEnv    []string   `json:"passEnv,omitempty"`
}

// DockerfileName is the Dockerfile the component names, or the default.
func (component Component) DockerfileName() string {
	if component.Dockerfile == "" {
		return DefaultDockerfile
	}

	return component.Dockerfile
}

// DockerfilePath is the Dockerfile's path from the project root.
func (component Component) DockerfilePath() string {
	return path.Join(component.Path, component.DockerfileName())
}

// Validate checks the component's own fields. Whether its files exist is
// the planner's and the runner's check, not a lexical one.
func (component Component) Validate() error {
	pathRefused := CheckComponentPath(component.Path) != nil ||
		IsPlatformPath(component.Path)

	if pathRefused {
		return invalid(fmt.Sprintf("component path %q must be a clean relative path inside the checkout and outside .langdrift", component.Path))
	}

	if err := CheckIdentifier(component.Name); err != nil {
		return invalid(fmt.Sprintf("component %s: %v", component.Path, err))
	}

	if !isFileName(component.DockerfileName()) {
		return invalid(fmt.Sprintf("component %s: dockerfile %q must be a file name inside the component folder", component.Name, component.Dockerfile))
	}

	portInvalid := component.Port != 0 && !ValidPort(component.Port)

	if portInvalid {
		return invalid(fmt.Sprintf("component %s: port %d is outside 1 to 65535", component.Name, component.Port))
	}

	if component.Readiness != nil {
		if err := component.Readiness.Validate(); err != nil {
			return invalid(fmt.Sprintf("component %s: %v", component.Name, err))
		}
	}

	if err := component.validateEnvFile(); err != nil {
		return err
	}

	return component.validatePassEnv()
}

// CheckEnv refuses a pass-through name the env file also sets, so a value
// never has two sources.
func (component Component) CheckEnv(variables map[string]string) error {
	for _, name := range component.PassEnv {
		if _, found := variables[name]; found {
			return invalid(fmt.Sprintf("component %s: %s is both in %s and passed through", component.Name, name, component.EnvFile))
		}
	}

	return nil
}

func (component Component) validateEnvFile() error {
	if component.EnvFile == "" {
		return nil
	}

	envFileRefused := CheckRelativePath(component.EnvFile) != nil ||
		IsPlatformPath(component.EnvFile)

	if envFileRefused {
		return invalid(fmt.Sprintf("component %s: envFile %q must be a clean relative path inside the checkout and outside .langdrift", component.Name, component.EnvFile))
	}

	return nil
}

func (component Component) validatePassEnv() error {
	if len(component.PassEnv) > maximumPassEnv {
		return invalid(fmt.Sprintf("component %s: at most %d pass-through variables", component.Name, maximumPassEnv))
	}

	seen := map[string]bool{}
	for _, name := range component.PassEnv {
		if err := CheckVariableName(name); err != nil {
			return invalid(fmt.Sprintf("component %s: %v", component.Name, err))
		}

		if seen[name] {
			return invalid(fmt.Sprintf("component %s: pass-through %s is named twice", component.Name, name))
		}

		seen[name] = true
	}

	return nil
}

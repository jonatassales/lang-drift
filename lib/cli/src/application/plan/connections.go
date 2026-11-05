package plan

import "github.com/jonatassales/tao/lib/cli/src/domain/project"

// Connections reads the console a project is connected to. A project that
// never connected reads as a Connection that is not recorded, not as an
// error.
type Connections interface {
	Read(owner project.Project) (Connection, error)
}

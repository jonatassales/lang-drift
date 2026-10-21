package support

import (
	"os"
	"path/filepath"

	"github.com/exphere/ex/src/application"
)

// Selection maps the project flags to a project selection. An explicit
// --project names the project root itself; otherwise the nearest project at
// or above the fallback or working directory is used.
//
// Resolving a path is presentation's job. Deciding whether it holds a valid
// project is the application's.
func (globals *Globals) Selection() (application.ProjectSelection, error) {
	switch {
	case globals.Project != "":
		directory, err := filepath.Abs(globals.Project)

		return application.ProjectSelection{Directory: directory, Explicit: true}, err

	case globals.Fallback != "":
		directory, err := filepath.Abs(globals.Fallback)

		return application.ProjectSelection{Directory: directory}, err
	}

	directory, err := os.Getwd()

	return application.ProjectSelection{Directory: directory}, err
}

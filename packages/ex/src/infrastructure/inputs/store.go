// Package inputs keeps the guardrail inputs privately, one record per kind
// under the project's private state, each with the record it replaced.
package inputs

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"

	"github.com/jonatassales/tao/packages/ex/src/application"
	app "github.com/jonatassales/tao/packages/ex/src/application/inputs"
	"github.com/jonatassales/tao/packages/ex/src/domain/project"
	"github.com/jonatassales/tao/packages/ex/src/infrastructure/files"
)

// The records live under the private state; the lock serializes writers.
const (
	directory = ".langdrift/agent/state/inputs"
	lockFile  = ".langdrift/agent/state/inputs.lock"
)

var kindPattern = regexp.MustCompile(`^[a-z]+$`)

// Store is the file store of the inputs.
type Store struct{}

// Load reads the current record of one kind, or reports application.ErrNotFound.
func (Store) Load(owner project.Project, kind string) (app.Record, error) {
	path, err := recordPath(owner, kind)
	if err != nil {
		return app.Record{}, err
	}

	var record app.Record

	err = files.ReadJSON(path, &record)

	switch {
	case errors.Is(err, fs.ErrNotExist):
		return app.Record{}, application.NotFound("no %s input is recorded; set it with ex harness or ex finops", kind)

	case err != nil:
		return app.Record{}, err

	case record.Kind != kind:
		return app.Record{}, fmt.Errorf("the %s input record belongs to %s", kind, record.Kind)
	}

	return record, nil
}

// Save replaces the record of its kind, under the inputs lock.
func (Store) Save(owner project.Project, record app.Record) error {
	path, err := recordPath(owner, record.Kind)
	if err != nil {
		return err
	}

	lock, err := files.Within(owner.Root, lockFile)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(lock), 0o700); err != nil {
		return err
	}

	unlock, err := files.Lock(lock)
	if err != nil {
		return err
	}
	defer unlock()

	return files.WriteJSON(path, record)
}

func recordPath(owner project.Project, kind string) (string, error) {
	if !kindPattern.MatchString(kind) {
		return "", fmt.Errorf("input kind %q is not a word", kind)
	}

	return files.Within(owner.Root, directory+"/"+kind+".json")
}

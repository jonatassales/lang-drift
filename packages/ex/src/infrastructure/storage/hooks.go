package storage

import (
	"path/filepath"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/domain/project"
	"github.com/exphere/ex/src/infrastructure/files"
)

// RecordHook replaces a hook's intent or outcome record without keeping its
// output, arguments, or environment in durable evidence.
func (Evidence) RecordHook(owner project.Project, outcome application.HookOutcome) error {
	return files.WriteJSON(filepath.Join(LayoutOf(owner).evidence(outcome.RunID), "hook.json"), outcome)
}

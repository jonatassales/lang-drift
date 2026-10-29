package storage

import (
	"fmt"

	"github.com/jonatassales/tao/packages/ex/src/application"
	"github.com/jonatassales/tao/packages/ex/src/domain/microruntime"
	"github.com/jonatassales/tao/packages/ex/src/domain/project"
	"github.com/jonatassales/tao/packages/ex/src/infrastructure/files"
)

// RecordHook replaces a hook's intent or outcome record without keeping its
// output, arguments, or environment in durable evidence.
func (Evidence) RecordHook(owner project.Project, outcome application.HookOutcome) error {
	runtime, err := microruntime.ParseIdentifier(outcome.Runtime)
	if err != nil {
		return fmt.Errorf("hook evidence for runtime %q: %w", outcome.Runtime, err)
	}

	layout := LayoutOf(owner)

	if err := ensureRunRecord(layout, runtime); err != nil {
		return err
	}

	return files.WriteJSON(layout.hookEvidence(runtime, outcome.RunID), outcome)
}

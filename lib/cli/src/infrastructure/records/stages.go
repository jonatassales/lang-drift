package records

import (
	"path/filepath"

	"github.com/jonatassales/tao/lib/cli/src/application/execution"
	"github.com/jonatassales/tao/lib/cli/src/domain/microruntime"
	"github.com/jonatassales/tao/lib/cli/src/domain/project"
	"github.com/jonatassales/tao/lib/cli/src/infrastructure/files"
)

// WriteStageEvidence writes the token-free record of one stage call under
// the lap's evidence in the runtime folder.
func (Runtimes) WriteStageEvidence(owner project.Project, runtime microruntime.Identifier, evidence execution.StageEvidence) error {
	path, err := runtimeFile(owner, runtime, filepath.Join("evidence", "laps", evidence.Lap, "stages", evidence.Stage+"-"+evidence.Skill+".json"))
	if err != nil {
		return err
	}

	data, err := Format(evidence, FileProfile)
	if err != nil {
		return err
	}

	return files.WriteDocument(path, data, documentMode)
}

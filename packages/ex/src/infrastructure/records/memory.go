package records

import (
	"strings"

	"github.com/jonatassales/tao/packages/ex/src/application/execution"
	"github.com/jonatassales/tao/packages/ex/src/domain/microruntime"
	"github.com/jonatassales/tao/packages/ex/src/domain/project"
	"github.com/jonatassales/tao/packages/ex/src/infrastructure/files"
)

// WriteMemory renders and writes the run's memory.
func (Runtimes) WriteMemory(owner project.Project, runtime microruntime.Identifier, memory execution.Memory) error {
	path, err := runtimeFile(owner, runtime, MemoryFile)
	if err != nil {
		return err
	}

	markdown, err := RenderMemory(memory)
	if err != nil {
		return err
	}

	return files.WriteDocument(path, markdown, documentMode)
}

// RenderMemory renders what the run teaches: its headline, its outcome,
// the laps' story, the lessons and the langdrift-run fence the RAG indexes.
func RenderMemory(memory execution.Memory) ([]byte, error) {
	var out strings.Builder

	out.WriteString("# Run " + memory.Runtime + ": " + memory.Name + "\n\n")
	out.WriteString("**Ticket:** " + memory.Ticket + "\n")
	out.WriteString("**Kind:** " + memory.Kind + "\n")
	out.WriteString("**Outcome:** " + memory.Outcome + "\n")
	out.WriteString("**Proposed ticket status:** " + memory.ProposedStatus + "\n\n")
	if memory.Story != "" {
		out.WriteString(memory.Story + "\n\n")
	}

	out.WriteString("## Worth remembering\n\n")
	if len(memory.Lessons) == 0 {
		out.WriteString("- The run measured nothing worth keeping yet.\n")
	}
	for _, lesson := range memory.Lessons {
		out.WriteString("- " + lesson + "\n")
	}

	fence, err := Format(memory.Fence, FenceProfile)
	if err != nil {
		return nil, err
	}

	out.WriteString("\n```langdrift-run\n" + string(fence) + "```\n")

	return []byte(out.String()), nil
}

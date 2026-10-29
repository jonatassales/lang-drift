package records

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/jonatassales/tao/packages/ex/src/application"
	"github.com/jonatassales/tao/packages/ex/src/application/execution"
	"github.com/jonatassales/tao/packages/ex/src/domain/finops"
	"github.com/jonatassales/tao/packages/ex/src/domain/microruntime"
	"github.com/jonatassales/tao/packages/ex/src/domain/project"
	"github.com/jonatassales/tao/packages/ex/src/domain/run"
	"github.com/jonatassales/tao/packages/ex/src/infrastructure/files"
)

// LapStore keeps each lap's private state under the project's private
// state, from which the lap file is rendered.
type LapStore struct{}

var lapNumber = regexp.MustCompile(`^\d{4}\.json$`)

func lapStatePath(owner project.Project, runtime microruntime.Identifier, number string) (string, error) {
	if !regexp.MustCompile(`^\d{4}$`).MatchString(number) {
		return "", fmt.Errorf("lap number %q is not four digits", number)
	}

	return files.Within(owner.Root, filepath.Join(".langdrift", "agent", "state", "runs", runtime.String(), "laps", number+".json"))
}

// Load reads one lap's state, or reports application.ErrNotFound.
func (LapStore) Load(owner project.Project, runtime microruntime.Identifier, number string) (execution.Lap, error) {
	path, err := lapStatePath(owner, runtime, number)
	if err != nil {
		return execution.Lap{}, err
	}

	var lap execution.Lap

	err = files.ReadJSON(path, &lap)
	if errors.Is(err, fs.ErrNotExist) {
		return execution.Lap{}, application.NotFound("lap %s of %s does not exist", number, runtime)
	}

	return lap, err
}

// Save writes one lap's state.
func (LapStore) Save(owner project.Project, runtime microruntime.Identifier, lap execution.Lap) error {
	path, err := lapStatePath(owner, runtime, lap.Number)
	if err != nil {
		return err
	}

	return files.WriteJSON(path, lap)
}

// Open finds the lap that has not ended, if any.
func (store LapStore) Open(owner project.Project, runtime microruntime.Identifier) (execution.Lap, bool, error) {
	directory, err := files.Within(owner.Root, filepath.Join(".langdrift", "agent", "state", "runs", runtime.String(), "laps"))
	if err != nil {
		return execution.Lap{}, false, err
	}

	entries, err := os.ReadDir(directory)
	if errors.Is(err, fs.ErrNotExist) {
		return execution.Lap{}, false, nil
	}
	if err != nil {
		return execution.Lap{}, false, err
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if lapNumber.MatchString(entry.Name()) {
			names = append(names, strings.TrimSuffix(entry.Name(), ".json"))
		}
	}
	sort.Strings(names)

	for index := len(names) - 1; index >= 0; index-- {
		lap, err := store.Load(owner, runtime, names[index])
		if err != nil {
			return execution.Lap{}, false, err
		}
		if !lap.Over() {
			return lap, true, nil
		}
	}

	return execution.Lap{}, false, nil
}

// WriteLap renders a lap's file from its state.
func (Runtimes) WriteLap(owner project.Project, runtime microruntime.Identifier, lap execution.Lap) error {
	path, err := runtimeFile(owner, runtime, filepath.Join("laps", lap.Number+".lap.md"))
	if err != nil {
		return err
	}

	markdown, err := RenderLap(lap)
	if err != nil {
		return err
	}

	return files.WriteDocument(path, markdown, documentMode)
}

// RenderLap renders one lap: its headline and choice, its steps, what it
// observed, and the langdrift-lap fence.
func RenderLap(lap execution.Lap) ([]byte, error) {
	var out strings.Builder

	out.WriteString("# Lap " + lap.Number + ": " + headline(lap) + "\n\n")
	out.WriteString("**Run:** " + lap.Run + "\n")
	out.WriteString("**Choice:** " + orDefault(lap.ModelName, lap.Choice.Model) + " at " + lap.Choice.Effort + " effort, " + chosen(lap) + "\n")
	out.WriteString("**Outcome:** " + outcomeWords(lap) + "\n")
	if lap.Runtime != nil {
		out.WriteString("**Runtime:** `" + lap.Runtime.ID + "`, revision " + fmt.Sprint(lap.Runtime.Revision) + ", `runtime.md` sha256 " + lap.Runtime.MarkdownSHA256 + "\n")
	}
	if lap.Execution != "" {
		out.WriteString("**Evidence:** `evidence/" + lap.Execution + "/` in this runtime\n")
	}

	out.WriteString("\n| Step | At | Note |\n| --- | --- | --- |\n")
	for _, step := range lap.Steps {
		out.WriteString("| " + step.Name + " | " + step.At.UTC().Format("15:04:05") + " | " + step.Note + " |\n")
	}

	if lap.NothingPaid() {
		out.WriteString("\nNo runtime revision was materialized and no execution was claimed, so nothing was observed.\n")
	}
	if lap.Unmeasured {
		out.WriteString("\nA call reported no usage, so the observed spend is a floor, not a measurement.\n")
	}

	fence, err := Format(fenceOf(lap), FenceProfile)
	if err != nil {
		return nil, err
	}

	out.WriteString("\n```langdrift-lap\n" + string(fence) + "```\n")

	return []byte(out.String()), nil
}

func fenceOf(lap execution.Lap) run.LapFence {
	fence := run.LapFence{Run: lap.Run, Lap: lap.Number, Choice: run.ChoiceFence{Model: lap.Choice.Model, Effort: lap.Choice.Effort}, Outcome: lap.Outcome, Runtime: lap.Runtime, Observed: lap.Observed}
	if lap.Gate != "" {
		fence.Gate = &lap.Gate
	}
	if lap.Execution != "" {
		execution := lap.Execution
		fence.Execution = &execution
	}

	return fence
}

func headline(lap execution.Lap) string {
	switch lap.Outcome {
	case execution.PausedAtGate:
		return "pit stop at the gate"
	case execution.BlockedAt:
		return "blocked at the gate"
	case execution.Open:
		return "the lap is open"
	case execution.Completed:
		return lap.Run + " ran"
	case execution.Failed:
		return "the lap failed"
	}

	return "outcome unknown"
}

func chosen(lap execution.Lap) string {
	if lap.Choice.Overrode != nil {
		return "chosen by " + lap.Choice.By + " over the " + *lap.Choice.Overrode
	}

	return "the suggestion"
}

func outcomeWords(lap execution.Lap) string {
	switch {
	case lap.NothingPaid() && lap.Outcome == execution.PausedAtGate:
		return "paused at gate " + lap.Gate + ", nothing paid"
	case lap.NothingPaid():
		return "blocked at gate " + lap.Gate + ", nothing paid"
	case lap.Outcome == execution.PausedAtGate:
		return "paused at gate " + lap.Gate + ", " + finops.Dollars(lap.ObservedCents) + " observed of " + finops.Dollars(lap.Reserved) + " reserved"
	case lap.Outcome == execution.Open && lap.StagesCompleted != nil:
		return "open, the stages done, " + finops.Dollars(lap.ObservedCents) + " observed of " + finops.Dollars(lap.Reserved) + " reserved"
	case lap.Outcome == execution.Open && (lap.ObservedCents > 0 || lap.Unmeasured):
		return "open, " + finops.Dollars(lap.ObservedCents) + " observed so far of " + finops.Dollars(lap.Reserved) + " reserved"
	case lap.Outcome == execution.Open:
		return "open, " + finops.Dollars(lap.Reserved) + " reserved"
	}

	return lap.Outcome
}

func orDefault(value, fallback string) string {
	if value == "" {
		return fallback
	}

	return value
}

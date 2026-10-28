// Package storage keeps EX's records inside the consumer project: the project
// identity, the private state under .langdrift/agent/state (cluster and
// execution ledgers, the project lock, preparation state), and the runtimes
// under .langdrift/microruntimes with their evidence. Records are versioned
// JSON documents written atomically.
package storage

import (
	"path/filepath"

	"github.com/jonatassales/tao/packages/ex/src/domain/microruntime"
	"github.com/jonatassales/tao/packages/ex/src/domain/project"
)

// Layout names EX's files inside one project.
type Layout struct{ root string }

// LayoutOf returns the layout of owner's project.
func LayoutOf(owner project.Project) Layout { return Layout{root: owner.Root} }

func layoutAt(root string) Layout { return Layout{root: root} }

// RuntimeState is the private directory of one runtime: its execution ledger
// and its preparation state.
func (layout Layout) RuntimeState(runtime microruntime.Identifier) string {
	return filepath.Join(layout.runtimes(), runtime.String())
}

// Runtime is the runtime's own folder: its authored and materialized
// configuration, its records and its evidence.
func (layout Layout) Runtime(runtime microruntime.Identifier) string {
	return filepath.Join(layout.root, filepath.FromSlash(microruntime.ConfigurationDirectory(runtime).String()))
}

func (layout Layout) metadataDirectory() string {
	return filepath.Join(layout.root, ".langdrift")
}

func (layout Layout) metadata() string {
	return filepath.Join(layout.metadataDirectory(), "config.json")
}

func (layout Layout) state() string {
	return filepath.Join(layout.metadataDirectory(), "agent", "state")
}

func (layout Layout) stateIgnore() string {
	return filepath.Join(layout.state(), ".gitignore")
}

func (layout Layout) initializationLock() string {
	return filepath.Join(layout.state(), "init.lock")
}

func (layout Layout) projectLock() string {
	return filepath.Join(layout.state(), "cluster.lock")
}

func (layout Layout) clusterRecord() string {
	return filepath.Join(layout.state(), "cluster.json")
}

func (layout Layout) kubeconfig() string {
	return filepath.Join(layout.state(), "cluster.kubeconfig")
}

func (layout Layout) runtimes() string {
	return filepath.Join(layout.state(), "runtimes")
}

func (layout Layout) executionRecord(runtime microruntime.Identifier) string {
	return filepath.Join(layout.RuntimeState(runtime), "state.json")
}

// runRecord is the runtime's execution record, run.json.
func (layout Layout) runRecord(runtime microruntime.Identifier) string {
	return filepath.Join(layout.Runtime(runtime), "run.json")
}

// evidence is where the records of one execution land, inside its runtime.
func (layout Layout) evidence(runtime microruntime.Identifier, executionRunID string) string {
	return filepath.Join(layout.Runtime(runtime), "evidence", executionRunID)
}

// hookEvidence is where a hook's record lands, beside the executions' evidence.
func (layout Layout) hookEvidence(runtime microruntime.Identifier, runID string) string {
	return filepath.Join(layout.Runtime(runtime), "evidence", "hooks", runID+".json")
}

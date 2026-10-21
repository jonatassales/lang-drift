// Package storage keeps EX's private records inside the consumer project:
// project identity, the cluster and execution ledgers, the project lock, and
// evidence. Records are versioned JSON documents written atomically.
package storage

import (
	"path/filepath"

	"github.com/exphere/ex/src/domain/microruntime"
	"github.com/exphere/ex/src/domain/project"
)

// Layout names EX's private files inside one project.
type Layout struct{ root string }

// LayoutOf returns the layout of owner's project.
func LayoutOf(owner project.Project) Layout { return Layout{root: owner.Root} }

func layoutAt(root string) Layout { return Layout{root: root} }

// RuntimeState is the private directory of one runtime: its execution ledger
// and its preparation state.
func (layout Layout) RuntimeState(runtime microruntime.Identifier) string {
	return filepath.Join(layout.runtimes(), runtime.String())
}

func (layout Layout) metadataDirectory() string {
	return filepath.Join(layout.root, ".exphere")
}

func (layout Layout) metadata() string {
	return filepath.Join(layout.metadataDirectory(), "config.json")
}

func (layout Layout) state() string {
	return filepath.Join(layout.metadataDirectory(), "state")
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

func (layout Layout) evidence(runID string) string {
	return filepath.Join(layout.root, ".audits", "evidence", "runs", runID)
}

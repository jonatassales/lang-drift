package storage

import (
	"github.com/exphere/ex/src/domain/project"
	"github.com/exphere/ex/src/infrastructure/files"
)

// Locker serializes every cluster and runtime change of one project across
// processes, so deleting the cluster can never race preparing a runtime.
type Locker struct{}

// Lock takes the project lock, or fails at once when another operation holds it.
func (Locker) Lock(owner project.Project) (func(), error) {
	return files.Lock(LayoutOf(owner).projectLock())
}

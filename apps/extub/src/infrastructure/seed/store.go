// Package seed owns bounded, disposable fixture state and embedded seed data.
package seed

import (
	"embed"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/jonatassales/tao/apps/extub/src/domain"
)

//go:embed console.json
var files embed.FS

type Store struct {
	mu      sync.RWMutex
	catalog domain.Catalog
}

func New(now time.Time) (*Store, error) {
	data, err := files.ReadFile("console.json")
	if err != nil {
		return nil, err
	}
	var initial domain.Catalog
	if err = json.Unmarshal(data, &initial); err != nil {
		return nil, fmt.Errorf("decode embedded seeds: %w", err)
	}
	initial.Source = "stub"
	initial.Secrets = []domain.Secret{}
	initial.StartedAt = now.Add(-154 * time.Minute).UTC().Format(time.RFC3339Nano)
	for _, runtime := range initial.Runtimes {
		initial.RecordActivity(runtime.ID, "start", now.Add(-154*time.Minute))
	}
	return &Store{catalog: initial}, nil
}
func (store *Store) Snapshot() domain.Catalog {
	store.mu.RLock()
	defer store.mu.RUnlock()
	return store.catalog.Clone()
}

// Update commits only a successful isolated mutation; returned data cannot alias state.
func (store *Store) Update(change func(*domain.Catalog) error) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	next := store.catalog.Clone()
	if err := change(&next); err != nil {
		return err
	}
	store.catalog = next.Clone()
	return nil
}

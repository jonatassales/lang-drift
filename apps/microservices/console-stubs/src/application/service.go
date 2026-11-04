// Package application owns the seed use cases; it never executes product recipes.
package application

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/jonatassales/tao/apps/microservices/console-stubs/src/domain"
)

type Store interface {
	Snapshot() domain.Catalog
	Update(func(*domain.Catalog) error) error
}
type Service struct {
	Store     Store
	Now       func() time.Time
	LoadTests *LoadTester
}

func (service Service) Snapshot() domain.Catalog { return service.Store.Snapshot() }
func (service Service) Widgets(sequence uint64) []domain.Widget {
	return service.Snapshot().Widgets(service.Now(), sequence)
}
func (service Service) CreateProduct(folder, name string) (domain.Product, error) {
	if domain.ValidateFolder(folder) != nil || len(name) > 100 {
		return domain.Product{}, domain.ErrInvalid
	}
	folder = path.Clean(folder)
	name = strings.TrimSpace(name)
	if name == "" {
		name = path.Base(folder)
	}
	var result domain.Product
	err := service.Store.Update(func(catalog *domain.Catalog) error {
		for _, product := range catalog.Products {
			if product.Source.Path == folder {
				result = product
				return nil
			}
		}
		if len(catalog.Products) >= 16 {
			return domain.ErrLimit
		}
		result = domain.Product{ID: fmt.Sprintf("seed-product-%d", len(catalog.Products)+1), Name: name, Source: domain.Source{Kind: "local", Path: folder}, Components: []domain.Component{}, CreatedAt: service.Now().UTC().Format(time.RFC3339Nano)}
		catalog.Products = append([]domain.Product{result}, catalog.Products...)
		return nil
	})
	return result, err
}
func (service Service) SaveRuntime(id string, input domain.Runtime) (domain.Runtime, error) {
	if err := domain.ValidateRuntime(input); err != nil {
		return domain.Runtime{}, err
	}
	var result domain.Runtime
	err := service.Store.Update(func(catalog *domain.Catalog) error {
		if id != "" {
			index := catalog.RuntimeIndex(id)
			if index < 0 {
				return domain.ErrNotFound
			}
			current := catalog.Runtimes[index]
			if current.DesiredState != "stopped" || current.Revision != input.Revision {
				return domain.ErrConflict
			}
			current.Name = input.Name
			current.Components = input.Components
			current.Configuration = input.Configuration
			current.Revision++
			current.UpdatedAt = service.Now().UTC().Format(time.RFC3339Nano)
			catalog.Runtimes[index] = current
			result = current
			return nil
		}
		if _, err := catalog.Product(input.ProductID); err != nil {
			return err
		}
		if len(catalog.Runtimes) >= 16 {
			return domain.ErrLimit
		}
		input.ID = fmt.Sprintf("seed-runtime-%d", len(catalog.Runtimes)+1)
		input.LoadTest = nil
		input.NativeRuntimeID = input.ID
		input.DesiredState = "stopped"
		input.Revision = 1
		input.IntentVersion = 0
		input.CredentialRevoked = false
		input.Markdown = "# Synthetic microruntime\n\nSeed configuration; no workload is executed."
		input.CreatedAt = service.Now().UTC().Format(time.RFC3339Nano)
		input.UpdatedAt = input.CreatedAt
		catalog.Runtimes = append(catalog.Runtimes, input)
		result = input
		return nil
	})
	return result, err
}
func (service Service) Operate(id, action string) error {
	return service.Store.Update(func(catalog *domain.Catalog) error {
		index := catalog.RuntimeIndex(id)
		if index < 0 {
			return domain.ErrNotFound
		}
		runtime := &catalog.Runtimes[index]
		if runtime.DesiredState == "deleted" {
			return domain.ErrConflict
		}
		switch action {
		case "start":
			if runtime.CredentialRevoked {
				return domain.ErrConflict
			}
			runtime.DesiredState = "running"
			runtime.StartedAt = service.Now().UTC().Format(time.RFC3339Nano)
		case "stop":
			runtime.DesiredState = "stopped"
		case "revoke":
			runtime.CredentialRevoked = true
			runtime.DesiredState = "stopped"
		case "delete":
			runtime.CredentialRevoked = true
			runtime.DesiredState = "deleted"
		default:
			return domain.ErrInvalid
		}
		runtime.IntentVersion++
		runtime.UpdatedAt = service.Now().UTC().Format(time.RFC3339Nano)
		catalog.RecordActivity(id, action, service.Now())
		return nil
	})
}
func (service Service) CreateSecret(product, name string) (domain.Secret, error) {
	name = strings.TrimSpace(name)
	if len(name) > 100 {
		return domain.Secret{}, domain.ErrInvalid
	}
	var result domain.Secret
	err := service.Store.Update(func(catalog *domain.Catalog) error {
		if _, err := catalog.Product(product); err != nil {
			return err
		}
		if len(catalog.Secrets) >= 32 {
			return domain.ErrLimit
		}
		if name == "" {
			name = fmt.Sprintf("Connection %d", len(catalog.Secrets)+1)
		}
		result = domain.Secret{ID: fmt.Sprintf("seed-secret-%d", len(catalog.Secrets)+1), ProductID: product, Name: name, CreatedAt: service.Now().UTC().Format(time.RFC3339Nano)}
		catalog.Secrets = append(catalog.Secrets, result)
		return nil
	})
	return result, err
}
func (service Service) RevokeSecret(product, id string) error {
	return service.Store.Update(func(catalog *domain.Catalog) error {
		for i := range catalog.Secrets {
			secret := &catalog.Secrets[i]
			if secret.ID == id && secret.ProductID == product {
				now := service.Now().UTC().Format(time.RFC3339Nano)
				secret.RevokedAt = &now
				return nil
			}
		}
		return domain.ErrNotFound
	})
}

type Description struct {
	Type string `json:"type"`
	SDP  string `json:"sdp"`
}
type Signaler interface {
	Offer(context.Context, Description) (Description, error)
}
type FrameSink interface {
	Broadcast(context.Context, []byte) error
}

// Stream owns its ticker and exits on cancellation; the caller joins it and closes peers.
func (service Service) Stream(ctx context.Context, sink FrameSink) error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	var sequence uint64
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
		}
		sequence++
		frame, err := json.Marshal(struct {
			Source      string          `json:"source"`
			Kind        string          `json:"kind"`
			Sequence    uint64          `json:"sequence"`
			GeneratedAt string          `json:"generatedAt"`
			Widgets     []domain.Widget `json:"widgets"`
		}{"stub", "runtime-widgets", sequence, service.Now().UTC().Format(time.RFC3339Nano), service.Widgets(sequence)})
		if err != nil {
			return err
		}
		if err = sink.Broadcast(ctx, frame); err != nil {
			return err
		}
	}
}

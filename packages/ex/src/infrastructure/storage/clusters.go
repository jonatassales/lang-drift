package storage

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"time"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/domain/cluster"
	"github.com/exphere/ex/src/domain/project"
	"github.com/exphere/ex/src/infrastructure/files"
)

// clusterRecordVersion identifies the reviewed standalone refactor's record.
// The original format has no record version; inventoryVersion separately
// records which ownership evidence is safe to use for deletion.
const clusterRecordVersion = 1

type clusterFile struct {
	Version          *int                    `json:"version,omitempty"`
	ProjectID        string                  `json:"projectId"`
	Name             string                  `json:"name"`
	Context          string                  `json:"context"`
	Kubeconfig       string                  `json:"kubeconfig"`
	UID              string                  `json:"uid"`
	Status           string                  `json:"status"`
	Namespaces       map[string]string       `json:"namespaceUIDs,omitempty"`
	InventoryVersion *int                    `json:"inventoryVersion,omitempty"`
	Bootstrap        map[string]identityFile `json:"bootstrap,omitempty"`
	ResourceTypes    []string                `json:"resourceTypes,omitempty"`
	UpdatedAt        time.Time               `json:"updatedAt"`
	Error            string                  `json:"error,omitempty"`
}

type identityFile struct {
	UID         string   `json:"uid"`
	Controllers []string `json:"controllers,omitempty"`
}

// Clusters keeps the project's cluster record and its private kubeconfig.
type Clusters struct{}

// Load reads the record, or returns application.ErrNotFound.
func (Clusters) Load(owner project.Project) (cluster.Record, error) {
	layout := LayoutOf(owner)

	var file clusterFile

	err := files.ReadJSON(layout.clusterRecord(), &file)

	switch {
	case errors.Is(err, fs.ErrNotExist):
		return cluster.Record{}, application.NotFound("no cluster is recorded for this project")

	case err != nil:
		return cluster.Record{}, err

	case file.Version != nil && *file.Version != clusterRecordVersion:
		return cluster.Record{}, fmt.Errorf("unsupported cluster record version %d", *file.Version)

	case file.Kubeconfig != layout.kubeconfig():
		return cluster.Record{}, errors.New("the cluster kubeconfig must be the project's private file")
	}
	if err := files.RequireNoSymlinks(layout.kubeconfig()); err != nil {
		return cluster.Record{}, err
	}

	id, err := project.ParseID(file.ProjectID)
	if err != nil {
		return cluster.Record{}, err
	}

	return file.record(id), nil
}

// Save replaces the record atomically.
func (Clusters) Save(owner project.Project, record cluster.Record) error {
	return files.WriteJSON(LayoutOf(owner).clusterRecord(), clusterFileOf(record))
}

// Credentials returns the private kubeconfig location and whether a file exists there.
func (Clusters) Credentials(owner project.Project) (string, bool, error) {
	location := LayoutOf(owner).kubeconfig()

	if err := files.RequireNoSymlinks(location); err != nil {
		return "", false, err
	}

	_, err := os.Lstat(location)

	switch {
	case err == nil:
		return location, true, nil
	case errors.Is(err, fs.ErrNotExist):
		return location, false, nil
	}

	return "", false, err
}

// DiscardCredentials removes the private kubeconfig, if present.
func (Clusters) DiscardCredentials(owner project.Project) error {
	err := os.Remove(LayoutOf(owner).kubeconfig())
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}

	return err
}

func clusterFileOf(record cluster.Record) clusterFile {
	bootstrap := make(map[string]identityFile, len(record.Baseline.Inventory.Objects))
	for key, identity := range record.Baseline.Inventory.Objects {
		bootstrap[key] = identityFile(identity)
	}

	inventoryVersion := record.Baseline.InventoryVersion
	return clusterFile{
		ProjectID:        record.ProjectID.String(),
		Name:             record.Name,
		Context:          record.Context,
		Kubeconfig:       record.Kubeconfig,
		UID:              record.UID,
		Status:           string(record.Status),
		Namespaces:       record.Baseline.Namespaces,
		InventoryVersion: &inventoryVersion,
		Bootstrap:        bootstrap,
		ResourceTypes:    record.Baseline.Inventory.ResourceTypes,
		UpdatedAt:        record.UpdatedAt,
		Error:            record.Error,
	}
}

func (file clusterFile) record(id project.ID) cluster.Record {
	objects := make(map[string]cluster.ObjectIdentity, len(file.Bootstrap))
	for key, identity := range file.Bootstrap {
		objects[key] = cluster.ObjectIdentity(identity)
	}

	inventoryVersion := 0
	if file.InventoryVersion != nil {
		inventoryVersion = *file.InventoryVersion
	} else if file.Version != nil && *file.Version == clusterRecordVersion &&
		len(file.Bootstrap) > 0 && len(file.ResourceTypes) > 0 {
		// This known standalone format stored the same complete inventory
		// under its record version. Partial records still carry no proof.
		inventoryVersion = 1
	}

	return cluster.Record{
		ProjectID:  id,
		Name:       file.Name,
		Context:    file.Context,
		Kubeconfig: file.Kubeconfig,
		UID:        file.UID,
		Status:     cluster.Status(file.Status),
		Baseline: cluster.Baseline{
			Namespaces:       file.Namespaces,
			Inventory:        cluster.Inventory{Objects: objects, ResourceTypes: file.ResourceTypes},
			InventoryVersion: inventoryVersion,
		},
		UpdatedAt: file.UpdatedAt,
		Error:     file.Error,
	}
}

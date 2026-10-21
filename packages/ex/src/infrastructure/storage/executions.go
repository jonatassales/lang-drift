package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"time"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/domain/microruntime"
	"github.com/exphere/ex/src/domain/project"
	"github.com/exphere/ex/src/infrastructure/files"
)

// executionRecordVersion is the only execution record format EX reads.
const executionRecordVersion = 1

type executionFile struct {
	Version      int                        `json:"version"`
	ProjectID    string                     `json:"projectId"`
	Runtime      string                     `json:"id"`
	Status       string                     `json:"status"`
	Settled      bool                       `json:"settled"`
	Token        string                     `json:"token"`
	ClusterUID   string                     `json:"clusterUid"`
	Context      string                     `json:"context"`
	Namespace    string                     `json:"namespace"`
	NamespaceUID string                     `json:"namespaceUid,omitempty"`
	ConfigHash   string                     `json:"configHash"`
	RunID        string                     `json:"runId"`
	Operation    string                     `json:"operation"`
	UpdatedAt    time.Time                  `json:"updatedAt"`
	ToolVersions map[string]toolVersionFile `json:"toolVersions"`
	Images       []containerFile            `json:"images,omitempty"`
	Resources    []resourceFile             `json:"resources"`
	Units        map[string]string          `json:"units"`
	// Inventory is a historical obligation field. It is never produced by
	// this implementation, but must not be silently discarded before cleanup.
	Inventory json.RawMessage `json:"inventory,omitempty"`
}

type toolVersionFile struct {
	Version string `json:"version,omitempty"`
	Status  string `json:"status"`
	// Observed is read only for records written by the reviewed standalone
	// refactor. New writes retain the original execution-version-1 contract.
	Observed *bool `json:"observed,omitempty"`
}

type containerFile struct {
	Unit        string `json:"unit"`
	Workload    string `json:"workload"`
	Pod         string `json:"pod"`
	PodUID      string `json:"podUid"`
	Container   string `json:"container"`
	Image       string `json:"configuredImage"`
	ImageID     string `json:"observedImageId,omitempty"`
	ImageStatus string `json:"imageStatus"`
	Ready       bool   `json:"ready"`
	Restarts    int    `json:"restarts"`
	Phase       string `json:"podPhase"`
}

type resourceFile struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Name       string `json:"name"`
	Namespace  string `json:"namespace"`
	UID        string `json:"uid,omitempty"`
	Unit       string `json:"unit,omitempty"`
}

// Executions keeps each runtime's execution ledger in its private directory.
type Executions struct{}

// Load reads one runtime's ledger, or returns application.ErrNotFound.
func (Executions) Load(owner project.Project, runtime microruntime.Identifier) (*microruntime.Execution, error) {
	var file executionFile

	err := files.ReadJSON(LayoutOf(owner).executionRecord(runtime), &file)

	switch {
	case errors.Is(err, fs.ErrNotExist):
		return nil, application.NotFound("runtime %s has no execution record", runtime)

	case err != nil:
		return nil, err

	case file.Version != executionRecordVersion:
		return nil, fmt.Errorf("unsupported execution record version %d for runtime %s", file.Version, runtime)
	}

	return file.execution()
}

// Save replaces the ledger atomically.
func (Executions) Save(owner project.Project, execution *microruntime.Execution) error {
	return files.WriteJSON(LayoutOf(owner).executionRecord(execution.Runtime), executionFileOf(execution))
}

// List reads every runtime's ledger. Anything unexpected in the private
// runtime directory is an error, never silently skipped.
func (store Executions) List(owner project.Project) ([]*microruntime.Execution, error) {
	entries, err := os.ReadDir(LayoutOf(owner).runtimes())
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	executions := make([]*microruntime.Execution, 0, len(entries))

	for _, entry := range entries {
		runtime, parseErr := microruntime.ParseIdentifier(entry.Name())
		if parseErr != nil || !entry.IsDir() || entry.Type()&fs.ModeSymlink != 0 {
			return nil, fmt.Errorf("unexpected runtime state entry %q", entry.Name())
		}

		execution, err := store.Load(owner, runtime)
		if err != nil {
			return nil, fmt.Errorf("runtime %s has incomplete state: %w", runtime, err)
		}

		if execution.Runtime != runtime {
			return nil, fmt.Errorf("the record in runtime directory %s belongs to runtime %s", runtime, execution.Runtime)
		}

		executions = append(executions, execution)
	}

	return executions, nil
}

func executionFileOf(execution *microruntime.Execution) executionFile {
	tools := make(map[string]toolVersionFile, len(execution.ToolVersions))
	for name, version := range execution.ToolVersions {
		status := "unavailable"
		if version.Observed {
			status = "observed"
		}
		tools[name] = toolVersionFile{Version: version.Version, Status: status}
	}

	images := make([]containerFile, 0, len(execution.Images))
	for _, image := range execution.Images {
		images = append(images, containerFileOf(image))
	}

	var resources []resourceFile
	for _, resource := range execution.Resources {
		resources = append(resources, resourceFile{
			APIVersion: resource.APIVersion,
			Kind:       resource.Kind,
			Name:       resource.Name,
			Namespace:  resource.Namespace,
			UID:        resource.UID,
			Unit:       resource.Unit.String(),
		})
	}

	units := make(map[string]string, len(execution.Units))
	for unit, status := range execution.Units {
		units[unit.String()] = string(status)
	}

	return executionFile{
		Version:      executionRecordVersion,
		ProjectID:    execution.ProjectID.String(),
		Runtime:      execution.Runtime.String(),
		Status:       string(execution.Status),
		Settled:      execution.Settled,
		Token:        execution.Token,
		ClusterUID:   execution.ClusterUID,
		Context:      execution.Context,
		Namespace:    execution.Namespace.String(),
		NamespaceUID: execution.NamespaceUID,
		ConfigHash:   execution.ConfigHash,
		RunID:        execution.RunID,
		Operation:    string(execution.Operation),
		UpdatedAt:    execution.UpdatedAt,
		ToolVersions: tools,
		Images:       images,
		Resources:    resources,
		Units:        units,
	}
}

func (file executionFile) execution() (*microruntime.Execution, error) {
	if err := requireNoLegacyInventory(file.Inventory); err != nil {
		return nil, err
	}
	projectID, err := project.ParseID(file.ProjectID)
	if err != nil {
		return nil, err
	}

	var identifiers identifierParser

	execution := &microruntime.Execution{
		ProjectID:    projectID,
		Runtime:      identifiers.parse("runtime", file.Runtime),
		Status:       microruntime.ExecutionStatus(file.Status),
		Settled:      file.Settled,
		Token:        file.Token,
		ClusterUID:   file.ClusterUID,
		Context:      file.Context,
		Namespace:    identifiers.parse("namespace", file.Namespace),
		NamespaceUID: file.NamespaceUID,
		ConfigHash:   file.ConfigHash,
		RunID:        file.RunID,
		Operation:    microruntime.Operation(file.Operation),
		UpdatedAt:    file.UpdatedAt,
		ToolVersions: make(map[string]microruntime.ToolVersion, len(file.ToolVersions)),
		Units:        make(map[microruntime.Identifier]microruntime.UnitStatus, len(file.Units)),
	}

	for name, version := range file.ToolVersions {
		observation, err := version.observation()
		if err != nil {
			return nil, fmt.Errorf("%w: tool %s: %w", microruntime.ErrInconsistentExecution, name, err)
		}
		execution.ToolVersions[name] = observation
	}

	for _, image := range file.Images {
		observation := image.observation()
		if image.ImageStatus != "" && image.ImageStatus != observation.ImageStatus() {
			return nil, fmt.Errorf("%w: container %s image status conflicts with its observation", microruntime.ErrInconsistentExecution, image.Container)
		}
		execution.Images = append(execution.Images, observation)
	}

	for _, resource := range file.Resources {
		execution.Resources = append(execution.Resources, microruntime.Resource{
			APIVersion: resource.APIVersion,
			Kind:       resource.Kind,
			Name:       resource.Name,
			Namespace:  resource.Namespace,
			UID:        resource.UID,
			Unit:       identifiers.parse("resource unit", resource.Unit),
		})
	}

	for unit, status := range file.Units {
		execution.Units[identifiers.parse("unit", unit)] = microruntime.UnitStatus(status)
	}

	if identifiers.err != nil {
		return nil, fmt.Errorf("%w: %w", microruntime.ErrInconsistentExecution, identifiers.err)
	}

	return execution, nil
}

func requireNoLegacyInventory(data json.RawMessage) error {
	if len(data) == 0 {
		return nil
	}
	var inventory any
	if err := json.Unmarshal(data, &inventory); err != nil {
		return fmt.Errorf("%w: invalid historical inventory", microruntime.ErrInconsistentExecution)
	}
	switch value := inventory.(type) {
	case nil:
		return nil
	case []any:
		if len(value) == 0 {
			return nil
		}
	case map[string]any:
		if len(value) == 0 {
			return nil
		}
	}
	return fmt.Errorf("%w: historical inventory still records ownership obligations; explicit reconciliation is required", microruntime.ErrInconsistentExecution)
}

func (file toolVersionFile) observation() (microruntime.ToolVersion, error) {
	var observed bool
	switch file.Status {
	case "observed":
		observed = true
	case "unavailable":
	case "":
		if file.Observed == nil {
			return microruntime.ToolVersion{}, fmt.Errorf("missing observation status")
		}
		observed = *file.Observed
	default:
		return microruntime.ToolVersion{}, fmt.Errorf("unsupported observation status %q", file.Status)
	}

	if file.Observed != nil && *file.Observed != observed {
		return microruntime.ToolVersion{}, fmt.Errorf("conflicting observation status")
	}
	return microruntime.ToolVersion{Version: file.Version, Observed: observed}, nil
}

func containerFileOf(observation microruntime.ContainerObservation) containerFile {
	return containerFile{
		Unit: observation.Unit, Workload: observation.Workload,
		Pod: observation.Pod, PodUID: observation.PodUID,
		Container: observation.Container, Image: observation.Image,
		ImageID: observation.ImageID, ImageStatus: observation.ImageStatus(),
		Ready: observation.Ready, Restarts: observation.Restarts, Phase: observation.Phase,
	}
}

func (file containerFile) observation() microruntime.ContainerObservation {
	return microruntime.ContainerObservation{
		Unit: file.Unit, Workload: file.Workload,
		Pod: file.Pod, PodUID: file.PodUID,
		Container: file.Container, Image: file.Image,
		ImageID: file.ImageID, Ready: file.Ready, Restarts: file.Restarts, Phase: file.Phase,
	}
}

// identifierParser parses a series of identifiers, keeping the first error.
type identifierParser struct{ err error }

func (parser *identifierParser) parse(field, text string) microruntime.Identifier {
	identifier, err := microruntime.ParseIdentifier(text)
	if err != nil && parser.err == nil {
		parser.err = fmt.Errorf("%s: %w", field, err)
	}

	return identifier
}

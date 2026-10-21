package application

import (
	"time"

	"github.com/exphere/ex/src/domain/microruntime"
)

// ExecutionOutcome is the externally visible part of an execution record:
// everything except the private ownership token. It is also the evidence
// recorded for every lifecycle operation.
type ExecutionOutcome struct {
	Runtime      string                       `json:"id"`
	Status       string                       `json:"status"`
	Settled      bool                         `json:"settled"`
	Operation    string                       `json:"operation"`
	RunID        string                       `json:"runId"`
	ClusterUID   string                       `json:"clusterUid"`
	Context      string                       `json:"context"`
	Namespace    string                       `json:"namespace"`
	NamespaceUID string                       `json:"namespaceUid,omitempty"`
	ConfigHash   string                       `json:"configHash"`
	UpdatedAt    time.Time                    `json:"updatedAt"`
	Resources    []ResourceResult             `json:"resources"`
	Units        map[string]string            `json:"units"`
	ToolVersions map[string]ToolVersionResult `json:"toolVersions"`
	Images       []ContainerResult            `json:"images,omitempty"`
}

// ResourceResult identifies one recorded object.
type ResourceResult struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Name       string `json:"name"`
	Namespace  string `json:"namespace"`
	UID        string `json:"uid,omitempty"`
	Unit       string `json:"unit,omitempty"`
}

// ToolVersionResult reports one tool version, or that it was unavailable.
type ToolVersionResult struct {
	Version string `json:"version,omitempty"`
	Status  string `json:"status"`
}

// ContainerResult reports one container's configured and observed image.
type ContainerResult struct {
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

func outcomeOf(execution *microruntime.Execution) ExecutionOutcome {
	units := make(map[string]string, len(execution.Units))
	for unit, status := range execution.Units {
		units[unit.String()] = string(status)
	}

	return ExecutionOutcome{
		Runtime:      execution.Runtime.String(),
		Status:       string(execution.Status),
		Settled:      execution.Settled,
		Operation:    string(execution.Operation),
		RunID:        execution.RunID,
		ClusterUID:   execution.ClusterUID,
		Context:      execution.Context,
		Namespace:    execution.Namespace.String(),
		NamespaceUID: execution.NamespaceUID,
		ConfigHash:   execution.ConfigHash,
		UpdatedAt:    execution.UpdatedAt,
		Resources:    resourceResultsOf(execution.Resources),
		Units:        units,
		ToolVersions: toolVersionResultsOf(execution.ToolVersions),
		Images:       containerResultsOf(execution.Images),
	}
}

func resourceResultOf(resource microruntime.Resource) ResourceResult {
	return ResourceResult{
		APIVersion: resource.APIVersion,
		Kind:       resource.Kind,
		Name:       resource.Name,
		Namespace:  resource.Namespace,
		UID:        resource.UID,
		Unit:       resource.Unit.String(),
	}
}

func resourceResultsOf(resources []microruntime.Resource) []ResourceResult {
	var results []ResourceResult
	for _, resource := range resources {
		results = append(results, resourceResultOf(resource))
	}

	return results
}

func toolVersionResultsOf(versions map[string]microruntime.ToolVersion) map[string]ToolVersionResult {
	results := make(map[string]ToolVersionResult, len(versions))

	for tool, version := range versions {
		result := ToolVersionResult{Version: version.Version, Status: "unavailable"}
		if version.Observed {
			result = ToolVersionResult{Version: version.Version, Status: "observed"}
		}

		results[tool] = result
	}

	return results
}

func containerResultsOf(observations []microruntime.ContainerObservation) []ContainerResult {
	results := make([]ContainerResult, 0, len(observations))

	for _, observation := range observations {
		results = append(results, ContainerResult{
			Unit:        observation.Unit,
			Workload:    observation.Workload,
			Pod:         observation.Pod,
			PodUID:      observation.PodUID,
			Container:   observation.Container,
			Image:       observation.Image,
			ImageID:     observation.ImageID,
			ImageStatus: observation.ImageStatus(),
			Ready:       observation.Ready,
			Restarts:    observation.Restarts,
			Phase:       observation.Phase,
		})
	}

	return results
}

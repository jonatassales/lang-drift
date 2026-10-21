package application

import (
	"context"

	"github.com/exphere/ex/src/domain/cluster"
)

// Diagnose reports the host's prerequisites and, when a project and its
// cluster can be found, the cluster's verified state. It installs nothing.
type Diagnose struct {
	Prerequisites Prerequisites
	Projects      ProjectOpener
	Clusters      ClusterReader
	Observer      ClusterObserver
}

// DiagnoseRequest selects the project whose cluster is reported, if any.
type DiagnoseRequest struct {
	Project ProjectSelection
}

// Diagnosis is the doctor's report.
type Diagnosis struct {
	Tools           []ToolStatus   `json:"tools"`
	ToolsAvailable  bool           `json:"toolsAvailable"`
	ContainerEngine string         `json:"containerEngine"`
	Cluster         *ClusterResult `json:"cluster,omitempty"`
	ClusterNote     string         `json:"clusterNote,omitempty"`
	Installation    string         `json:"installation"`
}

// Execute never fails because of the host: missing tools and an absent
// cluster are findings, not errors.
func (useCase Diagnose) Execute(ctx context.Context, request DiagnoseRequest) (Diagnosis, error) {
	report := useCase.Prerequisites.Check(ctx)

	diagnosis := Diagnosis{
		Tools:           report.Tools,
		ToolsAvailable:  allAvailable(report.Tools),
		ContainerEngine: "unavailable",
		Installation:    "manual; no tools were installed",
	}

	if report.ContainerEngineReachable {
		diagnosis.ContainerEngine = "reachable"
	}

	record, err := useCase.verifiedCluster(ctx, request.Project)
	if err != nil {
		diagnosis.ClusterNote = "no verified project cluster: " + err.Error()
	} else {
		result := clusterResultOf(record)
		diagnosis.Cluster = &result
	}

	return diagnosis, nil
}

func (useCase Diagnose) verifiedCluster(ctx context.Context, selection ProjectSelection) (cluster.Record, error) {
	owner, err := useCase.Projects.Open(selection)
	if err != nil {
		return cluster.Record{}, err
	}

	return clusterGuard{records: useCase.Clusters, observer: useCase.Observer}.verified(ctx, owner)
}

func allAvailable(tools []ToolStatus) bool {
	for _, tool := range tools {
		if tool.Status != ToolAvailable {
			return false
		}
	}

	return true
}

package application_test

import (
	"context"
	"strings"
	"testing"

	"github.com/exphere/ex/src/application"
)

func (w *world) diagnose() application.Diagnose {
	return application.Diagnose{Prerequisites: w.tools, Projects: w.projects, Clusters: w.clusters, Observer: w.observer}
}

func TestDiagnose_ReportsToolsAndTheVerifiedCluster(t *testing.T) {
	t.Parallel()

	w := newWorld(t)

	diagnosis, err := w.diagnose().Execute(context.Background(), application.DiagnoseRequest{Project: w.clusterRequest().Project})
	if err != nil {
		t.Fatal(err)
	}

	if !diagnosis.ToolsAvailable || diagnosis.ContainerEngine != "reachable" || diagnosis.Cluster == nil ||
		diagnosis.Cluster.Status != "ready" || diagnosis.ClusterNote != "" || !strings.Contains(diagnosis.Installation, "no tools were installed") {
		t.Fatalf("diagnosis = %+v", diagnosis)
	}
}

func TestDiagnose_ReportsFindingsWithoutFailing(t *testing.T) {
	t.Parallel()

	w := newWorld(t)
	w.tools.report = application.PrerequisiteReport{Tools: []application.ToolStatus{
		{Name: "kind", Status: application.ToolAvailable},
		{Name: "terraform", Status: application.ToolMissing, Install: "https://developer.hashicorp.com/terraform/install"},
	}}
	w.projects.err = application.NotFound("no Exphere project at or above %s", projectRoot)

	diagnosis, err := w.diagnose().Execute(context.Background(), application.DiagnoseRequest{Project: w.clusterRequest().Project})
	if err != nil {
		t.Fatalf("Execute() error = %v, want findings instead", err)
	}

	if diagnosis.ToolsAvailable || diagnosis.ContainerEngine != "unavailable" || diagnosis.Cluster != nil ||
		!strings.Contains(diagnosis.ClusterNote, "no Exphere project") {
		t.Fatalf("diagnosis = %+v, want the missing tool, the engine, and the project noted", diagnosis)
	}
}

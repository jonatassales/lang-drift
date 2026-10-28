package domain_test

import (
	"testing"
	"time"

	"github.com/jonatassales/tao/apps/extub/src/domain"
	"github.com/jonatassales/tao/apps/extub/src/infrastructure/seed"
)

func TestDenseSeedTopologyAndStackContract(t *testing.T) {
	now := time.Now()
	store, err := seed.New(now)
	if err != nil {
		t.Fatal(err)
	}
	catalog := store.Snapshot()
	widget := catalog.Widgets(now, 0)[0]
	if len(widget.Runtime.Components) != 14 || len(widget.ComponentLinks) != 14 {
		t.Fatal("first fixture must exercise a dense, branching graph")
	}
	paths := map[string]bool{}
	for _, component := range widget.Runtime.Components {
		paths[component.Path] = true
		if len(component.Stacks) == 0 {
			t.Fatalf("missing authored stack for %s", component.Path)
		}
	}
	neighbors := map[string][]string{}
	for _, link := range widget.ComponentLinks {
		if !paths[link.From] || !paths[link.To] || link.From == link.To {
			t.Fatalf("invalid fixture edge: %+v", link)
		}
		neighbors[link.From] = append(neighbors[link.From], link.To)
		neighbors[link.To] = append(neighbors[link.To], link.From)
	}
	visited := map[string]bool{}
	queue := []string{widget.Runtime.Components[0].Path}
	for len(queue) > 0 {
		path := queue[0]
		queue = queue[1:]
		if visited[path] {
			continue
		}
		visited[path] = true
		queue = append(queue, neighbors[path]...)
	}
	if len(visited) != len(paths) || len(neighbors["gateway"]) < 4 {
		t.Fatal("fixture must connect all services and branch at the gateway")
	}
	if err := domain.ValidateRuntime(widget.Runtime); err != nil {
		t.Fatalf("dense fixture must remain editable: %v", err)
	}
	catalog.Runtimes[0].Components = catalog.Runtimes[0].Components[:2]
	if links := catalog.Widgets(now, 1)[0].ComponentLinks; len(links) != 0 {
		t.Fatalf("excluded endpoints must not appear: %+v", links)
	}
	if len(store.Snapshot().Runtimes[0].Components) != 14 {
		t.Fatal("snapshot mutation leaked into the store")
	}
	for len(widget.Runtime.Components) < 17 {
		widget.Runtime.Components = append(widget.Runtime.Components, widget.Runtime.Components[0])
	}
	if domain.ValidateRuntime(widget.Runtime) == nil {
		t.Fatal("unbounded component count accepted")
	}
}

func TestPerRuntimeActivationKeepsMissingDataHonest(t *testing.T) {
	now := time.Now()
	store, err := seed.New(now)
	if err != nil {
		t.Fatal(err)
	}
	catalog := store.Snapshot()
	widgets := catalog.Widgets(now, 0)
	for i, tone := range []string{"green", "yellow", "orange", "red"} {
		if widgets[i].Tone != tone {
			t.Fatalf("tone %d=%s", i, widgets[i].Tone)
		}
	}
	catalog.Runtimes[0].DesiredState = "stopped"
	widgets = catalog.Widgets(now, 1)
	if len(widgets[0].Executions) != 0 || widgets[0].Tone != "purple" || len(widgets[1].Executions) != 1 {
		t.Fatal("stopping one runtime affected another")
	}
	catalog.Runtimes[0].DesiredState = "running"
	catalog.Runtimes[0].Configuration.LLMBudgetUSD = nil
	widgets = catalog.Widgets(now, 2)
	if widgets[0].Executions[0].LLMCostUSD != nil || !widgets[0].Partial {
		t.Fatal("unbudgeted LLM cost was invented")
	}
}

func TestVisionProductFixtureAndReferentialIntegrity(t *testing.T) {
	store, err := seed.New(time.Now())
	if err != nil {
		t.Fatal(err)
	}
	catalog := store.Snapshot()
	if catalog.Products[0].VisionFixture != nil {
		t.Fatal("Microruntimes must not carry a Vision fixture")
	}
	vision, err := catalog.Product("seed-vision")
	if err != nil || vision.VisionFixture == nil {
		t.Fatalf("missing Vision fixture on seed-vision: %v", err)
	}
	fixture := vision.VisionFixture
	if fixture.FixtureVersion != "dashboard-v1" || !fixture.Illustrative {
		t.Fatalf("unexpected fixture metadata: %+v", fixture)
	}
	if fixture.Overview.ProductVision != 73 || fixture.Report.IntentionalEvolution != 14 || fixture.Report.UnexplainedDrift != 4 {
		t.Fatalf("unexpected overview/report figures: %+v %+v", fixture.Overview, fixture.Report)
	}
	if len(fixture.History.Points) != 5 || len(fixture.Products) != 3 || len(fixture.Decisions) != 2 || len(fixture.People) != 4 || len(fixture.Teams) != 3 {
		t.Fatal("missing dashboard records")
	}
	workspace := fixture.Workspace
	if workspace.Version != "workspace-v1" || len(workspace.Artifacts) != 10 || len(workspace.Drifts) != 3 || len(workspace.Targets) != 3 || len(workspace.Timeline) != 7 {
		t.Fatalf("unexpected workspace counts: %+v", workspace)
	}
	if workspace.PullRequest == nil || workspace.PullRequest.ID != "REC-108" || workspace.PullRequest.Number != 108 || len(workspace.PullRequest.Files) != 2 {
		t.Fatalf("unexpected pull request fixture: %+v", workspace.PullRequest)
	}
	identities := map[string]bool{}
	for _, artifact := range workspace.Artifacts {
		identities[artifact.ID] = true
	}
	for _, drift := range workspace.Drifts {
		identities[drift.ID] = true
	}
	for _, target := range workspace.Targets {
		identities[target.ID] = true
	}
	checkRef := func(id string) {
		if id != "" && !identities[id] {
			t.Fatalf("unresolved workspace reference: %s", id)
		}
	}
	for _, artifact := range workspace.Artifacts {
		for _, ref := range artifact.References {
			checkRef(ref)
		}
	}
	for _, drift := range workspace.Drifts {
		checkRef(drift.TargetID)
		checkRef(drift.DecisionID)
		for _, ref := range drift.ArtifactIDs {
			checkRef(ref)
		}
	}
	for _, target := range workspace.Targets {
		for _, ref := range target.ArtifactIDs {
			checkRef(ref)
		}
	}
	for _, event := range workspace.Timeline {
		checkRef(event.ArtifactID)
		checkRef(event.DriftID)
	}
}

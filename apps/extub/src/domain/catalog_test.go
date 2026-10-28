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

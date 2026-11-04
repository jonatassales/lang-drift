package domain_test

import (
	"strings"
	"testing"
	"time"

	"github.com/jonatassales/tao/apps/microservices/console-stubs/src/domain"
	"github.com/jonatassales/tao/apps/microservices/console-stubs/src/infrastructure/seed"
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

func TestComponentRunFieldsAreBounded(t *testing.T) {
	store, err := seed.New(time.Now())
	if err != nil {
		t.Fatal(err)
	}
	for _, check := range []struct {
		label  string
		change func(*domain.Component)
		valid  bool
	}{
		{"an unset port", func(component *domain.Component) { component.Port = 0 }, true},
		{"port 1", func(component *domain.Component) { component.Port = 1 }, true},
		{"port 65535", func(component *domain.Component) { component.Port = 65535 }, true},
		{"a negative port", func(component *domain.Component) { component.Port = -1 }, false},
		{"port 65536", func(component *domain.Component) { component.Port = 65536 }, false},
		{"no image", func(component *domain.Component) { component.Image = "" }, true},
		{"a 256-byte image", func(component *domain.Component) { component.Image = strings.Repeat("i", 256) }, true},
		{"a 257-byte image", func(component *domain.Component) { component.Image = strings.Repeat("i", 257) }, false},
		{"the root readiness path", func(component *domain.Component) { component.ReadinessPath = "/" }, true},
		{"a 256-byte readiness path", func(component *domain.Component) { component.ReadinessPath = "/" + strings.Repeat("r", 255) }, true},
		{"a 257-byte readiness path", func(component *domain.Component) { component.ReadinessPath = "/" + strings.Repeat("r", 256) }, false},
		{"a relative readiness path", func(component *domain.Component) { component.ReadinessPath = "health" }, false},
		{"no readiness path", func(component *domain.Component) { component.ReadinessPath = "" }, false},
	} {
		runtime := store.Snapshot().Runtimes[1]
		check.change(&runtime.Components[0])
		if valid := domain.ValidateRuntime(runtime) == nil; valid != check.valid {
			t.Fatalf("%s: valid is %v", check.label, valid)
		}
	}
}

func TestFoldersFollowTheProductPathContract(t *testing.T) {
	for folder, valid := range map[string]bool{
		"/":                             true,
		"/seeds/microruntimes":          true,
		"/" + strings.Repeat("a", 1023): true,
		"/" + strings.Repeat("a", 1024): false,
		"":                              false,
		"seeds":                         false,
		"./seeds":                       false,
		"/seeds\x00/microruntimes":      false,
	} {
		if accepted := domain.ValidateFolder(folder) == nil; accepted != valid {
			t.Fatalf("folder %.40q: valid is %v", folder, accepted)
		}
	}
}

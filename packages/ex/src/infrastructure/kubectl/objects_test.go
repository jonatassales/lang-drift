package kubectl_test

import (
	"context"
	"encoding/json"
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/exphere/ex/src/domain/microruntime"
	"github.com/exphere/ex/src/infrastructure/kubectl"
	"github.com/exphere/ex/src/infrastructure/process"
)

func deploymentResource() microruntime.Resource {
	return microruntime.Resource{APIVersion: "apps/v1", Kind: "Deployment", Name: "web", Namespace: "demo-local", UID: "uid-web"}
}

func serviceResource() microruntime.Resource {
	return microruntime.Resource{APIVersion: "v1", Kind: "Service", Name: "web", Namespace: "demo-local", UID: "uid-service"}
}

func TestClient_Get(t *testing.T) {
	t.Parallel()

	t.Run("reads one object of the explicit target", func(t *testing.T) {
		t.Parallel()

		runner := &recorder{respond: answer(`{"kind": "Service", "metadata": {"name": "web", "uid": "uid-web"}}`)}

		object, err := kubectl.Client{Runner: runner}.Get(context.Background(), target, "Service", "web")
		if err != nil || object.UID() != "uid-web" {
			t.Fatalf("Get() = %v, %v", object, err)
		}

		want := append(targetFlags("demo-local"), "get", "Service", "web", "--ignore-not-found=true", "-o", "json")
		if !slices.Equal(runner.calls[0].args, want) {
			t.Errorf("args = %v, want %v", runner.calls[0].args, want)
		}
	})

	t.Run("absent object", func(t *testing.T) {
		t.Parallel()

		object, err := kubectl.Client{Runner: &recorder{respond: answer("\n")}}.Get(context.Background(), target, "Service", "web")
		if object != nil || err != nil {
			t.Fatalf("Get() = %v, %v; want nil without an error", object, err)
		}
	})

	t.Run("undecodable output", func(t *testing.T) {
		t.Parallel()

		if _, err := (kubectl.Client{Runner: &recorder{respond: answer("{")}}).Get(context.Background(), target, "Service", "web"); err == nil {
			t.Fatal("Get() decoded invalid output")
		}
	})
}

func TestClient_List(t *testing.T) {
	t.Parallel()

	runner := &recorder{respond: answer(`{"items": [{"metadata": {"name": "a"}}, {"metadata": {"name": "b"}}]}`)}

	objects, err := kubectl.Client{Runner: runner}.List(context.Background(), target, "pods")
	if err != nil || len(objects) != 2 || objects[1].Name() != "b" {
		t.Fatalf("List() = %v, %v", objects, err)
	}

	if want := append(targetFlags("demo-local"), "get", "pods", "-o", "json"); !slices.Equal(runner.calls[0].args, want) {
		t.Errorf("args = %v, want %v", runner.calls[0].args, want)
	}

	if _, err := (kubectl.Client{Runner: &recorder{respond: answer("[]")}}).List(context.Background(), target, "pods"); err == nil {
		t.Error("List() decoded a document that is not a list")
	}
}

func TestClient_DryRunAndCreate(t *testing.T) {
	t.Parallel()

	object := microruntime.Object{"apiVersion": "v1", "kind": "Service", "metadata": map[string]any{"name": "web"}}
	runner := &recorder{respond: answer(`{"kind": "Service", "metadata": {"name": "web", "uid": "uid-created"}}`)}
	client := kubectl.Client{Runner: runner}

	if err := client.DryRun(context.Background(), target, object); err != nil {
		t.Fatal(err)
	}

	created, err := client.Create(context.Background(), target, object)
	if err != nil || created.UID() != "uid-created" {
		t.Fatalf("Create() = %v, %v", created, err)
	}

	wantDryRun := append(targetFlags("demo-local"), "apply", "--server-side", "--field-manager=exphere", "--dry-run=server", "-f", "-", "-o", "json")
	wantCreate := append(targetFlags("demo-local"), "create", "--field-manager=exphere", "-f", "-", "-o", "json")

	if !slices.Equal(runner.calls[0].args, wantDryRun) || !slices.Equal(runner.calls[1].args, wantCreate) {
		t.Errorf("args = %v and %v, want a server dry run, then a create that never adopts", runner.calls[0].args, runner.calls[1].args)
	}

	for _, recorded := range runner.calls {
		var submitted microruntime.Object
		if err := json.Unmarshal([]byte(recorded.stdin), &submitted); err != nil || submitted.Name() != "web" {
			t.Errorf("stdin = %q, want the object as JSON", recorded.stdin)
		}
	}

	if _, err := (kubectl.Client{Runner: &recorder{respond: answer("created")}}).Create(context.Background(), target, object); err == nil {
		t.Error("Create() accepted output that is not an object")
	}
}

func TestClient_Delete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		resource microruntime.Resource
		path     string
	}{
		{resource: serviceResource(), path: "/api/v1/namespaces/demo-local/services/web"},
		{resource: deploymentResource(), path: "/apis/apps/v1/namespaces/demo-local/deployments/web"},
		{
			resource: microruntime.Resource{APIVersion: "batch/v1", Kind: "Job", Name: "batch", Namespace: "demo-local", UID: "uid-job"},
			path:     "/apis/batch/v1/namespaces/demo-local/jobs/batch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.resource.Kind, func(t *testing.T) {
			t.Parallel()

			runner := &recorder{}
			if err := (kubectl.Client{Runner: runner}).Delete(context.Background(), target, tt.resource); err != nil {
				t.Fatal(err)
			}

			want := append(targetFlags("demo-local"), "delete", "--raw", tt.path, "-f", "-")
			if !slices.Equal(runner.calls[0].args, want) {
				t.Errorf("args = %v, want %v", runner.calls[0].args, want)
			}

			assertDeleteOptions(t, runner.calls[0].stdin, tt.resource.UID)
		})
	}

	runner := &recorder{}
	unsupported := microruntime.Resource{APIVersion: "v1", Kind: "ConfigMap", Name: "settings", Namespace: "demo-local", UID: "uid"}
	if err := (kubectl.Client{Runner: runner}).Delete(context.Background(), target, unsupported); err == nil || len(runner.calls) != 0 {
		t.Fatalf("Delete(ConfigMap) = %v after %d calls, want a refusal without any call", err, len(runner.calls))
	}
}

func assertDeleteOptions(t *testing.T, body, uid string) {
	t.Helper()

	var options struct {
		Kind          string            `json:"kind"`
		Preconditions map[string]string `json:"preconditions"`
		Propagation   string            `json:"propagationPolicy"`
	}

	if err := json.Unmarshal([]byte(body), &options); err != nil {
		t.Fatal(err)
	}

	if options.Kind != "DeleteOptions" || options.Preconditions["uid"] != uid || options.Propagation != "Foreground" {
		t.Errorf("delete options = %+v, want a UID precondition and foreground propagation", options)
	}
}

func TestClient_DeleteNamespace(t *testing.T) {
	t.Parallel()

	runner := &recorder{}
	client := kubectl.Client{Runner: runner}

	if err := client.DeleteNamespace(context.Background(), target, ""); err == nil || len(runner.calls) != 0 {
		t.Fatal("DeleteNamespace() deleted without an observed UID")
	}

	if err := client.DeleteNamespace(context.Background(), target, "uid-namespace"); err != nil {
		t.Fatal(err)
	}

	if want := append(targetFlags("demo-local"), "delete", "--raw", "/api/v1/namespaces/demo-local", "-f", "-"); !slices.Equal(runner.calls[0].args, want) {
		t.Errorf("args = %v, want %v", runner.calls[0].args, want)
	}

	assertDeleteOptions(t, runner.calls[0].stdin, "uid-namespace")
}

func TestClient_NamespaceInventory(t *testing.T) {
	t.Parallel()

	listing := map[string]string{
		"pods":     `{"items": [{"kind": "Pod", "metadata": {"name": "web-1"}}]}`,
		"services": `{"items": [{"kind": "Service", "metadata": {"name": "web"}}, {"kind": "Service", "metadata": {"name": "api"}}]}`,
	}

	runner := &recorder{respond: func(arguments []string) (process.Result, error) {
		switch arguments[0] {
		case "api-resources":
			return process.Result{Stdout: "pods\nservices\nevents\nevents.events.k8s.io\n"}, nil
		case "get":
			document, listed := listing[arguments[1]]
			if !listed {
				return process.Result{}, errors.New("listed " + arguments[1])
			}

			return process.Result{Stdout: document}, nil
		}

		return process.Result{}, errors.New("unexpected call")
	}}

	inventory, err := kubectl.Client{Runner: runner}.NamespaceInventory(context.Background(), target)
	if err != nil || len(inventory) != 3 {
		t.Fatalf("NamespaceInventory() = %v, %v; want every listed object except events", inventory, err)
	}

	if len(runner.calls) != 3 {
		t.Errorf("made %d calls, want discovery and one list per non-event resource", len(runner.calls))
	}

	// List-only views such as pods.metrics.k8s.io hold no state and carry no
	// UID; discovering them would make cleanup refuse a healthy namespace.
	wantDiscovery := append(targetFlags("demo-local"), "api-resources", "--verbs=list,delete", "--namespaced=true", "-o", "name")
	if !slices.Equal(runner.calls[0].args, wantDiscovery) {
		t.Errorf("discovery = %v, want %v", runner.calls[0].args, wantDiscovery)
	}

	failing := &recorder{respond: func(arguments []string) (process.Result, error) {
		if arguments[0] == "api-resources" {
			return process.Result{Stdout: "pods\nwidgets.example.com\n"}, nil
		}

		if arguments[1] == "widgets.example.com" {
			return process.Result{}, errors.New("the server is currently unable to handle the request")
		}

		return process.Result{Stdout: `{"items": []}`}, nil
	}}

	if _, err := (kubectl.Client{Runner: failing}).NamespaceInventory(context.Background(), target); err == nil ||
		!strings.Contains(err.Error(), "incomplete inventory (widgets.example.com)") {
		t.Fatalf("NamespaceInventory() error = %v, want a partial inventory refused", err)
	}

	unavailable := &recorder{respond: func([]string) (process.Result, error) { return process.Result{}, errors.New("unreachable") }}
	if _, err := (kubectl.Client{Runner: unavailable}).NamespaceInventory(context.Background(), target); err == nil {
		t.Fatal("NamespaceInventory() succeeded without API discovery")
	}
}

func TestClient_PodMetrics(t *testing.T) {
	t.Parallel()

	runner := &recorder{respond: answer(`{"items": [{"metadata": {"name": "web-1"}}]}`)}

	samples, err := kubectl.Client{Runner: runner}.PodMetrics(context.Background(), target)
	if err != nil || len(samples) != 1 {
		t.Fatalf("PodMetrics() = %v, %v", samples, err)
	}

	want := append(targetFlags("demo-local"), "get", "--raw", "/apis/metrics.k8s.io/v1beta1/namespaces/demo-local/pods")
	if !slices.Equal(runner.calls[0].args, want) {
		t.Errorf("args = %v, want %v", runner.calls[0].args, want)
	}
}

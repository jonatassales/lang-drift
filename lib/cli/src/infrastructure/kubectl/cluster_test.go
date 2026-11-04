package kubectl_test

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/jonatassales/tao/lib/cli/src/infrastructure/kubectl"
	"github.com/jonatassales/tao/lib/cli/src/infrastructure/process"
)

func TestClient_Namespaces(t *testing.T) {
	t.Parallel()

	runner := &recorder{respond: answer(`{"items": [
		{"metadata": {"name": "kube-system", "uid": "uid-kube-system"}},
		{"metadata": {"name": "default", "uid": "uid-default"}}]}`)}

	namespaces, err := kubectl.Client{Runner: runner}.Namespaces(context.Background(), target.ClusterAccess)
	if err != nil || len(namespaces) != 2 || namespaces["kube-system"] != "uid-kube-system" {
		t.Fatalf("Namespaces() = %v, %v", namespaces, err)
	}

	if want := append(targetFlags(""), "get", "namespaces", "-o", "json"); !slices.Equal(runner.calls[0].args, want) {
		t.Errorf("args = %v, want no namespace flag: %v", runner.calls[0].args, want)
	}

	for _, document := range []string{
		`{"items": [{"metadata": {"name": "kube-system"}}]}`,
		`{"items": [{"metadata": {"uid": "uid-kube-system"}}]}`,
		`not json`,
	} {
		if _, err := (kubectl.Client{Runner: &recorder{respond: answer(document)}}).Namespaces(context.Background(), target.ClusterAccess); err == nil {
			t.Errorf("Namespaces() accepted %s", document)
		}
	}
}

// discovery answers API discovery and go-template listings for a cluster
// described by resource listings, keyed by resource.
type discovery struct {
	namespaced, clusterScoped string
	listings                  map[string]string
}

func (cluster discovery) respond(arguments []string) (process.Result, error) {
	switch {
	case arguments[0] == "api-resources" && arguments[1] == "--namespaced=true":
		return process.Result{Stdout: cluster.namespaced}, nil
	case arguments[0] == "api-resources" && arguments[1] == "--namespaced=false":
		return process.Result{Stdout: cluster.clusterScoped}, nil
	case arguments[0] == "get":
		listing, known := cluster.listings[arguments[1]]
		if !known {
			return process.Result{}, errors.New("cannot list " + arguments[1])
		}

		return process.Result{Stdout: listing}, nil
	}

	return process.Result{}, errors.New("unexpected call")
}

func freshCluster() discovery {
	return discovery{
		namespaced:    "pods\nevents\nconfigmaps\n",
		clusterScoped: "nodes\nnamespaces\n",
		listings: map[string]string{
			"pods":       "kube-system\tcoredns-a\tuid-pod\tuid-rs,\nkube-system\tetcd\tuid-etcd\t\n",
			"configmaps": "kube-system\tcoredns\tuid-cm\t<no value>\n",
			"nodes":      "\tex-control-plane\tuid-node\t\n",
			"namespaces": "\tdefault\tuid-default\t\n",
		},
	}
}

func TestClient_ClusterInventory(t *testing.T) {
	t.Parallel()

	runner := &recorder{respond: freshCluster().respond}

	inventory, err := kubectl.Client{Runner: runner}.ClusterInventory(context.Background(), target.ClusterAccess)
	if err != nil {
		t.Fatal(err)
	}

	if want := []string{"true/configmaps", "true/pods", "false/namespaces", "false/nodes"}; !slices.Equal(inventory.ResourceTypes, want) {
		t.Errorf("ResourceTypes = %v, want %v", inventory.ResourceTypes, want)
	}

	if len(inventory.Objects) != 5 {
		t.Fatalf("Objects = %v, want five identities", inventory.Objects)
	}

	if pod := inventory.Objects["pods/kube-system/coredns-a"]; pod.UID != "uid-pod" || !slices.Equal(pod.Controllers, []string{"uid-rs"}) {
		t.Errorf("pod identity = %+v, want its UID and controller", pod)
	}

	if configMap := inventory.Objects["configmaps/kube-system/coredns"]; len(configMap.Controllers) != 0 {
		t.Errorf("an absent controller was recorded: %+v", configMap)
	}

	if node := inventory.Objects["nodes//ex-control-plane"]; node.UID != "uid-node" {
		t.Errorf("node identity = %+v", node)
	}

	for _, recorded := range runner.calls {
		arguments := operation(recorded.args)

		switch arguments[0] {
		case "api-resources":
			if !slices.Contains(arguments, "--verbs=list,delete") {
				t.Errorf("discovery %v must exclude views that cannot be listed and deleted", arguments)
			}
		case "get":
			namespaced := arguments[1] == "pods" || arguments[1] == "configmaps"
			if slices.Contains(arguments, "--all-namespaces") != namespaced || !strings.HasPrefix(arguments[3], "go-template=") {
				t.Errorf("listing %v does not project identities across the right scope", arguments)
			}
		}

		if slices.Contains(recorded.args, "--namespace") {
			t.Errorf("cluster inventory call %v named a namespace", recorded.args)
		}
	}
}

func TestClient_ClusterInventoryRefusesIncompleteEvidence(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		edit func(*discovery)
	}{
		{name: "object without a UID", edit: func(d *discovery) { d.listings["nodes"] = "\tex-control-plane\t<no value>\t\n" }},
		{name: "object with an empty UID", edit: func(d *discovery) { d.listings["nodes"] = "\tex-control-plane\t\t\n" }},
		{name: "namespaced object without a namespace", edit: func(d *discovery) { d.listings["pods"] = "\tcoredns-a\tuid-pod\t\n" }},
		{name: "malformed projection", edit: func(d *discovery) { d.listings["pods"] = "kube-system\tcoredns-a\tuid-pod\n" }},
		{name: "duplicate identity", edit: func(d *discovery) { d.listings["pods"] = "kube-system\ta\tuid-1\t\nkube-system\ta\tuid-2\t\n" }},
		{name: "resource that cannot be listed", edit: func(d *discovery) { delete(d.listings, "configmaps") }},
		{name: "empty namespaced discovery", edit: func(d *discovery) { d.namespaced = "\n" }},
		{
			name: "no objects at all",
			edit: func(d *discovery) {
				d.namespaced, d.clusterScoped = "pods\n", "nodes\n"
				d.listings = map[string]string{"pods": "", "nodes": "\n"}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cluster := freshCluster()
			tt.edit(&cluster)

			if _, err := (kubectl.Client{Runner: &recorder{respond: cluster.respond}}).ClusterInventory(context.Background(), target.ClusterAccess); err == nil {
				t.Fatal("ClusterInventory() accepted incomplete ownership evidence")
			}
		})
	}
}

// BenchmarkClient_ClusterInventory measures parsing a whole-cluster inventory
// and reports how many kubectl invocations one inventory costs: one discovery
// per scope plus one listing per resource type.
func BenchmarkClient_ClusterInventory(b *testing.B) {
	sizes := []struct {
		name           string
		types, objects int
	}{
		{name: "fresh-kind", types: 40, objects: 8},
		{name: "large", types: 60, objects: 200},
	}

	for _, size := range sizes {
		b.Run(size.name, func(b *testing.B) {
			runner := &recorder{respond: syntheticCluster(size.types, size.objects).respond}
			client := kubectl.Client{Runner: runner}

			b.ReportAllocs()
			b.ResetTimer()

			for range b.N {
				runner.calls = runner.calls[:0]

				if _, err := client.ClusterInventory(context.Background(), target.ClusterAccess); err != nil {
					b.Fatal(err)
				}
			}

			b.ReportMetric(float64(len(runner.calls)), "kubectl-calls/op")
		})
	}
}

// syntheticCluster lists objects of each of types resource types; a quarter
// of the types are cluster-scoped, and every object has a controller.
func syntheticCluster(types, objects int) discovery {
	cluster := discovery{listings: map[string]string{}}

	var namespaced, clusterScoped strings.Builder

	for index := range types {
		resource := fmt.Sprintf("resource%02d", index)
		scoped := index%4 == 0

		var listing strings.Builder

		for object := range objects {
			namespace := ""
			if !scoped {
				namespace = fmt.Sprintf("namespace%d", object%5)
			}

			fmt.Fprintf(&listing, "%s\t%s-%d\tuid-%s-%d\tuid-owner-%d,\n", namespace, resource, object, resource, object, object)
		}

		cluster.listings[resource] = listing.String()

		if scoped {
			fmt.Fprintln(&clusterScoped, resource)
		} else {
			fmt.Fprintln(&namespaced, resource)
		}
	}

	cluster.namespaced, cluster.clusterScoped = namespaced.String(), clusterScoped.String()

	return cluster
}

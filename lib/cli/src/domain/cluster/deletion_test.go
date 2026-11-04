package cluster_test

import (
	"errors"
	"fmt"
	"maps"
	"testing"

	"github.com/jonatassales/tao/lib/cli/src/domain/cluster"
)

// bootstrapInventory is what a fresh cluster holds: a controller chain in
// kube-system and one cluster-scoped object.
func bootstrapInventory() cluster.Inventory {
	return cluster.Inventory{
		Objects: map[string]cluster.ObjectIdentity{
			"true/deployments.apps/kube-system/coredns":     {UID: "uid-deployment"},
			"true/replicasets.apps/kube-system/coredns-6f6": {UID: "uid-replicaset", Controllers: []string{"uid-deployment"}},
			"true/pods/kube-system/coredns-6f6-abcde":       {UID: "uid-pod", Controllers: []string{"uid-replicaset"}},
			"false/nodes//ex-control-plane":                 {UID: "uid-node"},
		},
		ResourceTypes: []string{"true/deployments.apps", "true/replicasets.apps", "true/pods", "false/nodes"},
	}
}

// observation copies the bootstrap so each case can change it independently.
func observation() (map[string]string, cluster.Inventory) {
	bootstrap := bootstrapInventory()

	return map[string]string{"kube-system": "uid-kube-system", "default": "uid-default"},
		cluster.Inventory{Objects: maps.Clone(bootstrap.Objects), ResourceTypes: append([]string(nil), bootstrap.ResourceTypes...)}
}

func TestRecord_AuthorizeDeletion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		edit    func(record *cluster.Record, namespaces map[string]string, observed *cluster.Inventory)
		wantErr bool
	}{
		{
			name: "unchanged bootstrap",
			edit: func(*cluster.Record, map[string]string, *cluster.Inventory) {},
		},
		{
			name: "deleted bootstrap object",
			edit: func(_ *cluster.Record, _ map[string]string, observed *cluster.Inventory) {
				delete(observed.Objects, "true/pods/kube-system/coredns-6f6-abcde")
			},
		},
		{
			name: "pod recreated by a bootstrap controller",
			edit: func(_ *cluster.Record, _ map[string]string, observed *cluster.Inventory) {
				delete(observed.Objects, "true/pods/kube-system/coredns-6f6-abcde")
				observed.Objects["true/pods/kube-system/coredns-6f6-fghij"] = cluster.ObjectIdentity{
					UID: "uid-new-pod", Controllers: []string{"uid-replicaset"},
				}
			},
		},
		{
			name: "rollout below a bootstrap controller, child listed before parent",
			edit: func(_ *cluster.Record, _ map[string]string, observed *cluster.Inventory) {
				observed.Objects["true/pods/kube-system/coredns-7a7-abcde"] = cluster.ObjectIdentity{
					UID: "uid-rollout-pod", Controllers: []string{"uid-rollout-replicaset"},
				}
				observed.Objects["true/replicasets.apps/kube-system/coredns-7a7"] = cluster.ObjectIdentity{
					UID: "uid-rollout-replicaset", Controllers: []string{"uid-deployment"},
				}
			},
		},
		{
			name: "new resource API without objects",
			edit: func(_ *cluster.Record, _ map[string]string, observed *cluster.Inventory) {
				observed.ResourceTypes = append(observed.ResourceTypes, "true/widgets.example.com")
			},
		},
		{
			name: "foreign object",
			edit: func(_ *cluster.Record, _ map[string]string, observed *cluster.Inventory) {
				observed.Objects["true/configmaps/default/consumer"] = cluster.ObjectIdentity{UID: "uid-foreign"}
			},
			wantErr: true,
		},
		{
			name: "object controlled by a foreign object",
			edit: func(_ *cluster.Record, _ map[string]string, observed *cluster.Inventory) {
				observed.Objects["true/pods/default/consumer"] = cluster.ObjectIdentity{
					UID: "uid-foreign-pod", Controllers: []string{"uid-foreign-controller"},
				}
			},
			wantErr: true,
		},
		{
			name: "replaced bootstrap object",
			edit: func(_ *cluster.Record, _ map[string]string, observed *cluster.Inventory) {
				observed.Objects["true/deployments.apps/kube-system/coredns"] = cluster.ObjectIdentity{UID: "uid-replacement"}
			},
			wantErr: true,
		},
		{
			name: "replaced child that a bootstrap controller recreated",
			edit: func(_ *cluster.Record, _ map[string]string, observed *cluster.Inventory) {
				observed.Objects["true/replicasets.apps/kube-system/coredns-6f6"] = cluster.ObjectIdentity{
					UID: "uid-recreated-replicaset", Controllers: []string{"uid-deployment"},
				}
				observed.Objects["true/pods/kube-system/coredns-6f6-abcde"] = cluster.ObjectIdentity{
					UID: "uid-recreated-pod", Controllers: []string{"uid-recreated-replicaset"},
				}
			},
		},
		{
			name: "vanished resource API",
			edit: func(_ *cluster.Record, _ map[string]string, observed *cluster.Inventory) {
				observed.ResourceTypes = observed.ResourceTypes[:len(observed.ResourceTypes)-1]
			},
			wantErr: true,
		},
		{
			name: "foreign namespace",
			edit: func(_ *cluster.Record, namespaces map[string]string, _ *cluster.Inventory) {
				namespaces["consumer"] = "uid-consumer"
			},
			wantErr: true,
		},
		{
			name: "foreign namespace without a UID",
			edit: func(_ *cluster.Record, namespaces map[string]string, _ *cluster.Inventory) {
				namespaces["consumer"] = ""
			},
			wantErr: true,
		},
		{
			name: "baseline namespace recorded without a UID",
			edit: func(record *cluster.Record, _ map[string]string, _ *cluster.Inventory) {
				record.Baseline.Namespaces["consumer"] = ""
			},
			wantErr: true,
		},
		{
			name: "objects without UIDs vouching for each other",
			edit: func(record *cluster.Record, _ map[string]string, observed *cluster.Inventory) {
				record.Baseline.Inventory.Objects["true/configmaps/default/unidentified"] = cluster.ObjectIdentity{}
				observed.Objects["true/configmaps/default/unidentified"] = cluster.ObjectIdentity{}
				observed.Objects["true/configmaps/default/consumer"] = cluster.ObjectIdentity{}
				observed.Objects["true/pods/default/consumer"] = cluster.ObjectIdentity{UID: "uid-foreign-pod", Controllers: []string{""}}
			},
			wantErr: true,
		},
		{
			name: "replaced namespace",
			edit: func(_ *cluster.Record, namespaces map[string]string, _ *cluster.Inventory) {
				namespaces["default"] = "uid-replacement"
			},
			wantErr: true,
		},
		{
			name: "missing baseline namespace",
			edit: func(_ *cluster.Record, namespaces map[string]string, _ *cluster.Inventory) {
				delete(namespaces, "default")
			},
			wantErr: true,
		},
		{
			name: "no namespace baseline",
			edit: func(record *cluster.Record, _ map[string]string, _ *cluster.Inventory) {
				record.Baseline.Namespaces = nil
			},
			wantErr: true,
		},
		{
			name: "record without versioned inventory proof",
			edit: func(record *cluster.Record, _ map[string]string, _ *cluster.Inventory) {
				record.Baseline.InventoryVersion = 0
			},
			wantErr: true,
		},
		{
			name: "empty bootstrap objects",
			edit: func(record *cluster.Record, _ map[string]string, _ *cluster.Inventory) {
				record.Baseline.Inventory.Objects = nil
			},
			wantErr: true,
		},
		{
			name: "empty bootstrap resource types",
			edit: func(record *cluster.Record, _ map[string]string, _ *cluster.Inventory) {
				record.Baseline.Inventory.ResourceTypes = nil
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			record := readyRecord(t)
			namespaces, observed := observation()
			tt.edit(&record, namespaces, &observed)

			err := record.AuthorizeDeletion(namespaces, observed)

			if tt.wantErr != errors.Is(err, cluster.ErrUnsafeDeletion) || (!tt.wantErr && err != nil) {
				t.Fatalf("AuthorizeDeletion() error = %v, want refusal %t", err, tt.wantErr)
			}
		})
	}
}

// BenchmarkRecord_AuthorizeDeletion measures reconciling an observed inventory
// against the bootstrap when every baseline object controls a chain of new
// objects, as a Deployment controls a new ReplicaSet and its Pods.
func BenchmarkRecord_AuthorizeDeletion(b *testing.B) {
	sizes := []struct {
		name         string
		roots, depth int
	}{
		{name: "fresh-kind", roots: 100, depth: 2},
		{name: "wide-2500", roots: 2500, depth: 2},
		{name: "wide-5000", roots: 5000, depth: 2},
		{name: "wide-10000", roots: 10000, depth: 2},
		{name: "roots-1000-depth-5", roots: 1000, depth: 5},
		{name: "chain-250", roots: 1, depth: 250},
		{name: "chain-500", roots: 1, depth: 500},
		{name: "chain-1000", roots: 1, depth: 1000},
		{name: "chain-2000", roots: 1, depth: 2000},
	}

	for _, size := range sizes {
		b.Run(size.name, func(b *testing.B) {
			record, namespaces, observed := reconciliation(size.roots, size.depth)

			b.ReportAllocs()
			b.ResetTimer()

			for range b.N {
				if err := record.AuthorizeDeletion(namespaces, observed); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// reconciliation builds a record whose bootstrap holds roots objects and an
// observation in which each root controls a new chain of depth objects.
func reconciliation(roots, depth int) (*cluster.Record, map[string]string, cluster.Inventory) {
	namespaces := map[string]string{"kube-system": "uid-kube-system"}
	resourceTypes := []string{"true/deployments.apps", "true/pods"}
	bootstrap := cluster.Inventory{Objects: map[string]cluster.ObjectIdentity{}, ResourceTypes: resourceTypes}
	observed := cluster.Inventory{Objects: map[string]cluster.ObjectIdentity{}, ResourceTypes: resourceTypes}

	for root := range roots {
		controller := fmt.Sprintf("uid-root-%d", root)
		key := fmt.Sprintf("true/deployments.apps/kube-system/root-%d", root)
		bootstrap.Objects[key] = cluster.ObjectIdentity{UID: controller}
		observed.Objects[key] = cluster.ObjectIdentity{UID: controller}

		for level := range depth {
			child := fmt.Sprintf("uid-child-%d-%d", root, level)
			observed.Objects[fmt.Sprintf("true/pods/kube-system/child-%d-%d", root, level)] = cluster.ObjectIdentity{
				UID: child, Controllers: []string{controller},
			}
			controller = child
		}
	}

	record := &cluster.Record{
		UID:      "uid-kube-system",
		Status:   cluster.StatusReady,
		Baseline: cluster.Baseline{Namespaces: namespaces, Inventory: bootstrap, InventoryVersion: 1},
	}

	return record, maps.Clone(namespaces), observed
}

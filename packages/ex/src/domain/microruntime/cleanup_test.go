package microruntime_test

import (
	"errors"
	"testing"

	"github.com/exphere/ex/src/domain/microruntime"
)

// namespaceObjects is the complete inventory of a running runtime's namespace,
// keyed by a short name so each case can change one object.
type namespaceObjects map[string]microruntime.Object

func (objects namespaceObjects) list() []microruntime.Object {
	list := make([]microruntime.Object, 0, len(objects))
	for _, object := range objects {
		list = append(list, object)
	}

	return list
}

func runtimeLabels(unit string) map[string]any {
	return map[string]any{
		microruntime.LabelRuntime:        "demo",
		microruntime.LabelUnit:           unit,
		microruntime.LabelOwnershipToken: token,
		"app":                            unit,
	}
}

func ownedBy(object microruntime.Object, uid string) microruntime.Object {
	object["metadata"].(map[string]any)["ownerReferences"] = []any{map[string]any{"uid": uid, "controller": true}}

	return object
}

func podReference(uid string) map[string]any {
	return map[string]any{"targetRef": map[string]any{"kind": "Pod", "name": "web-6f6-a", "namespace": "demo-local", "uid": uid}}
}

// runningNamespace is what runningExecution's namespace holds while it runs.
func runningNamespace() namespaceObjects {
	service := liveObject("v1", "Service", "web", "uid-Service-web", runtimeLabels("web"))
	service["spec"] = map[string]any{"selector": map[string]any{"app": "web"}}

	endpointsLabels := runtimeLabels("web")
	endpointsLabels["endpoints.kubernetes.io/managed-by"] = "endpoint-controller"
	endpoints := liveObject("v1", "Endpoints", "web", "uid-endpoints", endpointsLabels)
	endpoints["subsets"] = []any{map[string]any{"addresses": []any{podReference("uid-pod-web")}}}

	quota := liveObject("v1", "ResourceQuota", "runtime", "uid-quota", map[string]any{microruntime.LabelOwnershipToken: token})

	return namespaceObjects{
		"service account": liveObject("v1", "ServiceAccount", "default", "uid-sa", nil),
		"root CA":         liveObject("v1", "ConfigMap", "kube-root-ca.crt", "uid-ca", nil),
		"quota":           quota,
		"deployment":      liveObject("apps/v1", "Deployment", "web", "uid-Deployment-web", runtimeLabels("web")),
		"replica set":     ownedBy(liveObject("apps/v1", "ReplicaSet", "web-6f6", "uid-rs", runtimeLabels("web")), "uid-Deployment-web"),
		"web pod":         ownedBy(liveObject("v1", "Pod", "web-6f6-a", "uid-pod-web", runtimeLabels("web")), "uid-rs"),
		"service":         service,
		"endpoint slice":  ownedBy(liveObject("discovery.k8s.io/v1", "EndpointSlice", "web-abc", "uid-slice", nil), "uid-Service-web"),
		"endpoints":       endpoints,
		"job":             liveObject("batch/v1", "Job", "batch", "uid-Job-batch", runtimeLabels("batch")),
		"job pod":         ownedBy(liveObject("v1", "Pod", "batch-x", "uid-pod-batch", runtimeLabels("batch")), "uid-Job-batch"),
	}
}

func TestExecution_AuthorizeCleanupBeforeDeletion(t *testing.T) {
	t.Parallel()

	labels := func(object microruntime.Object) map[string]any {
		return object["metadata"].(map[string]any)["labels"].(map[string]any)
	}

	tests := []struct {
		name    string
		edit    func(namespaceObjects)
		wantErr bool
	}{
		{name: "owned objects, their descendants, and the Service's endpoints", edit: func(namespaceObjects) {}},
		{
			name: "endpoints before any pod is ready",
			edit: func(objects namespaceObjects) { delete(objects["endpoints"], "subsets") },
		},
		{
			name: "endpoints with not-ready owned pods",
			edit: func(objects namespaceObjects) {
				objects["endpoints"]["subsets"] = []any{map[string]any{"notReadyAddresses": []any{podReference("uid-pod-web")}}}
			},
		},
		{
			name: "foreign object",
			edit: func(objects namespaceObjects) {
				objects["consumer"] = liveObject("v1", "ConfigMap", "consumer", "uid-consumer", nil)
			},
			wantErr: true,
		},
		{
			name: "quota of another execution",
			edit: func(objects namespaceObjects) {
				objects["quota"] = liveObject("v1", "ResourceQuota", "runtime", "uid-quota", map[string]any{microruntime.LabelOwnershipToken: "other"})
			},
			wantErr: true,
		},
		{
			name: "object without a UID, such as a metrics view",
			edit: func(objects namespaceObjects) {
				objects["metrics"] = microruntime.Object{
					"apiVersion": "metrics.k8s.io/v1beta1", "kind": "PodMetrics",
					"metadata": map[string]any{"name": "web-6f6-a", "namespace": "demo-local", "labels": runtimeLabels("web")},
				}
			},
			wantErr: true,
		},
		{
			name: "endpoints without the endpoint controller's label",
			edit: func(objects namespaceObjects) {
				delete(labels(objects["endpoints"]), "endpoints.kubernetes.io/managed-by")
			},
			wantErr: true,
		},
		{
			name:    "endpoints of another API version",
			edit:    func(objects namespaceObjects) { objects["endpoints"]["apiVersion"] = "v1beta1" },
			wantErr: true,
		},
		{
			name: "endpoints in another namespace",
			edit: func(objects namespaceObjects) {
				objects["endpoints"]["metadata"].(map[string]any)["namespace"] = "other"
			},
			wantErr: true,
		},
		{
			name:    "endpoints owned by a foreign object",
			edit:    func(objects namespaceObjects) { ownedBy(objects["endpoints"], "uid-foreign") },
			wantErr: true,
		},
		{
			name:    "endpoints of an unrecorded Service",
			edit:    func(objects namespaceObjects) { objects["endpoints"]["metadata"].(map[string]any)["name"] = "api" },
			wantErr: true,
		},
		{
			name:    "endpoints of a Service that lost its ownership",
			edit:    func(objects namespaceObjects) { delete(labels(objects["service"]), microruntime.LabelOwnershipToken) },
			wantErr: true,
		},
		{
			name:    "endpoints claimed for another unit",
			edit:    func(objects namespaceObjects) { labels(objects["endpoints"])[microruntime.LabelUnit] = "batch" },
			wantErr: true,
		},
		{
			name:    "Service without a selector",
			edit:    func(objects namespaceObjects) { objects["service"]["spec"] = map[string]any{} },
			wantErr: true,
		},
		{
			name: "Service selector with a non-text value",
			edit: func(objects namespaceObjects) {
				objects["service"]["spec"] = map[string]any{"selector": map[string]any{"app": true}}
			},
			wantErr: true,
		},
		{
			name: "endpoints advertising an unknown pod",
			edit: func(objects namespaceObjects) {
				objects["endpoints"]["subsets"] = []any{map[string]any{"addresses": []any{podReference("uid-unknown")}}}
			},
			wantErr: true,
		},
		{
			name: "endpoints advertising a not-ready unknown pod",
			edit: func(objects namespaceObjects) {
				objects["endpoints"]["subsets"] = []any{map[string]any{
					"addresses":         []any{podReference("uid-pod-web")},
					"notReadyAddresses": []any{podReference("uid-unknown")},
				}}
			},
			wantErr: true,
		},
		{
			name:    "endpoints advertising an owned pod the selector does not match",
			edit:    func(objects namespaceObjects) { labels(objects["web pod"])["app"] = "other" },
			wantErr: true,
		},
		{
			name: "endpoints targeting a node",
			edit: func(objects namespaceObjects) {
				reference := podReference("uid-pod-web")
				reference["targetRef"].(map[string]any)["kind"] = "Node"
				objects["endpoints"]["subsets"] = []any{map[string]any{"addresses": []any{reference}}}
			},
			wantErr: true,
		},
		{
			name:    "malformed endpoint subsets",
			edit:    func(objects namespaceObjects) { objects["endpoints"]["subsets"] = map[string]any{"addresses": []any{}} },
			wantErr: true,
		},
		{
			name:    "endpoint subset that is not an object",
			edit:    func(objects namespaceObjects) { objects["endpoints"]["subsets"] = []any{"uid-pod-web"} },
			wantErr: true,
		},
		{
			name: "malformed endpoint addresses",
			edit: func(objects namespaceObjects) {
				objects["endpoints"]["subsets"] = []any{map[string]any{"addresses": "uid-pod-web"}}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			objects := runningNamespace()
			tt.edit(objects)

			err := runningExecution(t).AuthorizeCleanup(objects.list(), false)

			if tt.wantErr != errors.Is(err, microruntime.ErrCleanupBlocked) || (!tt.wantErr && err != nil) {
				t.Fatalf("AuthorizeCleanup() error = %v, want refusal %t", err, tt.wantErr)
			}
		})
	}
}

func TestExecution_AuthorizeCleanupAfterDeletion(t *testing.T) {
	t.Parallel()

	remaining := func(names ...string) []microruntime.Object {
		objects := runningNamespace()

		var list []microruntime.Object
		for _, name := range names {
			list = append(list, objects[name])
		}

		return list
	}

	tests := []struct {
		name      string
		inventory []microruntime.Object
		wantErr   bool
	}{
		{name: "namespace defaults and the runtime's quota", inventory: remaining("service account", "root CA", "quota")},
		{name: "empty namespace", inventory: nil},
		{name: "remaining recorded object", inventory: remaining("quota", "deployment"), wantErr: true},
		{name: "remaining descendant", inventory: remaining("quota", "job pod"), wantErr: true},
		{name: "remaining endpoints", inventory: remaining("quota", "endpoints"), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := runningExecution(t).AuthorizeCleanup(tt.inventory, true)

			if tt.wantErr != errors.Is(err, microruntime.ErrCleanupBlocked) || (!tt.wantErr && err != nil) {
				t.Fatalf("AuthorizeCleanup() error = %v, want refusal %t", err, tt.wantErr)
			}
		})
	}
}

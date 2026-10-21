package microruntime_test

import (
	"slices"
	"testing"

	"github.com/exphere/ex/src/domain/microruntime"
)

const token = "0123456789abcdef0123456789abcdef"

// liveObject is an object as the API server reports it: JSON-decoded, with a
// UID and the given labels.
func liveObject(apiVersion, kind, name, uid string, labels map[string]any) microruntime.Object {
	metadata := map[string]any{"name": name, "namespace": "demo-local", "uid": uid}
	if labels != nil {
		metadata["labels"] = labels
	}

	return microruntime.Object{"apiVersion": apiVersion, "kind": kind, "metadata": metadata}
}

func TestObject_ReadsMistypedFieldsAsZeroValues(t *testing.T) {
	t.Parallel()

	object := microruntime.Object{
		"kind":     42,
		"metadata": "not an object",
		"spec":     map[string]any{"replicas": "two", "paused": "yes", "ports": "80"},
	}

	if object.Kind() != "" || object.Name() != "" || object.UID() != "" || object.Labels() != nil {
		t.Errorf("mistyped kind or metadata read as %q, %q, %q, %v", object.Kind(), object.Name(), object.UID(), object.Labels())
	}

	spec := object.Spec()
	if spec.Number("replicas") != 0 || spec.Flag("paused") || len(spec.Children("ports")) != 0 {
		t.Errorf("mistyped spec fields were not read as zero values: %v", spec)
	}

	var absent microruntime.Object
	if absent.Kind() != "" || absent.Spec() != nil || absent.PodSpec() != nil || len(absent.Containers()) != 0 {
		t.Error("a nil object did not read as empty")
	}
}

func TestObject_Children(t *testing.T) {
	t.Parallel()

	object := microruntime.Object{"items": []any{map[string]any{"name": "first"}, "not an object", microruntime.Object{"name": "third"}}}

	children := object.Children("items")
	if len(children) != 3 {
		t.Fatalf("Children() returned %d entries, want one per element", len(children))
	}

	if children[0].Text("name") != "first" || children[1] != nil || children[2].Text("name") != "third" {
		t.Fatalf("Children() = %v, want objects with nil for the non-object element", children)
	}
}

func TestObject_Containers(t *testing.T) {
	t.Parallel()

	pod := microruntime.Object{
		"containers":     []any{map[string]any{"name": "app"}},
		"initContainers": []any{map[string]any{"name": "setup"}},
	}

	var names []string
	for _, container := range pod.Containers() {
		names = append(names, container.Text("name"))
	}

	if !slices.Equal(names, []string{"app", "setup"}) {
		t.Fatalf("Containers() = %v, want containers followed by init containers", names)
	}
}

func TestObject_OwnerUIDs(t *testing.T) {
	t.Parallel()

	object := microruntime.Object{"metadata": map[string]any{"ownerReferences": []any{
		map[string]any{"uid": "uid-controller", "controller": true},
		map[string]any{"uid": "uid-owner"},
		map[string]any{"uid": "uid-not-controller", "controller": false},
	}}}

	if got := object.OwnerUIDs(); !slices.Equal(got, []string{"uid-controller", "uid-owner", "uid-not-controller"}) {
		t.Errorf("OwnerUIDs() = %v", got)
	}

	if got := object.ControllerUIDs(); !slices.Equal(got, []string{"uid-controller"}) {
		t.Errorf("ControllerUIDs() = %v", got)
	}
}

func TestObject_IsOwnedBy(t *testing.T) {
	t.Parallel()

	owned := liveObject("v1", "Service", "web", "uid-web", map[string]any{microruntime.LabelOwnershipToken: token})

	tests := []struct {
		name   string
		object microruntime.Object
		uid    string
		token  string
		want   bool
	}{
		{name: "recorded object with the token", object: owned, uid: "uid-web", token: token, want: true},
		{name: "replacement under the same name", object: owned, uid: "uid-earlier", token: token},
		{name: "another execution's token", object: owned, uid: "uid-web", token: "another-token"},
		{name: "no recorded UID", object: owned, uid: "", token: token},
		{name: "unlabelled object", object: liveObject("v1", "Service", "web", "uid-web", nil), uid: "uid-web", token: token},
		{name: "absent object", object: nil, uid: "uid-web", token: token},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.object.IsOwnedBy(tt.uid, tt.token); got != tt.want {
				t.Fatalf("IsOwnedBy(%q, %q) = %t, want %t", tt.uid, tt.token, got, tt.want)
			}
		})
	}
}

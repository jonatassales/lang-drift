package terraform

import (
	"encoding/json"
	"testing"

	"github.com/exphere/ex/src/domain/microruntime"
)

const testToken = "0123456789abcdef0123456789abcdef"

type change struct {
	kind    string
	mode    string
	before  map[string]any
	after   map[string]any
	actions []string
}

func metadata(name, namespace, token string) map[string]any {
	block := map[string]any{"labels": map[string]any{microruntime.LabelOwnershipToken: token}}
	if name != "" {
		block["name"] = name
	}

	if namespace != "" {
		block["namespace"] = namespace
	}

	return map[string]any{"metadata": []any{block}}
}

func plan(t *testing.T, changes ...change) []byte {
	t.Helper()

	var resources []any
	for _, c := range changes {
		mode := c.mode
		if mode == "" {
			mode = "managed"
		}

		resources = append(resources, map[string]any{
			"type": c.kind, "mode": mode,
			"change": map[string]any{"before": c.before, "after": c.after, "actions": c.actions},
		})
	}

	data, err := json.Marshal(map[string]any{"resource_changes": resources})
	if err != nil {
		t.Fatal(err)
	}

	return data
}

func TestValidatePlan(t *testing.T) {
	t.Parallel()

	namespace := metadata("demo-local", "", testToken)
	quota := metadata("runtime", "demo-local", testToken)

	tests := []struct {
		name    string
		plan    func(*testing.T) []byte
		wantErr bool
	}{
		{name: "no changes", plan: func(t *testing.T) []byte { return plan(t) }},
		{
			name: "create namespace and quota",
			plan: func(t *testing.T) []byte {
				return plan(t,
					change{kind: NamespaceResource, after: namespace, actions: []string{"create"}},
					change{kind: QuotaResource, after: quota, actions: []string{"create"}})
			},
		},
		{
			name: "update quota in place",
			plan: func(t *testing.T) []byte {
				return plan(t, change{kind: QuotaResource, before: quota, after: quota, actions: []string{"update"}})
			},
		},
		{
			name: "destroy namespace and quota",
			plan: func(t *testing.T) []byte {
				return plan(t,
					change{kind: QuotaResource, before: quota, actions: []string{"delete"}},
					change{kind: NamespaceResource, before: namespace, actions: []string{"delete"}})
			},
		},
		{
			name: "replace namespace",
			plan: func(t *testing.T) []byte {
				return plan(t, change{kind: NamespaceResource, before: namespace, after: namespace, actions: []string{"delete", "create"}})
			},
			wantErr: true,
		},
		{
			name: "data source",
			plan: func(t *testing.T) []byte {
				return plan(t, change{kind: NamespaceResource, mode: "data", after: namespace})
			},
			wantErr: true,
		},
		{
			name:    "other resource type",
			plan:    func(t *testing.T) []byte { return plan(t, change{kind: "kubernetes_secret_v1", after: namespace}) },
			wantErr: true,
		},
		{
			name: "namespace without the ownership token",
			plan: func(t *testing.T) []byte {
				return plan(t, change{kind: NamespaceResource, after: metadata("demo-local", "", "other")})
			},
			wantErr: true,
		},
		{
			name: "another namespace",
			plan: func(t *testing.T) []byte {
				return plan(t, change{kind: NamespaceResource, after: metadata("default", "", testToken)})
			},
			wantErr: true,
		},
		{
			name: "quota in another namespace",
			plan: func(t *testing.T) []byte {
				return plan(t, change{kind: QuotaResource, after: metadata("runtime", "default", testToken)})
			},
			wantErr: true,
		},
		{
			name:    "unknown metadata",
			plan:    func(t *testing.T) []byte { return plan(t, change{kind: NamespaceResource, after: map[string]any{}}) },
			wantErr: true,
		},
		{
			name: "malformed metadata",
			plan: func(t *testing.T) []byte {
				return plan(t, change{kind: NamespaceResource, after: map[string]any{"metadata": []any{"demo-local"}}})
			},
			wantErr: true,
		},
		{name: "undecodable plan", plan: func(*testing.T) []byte { return []byte("{") }, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := validatePlan(tt.plan(t), "demo-local", testToken)

			if tt.wantErr != (err != nil) {
				t.Fatalf("validatePlan() error = %v, want refusal %t", err, tt.wantErr)
			}
		})
	}
}

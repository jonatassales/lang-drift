package microruntime_test

import (
	"errors"
	"slices"
	"testing"

	"github.com/exphere/ex/src/domain/microruntime"
)

// The authored objects below are shaped as a YAML decoder returns them:
// nested map[string]any values and integer numbers.

func authoredDeployment() microruntime.Object {
	return microruntime.Object{
		"apiVersion": "apps/v1",
		"kind":       "Deployment",
		"metadata":   map[string]any{"name": "web", "labels": map[string]any{"app": "web"}},
		"spec": map[string]any{
			"replicas": 2,
			"selector": map[string]any{"matchLabels": map[string]any{"app": "web"}},
			"template": map[string]any{
				"metadata": map[string]any{"labels": map[string]any{"app": "web"}},
				"spec": map[string]any{
					"containers": []any{map[string]any{"name": "web", "image": "web:dev", "imagePullPolicy": "Never"}},
				},
			},
		},
	}
}

func authoredService() microruntime.Object {
	return microruntime.Object{
		"apiVersion": "v1",
		"kind":       "Service",
		"metadata":   map[string]any{"name": "web"},
		"spec": map[string]any{
			"selector": map[string]any{"app": "web"},
			"ports":    []any{map[string]any{"port": 8080, "targetPort": "http"}},
		},
	}
}

func authoredJob() microruntime.Object {
	return microruntime.Object{
		"apiVersion": "batch/v1",
		"kind":       "Job",
		"metadata":   map[string]any{"name": "batch"},
		"spec": map[string]any{
			"activeDeadlineSeconds": 120,
			"template": map[string]any{
				"spec": map[string]any{
					"restartPolicy": "Never",
					"containers": []any{map[string]any{
						"name": "batch", "image": "batch:dev", "imagePullPolicy": "Never",
						"env": []any{map[string]any{"name": "MODE", "value": "full"}},
					}},
				},
			},
		},
	}
}

// validManifests are the authored objects of validDescriptor's units.
func validManifests(t testing.TB) []microruntime.Manifest {
	t.Helper()

	web, batch := identifier(t, "web"), identifier(t, "batch")

	return []microruntime.Manifest{
		{Unit: web, Object: authoredDeployment()},
		{Unit: web, Object: authoredService()},
		{Unit: batch, Object: authoredJob()},
	}
}

func TestManifest_Key(t *testing.T) {
	t.Parallel()

	manifest := microruntime.Manifest{Object: authoredService()}
	if got := manifest.Key(); got != "Service/web" {
		t.Fatalf("Key() = %q, want Service/web", got)
	}
}

func TestDescriptor_ValidateManifests(t *testing.T) {
	t.Parallel()

	stranger := identifier(t, "stranger")

	tests := []struct {
		name    string
		edit    func([]microruntime.Manifest) []microruntime.Manifest
		wantErr bool
	}{
		{name: "every unit contributes supported objects", edit: func(m []microruntime.Manifest) []microruntime.Manifest { return m }},
		{
			name: "object of an undeclared unit",
			edit: func(m []microruntime.Manifest) []microruntime.Manifest {
				return append(m, microruntime.Manifest{Unit: stranger, Object: authoredService()})
			},
			wantErr: true,
		},
		{
			name: "object declared twice across units",
			edit: func(m []microruntime.Manifest) []microruntime.Manifest {
				return append(m, microruntime.Manifest{Unit: m[2].Unit, Object: authoredService()})
			},
			wantErr: true,
		},
		{
			name:    "unit without objects",
			edit:    func(m []microruntime.Manifest) []microruntime.Manifest { return m[:2] },
			wantErr: true,
		},
		{
			name: "object violating the policy",
			edit: func(m []microruntime.Manifest) []microruntime.Manifest {
				m[0].Object["kind"] = "StatefulSet"
				return m
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := validDescriptor(t).ValidateManifests(tt.edit(validManifests(t)))

			if tt.wantErr != errors.Is(err, microruntime.ErrUnsupportedManifest) || (!tt.wantErr && err != nil) {
				t.Fatalf("ValidateManifests() error = %v, want refusal %t", err, tt.wantErr)
			}
		})
	}
}

// TestManifestPolicy checks each rule of the stateless pilot policy on one
// object, which is otherwise valid.
func TestManifestPolicy(t *testing.T) {
	t.Parallel()

	podSpec := func(object microruntime.Object) map[string]any {
		return object["spec"].(map[string]any)["template"].(map[string]any)["spec"].(map[string]any)
	}
	container := func(object microruntime.Object) map[string]any {
		return podSpec(object)["containers"].([]any)[0].(map[string]any)
	}
	metadata := func(object microruntime.Object) map[string]any {
		return object["metadata"].(map[string]any)
	}
	spec := func(object microruntime.Object) map[string]any {
		return object["spec"].(map[string]any)
	}

	tests := []struct {
		name    string
		object  func() microruntime.Object
		edit    func(microruntime.Object)
		wantErr bool
	}{
		{name: "namespace equal to the target", object: authoredDeployment, edit: func(o microruntime.Object) { metadata(o)["namespace"] = "demo-local" }},
		{name: "headless ClusterIP Service", object: authoredService, edit: func(o microruntime.Object) { spec(o)["clusterIP"] = "None" }},
		{name: "explicit ClusterIP Service", object: authoredService, edit: func(o microruntime.Object) { spec(o)["type"] = "ClusterIP" }},
		{
			name:   "stateless volumes",
			object: authoredDeployment,
			edit: func(o microruntime.Object) {
				podSpec(o)["volumes"] = []any{
					map[string]any{"name": "scratch", "emptyDir": map[string]any{}},
					map[string]any{"name": "settings", "configMap": map[string]any{"name": "settings"}},
					map[string]any{"name": "credentials", "secret": map[string]any{"secretName": "credentials"}},
				}
			},
		},
		{name: "Job deadline decoded from JSON", object: authoredJob, edit: func(o microruntime.Object) { spec(o)["activeDeadlineSeconds"] = 120.0 }},
		{name: "Job deadline as int64", object: authoredJob, edit: func(o microruntime.Object) { spec(o)["activeDeadlineSeconds"] = int64(120) }},
		{name: "Job deadline as uint64", object: authoredJob, edit: func(o microruntime.Object) { spec(o)["activeDeadlineSeconds"] = uint64(120) }},
		{name: "invalid name", object: authoredDeployment, edit: func(o microruntime.Object) { metadata(o)["name"] = "Web_API" }, wantErr: true},
		{name: "missing name", object: authoredService, edit: func(o microruntime.Object) { delete(metadata(o), "name") }, wantErr: true},
		{name: "foreign namespace", object: authoredDeployment, edit: func(o microruntime.Object) { metadata(o)["namespace"] = "default" }, wantErr: true},
		{name: "generated name", object: authoredDeployment, edit: func(o microruntime.Object) { metadata(o)["generateName"] = "web-" }, wantErr: true},
		{name: "authored UID", object: authoredDeployment, edit: func(o microruntime.Object) { metadata(o)["uid"] = "uid-web" }, wantErr: true},
		{name: "unsupported kind", object: authoredService, edit: func(o microruntime.Object) { o["kind"] = "ConfigMap" }, wantErr: true},
		{name: "unsupported API version", object: authoredDeployment, edit: func(o microruntime.Object) { o["apiVersion"] = "apps/v1beta1" }, wantErr: true},
		{name: "NodePort Service", object: authoredService, edit: func(o microruntime.Object) { spec(o)["type"] = "NodePort" }, wantErr: true},
		{name: "LoadBalancer Service", object: authoredService, edit: func(o microruntime.Object) { spec(o)["type"] = "LoadBalancer" }, wantErr: true},
		{name: "external IPs", object: authoredService, edit: func(o microruntime.Object) { spec(o)["externalIPs"] = []any{"192.0.2.1"} }, wantErr: true},
		{name: "external name", object: authoredService, edit: func(o microruntime.Object) { spec(o)["externalName"] = "example.com" }, wantErr: true},
		{name: "Job without deadline", object: authoredJob, edit: func(o microruntime.Object) { delete(spec(o), "activeDeadlineSeconds") }, wantErr: true},
		{name: "Job with zero deadline", object: authoredJob, edit: func(o microruntime.Object) { spec(o)["activeDeadlineSeconds"] = 0 }, wantErr: true},
		{name: "Job with fractional deadline", object: authoredJob, edit: func(o microruntime.Object) { spec(o)["activeDeadlineSeconds"] = 1.5 }, wantErr: true},
		{name: "Job with textual deadline", object: authoredJob, edit: func(o microruntime.Object) { spec(o)["activeDeadlineSeconds"] = "120" }, wantErr: true},
		{name: "suspended Job", object: authoredJob, edit: func(o microruntime.Object) { spec(o)["suspend"] = true }, wantErr: true},
		{name: "host network", object: authoredDeployment, edit: func(o microruntime.Object) { podSpec(o)["hostNetwork"] = true }, wantErr: true},
		{name: "host PID", object: authoredJob, edit: func(o microruntime.Object) { podSpec(o)["hostPID"] = true }, wantErr: true},
		{name: "host IPC", object: authoredDeployment, edit: func(o microruntime.Object) { podSpec(o)["hostIPC"] = true }, wantErr: true},
		{
			name:   "host path volume",
			object: authoredDeployment,
			edit: func(o microruntime.Object) {
				podSpec(o)["volumes"] = []any{map[string]any{"name": "host", "hostPath": map[string]any{"path": "/"}}}
			},
			wantErr: true,
		},
		{
			name:   "persistent volume claim",
			object: authoredDeployment,
			edit: func(o microruntime.Object) {
				podSpec(o)["volumes"] = []any{map[string]any{"name": "data", "persistentVolumeClaim": map[string]any{"claimName": "data"}}}
			},
			wantErr: true,
		},
		{
			name:   "volume with two sources",
			object: authoredDeployment,
			edit: func(o microruntime.Object) {
				podSpec(o)["volumes"] = []any{map[string]any{"name": "mixed", "emptyDir": map[string]any{}, "secret": map[string]any{}}}
			},
			wantErr: true,
		},
		{
			name:    "volume without a source",
			object:  authoredDeployment,
			edit:    func(o microruntime.Object) { podSpec(o)["volumes"] = []any{map[string]any{"name": "empty"}} },
			wantErr: true,
		},
		{
			name:   "volume source that is not an object",
			object: authoredDeployment,
			edit: func(o microruntime.Object) {
				podSpec(o)["volumes"] = []any{map[string]any{"name": "scratch", "emptyDir": "yes"}}
			},
			wantErr: true,
		},
		{
			name:   "volume without a name",
			object: authoredDeployment,
			edit: func(o microruntime.Object) {
				podSpec(o)["volumes"] = []any{map[string]any{"emptyDir": map[string]any{}}}
			},
			wantErr: true,
		},
		{
			name:    "volumes that are not an array",
			object:  authoredDeployment,
			edit:    func(o microruntime.Object) { podSpec(o)["volumes"] = map[string]any{"name": "scratch"} },
			wantErr: true,
		},
		{name: "no containers", object: authoredDeployment, edit: func(o microruntime.Object) { delete(podSpec(o), "containers") }, wantErr: true},
		{name: "blank image", object: authoredJob, edit: func(o microruntime.Object) { container(o)["image"] = " " }, wantErr: true},
		{
			name:    "privileged container",
			object:  authoredDeployment,
			edit:    func(o microruntime.Object) { container(o)["securityContext"] = map[string]any{"privileged": true} },
			wantErr: true,
		},
		{name: "registry pull policy", object: authoredDeployment, edit: func(o microruntime.Object) { container(o)["imagePullPolicy"] = "IfNotPresent" }, wantErr: true},
		{name: "default pull policy", object: authoredJob, edit: func(o microruntime.Object) { delete(container(o), "imagePullPolicy") }, wantErr: true},
		{
			name:   "init container pulling from a registry",
			object: authoredDeployment,
			edit: func(o microruntime.Object) {
				podSpec(o)["initContainers"] = []any{map[string]any{"name": "setup", "image": "setup:dev"}}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			descriptor := validDescriptor(t)
			manifests := validManifests(t)

			for index := range manifests {
				if manifests[index].Object.Kind() == tt.object().Kind() {
					tt.edit(manifests[index].Object)
				}
			}

			err := descriptor.ValidateManifests(manifests)

			if tt.wantErr != errors.Is(err, microruntime.ErrUnsupportedManifest) || (!tt.wantErr && err != nil) {
				t.Fatalf("ValidateManifests() error = %v, want refusal %t", err, tt.wantErr)
			}
		})
	}
}

func TestClaim_Apply(t *testing.T) {
	t.Parallel()

	claim := microruntime.Claim{
		Runtime:   identifier(t, "demo"),
		Namespace: identifier(t, "demo-local"),
		Token:     token,
		RunID:     "run-1",
	}
	web, batch := identifier(t, "web"), identifier(t, "batch")

	wantLabels := func(t *testing.T, labels microruntime.Object, unit string) {
		t.Helper()

		if labels.Text(microruntime.LabelRuntime) != "demo" ||
			labels.Text(microruntime.LabelUnit) != unit ||
			labels.Text(microruntime.LabelOwnershipToken) != token {
			t.Fatalf("labels = %v, want the runtime, unit %s, and ownership token", labels, unit)
		}
	}

	t.Run("Service carries the claim on its metadata", func(t *testing.T) {
		t.Parallel()

		service := authoredService()
		claim.Apply(microruntime.Manifest{Unit: web, Object: service})

		if service.Namespace() != "demo-local" {
			t.Errorf("namespace = %q, want demo-local", service.Namespace())
		}

		wantLabels(t, service.Labels(), "web")

		if service.Spec().Child("template") != nil {
			t.Error("a Service gained a pod template")
		}
	})

	t.Run("Deployment keeps authored labels and labels its pod template", func(t *testing.T) {
		t.Parallel()

		deployment := authoredDeployment()
		claim.Apply(microruntime.Manifest{Unit: web, Object: deployment})

		wantLabels(t, deployment.Labels(), "web")
		wantLabels(t, deployment.Spec().Child("template").Child("metadata").Child("labels"), "web")

		if deployment.Labels().Text("app") != "web" {
			t.Error("an authored label was dropped")
		}
	})

	t.Run("Job containers receive fresh correlation variables", func(t *testing.T) {
		t.Parallel()

		job := authoredJob()
		container := job.PodSpec().Children("containers")[0]
		container["env"] = append(container["env"].([]any),
			map[string]any{"name": microruntime.EnvironmentRunID, "value": "authored"},
		)

		claim.Apply(microruntime.Manifest{Unit: batch, Object: job})

		environment := map[string]string{}
		for _, variable := range job.PodSpec().Children("containers")[0].Children("env") {
			if _, duplicate := environment[variable.Text("name")]; duplicate {
				t.Fatalf("variable %s is declared twice", variable.Text("name"))
			}

			environment[variable.Text("name")] = variable.Text("value")
		}

		want := map[string]string{"MODE": "full", microruntime.EnvironmentRuntimeID: "demo", microruntime.EnvironmentRunID: "run-1"}
		if len(environment) != len(want) {
			t.Fatalf("env = %v, want %v", environment, want)
		}

		for name, value := range want {
			if environment[name] != value {
				t.Errorf("env %s = %q, want %q", name, environment[name], value)
			}
		}

		wantLabels(t, job.Spec().Child("template").Child("metadata").Child("labels"), "batch")
	})

	t.Run("malformed Job container is left for the policy to reject", func(t *testing.T) {
		t.Parallel()

		job := authoredJob()
		job.PodSpec()["containers"] = []any{"not an object", job.PodSpec().Children("containers")[0]}

		claim.Apply(microruntime.Manifest{Unit: batch, Object: job})

		containers := job.PodSpec()["containers"].([]any)
		if containers[0] != "not an object" || len(job.PodSpec().Children("containers")[1].Children("env")) != 3 {
			t.Fatalf("containers = %v, want the malformed entry unchanged and the valid one correlated", containers)
		}
	})

	t.Run("object without metadata gains it", func(t *testing.T) {
		t.Parallel()

		object := microruntime.Object{"apiVersion": "v1", "kind": "Service"}
		claim.Apply(microruntime.Manifest{Unit: web, Object: object})

		if object.Namespace() != "demo-local" {
			t.Fatalf("namespace = %q, want demo-local", object.Namespace())
		}

		wantLabels(t, object.Labels(), "web")
	})
}

func TestImages(t *testing.T) {
	t.Parallel()

	manifests := validManifests(t)
	second := authoredDeployment()
	second["metadata"].(map[string]any)["name"] = "admin"
	manifests = append(manifests, microruntime.Manifest{Unit: manifests[0].Unit, Object: second})

	if got := microruntime.Images(manifests); !slices.Equal(got, []string{"batch:dev", "web:dev"}) {
		t.Fatalf("Images() = %v, want distinct images in sorted order", got)
	}
}

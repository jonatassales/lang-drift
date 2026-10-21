package microruntime_test

import (
	"testing"

	"github.com/exphere/ex/src/domain/microruntime"
)

func TestWorkload_NativeObjects(t *testing.T) {
	t.Parallel()

	runtime, namespace := identifier(t, "demo"), identifier(t, "demo-local")

	t.Run("generated objects satisfy the manifest policy", func(t *testing.T) {
		t.Parallel()

		blueprint := validBlueprint(t)
		descriptor := blueprint.Descriptor(runtime)

		var manifests []microruntime.Manifest
		for _, workload := range blueprint.Workloads {
			for _, object := range workload.NativeObjects(runtime, namespace) {
				manifests = append(manifests, microruntime.Manifest{Unit: workload.Name, Object: object})
			}
		}

		if err := descriptor.ValidateManifests(manifests); err != nil {
			t.Fatalf("generated manifests violate the policy: %v", err)
		}
	})

	t.Run("Deployment with its ClusterIP Service", func(t *testing.T) {
		t.Parallel()

		objects := deploymentWorkload(t).NativeObjects(runtime, namespace)
		if len(objects) != 2 || objects[0].Kind() != "Deployment" || objects[1].Kind() != "Service" {
			t.Fatalf("NativeObjects() = %v, want a Deployment and a Service", objects)
		}

		deployment, service := objects[0], objects[1]

		for _, object := range objects {
			if object.Name() != "web" || object.Namespace() != "demo-local" ||
				object.Labels().Text(microruntime.LabelRuntime) != "demo" || object.Labels().Text(microruntime.LabelUnit) != "web" {
				t.Errorf("%s metadata = %v, want name, namespace, runtime, and unit", object.Kind(), object.Metadata())
			}
		}

		if deployment.Spec()["replicas"] != 2 {
			t.Errorf("replicas = %v, want 2", deployment.Spec()["replicas"])
		}

		selector := deployment.Spec().Child("selector").Child("matchLabels")
		podLabels := deployment.Spec().Child("template").Child("metadata").Child("labels")
		for key, value := range selector {
			if podLabels[key] != value {
				t.Errorf("selector %s=%v does not match the pod template labels %v", key, value, podLabels)
			}
		}

		container := deployment.PodSpec().Children("containers")[0]
		if container.Text("image") != "web:dev" || container.Text("imagePullPolicy") != "Never" ||
			container.Child("readinessProbe").Child("httpGet").Text("path") != "/healthz" {
			t.Errorf("container = %v, want the prebuilt image, pull policy Never, and the readiness path", container)
		}

		if deployment.PodSpec().Flag("automountServiceAccountToken") {
			t.Error("the pod mounts a service account token")
		}

		port := service.Spec().Children("ports")[0]
		if service.Spec().Text("type") != "ClusterIP" || port["port"] != 8080 || port.Text("targetPort") != "http" {
			t.Errorf("service spec = %v, want a ClusterIP Service on 8080 targeting the http port", service.Spec())
		}
	})

	t.Run("finite Job without retries", func(t *testing.T) {
		t.Parallel()

		objects := jobWorkload(t).NativeObjects(runtime, namespace)
		if len(objects) != 1 || objects[0].Kind() != "Job" {
			t.Fatalf("NativeObjects() = %v, want one Job", objects)
		}

		spec := objects[0].Spec()
		if spec["backoffLimit"] != 0 || spec["activeDeadlineSeconds"] != 120 || spec["completions"] != 1 || spec["parallelism"] != 1 {
			t.Errorf("job spec = %v, want one completion, one worker, no retries, and a 120-second deadline", spec)
		}

		if objects[0].PodSpec().Text("restartPolicy") != "Never" {
			t.Errorf("restartPolicy = %q, want Never", objects[0].PodSpec().Text("restartPolicy"))
		}
	})

	t.Run("objects share no mutable state", func(t *testing.T) {
		t.Parallel()

		objects := deploymentWorkload(t).NativeObjects(runtime, namespace)
		objects[0].Labels()["edited"] = "yes"
		objects[0].Spec().Child("template").Child("metadata").Child("labels")["edited"] = "yes"

		if _, shared := objects[1].Labels()["edited"]; shared {
			t.Error("editing the Deployment labels changed the Service")
		}

		if _, shared := objects[1].Spec().Child("selector")["edited"]; shared {
			t.Error("editing the pod template labels changed the Service selector")
		}
	})
}

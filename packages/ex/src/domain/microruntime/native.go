package microruntime

// Defaults that seed the native objects generated for a workload. The
// generated files are consumer-owned; these values are only a starting point.
const (
	jobDeadlineSeconds      = 120
	terminationGraceSeconds = 30
	readinessPeriodSeconds  = 5
	readinessTimeoutSeconds = 2
	readinessFailures       = 12
	httpPortName            = "http"
)

// NativeObjects returns the Kubernetes objects that start a workload's native
// configuration: a Deployment and its ClusterIP Service, or one finite Job
// with no retries. Every object is built independently, so changing one never
// changes another. Call Validate first.
func (workload Workload) NativeObjects(runtime, namespace Identifier) []Object {
	if workload.Kind == JobWorkload {
		return []Object{workload.job(runtime, namespace)}
	}

	return []Object{
		workload.deployment(runtime, namespace),
		workload.service(runtime, namespace),
	}
}

func (workload Workload) job(runtime, namespace Identifier) Object {
	pod := workload.podSpec(workload.container())
	pod["restartPolicy"] = "Never"

	return Object{
		"apiVersion": "batch/v1",
		"kind":       "Job",
		"metadata":   workload.metadata(runtime, namespace),
		"spec": map[string]any{
			"backoffLimit":          0,
			"activeDeadlineSeconds": jobDeadlineSeconds,
			"completions":           1,
			"parallelism":           1,
			"template":              workload.podTemplate(runtime, pod),
		},
	}
}

func (workload Workload) deployment(runtime, namespace Identifier) Object {
	container := workload.container()
	container["ports"] = []any{map[string]any{"name": httpPortName, "containerPort": workload.Service.Port}}
	container["readinessProbe"] = map[string]any{
		"httpGet":          map[string]any{"path": workload.Service.ReadinessPath, "port": httpPortName},
		"periodSeconds":    readinessPeriodSeconds,
		"timeoutSeconds":   readinessTimeoutSeconds,
		"failureThreshold": readinessFailures,
	}

	return Object{
		"apiVersion": "apps/v1",
		"kind":       "Deployment",
		"metadata":   workload.metadata(runtime, namespace),
		"spec": map[string]any{
			"replicas": workload.Service.Replicas,
			"selector": map[string]any{"matchLabels": workload.labels(runtime)},
			"template": workload.podTemplate(runtime, workload.podSpec(container)),
		},
	}
}

func (workload Workload) service(runtime, namespace Identifier) Object {
	return Object{
		"apiVersion": "v1",
		"kind":       "Service",
		"metadata":   workload.metadata(runtime, namespace),
		"spec": map[string]any{
			"type":     "ClusterIP",
			"selector": workload.labels(runtime),
			"ports": []any{map[string]any{
				"name":       httpPortName,
				"port":       workload.Service.Port,
				"targetPort": httpPortName,
			}},
		},
	}
}

func (workload Workload) metadata(runtime, namespace Identifier) map[string]any {
	return map[string]any{
		"name":      workload.Name.String(),
		"namespace": namespace.String(),
		"labels":    workload.labels(runtime),
	}
}

// labels returns a new map of the labels that tie generated objects to the
// runtime and the unit.
func (workload Workload) labels(runtime Identifier) map[string]any {
	return map[string]any{
		LabelRuntime: runtime.String(),
		LabelUnit:    workload.Name.String(),
	}
}

func (workload Workload) container() map[string]any {
	return map[string]any{
		"name":            workload.Name.String(),
		"image":           workload.Image,
		"imagePullPolicy": "Never",
		"resources": map[string]any{
			"requests": map[string]any{
				"cpu":    workload.Resources.CPURequest.String(),
				"memory": workload.Resources.MemoryRequest.String(),
			},
			"limits": map[string]any{
				"cpu":    workload.Resources.CPULimit.String(),
				"memory": workload.Resources.MemoryLimit.String(),
			},
		},
	}
}

func (workload Workload) podSpec(container map[string]any) map[string]any {
	return map[string]any{
		"automountServiceAccountToken":  false,
		"containers":                    []any{container},
		"terminationGracePeriodSeconds": terminationGraceSeconds,
	}
}

func (workload Workload) podTemplate(runtime Identifier, pod map[string]any) map[string]any {
	return map[string]any{
		"metadata": map[string]any{"labels": workload.labels(runtime)},
		"spec":     pod,
	}
}

package application

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/exphere/ex/src/domain/microruntime"
)

// InspectRuntime observes a runtime as recorded: its objects, containers,
// quotas, events, and optional usage. It keeps only allowlisted fields, writes
// the report as evidence, and never changes the ledger or the cluster.
type InspectRuntime struct {
	Projects   ProjectOpener
	Clusters   ClusterReader
	Observer   ClusterObserver
	Executions ExecutionReader
	Objects    WorkloadObserver
	Evidence   InspectionEvidence
	Tools      ToolVersions
}

// Inspection is a sanitized observation of one runtime. Missing metrics and
// missing image observations are reported as unavailable, never as zero.
type Inspection struct {
	Runtime        string                       `json:"runtimeId"`
	ExecutionRunID string                       `json:"executionRunId"`
	ObservationID  string                       `json:"observationId"`
	ObservedAt     time.Time                    `json:"observedAt"`
	Execution      ExecutionOutcome             `json:"state"`
	Resources      []ResourceObservation        `json:"resources"`
	Containers     []ContainerResult            `json:"containers"`
	Quotas         []QuotaObservation           `json:"quotas"`
	Events         []EventObservation           `json:"events"`
	Usage          UsageObservation             `json:"usage"`
	ToolVersions   map[string]ToolVersionResult `json:"toolVersions"`
	Cost           string                       `json:"cost"`
}

// ResourceObservation reports one recorded object as found in the cluster.
type ResourceObservation struct {
	Resource           ResourceResult        `json:"resource"`
	Observed           string                `json:"observed"`
	ObservedUID        string                `json:"observedUID,omitempty"`
	ConfiguredReplicas any                   `json:"configuredReplicas,omitempty"`
	Status             map[string]any        `json:"status,omitempty"`
	Containers         []ConfiguredContainer `json:"containers,omitempty"`
}

// MarshalJSON retains the original observation contract: an absent resource
// has only its reference and absence marker; a present resource has details.
func (observation ResourceObservation) MarshalJSON() ([]byte, error) {
	if observation.Observed == "absent" {
		return json.Marshal(struct {
			Resource ResourceResult `json:"resource"`
			Observed string         `json:"observed"`
		}{observation.Resource, observation.Observed})
	}
	return json.Marshal(struct {
		Resource           ResourceResult        `json:"resource"`
		ObservedUID        string                `json:"observedUID"`
		ConfiguredReplicas any                   `json:"configuredReplicas"`
		Status             map[string]any        `json:"status"`
		Containers         []ConfiguredContainer `json:"containers"`
	}{observation.Resource, observation.ObservedUID, observation.ConfiguredReplicas, observation.Status, observation.Containers})
}

// ConfiguredContainer is a container as declared on a live workload.
type ConfiguredContainer struct {
	Name      string `json:"name"`
	Image     string `json:"configuredImage"`
	Resources any    `json:"resources"`
}

// QuotaObservation reports one runtime-owned resource quota.
type QuotaObservation struct {
	Name           string `json:"name"`
	ConfiguredHard any    `json:"configuredHard"`
	ObservedHard   any    `json:"observedHard"`
	ObservedUsed   any    `json:"observedUsed"`
}

// EventObservation classifies an event about an owned object. Free-form
// messages can echo environment values, so they stay in Kubernetes.
type EventObservation struct {
	Name           string       `json:"name"`
	Reason         string       `json:"reason"`
	Type           string       `json:"type"`
	Count          any          `json:"count"`
	FirstTimestamp any          `json:"firstTimestamp"`
	LastTimestamp  any          `json:"lastTimestamp"`
	Object         EventSubject `json:"object"`
}

// EventSubject identifies the object an event is about.
type EventSubject struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
	UID  string `json:"uid"`
}

// UsageObservation reports pod metrics when a metrics API answered.
type UsageObservation struct {
	Status string     `json:"status"`
	Reason string     `json:"reason,omitempty"`
	Pods   []PodUsage `json:"pods,omitempty"`
}

// MarshalJSON keeps an observed empty sample distinct from unavailable data.
func (usage UsageObservation) MarshalJSON() ([]byte, error) {
	if usage.Status == "observed" {
		return json.Marshal(struct {
			Status string     `json:"status"`
			Pods   []PodUsage `json:"pods"`
		}{usage.Status, usage.Pods})
	}
	return json.Marshal(struct {
		Status string `json:"status"`
		Reason string `json:"reason"`
	}{usage.Status, usage.Reason})
}

// PodUsage is one pod's metrics sample.
type PodUsage struct {
	Pod        string           `json:"pod"`
	Timestamp  any              `json:"timestamp"`
	Window     any              `json:"window"`
	Containers []ContainerUsage `json:"containers"`
}

// ContainerUsage is one container's metrics sample.
type ContainerUsage struct {
	Name  string `json:"name"`
	Usage any    `json:"usage"`
}

// workloadStatusFields are the live status fields safe to keep as evidence.
var workloadStatusFields = []string{
	"replicas", "readyReplicas", "availableReplicas", "updatedReplicas",
	"observedGeneration", "active", "succeeded", "failed",
	"startTime", "completionTime",
}

// Execute observes the runtime and records the report.
func (useCase InspectRuntime) Execute(ctx context.Context, request RuntimeRequest) (Inspection, error) {
	runtime, err := parseRuntime(request.Runtime)
	if err != nil {
		return Inspection{}, err
	}

	owner, err := useCase.Projects.Open(request.Project)
	if err != nil {
		return Inspection{}, err
	}

	execution, err := loadExecution(useCase.Executions, owner, runtime)
	if err != nil {
		return Inspection{}, err
	}

	record, err := clusterGuard{records: useCase.Clusters, observer: useCase.Observer}.ready(ctx, owner)
	if err != nil {
		return Inspection{}, err
	}

	if record.UID != execution.ClusterUID || record.Context != execution.Context {
		return Inspection{}, fmt.Errorf("%w: inspection must target the cluster the execution used", ErrTargetMismatch)
	}

	observationID, err := newID()
	if err != nil {
		return Inspection{}, err
	}

	report := Inspection{
		Runtime:        runtime.String(),
		ExecutionRunID: execution.RunID,
		ObservationID:  observationID,
		ObservedAt:     time.Now().UTC(),
		Execution:      outcomeOf(execution),
		ToolVersions:   toolVersionResultsOf(useCase.Tools.Versions(ctx)),
		Cost:           "unavailable: no billing source",
		Usage:          UsageObservation{Status: "unavailable", Reason: "metrics API not observed"},
	}

	inspection := inspector{
		objects:   useCase.Objects,
		target:    Target{ClusterAccess: accessOf(record), Namespace: execution.Namespace.String()},
		execution: execution,
	}

	if err := inspection.observe(ctx, &report); err != nil {
		return Inspection{}, err
	}

	if err := useCase.Evidence.RecordInspection(owner, report); err != nil {
		return Inspection{}, err
	}

	return report, nil
}

// inspector gathers one inspection's observations.
type inspector struct {
	objects   WorkloadObserver
	target    Target
	execution *microruntime.Execution
}

func (inspector inspector) observe(ctx context.Context, report *Inspection) error {
	namespace, err := inspector.objects.Get(ctx, inspector.target, "namespace", inspector.target.Namespace)
	if err != nil {
		return err
	}

	if namespace == nil {
		for _, resource := range inspector.execution.Resources {
			report.Resources = append(report.Resources, ResourceObservation{Resource: resourceResultOf(resource), Observed: "absent"})
		}

		return nil
	}

	namespaceUID := inspector.execution.NamespaceUID
	if namespaceUID == "" {
		namespaceUID = namespace.UID()
	}

	if !namespace.IsOwnedBy(namespaceUID, inspector.execution.Token) {
		return fmt.Errorf("%w: namespace %s differs from the recorded runtime", microruntime.ErrOwnershipChanged, inspector.target.Namespace)
	}

	if report.Resources, err = inspector.resources(ctx); err != nil {
		return err
	}

	containers, owned, err := observeContainers(ctx, inspector.objects, inspector.target, inspector.execution)
	if err != nil {
		return err
	}

	report.Containers = containerResultsOf(containers)

	if report.Quotas, err = inspector.quotas(ctx); err != nil {
		return err
	}

	owned[namespaceUID] = true

	if report.Events, err = inspector.events(ctx, owned); err != nil {
		return err
	}

	if usage, observed := inspector.usage(ctx, containers); observed {
		report.Usage = usage
	}

	return nil
}

func (inspector inspector) resources(ctx context.Context) ([]ResourceObservation, error) {
	var observations []ResourceObservation

	for _, resource := range inspector.execution.Resources {
		live, err := inspector.objects.Get(ctx, inspector.target, resource.Kind, resource.Name)
		if err != nil {
			return nil, err
		}

		if live == nil {
			observations = append(observations, ResourceObservation{Resource: resourceResultOf(resource), Observed: "absent"})

			continue
		}

		expected := resource.UID
		if expected == "" {
			expected = live.UID()
		}

		if !live.IsOwnedBy(expected, inspector.execution.Token) {
			return nil, fmt.Errorf("%w: %s %s differs from the recorded state", microruntime.ErrOwnershipChanged, resource.Kind, resource.Name)
		}

		observations = append(observations, ResourceObservation{
			Resource:           resourceResultOf(resource),
			Observed:           "present",
			ObservedUID:        live.UID(),
			ConfiguredReplicas: live.Spec()["replicas"],
			Status:             pick(live.Status(), workloadStatusFields),
			Containers:         configuredContainersOf(live),
		})
	}

	return observations, nil
}

func (inspector inspector) quotas(ctx context.Context) ([]QuotaObservation, error) {
	quotas, err := inspector.objects.List(ctx, inspector.target, "resourcequotas")
	if err != nil {
		return nil, err
	}

	var observations []QuotaObservation

	for _, quota := range quotas {
		if !quota.IsOwnedBy(quota.UID(), inspector.execution.Token) {
			return nil, fmt.Errorf("%w: a foreign resource quota prevents trusted inspection", microruntime.ErrOwnershipChanged)
		}

		observations = append(observations, QuotaObservation{
			Name:           quota.Name(),
			ConfiguredHard: quota.Spec()["hard"],
			ObservedHard:   quota.Status()["hard"],
			ObservedUsed:   quota.Status()["used"],
		})
	}

	return observations, nil
}

// events keeps the classification of events about owned objects only.
func (inspector inspector) events(ctx context.Context, owned map[string]bool) ([]EventObservation, error) {
	events, err := inspector.objects.List(ctx, inspector.target, "events")
	if err != nil {
		return nil, err
	}

	var observations []EventObservation

	for _, event := range events {
		subject := event.Child("involvedObject")
		if !owned[subject.Text("uid")] {
			continue
		}

		observations = append(observations, EventObservation{
			Name:           event.Name(),
			Reason:         event.Text("reason"),
			Type:           event.Text("type"),
			Count:          event["count"],
			FirstTimestamp: event["firstTimestamp"],
			LastTimestamp:  event["lastTimestamp"],
			Object: EventSubject{
				Kind: subject.Text("kind"),
				Name: subject.Text("name"),
				UID:  subject.Text("uid"),
			},
		})
	}

	return observations, nil
}

// usage reads optional metrics for owned pods. An absent metrics API is not
// an error; it leaves usage reported as unavailable.
func (inspector inspector) usage(ctx context.Context, containers []microruntime.ContainerObservation) (UsageObservation, bool) {
	metrics, err := inspector.objects.PodMetrics(ctx, inspector.target)
	if err != nil {
		return UsageObservation{}, false
	}

	ownedPods := map[string]bool{}
	for _, container := range containers {
		ownedPods[container.Pod] = true
	}

	usage := UsageObservation{Status: "observed", Pods: []PodUsage{}}

	for _, sample := range metrics {
		if !ownedPods[sample.Name()] {
			continue
		}

		pod := PodUsage{Pod: sample.Name(), Timestamp: sample["timestamp"], Window: sample["window"], Containers: []ContainerUsage{}}
		for _, container := range sample.Children("containers") {
			pod.Containers = append(pod.Containers, ContainerUsage{Name: container.Text("name"), Usage: container["usage"]})
		}

		usage.Pods = append(usage.Pods, pod)
	}

	return usage, true
}

func configuredContainersOf(workload microruntime.Object) []ConfiguredContainer {
	containers := []ConfiguredContainer{}

	for _, container := range workload.PodSpec().Containers() {
		containers = append(containers, ConfiguredContainer{
			Name:      container.Text("name"),
			Image:     container.Text("image"),
			Resources: container["resources"],
		})
	}

	return containers
}

// pick copies the listed fields that are present in object.
func pick(object microruntime.Object, fields []string) map[string]any {
	picked := map[string]any{}

	for _, field := range fields {
		if value, present := object[field]; present {
			picked[field] = value
		}
	}

	return picked
}

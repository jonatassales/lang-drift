package application_test

// In-memory fakes of the ports in ports.go. The Kubernetes fake imitates the
// controllers a runtime relies on, so use cases run their real flows against
// objects shaped like kubectl output, without any native tool.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/domain/cluster"
	"github.com/exphere/ex/src/domain/microruntime"
	"github.com/exphere/ex/src/domain/project"
)

const (
	projectRoot    = "/workspace"
	kubeconfigPath = "/workspace/.exphere/state/cluster.kubeconfig"
	namespaceName  = "demo-local"
	kubeSystemUID  = "uid-kube-system"
)

// errInjected is the failure a fake returns when a test asks it to fail.
var errInjected = errors.New("injected failure")

func mustIdentifier(t testing.TB, text string) microruntime.Identifier {
	t.Helper()

	identifier, err := microruntime.ParseIdentifier(text)
	if err != nil {
		t.Fatal(err)
	}

	return identifier
}

func mustPath(t testing.TB, text string) microruntime.RelativePath {
	t.Helper()

	path, err := microruntime.ParseRelativePath(text)
	if err != nil {
		t.Fatal(err)
	}

	return path
}

// cloneObject copies an object the way a JSON API round trip would, so stored
// state never aliases what a use case holds and numbers become float64.
func cloneObject(object microruntime.Object) microruntime.Object {
	if object == nil {
		return nil
	}

	data, err := json.Marshal(object)
	if err != nil {
		panic(err)
	}

	var clone microruntime.Object
	if err := json.Unmarshal(data, &clone); err != nil {
		panic(err)
	}

	return clone
}

func cloneExecution(execution *microruntime.Execution) *microruntime.Execution {
	clone := *execution
	clone.Resources = slices.Clone(execution.Resources)
	clone.Images = slices.Clone(execution.Images)
	clone.Units = maps.Clone(execution.Units)
	clone.ToolVersions = maps.Clone(execution.ToolVersions)

	return &clone
}

// --- Projects ---------------------------------------------------------------

type fakeProjects struct {
	project project.Project
	err     error
	opened  []application.ProjectSelection
}

func (projects *fakeProjects) Open(selection application.ProjectSelection) (project.Project, error) {
	projects.opened = append(projects.opened, selection)
	if projects.err != nil {
		return project.Project{}, projects.err
	}

	return projects.project, nil
}

type fakeLocks struct {
	err      error
	locked   int
	unlocked int
}

func (locks *fakeLocks) Lock(project.Project) (func(), error) {
	if locks.err != nil {
		return nil, locks.err
	}

	locks.locked++

	return func() { locks.unlocked++ }, nil
}

// --- Cluster ----------------------------------------------------------------

type fakeClusters struct {
	record       *cluster.Record
	loadErr      error
	saveErr      error
	saveErrAfter int
	saved        []cluster.Record
}

func (clusters *fakeClusters) Load(project.Project) (cluster.Record, error) {
	switch {
	case clusters.loadErr != nil:
		return cluster.Record{}, clusters.loadErr
	case clusters.record == nil:
		return cluster.Record{}, application.NotFound("no cluster is recorded for this project")
	}

	return *clusters.record, nil
}

func (clusters *fakeClusters) Save(_ project.Project, record cluster.Record) error {
	if clusters.saveErr != nil || (clusters.saveErrAfter > 0 && len(clusters.saved) >= clusters.saveErrAfter) {
		return errors.Join(clusters.saveErr, errInjected)
	}

	clusters.saved = append(clusters.saved, record)
	clusters.record = &record

	return nil
}

func (clusters *fakeClusters) statuses() []cluster.Status {
	var statuses []cluster.Status
	for _, record := range clusters.saved {
		statuses = append(statuses, record.Status)
	}

	return statuses
}

type fakeCredentials struct {
	exists     bool
	err        error
	discardErr error
	discarded  int
}

func (credentials *fakeCredentials) Credentials(project.Project) (string, bool, error) {
	return kubeconfigPath, credentials.exists, credentials.err
}

func (credentials *fakeCredentials) DiscardCredentials(project.Project) error {
	if credentials.discardErr != nil {
		return credentials.discardErr
	}

	credentials.discarded++
	credentials.exists = false

	return nil
}

type fakeProvisioner struct {
	clusters        []string
	createErr       error
	deleteErr       error
	listErr         error
	keepAfterDelete bool
	blockCreate     bool
	created         []string
	deleted         []string
}

func (provisioner *fakeProvisioner) Clusters(context.Context) ([]string, error) {
	return slices.Clone(provisioner.clusters), provisioner.listErr
}

func (provisioner *fakeProvisioner) Create(ctx context.Context, name, _ string) (string, error) {
	provisioner.created = append(provisioner.created, name)

	if provisioner.blockCreate {
		<-ctx.Done()

		return "", ctx.Err()
	}

	if provisioner.createErr != nil {
		return "", provisioner.createErr
	}

	provisioner.clusters = append(provisioner.clusters, name)

	return "kind-" + name, nil
}

func (provisioner *fakeProvisioner) Delete(_ context.Context, name string) error {
	provisioner.deleted = append(provisioner.deleted, name)
	if provisioner.deleteErr != nil {
		return provisioner.deleteErr
	}

	if !provisioner.keepAfterDelete {
		provisioner.clusters = slices.DeleteFunc(provisioner.clusters, func(existing string) bool { return existing == name })
	}

	return nil
}

// fakeObserver reports the cluster's identity and inventory.
type fakeObserver struct {
	namespaces    map[string]string
	inventory     cluster.Inventory
	namespacesErr error
	inventoryErr  error
}

func (observer *fakeObserver) Namespaces(context.Context, application.ClusterAccess) (map[string]string, error) {
	return maps.Clone(observer.namespaces), observer.namespacesErr
}

func (observer *fakeObserver) ClusterInventory(context.Context, application.ClusterAccess) (cluster.Inventory, error) {
	return observer.inventory, observer.inventoryErr
}

type fakeImages struct {
	missing error
	loadErr error
	checked [][]string
	loaded  []string
}

func (images *fakeImages) RequireImages(_ context.Context, _ string, required []string) error {
	images.checked = append(images.checked, required)

	return images.missing
}

func (images *fakeImages) LoadImage(_ context.Context, clusterName, image string) error {
	images.loaded = append(images.loaded, clusterName+" "+image)

	return images.loadErr
}

// --- Configuration ----------------------------------------------------------

type fakeConfiguration struct {
	descriptor    microruntime.Descriptor
	manifests     func() []microruntime.Manifest
	hash          string
	compatible    []string
	loadErr       error
	snapshotErr   error
	referencesErr error
	hook          application.HookScript
	hookErr       error
	resolved      []microruntime.RelativePath
}

func (configuration *fakeConfiguration) Load(project.Project, microruntime.Identifier) (microruntime.Descriptor, error) {
	return configuration.descriptor, configuration.loadErr
}

func (configuration *fakeConfiguration) Snapshot(project.Project, microruntime.Descriptor) (application.Snapshot, error) {
	if configuration.snapshotErr != nil {
		return application.Snapshot{}, configuration.snapshotErr
	}

	return application.Snapshot{
		Manifests:            configuration.manifests(),
		PreparationDirectory: projectRoot + "/.exphere/runtimes/demo/terraform",
		Hash:                 configuration.hash,
		CompatibleHashes:     slices.Clone(configuration.compatible),
	}, nil
}

func (configuration *fakeConfiguration) CheckReferences(project.Project, microruntime.Descriptor) error {
	return configuration.referencesErr
}

func (configuration *fakeConfiguration) ResolveHook(_ project.Project, _ microruntime.Identifier, script microruntime.RelativePath) (application.HookScript, error) {
	configuration.resolved = append(configuration.resolved, script)

	return configuration.hook, configuration.hookErr
}

// --- Execution ledger and evidence ------------------------------------------

type fakeLedger struct {
	executions map[microruntime.Identifier]*microruntime.Execution
	saveErr    error
	failAfter  int
	listErr    error
	saved      []*microruntime.Execution
}

func (ledger *fakeLedger) Load(_ project.Project, runtime microruntime.Identifier) (*microruntime.Execution, error) {
	execution, recorded := ledger.executions[runtime]
	if !recorded {
		return nil, application.NotFound("runtime %s has no execution record", runtime)
	}

	return cloneExecution(execution), nil
}

func (ledger *fakeLedger) Save(_ project.Project, execution *microruntime.Execution) error {
	if ledger.saveErr != nil {
		return ledger.saveErr
	}

	if ledger.failAfter > 0 && len(ledger.saved) >= ledger.failAfter {
		return errInjected
	}

	if ledger.executions == nil {
		ledger.executions = map[microruntime.Identifier]*microruntime.Execution{}
	}

	ledger.executions[execution.Runtime] = cloneExecution(execution)
	ledger.saved = append(ledger.saved, cloneExecution(execution))

	return nil
}

func (ledger *fakeLedger) List(project.Project) ([]*microruntime.Execution, error) {
	if ledger.listErr != nil {
		return nil, ledger.listErr
	}

	var executions []*microruntime.Execution
	for _, execution := range ledger.executions {
		executions = append(executions, cloneExecution(execution))
	}

	return executions, nil
}

// last returns the most recently saved record of the runtime.
func (ledger *fakeLedger) last(t *testing.T, runtime string) *microruntime.Execution {
	t.Helper()

	execution, recorded := ledger.executions[mustIdentifier(t, runtime)]
	if !recorded {
		t.Fatalf("runtime %s has no saved execution record", runtime)
	}

	return execution
}

type fakeEvidence struct {
	err         error
	executions  []application.ExecutionOutcome
	inspections []application.Inspection
	hooks       []application.HookOutcome
	failHookAt  int
}

func (evidence *fakeEvidence) RecordExecution(_ project.Project, outcome application.ExecutionOutcome) error {
	evidence.executions = append(evidence.executions, outcome)

	return evidence.err
}

func (evidence *fakeEvidence) RecordInspection(_ project.Project, report application.Inspection) error {
	evidence.inspections = append(evidence.inspections, report)

	return evidence.err
}

func (evidence *fakeEvidence) RecordHook(_ project.Project, outcome application.HookOutcome) error {
	evidence.hooks = append(evidence.hooks, outcome)
	if evidence.failHookAt == len(evidence.hooks) {
		return errInjected
	}

	return nil
}

type fakeTools struct {
	report application.PrerequisiteReport
}

func (fakeTools) Versions(context.Context) map[string]microruntime.ToolVersion {
	return map[string]microruntime.ToolVersion{
		"kubectl":   {Version: "v1.32.0", Observed: true},
		"terraform": {Observed: false},
	}
}

func (tools fakeTools) Check(context.Context) application.PrerequisiteReport { return tools.report }

// --- Kubernetes -------------------------------------------------------------

// fakeKubernetes holds one runtime namespace and imitates the controllers a
// runtime relies on: Deployments create a ReplicaSet and a Pod, Jobs create a
// Pod and complete, Services gain controller-generated Endpoints and an
// EndpointSlice, and foreground deletion removes descendants first.
type fakeKubernetes struct {
	mu      sync.Mutex
	objects []microruntime.Object
	uids    int

	getErr       error
	dryRunErr    error
	createErr    error
	inventoryErr error
	deleteErr    error
	unready      bool
	jobFails     bool
	metrics      []microruntime.Object
	metricsErr   error
	events       []microruntime.Object
	onDelete     func(*fakeKubernetes)

	dryRuns []string
	creates []string
	deletes []string
}

func (kubernetes *fakeKubernetes) nextUID(kind string) string {
	kubernetes.uids++

	return fmt.Sprintf("uid-%s-%d", strings.ToLower(kind), kubernetes.uids)
}

func (kubernetes *fakeKubernetes) find(kind, name string) (int, microruntime.Object) {
	for index, object := range kubernetes.objects {
		if strings.EqualFold(object.Kind(), kind) && object.Name() == name {
			return index, object
		}
	}

	return -1, nil
}

// add stores a new object with a fresh UID in the runtime namespace.
func (kubernetes *fakeKubernetes) add(object microruntime.Object) microruntime.Object {
	object = cloneObject(object)
	metadata := object.Metadata()
	if metadata == nil {
		metadata = map[string]any{}
		object["metadata"] = metadata
	}

	metadata["uid"] = kubernetes.nextUID(object.Kind())
	if object.Kind() != "Namespace" {
		metadata["namespace"] = namespaceName
	}

	kubernetes.objects = append(kubernetes.objects, object)

	return object
}

// put adds an object as given, for objects a test or another actor created.
func (kubernetes *fakeKubernetes) put(object microruntime.Object) {
	kubernetes.mu.Lock()
	defer kubernetes.mu.Unlock()

	kubernetes.objects = append(kubernetes.objects, cloneObject(object))
}

// live returns a copy of the stored object, or nil.
func (kubernetes *fakeKubernetes) live(kind, name string) microruntime.Object {
	kubernetes.mu.Lock()
	defer kubernetes.mu.Unlock()

	_, object := kubernetes.find(kind, name)

	return cloneObject(object)
}

// edit changes the stored object in place.
func (kubernetes *fakeKubernetes) edit(kind, name string, change func(microruntime.Object)) {
	kubernetes.mu.Lock()
	defer kubernetes.mu.Unlock()

	if _, object := kubernetes.find(kind, name); object != nil {
		change(object)
	}
}

func (kubernetes *fakeKubernetes) Get(_ context.Context, _ application.Target, kind, name string) (microruntime.Object, error) {
	kubernetes.mu.Lock()
	defer kubernetes.mu.Unlock()

	if kubernetes.getErr != nil {
		return nil, kubernetes.getErr
	}

	_, object := kubernetes.find(kind, name)

	return cloneObject(object), nil
}

var resourceKinds = map[string]string{
	"pods": "Pod", "replicasets": "ReplicaSet", "resourcequotas": "ResourceQuota",
}

func (kubernetes *fakeKubernetes) List(_ context.Context, _ application.Target, resource string) ([]microruntime.Object, error) {
	kubernetes.mu.Lock()
	defer kubernetes.mu.Unlock()

	if resource == "events" {
		return slices.Clone(kubernetes.events), nil
	}

	kind, known := resourceKinds[resource]
	if !known {
		return nil, fmt.Errorf("fake cannot list %s", resource)
	}

	var objects []microruntime.Object
	for _, object := range kubernetes.objects {
		if object.Kind() == kind {
			objects = append(objects, cloneObject(object))
		}
	}

	return objects, nil
}

func (kubernetes *fakeKubernetes) DryRun(_ context.Context, _ application.Target, object microruntime.Object) error {
	kubernetes.mu.Lock()
	defer kubernetes.mu.Unlock()

	kubernetes.dryRuns = append(kubernetes.dryRuns, object.Kind()+"/"+object.Name())

	return kubernetes.dryRunErr
}

func (kubernetes *fakeKubernetes) Create(_ context.Context, _ application.Target, object microruntime.Object) (microruntime.Object, error) {
	kubernetes.mu.Lock()
	defer kubernetes.mu.Unlock()

	kubernetes.creates = append(kubernetes.creates, object.Kind()+"/"+object.Name())

	if kubernetes.createErr != nil {
		return nil, kubernetes.createErr
	}

	if _, existing := kubernetes.find(object.Kind(), object.Name()); existing != nil {
		return nil, fmt.Errorf("%s %s already exists", object.Kind(), object.Name())
	}

	created := kubernetes.add(object)
	created.Metadata()["generation"] = 1.0

	switch created.Kind() {
	case "Deployment":
		kubernetes.runDeployment(created)
	case "Job":
		kubernetes.runJob(created)
	}

	kubernetes.reconcileEndpoints()

	return cloneObject(created), nil
}

func controlledBy(owner microruntime.Object) []any {
	return []any{map[string]any{"apiVersion": owner.APIVersion(), "kind": owner.Kind(), "name": owner.Name(), "uid": owner.UID(), "controller": true}}
}

func (kubernetes *fakeKubernetes) runDeployment(deployment microruntime.Object) {
	replicas := 1.0
	if deployment.Spec()["replicas"] != nil {
		replicas = deployment.Spec().Number("replicas")
	}

	if !kubernetes.unready {
		deployment["status"] = map[string]any{
			"observedGeneration": 1.0, "replicas": replicas, "updatedReplicas": replicas, "availableReplicas": replicas,
		}
	}

	template := deployment.Spec().Child("template")
	replicaSet := kubernetes.add(microruntime.Object{
		"apiVersion": "apps/v1",
		"kind":       "ReplicaSet",
		"metadata": map[string]any{
			"name":            deployment.Name() + "-6f6",
			"labels":          template.Child("metadata").Child("labels"),
			"ownerReferences": controlledBy(deployment),
		},
	})

	kubernetes.runPod(replicaSet.Name()+"-a", template, replicaSet, "Running")
}

func (kubernetes *fakeKubernetes) runJob(job microruntime.Object) {
	condition := map[string]any{"type": "Complete", "status": "True"}
	if kubernetes.jobFails {
		condition = map[string]any{"type": "Failed", "status": "True", "reason": "BackoffLimitExceeded"}
	}

	if !kubernetes.unready {
		job["status"] = map[string]any{"conditions": []any{condition}}
	}

	kubernetes.runPod(job.Name()+"-x", job.Spec().Child("template"), job, "Succeeded")
}

func (kubernetes *fakeKubernetes) runPod(name string, template, owner microruntime.Object, phase string) {
	var statuses []any
	for _, container := range template.Child("spec").Children("containers") {
		statuses = append(statuses, map[string]any{
			"name":         container.Text("name"),
			"imageID":      "docker.io/library/" + container.Text("image") + "@sha256:feed",
			"ready":        phase == "Running",
			"restartCount": 0.0,
		})
	}

	kubernetes.add(microruntime.Object{
		"apiVersion": "v1",
		"kind":       "Pod",
		"metadata": map[string]any{
			"name":            name,
			"labels":          template.Child("metadata").Child("labels"),
			"ownerReferences": controlledBy(owner),
		},
		"spec":   template.Child("spec"),
		"status": map[string]any{"phase": phase, "containerStatuses": statuses},
	})
}

// reconcileEndpoints keeps one Endpoints object and one EndpointSlice per
// Service, as the endpoint controllers do.
func (kubernetes *fakeKubernetes) reconcileEndpoints() {
	services := map[string]bool{}
	for _, object := range kubernetes.objects {
		if object.Kind() == "Service" {
			services[object.Name()] = true
		}
	}

	kubernetes.objects = slices.DeleteFunc(kubernetes.objects, func(object microruntime.Object) bool {
		derived := object.Kind() == "Endpoints" || object.Kind() == "EndpointSlice"

		return derived && !services[strings.TrimSuffix(object.Name(), "-abc")]
	})

	for _, service := range slices.Clone(kubernetes.objects) {
		if service.Kind() != "Service" {
			continue
		}

		var addresses []any
		for _, pod := range kubernetes.objects {
			if pod.Kind() == "Pod" && selects(service.Spec().Child("selector"), pod.Labels()) {
				addresses = append(addresses, map[string]any{"targetRef": map[string]any{
					"kind": "Pod", "name": pod.Name(), "namespace": namespaceName, "uid": pod.UID(),
				}})
			}
		}

		labels := maps.Clone(service.Labels())
		labels["endpoints.kubernetes.io/managed-by"] = "endpoint-controller"

		subsets := []any{}
		if len(addresses) > 0 {
			subsets = []any{map[string]any{"addresses": addresses}}
		}

		if index, endpoints := kubernetes.find("Endpoints", service.Name()); endpoints != nil {
			kubernetes.objects[index]["subsets"] = subsets
		} else {
			kubernetes.add(microruntime.Object{
				"apiVersion": "v1", "kind": "Endpoints",
				"metadata": map[string]any{"name": service.Name(), "labels": labels},
				"subsets":  subsets,
			})
		}

		if _, slice := kubernetes.find("EndpointSlice", service.Name()+"-abc"); slice == nil {
			kubernetes.add(microruntime.Object{
				"apiVersion": "discovery.k8s.io/v1", "kind": "EndpointSlice",
				"metadata": map[string]any{"name": service.Name() + "-abc", "ownerReferences": controlledBy(service)},
			})
		}
	}
}

func selects(selector, labels microruntime.Object) bool {
	if len(selector) == 0 {
		return false
	}

	for key, value := range selector {
		if labels[key] != value {
			return false
		}
	}

	return true
}

func (kubernetes *fakeKubernetes) Delete(_ context.Context, _ application.Target, resource microruntime.Resource) error {
	kubernetes.mu.Lock()
	defer kubernetes.mu.Unlock()

	kubernetes.deletes = append(kubernetes.deletes, resource.Kind+"/"+resource.Name)

	if kubernetes.deleteErr != nil {
		return kubernetes.deleteErr
	}

	_, object := kubernetes.find(resource.Kind, resource.Name)
	if object == nil || object.UID() != resource.UID {
		return fmt.Errorf("precondition failed: UID in precondition %s does not match", resource.UID)
	}

	kubernetes.removeWithDescendants(object.UID())
	kubernetes.reconcileEndpoints()

	if kubernetes.onDelete != nil {
		kubernetes.onDelete(kubernetes)
	}

	return nil
}

// removeWithDescendants imitates foreground deletion.
func (kubernetes *fakeKubernetes) removeWithDescendants(uid string) {
	doomed := map[string]bool{uid: true}

	for changed := true; changed; {
		changed = false

		for _, object := range kubernetes.objects {
			if doomed[object.UID()] {
				continue
			}

			for _, owner := range object.OwnerUIDs() {
				if doomed[owner] {
					doomed[object.UID()], changed = true, true
				}
			}
		}
	}

	kubernetes.objects = slices.DeleteFunc(kubernetes.objects, func(object microruntime.Object) bool {
		return doomed[object.UID()]
	})
}

func (kubernetes *fakeKubernetes) DeleteNamespace(_ context.Context, _ application.Target, uid string) error {
	kubernetes.mu.Lock()
	defer kubernetes.mu.Unlock()

	kubernetes.deletes = append(kubernetes.deletes, "Namespace/"+namespaceName)

	if kubernetes.deleteErr != nil {
		return kubernetes.deleteErr
	}

	_, namespace := kubernetes.find("Namespace", namespaceName)
	if namespace == nil || namespace.UID() != uid {
		return fmt.Errorf("precondition failed: namespace UID %s does not match", uid)
	}

	kubernetes.objects = nil

	return nil
}

func (kubernetes *fakeKubernetes) NamespaceInventory(context.Context, application.Target) ([]microruntime.Object, error) {
	kubernetes.mu.Lock()
	defer kubernetes.mu.Unlock()

	if kubernetes.inventoryErr != nil {
		return nil, kubernetes.inventoryErr
	}

	var inventory []microruntime.Object
	for _, object := range kubernetes.objects {
		if object.Kind() != "Namespace" {
			inventory = append(inventory, cloneObject(object))
		}
	}

	return inventory, nil
}

func (kubernetes *fakeKubernetes) PodMetrics(context.Context, application.Target) ([]microruntime.Object, error) {
	return kubernetes.metrics, kubernetes.metricsErr
}

// fakePreparation creates the namespace, its defaults, and the quota, as the
// generated Terraform root does.
type fakePreparation struct {
	kubernetes *fakeKubernetes
	applyErr   error
	destroyErr error
	recreate   bool
	extra      []microruntime.Object
	state      bool
	applied    []application.PreparationRequest
	destroyed  []application.PreparationRequest
}

func (preparation *fakePreparation) Apply(_ context.Context, request application.PreparationRequest) error {
	preparation.applied = append(preparation.applied, request)
	if preparation.applyErr != nil {
		return preparation.applyErr
	}

	kubernetes := preparation.kubernetes
	kubernetes.mu.Lock()
	defer kubernetes.mu.Unlock()

	owned := map[string]any{microruntime.LabelOwnershipToken: request.Token, microruntime.LabelRuntime: request.Runtime.String()}

	kubernetes.add(microruntime.Object{"apiVersion": "v1", "kind": "Namespace", "metadata": map[string]any{"name": request.Namespace, "labels": owned}})
	kubernetes.add(microruntime.Object{"apiVersion": "v1", "kind": "ServiceAccount", "metadata": map[string]any{"name": "default"}})
	kubernetes.add(microruntime.Object{"apiVersion": "v1", "kind": "ConfigMap", "metadata": map[string]any{"name": "kube-root-ca.crt"}})
	kubernetes.add(microruntime.Object{
		"apiVersion": "v1", "kind": "ResourceQuota",
		"metadata": map[string]any{"name": "runtime", "labels": owned},
		"spec":     map[string]any{"hard": map[string]any{"pods": "10"}},
		"status":   map[string]any{"hard": map[string]any{"pods": "10"}, "used": map[string]any{"pods": "0"}},
	})

	for _, object := range preparation.extra {
		kubernetes.add(object)
	}

	preparation.state = true

	return nil
}

func (preparation *fakePreparation) HasState(project.Project, microruntime.Identifier) (bool, error) {
	return preparation.state, nil
}

func (preparation *fakePreparation) Destroy(_ context.Context, request application.PreparationRequest) error {
	preparation.destroyed = append(preparation.destroyed, request)
	if preparation.destroyErr != nil {
		return preparation.destroyErr
	}

	kubernetes := preparation.kubernetes
	kubernetes.mu.Lock()
	defer kubernetes.mu.Unlock()

	kubernetes.objects = slices.DeleteFunc(kubernetes.objects, func(object microruntime.Object) bool {
		return object.Kind() == "Namespace" || object.Kind() == "ResourceQuota"
	})
	preparation.state = false

	if preparation.recreate {
		kubernetes.add(microruntime.Object{"apiVersion": "v1", "kind": "Namespace", "metadata": map[string]any{"name": request.Namespace}})
	}

	return nil
}

// --- Diagnostics streams ----------------------------------------------------

type fakeLogs struct {
	err      error
	requests []string
}

func (logs *fakeLogs) Logs(_ context.Context, target application.Target, resource microruntime.Resource, follow bool, output, _ io.Writer) error {
	logs.requests = append(logs.requests, fmt.Sprintf("%s %s/%s follow=%t", target.Namespace, resource.Kind, resource.Name, follow))
	_, _ = io.WriteString(output, "log line\n")

	return logs.err
}

type fakeForwarder struct {
	err      error
	requests []string
}

func (forwarder *fakeForwarder) Forward(_ context.Context, target application.Target, resource microruntime.Resource, localPort, remotePort int, _ io.Writer) error {
	forwarder.requests = append(forwarder.requests, fmt.Sprintf("%s %s/%s %d:%d", target.Namespace, resource.Kind, resource.Name, localPort, remotePort))

	return forwarder.err
}

// --- World ------------------------------------------------------------------

// world is one consumer project with a ready, verified cluster and a valid
// runtime configuration named demo, built entirely from fakes.
type world struct {
	owner         project.Project
	projects      *fakeProjects
	locks         *fakeLocks
	clusters      *fakeClusters
	credentials   *fakeCredentials
	provisioner   *fakeProvisioner
	observer      *fakeObserver
	images        *fakeImages
	configuration *fakeConfiguration
	ledger        *fakeLedger
	kubernetes    *fakeKubernetes
	preparation   *fakePreparation
	evidence      *fakeEvidence
	tools         fakeTools
	logs          *fakeLogs
	forwarder     *fakeForwarder
}

func newWorld(t *testing.T) *world {
	t.Helper()

	id, err := project.ParseID("0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatal(err)
	}

	owner := project.Project{ID: id, Root: projectRoot}
	namespaces := map[string]string{"kube-system": kubeSystemUID, "default": "uid-default"}
	bootstrap := cluster.Inventory{
		Objects:       map[string]cluster.ObjectIdentity{"true/pods/kube-system/coredns": {UID: "uid-coredns"}},
		ResourceTypes: []string{"true/pods"},
	}

	record := cluster.Begin(owner, kubeconfigPath)
	if err := record.Identify(maps.Clone(namespaces)); err != nil {
		t.Fatal(err)
	}

	if err := record.Complete(bootstrap); err != nil {
		t.Fatal(err)
	}

	kubernetes := &fakeKubernetes{}

	return &world{
		owner:       owner,
		projects:    &fakeProjects{project: owner},
		locks:       &fakeLocks{},
		clusters:    &fakeClusters{record: &record},
		credentials: &fakeCredentials{exists: true},
		provisioner: &fakeProvisioner{clusters: []string{owner.ClusterName()}},
		observer:    &fakeObserver{namespaces: namespaces, inventory: bootstrap},
		images:      &fakeImages{},
		configuration: &fakeConfiguration{
			descriptor: demoDescriptor(t, record.Context),
			manifests:  func() []microruntime.Manifest { return demoManifests(t) },
			hash:       "hash-current",
			compatible: []string{"hash-previous"},
		},
		ledger:      &fakeLedger{},
		kubernetes:  kubernetes,
		preparation: &fakePreparation{kubernetes: kubernetes},
		evidence:    &fakeEvidence{},
		tools: fakeTools{report: application.PrerequisiteReport{
			Tools:                    []application.ToolStatus{{Name: "kind", Status: application.ToolAvailable}},
			ContainerEngineReachable: true,
		}},
		logs:      &fakeLogs{},
		forwarder: &fakeForwarder{},
	}
}

func demoDescriptor(t testing.TB, context string) microruntime.Descriptor {
	t.Helper()

	return microruntime.Descriptor{
		Version:            microruntime.SupportedVersion,
		ID:                 mustIdentifier(t, "demo"),
		Environment:        microruntime.LocalEnvironment,
		Owner:              "team-payments",
		Target:             microruntime.KubernetesTarget{Context: context, Namespace: mustIdentifier(t, namespaceName)},
		TerraformDirectory: mustPath(t, "terraform"),
		Units: map[microruntime.Identifier]microruntime.Unit{
			mustIdentifier(t, "web"):   {Component: mustPath(t, "services/web"), Manifests: mustPath(t, "units/web/kubernetes")},
			mustIdentifier(t, "batch"): {Component: mustPath(t, "jobs/batch"), Manifests: mustPath(t, "units/batch/kubernetes")},
		},
		Hooks: map[microruntime.Identifier]microruntime.RelativePath{
			mustIdentifier(t, "seed"): mustPath(t, "hooks/seed.sh"),
		},
	}
}

// demoManifests are fresh copies of the native objects, as a snapshot reads them.
func demoManifests(t testing.TB) []microruntime.Manifest {
	t.Helper()

	web, batch := mustIdentifier(t, "web"), mustIdentifier(t, "batch")
	container := func(name, image string) map[string]any {
		return map[string]any{"name": name, "image": image, "imagePullPolicy": "Never"}
	}

	return []microruntime.Manifest{
		{Unit: web, Object: microruntime.Object{
			"apiVersion": "apps/v1", "kind": "Deployment",
			"metadata": map[string]any{"name": "web"},
			"spec": map[string]any{
				"replicas": 1,
				"selector": map[string]any{"matchLabels": map[string]any{"app": "web"}},
				"template": map[string]any{
					"metadata": map[string]any{"labels": map[string]any{"app": "web"}},
					"spec":     map[string]any{"containers": []any{container("web", "web:dev")}},
				},
			},
		}},
		{Unit: web, Object: microruntime.Object{
			"apiVersion": "v1", "kind": "Service",
			"metadata": map[string]any{"name": "web"},
			"spec": map[string]any{
				"selector": map[string]any{"app": "web"},
				"ports":    []any{map[string]any{"port": 8080, "targetPort": "http"}},
			},
		}},
		{Unit: batch, Object: microruntime.Object{
			"apiVersion": "batch/v1", "kind": "Job",
			"metadata": map[string]any{"name": "batch"},
			"spec": map[string]any{
				"activeDeadlineSeconds": 120,
				"template": map[string]any{
					"spec": map[string]any{"restartPolicy": "Never", "containers": []any{container("batch", "batch:dev")}},
				},
			},
		}},
	}
}

func (w *world) runtimeRequest() application.RuntimeRequest {
	return application.RuntimeRequest{Project: application.ProjectSelection{Directory: projectRoot}, Runtime: "demo"}
}

func (w *world) clusterRequest() application.ClusterRequest {
	return application.ClusterRequest{Project: application.ProjectSelection{Directory: projectRoot}}
}

func (w *world) startRuntime() application.StartRuntime {
	return application.StartRuntime{
		Projects: w.projects, Locks: w.locks, Configuration: w.configuration, Clusters: w.clusters,
		Observer: w.observer, Images: w.images, Executions: w.ledger, Objects: w.kubernetes,
		Preparation: w.preparation, Evidence: w.evidence, Tools: w.tools,
	}
}

func (w *world) stopRuntime() application.StopRuntime {
	return application.StopRuntime{
		Projects: w.projects, Locks: w.locks, Configuration: w.configuration, Clusters: w.clusters,
		Observer: w.observer, Executions: w.ledger, Objects: w.kubernetes,
		Preparation: w.preparation, Evidence: w.evidence, Tools: w.tools,
	}
}

// up starts the demo runtime and fails the test unless it is running.
func (w *world) up(t *testing.T) application.ExecutionOutcome {
	t.Helper()

	outcome, err := w.startRuntime().Execute(context.Background(), w.runtimeRequest())
	if err != nil {
		t.Fatalf("runtime up: %v", err)
	}

	return outcome
}

// requireLockBalanced checks that every acquired lock was released.
func (w *world) requireLockBalanced(t *testing.T) {
	t.Helper()

	if w.locks.locked != w.locks.unlocked {
		t.Fatalf("locked %d times but unlocked %d times", w.locks.locked, w.locks.unlocked)
	}
}

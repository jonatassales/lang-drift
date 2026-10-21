package connected

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jonatassales/tao/packages/ex/src/application"
	app "github.com/jonatassales/tao/packages/ex/src/application/connected"
	"github.com/jonatassales/tao/packages/ex/src/domain/connection"
	"github.com/jonatassales/tao/packages/ex/src/domain/microruntime"
	"github.com/jonatassales/tao/packages/ex/src/domain/project"
	"github.com/jonatassales/tao/packages/ex/src/infrastructure/configuration"
	"github.com/jonatassales/tao/packages/ex/src/infrastructure/files"
	"github.com/jonatassales/tao/packages/ex/src/infrastructure/process"
)

type clusterCreator interface {
	Execute(context.Context, application.ClusterRequest) (application.ClusterResult, error)
}
type imageLoader interface {
	Execute(context.Context, application.LoadClusterImageRequest) (application.ImageResult, error)
}
type starter interface {
	Execute(context.Context, application.RuntimeRequest) (application.ExecutionOutcome, error)
}
type stopper interface {
	Execute(context.Context, application.StopRuntimeRequest) (application.ExecutionOutcome, error)
}
type inspector interface {
	Execute(context.Context, application.RuntimeRequest) (application.Inspection, error)
}
type runner interface {
	Run(context.Context, process.Command) (process.Result, error)
}

// Native explicitly reuses the reviewed ownership-aware local lifecycle.
type Native struct {
	Create       clusterCreator
	Load         imageLoader
	StartRuntime starter
	StopRuntime  stopper
	Inspect      inspector
	Executions   application.ExecutionReader
	Runner       runner
}

func request(owner project.Project, id string) application.RuntimeRequest {
	return application.RuntimeRequest{Project: application.ProjectSelection{Directory: owner.Root, Explicit: true}, Runtime: id}
}
func digest(data []byte) string { sum := sha256.Sum256(data); return hex.EncodeToString(sum[:]) }

func (native Native) VerifyRepository(ctx context.Context, owner project.Project, repository app.Repository) error {
	if repository.Kind == "local" {
		root, err := files.CanonicalDirectory(repository.Path)
		if err != nil || root != owner.Root {
			return errors.New("local folder differs from the connected project")
		}
		return nil
	}
	if repository.FullName == "" {
		return errors.New("console repository identity is missing")
	}
	result, err := native.Runner.Run(ctx, process.Command{Name: "git", Args: []string{"config", "--get", "remote.origin.url"}, Dir: owner.Root})
	if err != nil {
		return fmt.Errorf("project requires the connected GitHub repository as origin: %w", err)
	}
	origin := strings.TrimSpace(result.Stdout)
	origin = strings.TrimSuffix(origin, ".git")
	origin = strings.TrimPrefix(origin, "git@github.com:")
	origin = strings.TrimPrefix(origin, "https://github.com/")
	origin = strings.TrimPrefix(origin, "ssh://git@github.com/")
	if !strings.EqualFold(origin, repository.FullName) {
		return errors.New("local origin differs from the console repository; no local files were changed")
	}
	head, err := native.Runner.Run(ctx, process.Command{Name: "git", Args: []string{"rev-parse", "HEAD"}, Dir: owner.Root})
	if err != nil {
		return fmt.Errorf("local source revision is unavailable: %w", err)
	}
	if repository.Ref != "" && strings.TrimSpace(head.Stdout) != repository.Ref {
		return errors.New("local HEAD differs from the imported repository revision")
	}
	dirty, err := native.Runner.Run(ctx, process.Command{Name: "git", Args: []string{"status", "--porcelain", "--untracked-files=normal", "--", ".", ":(exclude).langdrift"}, Dir: owner.Root})
	if err != nil || strings.TrimSpace(dirty.Stdout) != "" {
		return errors.New("source changes must be committed and imported before execution; generated .langdrift files are excluded")
	}
	return nil
}

func blueprint(owner project.Project, config connection.Configuration) (microruntime.Blueprint, error) {
	runtime, err := microruntime.ParseIdentifier(config.RuntimeID)
	if err != nil {
		return microruntime.Blueprint{}, err
	}
	cpu, _ := microruntime.ParseQuantity(strconv.Itoa(config.CPUMillis) + "m")
	memory, _ := microruntime.ParseQuantity(strconv.Itoa(config.MemoryMiB) + "Mi")
	plan := microruntime.Blueprint{Owner: config.ProjectID, Target: microruntime.KubernetesTarget{Context: "kind-" + owner.ClusterName(), Namespace: runtime}, Quotas: microruntime.ResourceQuotas{microruntime.QuotaRequestsCPU: cpu, microruntime.QuotaLimitsCPU: cpu, microruntime.QuotaRequestsMemory: memory, microruntime.QuotaLimitsMemory: memory}}
	for i, component := range config.Components {
		name, _ := microruntime.ParseIdentifier("c-" + digest([]byte(component.Path))[:12])
		relative, err := microruntime.ParseComponentPath(component.Path)
		if err != nil {
			return plan, err
		}
		cpuShare := config.CPUMillis / len(config.Components)
		memoryShare := config.MemoryMiB / len(config.Components)
		if i < config.CPUMillis%len(config.Components) {
			cpuShare++
		}
		if i < config.MemoryMiB%len(config.Components) {
			memoryShare++
		}
		cpu, _ := microruntime.ParseQuantity(strconv.Itoa(cpuShare) + "m")
		memory, _ := microruntime.ParseQuantity(strconv.Itoa(memoryShare) + "Mi")
		plan.Workloads = append(plan.Workloads, microruntime.Workload{Name: name, Component: relative, Image: component.Image, Kind: microruntime.DeploymentWorkload, Resources: microruntime.Resources{CPURequest: cpu, CPULimit: cpu, MemoryRequest: memory, MemoryLimit: memory}, Service: microruntime.ServiceSpec{Port: component.Port, ReadinessPath: component.ReadinessPath, Replicas: 1}})
	}
	return plan, plan.Validate()
}

// Materialize refuses changed authored/generated files. Native v2 migration
// retains all original paths, descriptors and evidence; only Markdown is added.
func (native Native) Materialize(ctx context.Context, owner project.Project, record app.Connection, config connection.Configuration, markdown string) (map[string]string, string, error) {
	if err := native.VerifyRepository(ctx, owner, record.Repository); err != nil {
		return nil, "", err
	}
	if err := native.requireStopped(owner, config.RuntimeID); err != nil {
		return nil, "", err
	}
	runtime, _ := microruntime.ParseIdentifier(config.RuntimeID)
	rendered := map[string][]byte{"runtime.md": []byte(markdown)}
	plan, err := blueprint(owner, config)
	if err != nil {
		return nil, "", err
	}
	migrating := false
	var originalDescriptor *microruntime.Descriptor
	if config.Native != nil {
		original, fresh, err := native.migrationSource(owner, runtime)
		if err != nil {
			return nil, "", err
		}
		migrating = fresh
		originalDescriptor = &original
		plan.Target = original.Target
		for index := range plan.Workloads {
			for name, unit := range original.Units {
				if unit.Component == plan.Workloads[index].Component {
					plan.Workloads[index].Name = name
				}
			}
		}
	}
	generated, err := configuration.RenderBlueprint(runtime, plan)
	if err != nil {
		return nil, "", err
	}
	if config.Native != nil {
		if err := native.preserveNative(owner, *originalDescriptor, plan, generated); err != nil {
			return nil, "", err
		}
		generated, err = relocateMigration(generated)
		if err != nil {
			return nil, "", err
		}
	}
	for name, data := range generated {
		rendered[name] = data
	}

	directory, err := files.Within(owner.Root, microruntime.ConfigurationDirectory(runtime).String())
	if err != nil {
		return nil, "", err
	}
	// Validate all writes and removals before changing any file. A file left
	// by an interrupted materialization matches the recorded intent and is
	// adopted rather than reported as a conflict.
	names := make([]string, 0, len(rendered))
	for name := range rendered {
		names = append(names, name)
	}
	sort.Strings(names)
	intent, err := native.readIntent(owner, config.RuntimeID)
	if err != nil {
		return nil, "", err
	}
	ours := func(name string, data []byte) bool {
		return record.Managed[name] == digest(data) || intent[name] == digest(data)
	}
	for name, prior := range record.Managed {
		target, err := files.Existing(directory, name)
		if err != nil {
			return nil, "", fmt.Errorf("local configuration conflict: %w", err)
		}
		data, err := files.ReadBounded(target, 4<<20)
		if err != nil || (digest(data) != prior && intent[name] != digest(data)) {
			return nil, "", fmt.Errorf("local configuration conflict at %s", name)
		}
	}
	for _, name := range names {
		target, err := files.Within(directory, name)
		if err != nil {
			return nil, "", err
		}
		data, err := files.ReadBounded(target, 4<<20)
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			return nil, "", err
		}
		if err == nil && record.Managed[name] == "" && digest(data) != digest(rendered[name]) && !ours(name, data) && !(migrating && name == "schema.json") {
			return nil, "", fmt.Errorf("unmanaged local configuration conflicts at %s", name)
		}
	}
	// The intent persists before the first write, so a crash between writes
	// leaves files the next attempt recognizes as its own.
	managed := map[string]string{}
	for _, name := range names {
		managed[name] = digest(rendered[name])
	}
	if err := native.writeIntent(owner, config.RuntimeID, managed); err != nil {
		return nil, "", err
	}
	for _, name := range names {
		target, _ := files.Within(directory, name)
		if err := atomicWrite(target, rendered[name], 0o644); err != nil {
			return nil, "", err
		}
	}
	for name := range record.Managed {
		if _, keep := rendered[name]; !keep {
			target, err := files.Existing(directory, name)
			if err != nil {
				return nil, "", err
			}
			if err := os.Remove(target); err != nil {
				return nil, "", err
			}
		}
	}
	snapshot, err := native.snapshot(owner, config)
	if err != nil {
		return nil, "", err
	}
	return managed, snapshot.Hash, nil
}
func (native Native) requireStopped(owner project.Project, id string) error {
	runtime, err := microruntime.ParseIdentifier(id)
	if err != nil {
		return err
	}
	execution, err := native.Executions.Load(owner, runtime)
	if errors.Is(err, application.ErrNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	if err := execution.Validate(owner.ID, runtime); err != nil {
		return err
	}
	if err := execution.RequireClosed(); err != nil {
		return fmt.Errorf("%w; explicitly recover the legacy runtime with ex agent migrate %s --recover", err, id)
	}
	return nil
}
func (native Native) snapshot(owner project.Project, config connection.Configuration) (application.Snapshot, error) {
	runtime, _ := microruntime.ParseIdentifier(config.RuntimeID)
	source := configuration.Files{}
	descriptor, err := source.Load(owner, runtime)
	if err != nil {
		return application.Snapshot{}, err
	}
	if err = descriptor.Validate(); err != nil {
		return application.Snapshot{}, err
	}
	snapshot, err := source.Snapshot(owner, descriptor)
	if err != nil {
		return snapshot, err
	}
	if err = descriptor.ValidateManifests(snapshot.Manifests); err != nil {
		return snapshot, err
	}
	if len(descriptor.Units) != len(config.Components) {
		return snapshot, errors.New("native component selection differs from the console")
	}
	selected := map[string]connection.Component{}
	for _, component := range config.Components {
		selected[component.Path] = component
	}
	for _, unit := range descriptor.Units {
		if _, ok := selected[unit.Component.String()]; !ok {
			return snapshot, errors.New("native component selection differs from the console")
		}
	}
	var cpu, memory float64
	for _, manifest := range snapshot.Manifests {
		if manifest.Object.Kind() == "Service" {
			continue
		}
		replicas := manifest.Object.Spec().Number("replicas")
		if replicas == 0 {
			replicas = 1
		}
		for _, container := range manifest.Object.PodSpec().Containers() {
			limits := container.Child("resources").Child("limits")
			c, e := quantity(limits.Text("cpu"), false)
			if e != nil {
				return snapshot, e
			}
			m, e := quantity(limits.Text("memory"), true)
			if e != nil {
				return snapshot, e
			}
			cpu += c * replicas
			memory += m * replicas
		}
	}
	if cpu*1000 > float64(config.CPUMillis)+0.0001 || memory > float64(config.MemoryMiB)+0.0001 {
		return snapshot, errors.New("native workloads exceed the configured per-execution CPU or memory envelope")
	}
	return snapshot, nil
}

// Start runs only explicit recipes, then loads the images into the private
// project cluster and waits for native readiness. Diagnostics stay local.
func (native Native) Start(ctx context.Context, owner project.Project, record app.Connection, config connection.Configuration, report func(string, string) error) error {
	if err := native.VerifyRepository(ctx, owner, record.Repository); err != nil {
		return err
	}
	if err := verifyManaged(owner, record); err != nil {
		return err
	}
	snapshot, err := native.snapshot(owner, config)
	if err != nil {
		return err
	}
	if snapshot.Hash != record.ConcreteHash {
		return errors.New("native configuration changed after materialization")
	}
	var buildRoot, revision string
	var cleanup func()
	if record.Repository.Kind == "local" {
		buildRoot, cleanup, revision, err = native.localBuildTree(ctx, owner)
	} else {
		var source process.Result
		source, err = native.Runner.Run(ctx, process.Command{Name: "git", Args: []string{"rev-parse", "HEAD"}, Dir: owner.Root})
		if err == nil {
			revision = strings.TrimSpace(source.Stdout)
			buildRoot, cleanup, err = native.buildTree(ctx, owner)
		}
	}
	if err != nil {
		return err
	}
	defer cleanup()
	if err := report("build_started", revision); err != nil {
		return err
	}
	// The whole build's output streams to one private log, so a verbose build
	// never fails on capture limits and a failed one names its evidence.
	buildLog, logPath, err := native.buildLog(owner, record.RuntimeID)
	if err != nil {
		return err
	}
	defer buildLog.Close()
	logWriter := process.SyncWriter(buildLog)
	for _, component := range config.Components {
		directory, err := files.Existing(buildRoot, component.Path)
		if err != nil {
			return err
		}
		// A Docker recipe must refer to an existing Dockerfile. No tool installation
		// or implicit fallback recipe is attempted.
		if component.Build.Executable == "docker" {
			if _, err := files.Existing(directory, "Dockerfile"); err != nil {
				return errors.New("declared Docker build requires an existing Dockerfile")
			}
		}
		command := process.Command{
			Name: component.Build.Executable,
			Args: component.Build.Args,
			Dir:  directory,
			// Recipes see the minimal base plus the container engine's own
			// selection variables, never the agent's credentials.
			PassEnv: []string{"DOCKER_*"},
			Stdout:  logWriter,
			Stderr:  logWriter,
		}
		// A committed script resolves inside the component's checkout directory,
		// through the symlink-refusing resolver, and runs by absolute path.
		if strings.Contains(component.Build.Executable, "/") {
			resolved, err := files.Existing(directory, component.Build.Executable)
			if err != nil {
				return fmt.Errorf("committed build script: %w", err)
			}
			command.Name = resolved
		}
		if result, err := native.Runner.Run(ctx, command); err != nil {
			_ = report("build_failed", "build_failed")
			failure := fmt.Errorf("build recipe %s failed for %s; no runtime was started (full output: %s)", component.Build.Executable, component.Path, logPath)
			if tail := process.Tail(result.Stderr, 2048); tail != "" {
				return fmt.Errorf("%w\n%s", failure, tail)
			}
			return failure
		}
	}
	if err := report("build_succeeded", "build_completed"); err != nil {
		return err
	}
	if _, err := native.Create.Execute(ctx, application.ClusterRequest{Project: request(owner, record.RuntimeID).Project}); err != nil {
		return err
	}
	for _, component := range config.Components {
		if _, err := native.Load.Execute(ctx, application.LoadClusterImageRequest{Project: request(owner, record.RuntimeID).Project, Image: component.Image}); err != nil {
			return err
		}
	}
	if err := verifyManaged(owner, record); err != nil {
		return err
	}
	_, err = native.StartRuntime.Execute(ctx, request(owner, record.RuntimeID))
	return err
}
func (native Native) Stop(ctx context.Context, owner project.Project, record app.Connection) error {
	_, err := native.StopRuntime.Execute(ctx, application.StopRuntimeRequest{RuntimeRequest: request(owner, record.RuntimeID), Recover: true})
	if errors.Is(err, application.ErrNotFound) {
		return nil
	}
	return err
}
func (native Native) Observe(ctx context.Context, owner project.Project, record app.Connection) (*app.Observation, error) {
	result, err := native.Inspect.Execute(ctx, request(owner, record.RuntimeID))
	if err != nil {
		return nil, err
	}
	observation := &app.Observation{ObservedAt: result.ObservedAt, State: "running"}
	if result.Execution.Status != "running" || !result.Execution.Settled || len(result.Resources) == 0 {
		observation.State = "unknown"
	}
	for _, resource := range result.Resources {
		if resource.Observed == "absent" {
			observation.State = "unknown"
		}
	}
	for _, container := range result.Containers {
		if !container.Ready && container.Phase != "Succeeded" {
			observation.State = "unknown"
		}
	}
	if len(result.Containers) == 0 {
		observation.State = "unknown"
	}
	if observation.State != "running" || result.Usage.Status != "observed" || len(result.Usage.Pods) == 0 {
		return observation, nil
	}
	// Match every running pod/container before summing. A successful metrics API
	// response may contain only a subset of the runtime's current containers.
	expected := map[string]map[string]bool{}
	for _, container := range result.Containers {
		if container.Phase != "Running" {
			continue
		}
		if container.Pod == "" || container.Container == "" {
			return observation, nil
		}
		if expected[container.Pod] == nil {
			expected[container.Pod] = map[string]bool{}
		}
		expected[container.Pod][container.Container] = true
	}
	if len(expected) == 0 || len(expected) > 128 {
		return observation, nil
	}
	seen := map[string]bool{}
	var cpu, memory float64
	var oldest, measuredFrom time.Time
	windows := make([]app.MeasurementWindow, 0, len(expected))
	for _, pod := range result.Usage.Pods {
		containers, relevant := expected[pod.Pod]
		if !relevant {
			continue
		}
		if seen[pod.Pod] || len(pod.Containers) != len(containers) {
			return observation, nil
		}
		seen[pod.Pod] = true
		timestamp, timeOK := pod.Timestamp.(string)
		window, windowOK := pod.Window.(string)
		observedAt, timeErr := time.Parse(time.RFC3339Nano, timestamp)
		duration, windowErr := time.ParseDuration(window)
		if !timeOK || !windowOK || timeErr != nil || windowErr != nil || observedAt.IsZero() || duration <= 0 {
			return observation, nil
		}
		from := observedAt.Add(-duration)
		if oldest.IsZero() || observedAt.Before(oldest) {
			oldest = observedAt
		}
		if measuredFrom.IsZero() || from.Before(measuredFrom) {
			measuredFrom = from
		}
		windows = append(windows, app.MeasurementWindow{Pod: pod.Pod, ObservedAt: observedAt, Window: window})
		matched := map[string]bool{}
		for _, container := range pod.Containers {
			if !containers[container.Name] || matched[container.Name] {
				return observation, nil
			}
			matched[container.Name] = true
			usage, ok := container.Usage.(map[string]any)
			if typed, valid := container.Usage.(microruntime.Object); valid {
				usage = typed
				ok = true
			}
			if !ok {
				return observation, nil
			}
			c, err := quantity(fmt.Sprint(usage["cpu"]), false)
			if err != nil {
				return observation, nil
			}
			m, err := quantity(fmt.Sprint(usage["memory"]), true)
			if err != nil {
				return observation, nil
			}
			cpu += c
			memory += m
		}
	}
	if len(seen) != len(expected) {
		return observation, nil
	}
	observation.CPUUsed = &cpu
	observation.MemoryUsedMiB = &memory
	observation.ObservedAt = oldest
	observation.MeasuredFrom = &measuredFrom
	observation.MeasurementWindows = windows
	return observation, nil
}

// quantity accepts the native resource units supported by this adapter. Values
// outside the subset fail explicitly rather than inventing measurements.
func quantity(text string, memory bool) (float64, error) {
	factors := []struct {
		suffix string
		factor float64
	}{{"Ki", 1024}, {"Mi", 1024 * 1024}, {"Gi", 1024 * 1024 * 1024}, {"Ti", 1024 * 1024 * 1024 * 1024}, {"n", 1e-9}, {"u", 1e-6}, {"m", 1e-3}, {"k", 1e3}, {"M", 1e6}, {"G", 1e9}}
	factor := 1.0
	for _, entry := range factors {
		if strings.HasSuffix(text, entry.suffix) {
			text = strings.TrimSuffix(text, entry.suffix)
			factor = entry.factor
			break
		}
	}
	value, err := strconv.ParseFloat(text, 64)
	if err != nil || value < 0 || math.IsNaN(value) || math.IsInf(value, 0) {
		return 0, errors.New("unsupported native resource quantity")
	}
	value *= factor
	if memory {
		value /= 1024 * 1024
	}
	return value, nil
}

func (native Native) Migrate(ctx context.Context, owner project.Project, id string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	if err := native.requireStopped(owner, id); err != nil {
		return "", err
	}
	runtime, err := microruntime.ParseIdentifier(id)
	if err != nil {
		return "", err
	}
	source := configuration.Files{}
	descriptor, err := source.Load(owner, runtime)
	if err != nil {
		return "", err
	}
	if err = descriptor.Validate(); err != nil {
		return "", err
	}
	snapshot, err := source.Snapshot(owner, descriptor)
	if err != nil {
		return "", err
	}
	if err = descriptor.ValidateManifests(snapshot.Manifests); err != nil {
		return "", err
	}
	target, err := files.Within(owner.Root, microruntime.ConfigurationDirectory(runtime).String()+"/migration.md")
	if err != nil {
		return "", err
	}
	if _, err := os.Lstat(target); err == nil {
		return "", errors.New("migration.md already exists; preserving your draft")
	} else if !errors.Is(err, fs.ErrNotExist) {
		return "", err
	}
	text := "# Migrate " + id + "\n\nImport this existing JSON v2 runtime ID in the console: `" + id + "`.\n\nThe original descriptor, native configuration, ownership ledger and evidence remain unchanged. Select the same components below, provide their existing image/build recipes, and set CPU/RAM no lower than the sum of native workload limits. Configure the agent model and effort in the console. Connect this checkout only after all executions are confirmed stopped.\n\nNative configuration fingerprint: `" + snapshot.Hash + "`. This is configuration evidence, not an observed execution.\n\n"
	for _, unit := range descriptor.UnitNames() {
		text += "- Component: `" + descriptor.Units[unit].Component.String() + "`; native manifests: `" + descriptor.Units[unit].Manifests.String() + "`.\n"
	}
	text += "\nThe console will produce the authoritative runtime.md with `native.runtimeId` equal to `" + id + "`. Local conflicts block changes; no commit or infrastructure action is performed by migration.\n"
	return target, atomicWrite(target, []byte(text), 0o644)
}

// migrationSource stores the original descriptor before switching its active
// references. Existing native files, source paths and evidence remain intact.
func (native Native) migrationSource(owner project.Project, runtime microruntime.Identifier) (microruntime.Descriptor, bool, error) {
	source := configuration.Files{}
	target, err := files.Within(owner.Root, ".langdrift/state/migrations/"+runtime.String()+"/schema.json")
	if err != nil {
		return microruntime.Descriptor{}, false, err
	}
	original, err := files.ReadBounded(target, 4<<20)
	fresh := errors.Is(err, fs.ErrNotExist)
	if err != nil && !fresh {
		return microruntime.Descriptor{}, false, err
	}
	if fresh {
		descriptor, err := source.Load(owner, runtime)
		if err != nil {
			return descriptor, false, err
		}
		if err = descriptor.Validate(); err != nil {
			return descriptor, false, err
		}
		snapshot, err := source.Snapshot(owner, descriptor)
		if err != nil {
			return descriptor, false, err
		}
		if err = descriptor.ValidateManifests(snapshot.Manifests); err != nil {
			return descriptor, false, err
		}
		existing, err := files.Existing(owner.Root, microruntime.ConfigurationDirectory(runtime).String()+"/schema.json")
		if err != nil {
			return descriptor, false, err
		}
		original, err = files.ReadBounded(existing, 4<<20)
		if err != nil {
			return descriptor, false, err
		}
		if err = atomicWrite(target, original, 0o600); err != nil {
			return descriptor, false, err
		}
		return descriptor, true, nil
	}
	// Parse the preserved source without moving or rewriting original native files.
	descriptor, err := configuration.DecodeDescriptor(original)
	activePath, pathErr := files.Existing(owner.Root, microruntime.ConfigurationDirectory(runtime).String()+"/schema.json")
	if pathErr != nil {
		return descriptor, false, pathErr
	}
	active, readErr := files.ReadBounded(activePath, 4<<20)
	return descriptor, readErr == nil && digest(active) == digest(original), err
}

// readIntent loads the digests an earlier materialization intended to write.
func (native Native) readIntent(owner project.Project, runtime string) (map[string]string, error) {
	path, err := files.Within(owner.Root, ".langdrift/state/runtimes/"+runtime+"/materialize-intent.json")
	if err != nil {
		return nil, err
	}
	data, err := files.ReadBounded(path, 4<<20)
	if errors.Is(err, fs.ErrNotExist) {
		return map[string]string{}, nil
	}
	if err != nil {
		return nil, err
	}
	intent := map[string]string{}
	if err := json.Unmarshal(data, &intent); err != nil {
		return nil, fmt.Errorf("materialize intent is unreadable: %w", err)
	}
	return intent, nil
}

// writeIntent records what the current materialization is about to write.
func (native Native) writeIntent(owner project.Project, runtime string, managed map[string]string) error {
	path, err := files.Within(owner.Root, ".langdrift/state/runtimes/"+runtime+"/materialize-intent.json")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(managed, "", "  ")
	if err != nil {
		return err
	}
	return atomicWrite(path, append(data, '\n'), 0o600)
}

// buildLog opens the runtime's private build log for one start.
func (native Native) buildLog(owner project.Project, runtime string) (*os.File, string, error) {
	directory, err := files.Within(owner.Root, ".langdrift/state/runtimes/"+runtime)
	if err != nil {
		return nil, "", err
	}
	if err := os.MkdirAll(directory, 0o700); err != nil {
		return nil, "", err
	}
	path := filepath.Join(directory, "build.log")
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return nil, "", err
	}
	return file, path, nil
}

func relocateMigration(rendered map[string][]byte) (map[string][]byte, error) {
	var descriptor map[string]any
	if err := json.Unmarshal(rendered["schema.json"], &descriptor); err != nil {
		return nil, err
	}
	// The schema was rendered by this process, but a shape check beats a panic
	// if rendering and relocation ever drift apart.
	terraform, ok := descriptor["terraform"].(map[string]any)
	if !ok {
		return nil, errors.New("rendered schema lacks a terraform object")
	}
	terraform["directory"] = "console-terraform"
	units, ok := descriptor["units"].(map[string]any)
	if !ok {
		return nil, errors.New("rendered schema lacks a units object")
	}
	for name, value := range units {
		unit, ok := value.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("rendered unit %s is not an object", name)
		}
		kubernetes, ok := unit["kubernetes"].(map[string]any)
		if !ok {
			return nil, fmt.Errorf("rendered unit %s lacks a kubernetes object", name)
		}
		manifests, ok := kubernetes["manifests"].(string)
		if !ok {
			return nil, fmt.Errorf("rendered unit %s lacks a manifests path", name)
		}
		kubernetes["manifests"] = "console-" + manifests
	}
	data, err := json.MarshalIndent(descriptor, "", "  ")
	if err != nil {
		return nil, err
	}
	result := map[string][]byte{"schema.json": append(data, '\n')}
	for name, data := range rendered {
		if name != "schema.json" {
			result["console-"+name] = data
		}
	}
	return result, nil
}

// preserveNative retains consumer settings outside console-owned image and
// resource fields. Ambiguous multi-workload/container units require an explicit
// consumer migration instead of silently discarding native configuration.
func (native Native) preserveNative(owner project.Project, original microruntime.Descriptor, plan microruntime.Blueprint, rendered map[string][]byte) error {
	snapshot, err := (configuration.Files{}).Snapshot(owner, original)
	if err != nil {
		return err
	}
	for _, workload := range plan.Workloads {
		var objects []microruntime.Object
		for _, manifest := range snapshot.Manifests {
			if manifest.Unit == workload.Name {
				objects = append(objects, manifest.Object)
			}
		}
		if len(objects) == 0 {
			continue
		}
		workloads := 0
		for _, object := range objects {
			target := "units/" + workload.Name.String() + "/kubernetes/"
			if object.Kind() == "Service" {
				object.Spec()["ports"] = []any{map[string]any{"name": "http", "port": workload.Service.Port, "targetPort": "http"}}
				target += "service.json"
			} else {
				workloads++
				if object.Kind() == "Job" {
					delete(rendered, "units/"+workload.Name.String()+"/kubernetes/service.json")
				}
				if workloads > 1 || len(object.PodSpec().Containers()) != 1 {
					return errors.New("migration requires one workload and one container per existing unit")
				}
				target += "workload.json"
				container := object.PodSpec().Containers()[0]
				container["image"] = workload.Image
				container["resources"] = map[string]any{"requests": map[string]any{"cpu": workload.Resources.CPURequest.String(), "memory": workload.Resources.MemoryRequest.String()}, "limits": map[string]any{"cpu": workload.Resources.CPULimit.String(), "memory": workload.Resources.MemoryLimit.String()}}
				if object.Kind() == "Deployment" {
					object.Spec()["replicas"] = 1
					container["ports"] = []any{map[string]any{"name": "http", "containerPort": workload.Service.Port}}
					probe := container.Child("readinessProbe")
					if len(probe) == 0 {
						probe = microruntime.Object{}
						container["readinessProbe"] = probe
					}
					delete(probe, "exec")
					delete(probe, "tcpSocket")
					delete(probe, "grpc")
					probe["httpGet"] = map[string]any{"path": workload.Service.ReadinessPath, "port": "http"}
				}
			}
			data, err := json.MarshalIndent(object, "", "  ")
			if err != nil {
				return err
			}
			rendered[target] = append(data, '\n')
		}
	}
	return nil
}

// buildTree copies only committed project source, excluding platform state and
// dotenv credentials. Builds cannot accidentally send the agent's private .env
// through a root Docker context. This is a bounded context, not a code sandbox.
func (native Native) buildTree(ctx context.Context, owner project.Project) (string, func(), error) {
	result, err := native.Runner.Run(ctx, process.Command{Name: "git", Args: []string{"ls-files", "-z"}, Dir: owner.Root})
	if err != nil {
		return "", nil, fmt.Errorf("cannot enumerate committed build source: %w", err)
	}
	return copyBuildTree(ctx, owner, strings.Split(result.Stdout, "\x00"))
}

func copyBuildTree(ctx context.Context, owner project.Project, names []string) (string, func(), error) {
	parent, err := files.Within(owner.Root, ".langdrift/state/builds")
	if err != nil {
		return "", nil, err
	}
	if err = os.MkdirAll(parent, 0o700); err != nil {
		return "", nil, err
	}
	root, err := os.MkdirTemp(parent, "source-*")
	if err != nil {
		return "", nil, err
	}
	cleanup := func() { _ = os.RemoveAll(root) }
	var size int
	for _, name := range names {
		if err := ctx.Err(); err != nil {
			cleanup()
			return "", nil, err
		}
		if name == "" || name == ".langdrift" || strings.HasPrefix(name, ".langdrift/") {
			continue
		}
		base := filepath.Base(name)
		if base == ".env" || strings.HasPrefix(base, ".env.") {
			continue
		}
		source, err := files.Existing(owner.Root, name)
		if err != nil {
			cleanup()
			return "", nil, err
		}
		data, err := files.ReadBounded(source, 8<<20)
		if err != nil {
			cleanup()
			return "", nil, err
		}
		size += len(data)
		if size > 128<<20 {
			cleanup()
			return "", nil, errors.New("build source exceeds 128 MiB supported context")
		}
		target, err := files.Within(root, name)
		if err != nil {
			cleanup()
			return "", nil, err
		}
		info, err := os.Stat(source)
		if err != nil {
			cleanup()
			return "", nil, err
		}
		if err = atomicWrite(target, data, info.Mode().Perm()); err != nil {
			cleanup()
			return "", nil, err
		}
	}
	return root, cleanup, nil
}

func verifyManaged(owner project.Project, record app.Connection) error {
	root, err := files.Existing(owner.Root, ".langdrift/runtimes/"+record.RuntimeID)
	if err != nil {
		return err
	}
	if len(record.Managed) == 0 {
		return errors.New("configuration has not been materialized")
	}
	for name, expected := range record.Managed {
		target, err := files.Existing(root, name)
		if err != nil {
			return err
		}
		data, err := files.ReadBounded(target, 4<<20)
		if err != nil {
			return err
		}
		if digest(data) != expected {
			return fmt.Errorf("local configuration conflict at %s; synchronize before starting", name)
		}
	}
	return nil
}

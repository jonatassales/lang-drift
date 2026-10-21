package terraform

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/domain/microruntime"
	"github.com/exphere/ex/src/domain/project"
	"github.com/exphere/ex/src/infrastructure/process"
)

// recorder records Terraform invocations. While planning, it also captures
// the private variable file the plan reads.
type recorder struct {
	t         *testing.T
	commands  []process.Command
	variables map[string]string
	plan      []byte
	failOn    string
}

func (recorder *recorder) Run(_ context.Context, command process.Command) (process.Result, error) {
	recorder.commands = append(recorder.commands, command)
	operation := command.Args[0]

	if operation == recorder.failOn {
		return process.Result{Stderr: "Error: token=secret-provider-value"}, errors.New("exit status 1")
	}

	switch operation {
	case "plan":
		for _, argument := range command.Args {
			if path, found := strings.CutPrefix(argument, "-var-file="); found {
				data, err := os.ReadFile(path)
				if err != nil {
					recorder.t.Errorf("the variable file is unreadable during plan: %v", err)
				}

				if err := json.Unmarshal(data, &recorder.variables); err != nil {
					recorder.t.Errorf("the variable file is not JSON: %v", err)
				}
			}
		}
	case "show":
		return process.Result{Stdout: string(recorder.plan)}, nil
	}

	return process.Result{}, nil
}

func (recorder *recorder) operations() []string {
	var operations []string
	for _, command := range recorder.commands {
		operations = append(operations, command.Args[0])
	}

	return operations
}

func preparationRequest(t *testing.T) application.PreparationRequest {
	t.Helper()

	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	id, err := project.ParseID("0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatal(err)
	}

	runtime, err := microruntime.ParseIdentifier("demo")
	if err != nil {
		t.Fatal(err)
	}

	return application.PreparationRequest{
		Project:   project.Project{ID: id, Root: root},
		Runtime:   runtime,
		Directory: filepath.Join(root, ".exphere", "runtimes", "demo", "terraform"),
		Access:    application.ClusterAccess{Kubeconfig: filepath.Join(root, ".exphere", "state", "cluster.kubeconfig"), Context: "kind-ex-0123456789ab"},
		Namespace: "demo-local",
		Owner:     "team-payments",
		Token:     testToken,
	}
}

func validPlan(t *testing.T, destroy bool) []byte {
	t.Helper()

	namespace := change{kind: NamespaceResource, after: metadata("demo-local", "", testToken), actions: []string{"create"}}
	if destroy {
		namespace = change{kind: NamespaceResource, before: metadata("demo-local", "", testToken), actions: []string{"delete"}}
	}

	return plan(t, namespace)
}

func TestPreparation_Apply(t *testing.T) {
	request := preparationRequest(t)
	t.Setenv("KUBE_CONFIG_PATH", "/home/user/.kube/config")

	runner := &recorder{t: t, plan: validPlan(t, false)}
	if err := (Preparation{Runner: runner}).Apply(context.Background(), request); err != nil {
		t.Fatal(err)
	}

	if !slices.Equal(runner.operations(), []string{"init", "plan", "show", "apply"}) {
		t.Fatalf("operations = %v, want init, plan, inspection, then apply", runner.operations())
	}

	state := filepath.Join(request.Project.Root, ".exphere", "state", "runtimes", "demo")
	planFile := filepath.Join(state, "operation.tfplan")

	wantArguments := [][]string{
		{"init", "-input=false", "-no-color", "-reconfigure", "-backend-config=path=" + filepath.Join(state, "terraform.tfstate")},
		{"plan", "-input=false", "-no-color", "-out=" + planFile, "-var-file=" + filepath.Join(state, "operation.tfvars.json")},
		{"show", "-json", planFile},
		{"apply", "-input=false", "-no-color", planFile},
	}

	for index, command := range runner.commands {
		if command.Name != "terraform" || command.Dir != request.Directory || !slices.Equal(command.Args, wantArguments[index]) {
			t.Errorf("command %d = %s %v in %s, want %v in the preparation root", index, command.Name, command.Args, command.Dir, wantArguments[index])
		}

		for _, neutralized := range []string{
			"KUBECONFIG=", "KUBE_CONFIG_PATH=", "TF_CLI_ARGS=", "TF_CLI_ARGS_plan=", "TF_WORKSPACE=default",
			"TF_DATA_DIR=" + filepath.Join(state, "terraform-data"), "TF_IN_AUTOMATION=1", "TF_INPUT=0",
		} {
			if !slices.Contains(command.Env, neutralized) {
				t.Errorf("command %d environment lacks %s", index, neutralized)
			}
		}
	}

	wantVariables := map[string]string{
		VariableContext: request.Access.Context, VariableNamespace: "demo-local", VariableKubeconfig: request.Access.Kubeconfig,
		VariableOwner: "team-payments", VariableOwnershipToken: testToken,
	}
	if len(runner.variables) != len(wantVariables) {
		t.Fatalf("variables = %v, want %v", runner.variables, wantVariables)
	}

	for name, value := range wantVariables {
		if runner.variables[name] != value {
			t.Errorf("variable %s = %q, want %q", name, runner.variables[name], value)
		}
	}

	for _, transient := range []string{planFile, filepath.Join(state, "operation.tfvars.json")} {
		if _, err := os.Stat(transient); !os.IsNotExist(err) {
			t.Errorf("%s remains after the operation: %v", transient, err)
		}
	}

	if info, err := os.Stat(state); err != nil || info.Mode().Perm() != 0o700 {
		t.Errorf("state directory mode = %v, %v; want private", info.Mode().Perm(), err)
	}
}

func TestPreparation_Destroy(t *testing.T) {
	t.Parallel()

	request := preparationRequest(t)
	runner := &recorder{t: t, plan: validPlan(t, true)}

	if err := (Preparation{Runner: runner}).Destroy(context.Background(), request); err != nil {
		t.Fatal(err)
	}

	if planArguments := runner.commands[1].Args; planArguments[len(planArguments)-1] != "-destroy" {
		t.Fatalf("plan arguments = %v, want a destroy plan", planArguments)
	}

	if !slices.Equal(runner.operations(), []string{"init", "plan", "show", "apply"}) {
		t.Fatalf("operations = %v", runner.operations())
	}
}

func TestPreparation_StopsBeforeApplyingAnUnsafePlan(t *testing.T) {
	t.Parallel()

	request := preparationRequest(t)
	runner := &recorder{t: t, plan: plan(t, change{kind: NamespaceResource, after: metadata("default", "", testToken), actions: []string{"create"}})}

	if err := (Preparation{Runner: runner}).Apply(context.Background(), request); err == nil {
		t.Fatal("Apply() applied a plan for another namespace")
	}

	if slices.Contains(runner.operations(), "apply") {
		t.Fatalf("operations = %v, want no apply", runner.operations())
	}
}

func TestPreparation_WithholdsRawProviderOutput(t *testing.T) {
	t.Parallel()

	for _, step := range []string{"init", "plan", "show", "apply"} {
		t.Run(step, func(t *testing.T) {
			t.Parallel()

			request := preparationRequest(t)
			runner := &recorder{t: t, plan: validPlan(t, false), failOn: step}

			err := Preparation{Runner: runner}.Apply(context.Background(), request)
			if err == nil || !strings.Contains(err.Error(), "terraform "+step+" did not complete") || strings.Contains(err.Error(), "secret-provider-value") {
				t.Fatalf("Apply() error = %v, want the failed step without raw provider output", err)
			}

			if runner.operations()[len(runner.operations())-1] != step {
				t.Fatalf("operations = %v, want to stop at %s", runner.operations(), step)
			}
		})
	}
}

func TestPreparation_HasState(t *testing.T) {
	t.Parallel()

	request := preparationRequest(t)

	if prepared, err := (Preparation{}).HasState(request.Project, request.Runtime); err != nil || prepared {
		t.Fatalf("HasState() before preparation = %t, %v", prepared, err)
	}

	state := filepath.Join(request.Project.Root, ".exphere", "state", "runtimes", "demo", "terraform.tfstate")
	if err := os.MkdirAll(filepath.Dir(state), 0o700); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(state, []byte("{}"), 0o600); err != nil {
		t.Fatal(err)
	}

	if prepared, err := (Preparation{}).HasState(request.Project, request.Runtime); err != nil || !prepared {
		t.Fatalf("HasState() after preparation = %t, %v", prepared, err)
	}
}

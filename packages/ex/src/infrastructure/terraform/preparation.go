// Package terraform runs the consumer-installed Terraform to create and
// destroy a runtime's namespace and quota. Terraform's data, plans, and state
// stay in the project's private state, apart from authored configuration.
package terraform

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/domain/microruntime"
	"github.com/exphere/ex/src/domain/project"
	"github.com/exphere/ex/src/infrastructure/files"
	"github.com/exphere/ex/src/infrastructure/process"
	"github.com/exphere/ex/src/infrastructure/storage"
)

// Input variables every preparation root declares. EX supplies them on each
// run through a private variable file, never the command line or environment.
const (
	VariableContext        = "context"
	VariableNamespace      = "namespace"
	VariableKubeconfig     = "kubeconfig"
	VariableOwner          = "owner"
	VariableOwnershipToken = "ownership_token"
)

// Variables lists every input variable in declaration order.
var Variables = []string{
	VariableContext,
	VariableNamespace,
	VariableKubeconfig,
	VariableOwner,
	VariableOwnershipToken,
}

// Files Terraform keeps in a runtime's private state directory.
const (
	stateFile     = "terraform.tfstate"
	planFile      = "operation.tfplan"
	variablesFile = "operation.tfvars.json"
	dataDirectory = "terraform-data"
)

const stateDirectoryMode = 0o700

// commandRunner runs Terraform directly and propagates cancellation to its
// process group, including providers. Run returns after the invocation settles.
type commandRunner interface {
	Run(context.Context, process.Command) (process.Result, error)
}

// Preparation creates and destroys runtime namespaces and quotas.
type Preparation struct {
	Runner commandRunner
}

// Apply creates or converges the runtime's namespace and quota.
func (preparation Preparation) Apply(ctx context.Context, request application.PreparationRequest) error {
	return preparation.converge(ctx, request, false)
}

// Destroy removes what Apply created.
func (preparation Preparation) Destroy(ctx context.Context, request application.PreparationRequest) error {
	return preparation.converge(ctx, request, true)
}

// HasState reports whether Terraform recorded anything for the runtime.
func (Preparation) HasState(owner project.Project, runtime microruntime.Identifier) (bool, error) {
	_, err := os.Lstat(filepath.Join(storage.LayoutOf(owner).RuntimeState(runtime), stateFile))

	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, fs.ErrNotExist):
		return false, nil
	}

	return false, err
}

// converge plans, inspects the saved plan, and only then applies it. A
// failure leaves Terraform's state in place for inspection and recovery.
func (preparation Preparation) converge(ctx context.Context, request application.PreparationRequest, destroy bool) error {
	state := storage.LayoutOf(request.Project).RuntimeState(request.Runtime)
	if err := os.MkdirAll(state, stateDirectoryMode); err != nil {
		return err
	}

	run := invocation{runner: preparation.Runner, directory: request.Directory, state: state}

	if _, err := run.terraform(ctx, "init", "-input=false", "-no-color", "-reconfigure",
		"-backend-config=path="+filepath.Join(state, stateFile)); err != nil {
		return err
	}

	plan := filepath.Join(state, planFile)
	defer os.Remove(plan) // Plans may contain sensitive provider values.

	variables := filepath.Join(state, variablesFile)
	if err := files.WriteJSON(variables, inputs(request)); err != nil {
		return err
	}
	defer os.Remove(variables)

	arguments := []string{"plan", "-input=false", "-no-color", "-out=" + plan, "-var-file=" + variables}
	if destroy {
		arguments = append(arguments, "-destroy")
	}

	if _, err := run.terraform(ctx, arguments...); err != nil {
		return err
	}

	shown, err := run.terraform(ctx, "show", "-json", plan)
	if err != nil {
		return err
	}

	if err := validatePlan([]byte(shown.Stdout), request.Namespace, request.Token); err != nil {
		return err
	}

	_, err = run.terraform(ctx, "apply", "-input=false", "-no-color", plan)

	return err
}

func inputs(request application.PreparationRequest) map[string]string {
	return map[string]string{
		VariableContext:        request.Access.Context,
		VariableNamespace:      request.Namespace,
		VariableKubeconfig:     request.Access.Kubeconfig,
		VariableOwner:          request.Owner,
		VariableOwnershipToken: request.Token,
	}
}

// invocation runs Terraform for one runtime with an isolated environment.
type invocation struct {
	runner    commandRunner
	directory string
	state     string
}

// terraform runs one Terraform command. Raw output is withheld from errors
// because providers may print sensitive values.
func (invocation invocation) terraform(ctx context.Context, arguments ...string) (process.Result, error) {
	result, err := invocation.runner.Run(ctx, process.Command{
		Name: "terraform",
		Args: arguments,
		Dir:  invocation.directory,
		Env:  invocation.environment(),
	})
	if err != nil {
		return result, fmt.Errorf(
			"terraform %s did not complete; raw provider output is withheld because it may contain sensitive values: %w",
			arguments[0], err,
		)
	}

	return result, nil
}

// environment neutralizes every variable that could redirect Terraform or the
// Kubernetes provider away from the explicit runtime target.
func (invocation invocation) environment() []string {
	environment := []string{
		"TF_IN_AUTOMATION=1",
		"TF_INPUT=0",
		"TF_DATA_DIR=" + filepath.Join(invocation.state, dataDirectory),
		"TF_CLI_ARGS=",
		"TF_CLI_ARGS_init=",
		"TF_CLI_ARGS_plan=",
		"TF_CLI_ARGS_apply=",
		"TF_CLI_ARGS_show=",
		"TF_WORKSPACE=default",
		"TF_LOG=",
		"TF_LOG_PATH=",
		"TF_LOG_CORE=",
		"TF_LOG_PROVIDER=",
		"KUBECONFIG=",
		"KUBERNETES_MASTER=",
	}

	for _, entry := range os.Environ() {
		name, _, _ := strings.Cut(entry, "=")
		if strings.HasPrefix(name, "KUBE_") {
			environment = append(environment, name+"=")
		}
	}

	return environment
}

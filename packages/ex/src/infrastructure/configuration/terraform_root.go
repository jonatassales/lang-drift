package configuration

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/exphere/ex/src/domain/microruntime"
	"github.com/exphere/ex/src/infrastructure/files"
	"github.com/exphere/ex/src/infrastructure/terraform"
)

const (
	terraformSourceSuffix = ".tf.json"
	kubernetesProvider    = "kubernetes"
	kubernetesSource      = "hashicorp/kubernetes"
	kubernetesVersion     = "3.2.1"
	terraformVersions     = ">= 1.6, < 2.0"
)

// terraformBlocks are the top-level blocks a preparation root may declare.
var terraformBlocks = map[string]bool{
	"terraform": true,
	"provider":  true,
	"variable":  true,
	"resource":  true,
	"output":    true,
	"locals":    true,
}

// forbiddenFragments would let Terraform read files, data, or providers
// outside the preparation.
var forbiddenFragments = []string{"file(", "filebase64(", "templatefile(", "data.", "provider::"}

// forbiddenKeys would let Terraform run commands or reach other hosts.
var forbiddenKeys = map[string]bool{"provisioner": true, "connection": true, "exec": true}

// validateTerraformRoot accepts only a native .tf.json root that manages the
// runtime namespace and quota through one explicit Kubernetes provider bound
// to the runtime's kubeconfig and context variables. Terraform would
// otherwise fall back to an ambient target. It returns each source's digest.
func validateTerraformRoot(directory string) (map[string]string, error) {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	digests := map[string]string{}
	providers := 0

	for _, entry := range entries {
		name := entry.Name()

		switch {
		case entry.Type()&fs.ModeSymlink != 0:
			return nil, fmt.Errorf("symlinks are unsupported: %s", name)
		case entry.IsDir() && name == ".terraform":
			continue
		case entry.IsDir():
			return nil, fmt.Errorf("module directories are unsupported: %s", name)
		case strings.HasSuffix(name, ".tf"):
			return nil, fmt.Errorf("only native %s sources are accepted: %s", terraformSourceSuffix, name)
		case strings.Contains(name, ".tfvars"):
			return nil, fmt.Errorf("variable files are unsupported; runtime inputs are explicit: %s", name)
		case !strings.HasSuffix(name, terraformSourceSuffix):
			continue
		}

		data, err := files.ReadBounded(filepath.Join(directory, name), maxNativeFileBytes)
		if err != nil {
			return nil, err
		}

		declared, err := validateTerraformDocument(data)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}

		providers += declared
		digests[name] = digestOf(data)
	}

	if len(digests) == 0 {
		return nil, fmt.Errorf("no native %s preparation configuration", terraformSourceSuffix)
	}

	if providers != 1 {
		return nil, errors.New("exactly one explicit Kubernetes provider using the runtime kubeconfig and context variables is required")
	}

	return digests, nil
}

// validateTerraformDocument checks one source and returns how many provider
// configurations it declares.
func validateTerraformDocument(data []byte) (int, error) {
	var document map[string]any
	if err := decodeStrict(data, &document); err != nil {
		return 0, err
	}

	for block := range document {
		if !terraformBlocks[block] {
			return 0, fmt.Errorf("unsupported block %q", block)
		}
	}

	if containsForbidden(document) {
		return 0, errors.New("provisioners, imports, external data, and file functions are unsupported")
	}

	if err := validateSettings(asMap(document["terraform"])); err != nil {
		return 0, err
	}

	if err := validateResources(document); err != nil {
		return 0, err
	}

	return validateProviders(document)
}

func validateSettings(settings map[string]any) error {
	if settings["cloud"] != nil {
		return errors.New("remote Terraform execution is unsupported")
	}

	for backend := range asMap(settings["backend"]) {
		if backend != "local" {
			return errors.New("only local Terraform state is supported")
		}
	}

	for name, requirement := range asMap(settings["required_providers"]) {
		if name != kubernetesProvider || asMap(requirement)["source"] != kubernetesSource {
			return fmt.Errorf("only the %s provider is supported", kubernetesSource)
		}
	}

	return nil
}

func validateProviders(document map[string]any) (int, error) {
	raw, declared := document["provider"]
	if !declared {
		return 0, nil
	}

	providers, isObject := raw.(map[string]any)
	if !isObject {
		return 0, errors.New("provider blocks must use the object form")
	}

	for name, value := range providers {
		if name != kubernetesProvider {
			return 0, fmt.Errorf("unsupported provider %q", name)
		}

		configuration, isObject := value.(map[string]any)
		if !isObject {
			return 0, errors.New("the Kubernetes provider must be one unaliased configuration object")
		}

		if configuration["config_path"] != variableReference(terraform.VariableKubeconfig) ||
			configuration["config_context"] != variableReference(terraform.VariableContext) {
			return 0, errors.New("the Kubernetes provider must use the explicit kubeconfig and context variables")
		}

		for option := range configuration {
			if option != "config_path" && option != "config_context" {
				return 0, fmt.Errorf("unsupported provider option %s", option)
			}
		}
	}

	return len(providers), nil
}

func validateResources(document map[string]any) error {
	raw, declared := document["resource"]
	if !declared {
		return nil
	}

	resources, isObject := raw.(map[string]any)
	if !isObject {
		return errors.New("resources must use the object form")
	}

	for resourceType, declarations := range resources {
		if resourceType != terraform.NamespaceResource && resourceType != terraform.QuotaResource {
			return fmt.Errorf("unsupported preparation resource %s", resourceType)
		}

		named, isObject := declarations.(map[string]any)
		if !isObject {
			return fmt.Errorf("%s declarations must be named objects", resourceType)
		}

		for name, declaration := range named {
			resource, isObject := declaration.(map[string]any)
			if !isObject {
				return fmt.Errorf("resource %s.%s must be an object", resourceType, name)
			}

			// An omitted provider reference uses the one explicit default provider.
			if reference, present := resource["provider"]; present && reference != kubernetesProvider {
				return fmt.Errorf("resource %s.%s must use the default Kubernetes provider; aliases are unsupported", resourceType, name)
			}
		}
	}

	return nil
}

func containsForbidden(value any) bool {
	switch typed := value.(type) {
	case map[string]any:
		for key, nested := range typed {
			if forbiddenKeys[key] || containsForbidden(nested) {
				return true
			}
		}

	case []any:
		for _, nested := range typed {
			if containsForbidden(nested) {
				return true
			}
		}

	case string:
		for _, fragment := range forbiddenFragments {
			if strings.Contains(typed, fragment) {
				return true
			}
		}
	}

	return false
}

// generateTerraformRoot returns a preparation root that creates the runtime
// namespace and one quota, labelled with the execution's ownership token.
func generateTerraformRoot(runtime microruntime.Identifier, quotas microruntime.ResourceQuotas) map[string]any {
	variables := map[string]any{}
	for _, name := range terraform.Variables {
		variables[name] = map[string]any{"type": "string", "nullable": false}
	}

	hard := make(map[string]string, len(quotas))
	for resource, quantity := range quotas {
		hard[string(resource)] = quantity.String()
	}

	labels := map[string]any{
		microruntime.LabelRuntime:        runtime.String(),
		microruntime.LabelOwnershipToken: variableReference(terraform.VariableOwnershipToken),
	}

	return map[string]any{
		"terraform": map[string]any{
			"required_version": terraformVersions,
			"required_providers": map[string]any{
				kubernetesProvider: map[string]any{"source": kubernetesSource, "version": kubernetesVersion},
			},
			"backend": map[string]any{"local": map[string]any{}},
		},
		"variable": variables,
		"provider": map[string]any{
			kubernetesProvider: map[string]any{
				"config_path":    variableReference(terraform.VariableKubeconfig),
				"config_context": variableReference(terraform.VariableContext),
			},
		},
		"resource": map[string]any{
			terraform.NamespaceResource: map[string]any{
				"runtime": map[string]any{
					"metadata": []any{map[string]any{
						"name":        variableReference(terraform.VariableNamespace),
						"labels":      labels,
						"annotations": map[string]any{"exphere.dev/owner": variableReference(terraform.VariableOwner)},
					}},
				},
			},
			terraform.QuotaResource: map[string]any{
				"runtime": map[string]any{
					"metadata": []any{map[string]any{
						"name":      "runtime",
						"namespace": "${" + terraform.NamespaceResource + ".runtime.metadata[0].name}",
						"labels":    labels,
					}},
					"spec": []any{map[string]any{"hard": hard}},
				},
			},
		},
	}
}

func variableReference(name string) string { return "${var." + name + "}" }

func asMap(value any) map[string]any {
	object, _ := value.(map[string]any)

	return object
}

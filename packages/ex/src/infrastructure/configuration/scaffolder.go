package configuration

import (
	"fmt"
	"path"

	"github.com/exphere/ex/src/domain/microruntime"
	"github.com/exphere/ex/src/domain/project"
	"github.com/exphere/ex/src/infrastructure/scaffold"
)

const (
	readmeFileName    = "README.md"
	terraformFileName = "main.tf.json"
)

// nativeFileNames names the generated manifest file of each object kind.
var nativeFileNames = map[string]string{
	"Deployment": "workload.json",
	"Job":        "workload.json",
	"Service":    "service.json",
}

// Scaffolder publishes new runtime configuration directories.
type Scaffolder struct{}

// Scaffold writes a descriptor for the consumer to complete, with instructions.
func (Scaffolder) Scaffold(owner project.Project, descriptor microruntime.Descriptor) (string, error) {
	if err := rejectLegacyRuntime(owner.Root, descriptor.ID); err != nil {
		return "", err
	}

	encoded, err := encodeDescriptor(descriptor)
	if err != nil {
		return "", err
	}

	return scaffold.Publish(owner.Root, microruntime.ConfigurationDirectory(descriptor.ID).String(), func(stage string) error {
		if err := scaffold.WriteFile(stage, descriptorFileName, encoded); err != nil {
			return err
		}

		return scaffold.WriteFile(stage, readmeFileName, []byte(manualReadme(descriptor.ID)))
	})
}

// Generate writes the descriptor, native manifests, and Terraform root that a
// validated blueprint describes. Component sources are referenced, never copied.
func (Scaffolder) Generate(owner project.Project, runtime microruntime.Identifier, blueprint microruntime.Blueprint) (string, error) {
	if err := rejectLegacyRuntime(owner.Root, runtime); err != nil {
		return "", err
	}

	for _, workload := range blueprint.Workloads {
		if err := requireDirectory(owner.Root, workload.Component); err != nil {
			return "", fmt.Errorf("workload %s component: %w", workload.Name, err)
		}
	}

	descriptor := blueprint.Descriptor(runtime)

	encoded, err := encodeDescriptor(descriptor)
	if err != nil {
		return "", err
	}

	return scaffold.Publish(owner.Root, microruntime.ConfigurationDirectory(runtime).String(), func(stage string) error {
		for _, workload := range blueprint.Workloads {
			manifests := descriptor.Units[workload.Name].Manifests.String()

			for _, object := range workload.NativeObjects(runtime, blueprint.Target.Namespace) {
				target := path.Join(manifests, nativeFileNames[object.Kind()])
				if err := scaffold.WriteJSONFile(stage, target, object); err != nil {
					return err
				}
			}
		}

		root := path.Join(descriptor.TerraformDirectory.String(), terraformFileName)
		if err := scaffold.WriteJSONFile(stage, root, generateTerraformRoot(runtime, blueprint.Quotas)); err != nil {
			return err
		}

		if err := scaffold.WriteFile(stage, descriptorFileName, encoded); err != nil {
			return err
		}

		return scaffold.WriteFile(stage, readmeFileName, []byte(generatedReadme(runtime)))
	})
}

func manualReadme(runtime microruntime.Identifier) string {
	return "# " + runtime.String() + `

This manual scaffold is deliberately incomplete. In schema.json, set the owner,
the Kubernetes context and namespace, and at least one unit. Each unit name
must match its existing component directory; component paths are relative to
the project. Terraform and manifest references are relative to this runtime.
Add a native Terraform .tf.json namespace and quota configuration and the
units' Kubernetes manifests, then validate the runtime.

Native Kubernetes files own image, resources, replicas, probes, volumes, and
scheduling. Terraform owns the namespace and quotas, and receives the context,
namespace, kubeconfig, owner, and ownership_token input variables. Only the
hashicorp/kubernetes provider with namespace and quota resources is supported;
modules, provisioners, data sources, and remote backends are rejected. Use a
local backend; operational state belongs to .exphere/state, never here.
Quotas limit resources, not monetary spend. Application and agent sources stay
outside this runtime. Hooks are optional explicit .sh references and never run
implicitly. Only stateless Deployments, ClusterIP Services, and finite Jobs are
supported.

No migration, dependency installation, image build, or infrastructure
execution was performed.
`
}

func generatedReadme(runtime microruntime.Identifier) string {
	return "# " + runtime.String() + `

This configuration is consumer-owned and editable. Component sources are
referenced, never copied. The descriptor owns identity, references, and target
selection; Kubernetes owns workload configuration; Terraform owns the namespace
and quotas.

1. Review the native files and validate the runtime.
2. Check the local prerequisites, create the project cluster explicitly, and
   load your already-built images into it.
3. Start the runtime, inspect the result, and stop it again.

Terraform receives explicit context, namespace, kubeconfig, owner, and
ownership_token variables from the lifecycle. Do not store secrets here.
Operational state lives outside this configuration under .exphere/state and is
required for safe cleanup; deleting configuration is not cleanup. Stopping a
runtime preserves the project cluster.

Jobs are finite by default: one completion, one worker, no retries, and a
120-second deadline. Review those defaults for your use case. The image pull
policy is Never, for prebuilt local images. Generation verified neither network
access nor cluster execution. Quotas are resource limits; no financial cost or
saving is calculated.
`
}

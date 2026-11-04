// Package configuration reads and writes the runtime configuration consumers
// author under .langdrift/microruntimes/runtimes: the schema.json descriptor, native
// Kubernetes manifests, and the native Terraform preparation root.
package configuration

import (
	"encoding/json"
	"fmt"

	"github.com/jonatassales/tao/lib/cli/src/domain/microruntime"
)

// descriptorFile is the schema.json format. The domain descriptor is flatter;
// this type exists so the file format can evolve apart from the domain.
type descriptorFile struct {
	Version     int                 `json:"version"`
	ID          string              `json:"id"`
	Environment string              `json:"environment"`
	Owner       string              `json:"owner"`
	Units       map[string]unitFile `json:"units"`
	Terraform   terraformFile       `json:"terraform"`
	Kubernetes  kubernetesFile      `json:"kubernetes"`
	Hooks       map[string]string   `json:"hooks,omitempty"`
}

type unitFile struct {
	Component  string             `json:"component"`
	Kubernetes unitKubernetesFile `json:"kubernetes"`
}

type unitKubernetesFile struct {
	Manifests string `json:"manifests"`
}

type terraformFile struct {
	Directory string `json:"directory"`
}

type kubernetesFile struct {
	Context   string `json:"context"`
	Namespace string `json:"namespace"`
}

// decodeDescriptor reads schema.json strictly and maps it to the domain.
func decodeDescriptor(data []byte) (microruntime.Descriptor, error) {
	var file descriptorFile

	if err := decodeStrict(data, &file); err != nil {
		return microruntime.Descriptor{}, fmt.Errorf("%w: schema.json: %w", microruntime.ErrInvalidDescriptor, err)
	}

	if err := requireDescriptorFields(data); err != nil {
		return microruntime.Descriptor{}, fmt.Errorf("%w: schema.json: %w", microruntime.ErrInvalidDescriptor, err)
	}

	return file.descriptor()
}

// encodeDescriptor writes the domain descriptor in the schema.json format.
func encodeDescriptor(descriptor microruntime.Descriptor) ([]byte, error) {
	file := descriptorFile{
		Version:     descriptor.Version,
		ID:          descriptor.ID.String(),
		Environment: descriptor.Environment,
		Owner:       descriptor.Owner,
		Units:       make(map[string]unitFile, len(descriptor.Units)),
		Terraform:   terraformFile{Directory: descriptor.TerraformDirectory.String()},
		Kubernetes: kubernetesFile{
			Context:   descriptor.Target.Context,
			Namespace: descriptor.Target.Namespace.String(),
		},
	}

	for name, unit := range descriptor.Units {
		file.Units[name.String()] = unitFile{
			Component:  unit.Component.String(),
			Kubernetes: unitKubernetesFile{Manifests: unit.Manifests.String()},
		}
	}

	if len(descriptor.Hooks) > 0 {
		file.Hooks = make(map[string]string, len(descriptor.Hooks))
		for name, script := range descriptor.Hooks {
			file.Hooks[name.String()] = script.String()
		}
	}

	data, err := json.MarshalIndent(file, "", "  ")
	if err != nil {
		return nil, err
	}

	return append(data, '\n'), nil
}

func (file descriptorFile) descriptor() (microruntime.Descriptor, error) {
	var fields fieldParser

	descriptor := microruntime.Descriptor{
		Version:     file.Version,
		ID:          fields.identifier("id", file.ID),
		Environment: file.Environment,
		Owner:       file.Owner,
		Target: microruntime.KubernetesTarget{
			Context:   file.Kubernetes.Context,
			Namespace: fields.identifier("kubernetes.namespace", file.Kubernetes.Namespace),
		},
		TerraformDirectory: fields.path("terraform.directory", file.Terraform.Directory),
		Units:              make(map[microruntime.Identifier]microruntime.Unit, len(file.Units)),
	}

	for name, unit := range file.Units {
		descriptor.Units[fields.identifier("unit name", name)] = microruntime.Unit{
			Component: fields.componentPath("units."+name+".component", unit.Component),
			Manifests: fields.path("units."+name+".kubernetes.manifests", unit.Kubernetes.Manifests),
		}
	}

	if len(file.Hooks) > 0 {
		descriptor.Hooks = make(map[microruntime.Identifier]microruntime.RelativePath, len(file.Hooks))
		for name, script := range file.Hooks {
			descriptor.Hooks[fields.identifier("hook name", name)] = fields.path("hooks."+name, script)
		}
	}

	if fields.err != nil {
		return microruntime.Descriptor{}, fmt.Errorf("%w: schema.json: %w", microruntime.ErrInvalidDescriptor, fields.err)
	}

	return descriptor, nil
}

// fieldParser parses descriptor fields into domain values, keeping the first
// error. Empty text maps to the zero value, so the domain's own validation
// reports a missing field in its own terms.
type fieldParser struct{ err error }

func (parser *fieldParser) identifier(field, text string) microruntime.Identifier {
	if text == "" {
		return microruntime.Identifier{}
	}

	identifier, err := microruntime.ParseIdentifier(text)
	parser.keep(field, err)

	return identifier
}

func (parser *fieldParser) path(field, text string) microruntime.RelativePath {
	if text == "" {
		return microruntime.RelativePath{}
	}

	path, err := microruntime.ParseRelativePath(text)
	parser.keep(field, err)

	return path
}

func (parser *fieldParser) keep(field string, err error) {
	if err != nil && parser.err == nil {
		parser.err = fmt.Errorf("%s: %w", field, err)
	}
}

func (parser *fieldParser) componentPath(field, text string) microruntime.RelativePath {
	if text == "" {
		return microruntime.RelativePath{}
	}
	path, err := microruntime.ParseComponentPath(text)
	parser.keep(field, err)
	return path
}

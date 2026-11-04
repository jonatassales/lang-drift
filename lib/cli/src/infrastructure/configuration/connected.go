package configuration

import (
	"encoding/json"
	"path"

	"github.com/jonatassales/tao/lib/cli/src/domain/microruntime"
)

// RenderBlueprint produces the concrete files for an already validated console
// revision. Publication and conflict detection belong to the connected adapter.
func RenderBlueprint(runtime microruntime.Identifier, blueprint microruntime.Blueprint) (map[string][]byte, error) {
	if err := blueprint.Validate(); err != nil {
		return nil, err
	}
	descriptor := blueprint.Descriptor(runtime)
	encoded, err := encodeDescriptor(descriptor)
	if err != nil {
		return nil, err
	}
	result := map[string][]byte{descriptorFileName: encoded}
	put := func(name string, value any) error {
		data, err := json.MarshalIndent(value, "", "  ")
		if err != nil {
			return err
		}
		result[name] = append(data, '\n')
		return nil
	}
	for _, workload := range blueprint.Workloads {
		for _, object := range workload.NativeObjects(runtime, blueprint.Target.Namespace) {
			if err := put(path.Join(descriptor.Units[workload.Name].Manifests.String(), nativeFileNames[object.Kind()]), object); err != nil {
				return nil, err
			}
		}
	}
	if err := put(path.Join(descriptor.TerraformDirectory.String(), terraformFileName), generateTerraformRoot(runtime, blueprint.Quotas)); err != nil {
		return nil, err
	}
	return result, nil
}

// DecodeDescriptor reads a preserved migration source with the same strict
// validation used for the active descriptor.
func DecodeDescriptor(data []byte) (microruntime.Descriptor, error) { return decodeDescriptor(data) }

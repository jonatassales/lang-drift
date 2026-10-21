package configuration

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/exphere/ex/src/domain/microruntime"
)

func TestEncodeAndDecodeDescriptor(t *testing.T) {
	t.Parallel()

	descriptor := testBlueprint(t).Descriptor(identifier(t, "demo"))
	descriptor.Hooks = map[microruntime.Identifier]microruntime.RelativePath{identifier(t, "seed"): relativePath(t, "hooks/seed.sh")}

	encoded, err := encodeDescriptor(descriptor)
	if err != nil {
		t.Fatal(err)
	}

	decoded, err := decodeDescriptor(encoded)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(decoded, descriptor) {
		t.Fatalf("decodeDescriptor(encodeDescriptor(d)) = %+v, want %+v", decoded, descriptor)
	}

	withoutHooks := descriptor
	withoutHooks.Hooks = nil

	encoded, err = encodeDescriptor(withoutHooks)
	if err != nil {
		t.Fatal(err)
	}

	var document map[string]any
	if err := json.Unmarshal(encoded, &document); err != nil {
		t.Fatal(err)
	}

	if _, present := document["hooks"]; present {
		t.Error("a descriptor without hooks encoded a hooks field")
	}
}

func TestDecodeDescriptor_Refuses(t *testing.T) {
	t.Parallel()

	valid := func(t *testing.T) map[string]any {
		encoded, err := encodeDescriptor(testBlueprint(t).Descriptor(identifier(t, "demo")))
		if err != nil {
			t.Fatal(err)
		}

		var document map[string]any
		if err := json.Unmarshal(encoded, &document); err != nil {
			t.Fatal(err)
		}

		return document
	}
	unit := func(document map[string]any) map[string]any {
		return document["units"].(map[string]any)["web"].(map[string]any)
	}

	tests := []struct {
		name string
		edit func(map[string]any)
	}{
		{name: "missing owner", edit: func(d map[string]any) { delete(d, "owner") }},
		{name: "null units", edit: func(d map[string]any) { d["units"] = nil }},
		{name: "unknown top-level field", edit: func(d map[string]any) { d["cluster"] = "kind" }},
		{name: "null hooks", edit: func(d map[string]any) { d["hooks"] = nil }},
		{name: "missing Terraform directory", edit: func(d map[string]any) { d["terraform"] = map[string]any{} }},
		{name: "missing namespace", edit: func(d map[string]any) { delete(d["kubernetes"].(map[string]any), "namespace") }},
		{name: "unknown target field", edit: func(d map[string]any) { d["kubernetes"].(map[string]any)["cluster"] = "x" }},
		{name: "unit without manifests", edit: func(d map[string]any) { unit(d)["kubernetes"] = map[string]any{} }},
		{name: "unit without component", edit: func(d map[string]any) { delete(unit(d), "component") }},
		{name: "invalid runtime ID", edit: func(d map[string]any) { d["id"] = "Demo" }},
		{name: "invalid unit name", edit: func(d map[string]any) { d["units"].(map[string]any)["Web"] = unit(d) }},
		{name: "escaping component", edit: func(d map[string]any) { unit(d)["component"] = "../web" }},
		{name: "invalid hook name", edit: func(d map[string]any) { d["hooks"] = map[string]any{"Seed": "hooks/seed.sh"} }},
		{name: "escaping hook", edit: func(d map[string]any) { d["hooks"] = map[string]any{"seed": "/bin/seed.sh"} }},
		{name: "mistyped version", edit: func(d map[string]any) { d["version"] = "2" }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			document := valid(t)
			tt.edit(document)

			data, err := json.Marshal(document)
			if err != nil {
				t.Fatal(err)
			}

			if _, err := decodeDescriptor(data); !errors.Is(err, microruntime.ErrInvalidDescriptor) {
				t.Fatalf("decodeDescriptor() error = %v, want ErrInvalidDescriptor", err)
			}
		})
	}
}

package configuration

import (
	"strings"
	"testing"
)

func TestDecodeStrict(t *testing.T) {
	t.Parallel()

	type document struct {
		Name  string         `json:"name"`
		Items []any          `json:"items"`
		Extra map[string]any `json:"extra"`
	}

	tests := []struct {
		name    string
		data    string
		wantErr bool
	}{
		{name: "one document", data: `{"name": "demo", "items": [1, {"a": 1}], "extra": {"nested": {"x": 1}}}`},
		{name: "duplicate top-level key", data: `{"name": "a", "name": "b"}`, wantErr: true},
		{name: "duplicate nested key", data: `{"extra": {"nested": {"x": 1, "x": 2}}}`, wantErr: true},
		{name: "duplicate key inside an array", data: `{"items": [{"a": 1, "a": 1}]}`, wantErr: true},
		{name: "unknown field", data: `{"name": "demo", "unexpected": true}`, wantErr: true},
		{name: "second document", data: `{"name": "a"} {"name": "b"}`, wantErr: true},
		{name: "truncated document", data: `{"name": "a"`, wantErr: true},
		{name: "excessive nesting", data: `{"extra": ` + strings.Repeat(`{"x": `, 70) + `1` + strings.Repeat(`}`, 70) + `}`, wantErr: true},
		{name: "deep arrays", data: `{"items": ` + strings.Repeat(`[`, 70) + strings.Repeat(`]`, 70) + `}`, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var target document
			err := decodeStrict([]byte(tt.data), &target)

			if tt.wantErr != (err != nil) {
				t.Fatalf("decodeStrict() error = %v, want refusal %t", err, tt.wantErr)
			}
		})
	}
}

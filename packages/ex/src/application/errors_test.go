package application_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/exphere/ex/src/application"
)

func TestConditionErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		err       error
		condition error
		other     error
		message   string
	}{
		{
			name:      "not found",
			err:       application.NotFound("runtime %s has no execution record", "demo"),
			condition: application.ErrNotFound,
			other:     application.ErrAlreadyExists,
			message:   "runtime demo has no execution record",
		},
		{
			name:      "already exists",
			err:       application.AlreadyExists("project metadata already exists at %s", "/workspace"),
			condition: application.ErrAlreadyExists,
			other:     application.ErrNotFound,
			message:   "project metadata already exists at /workspace",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			wrapped := fmt.Errorf("open: %w", tt.err)

			if !errors.Is(wrapped, tt.condition) || errors.Is(wrapped, tt.other) {
				t.Errorf("errors.Is does not match only the reported condition")
			}

			if tt.err.Error() != tt.message {
				t.Errorf("Error() = %q, want %q without repeating the condition", tt.err.Error(), tt.message)
			}
		})
	}
}

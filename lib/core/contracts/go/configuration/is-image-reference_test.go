package configuration

import "testing"

func TestIsImageReference(t *testing.T) {
	acceptsTag := IsImageReference("console:local")
	refusesSpace := !IsImageReference("console local") && !IsImageReference(" ")

	if !acceptsTag || !refusesSpace {
		t.Error("IsImageReference accepts whitespace")
	}
}

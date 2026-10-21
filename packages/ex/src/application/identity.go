package application

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// newID returns an opaque 128-bit identifier in lowercase hex. Project IDs,
// ownership tokens, run IDs, and observation IDs all use it.
func newID() (string, error) {
	var random [16]byte
	if _, err := rand.Read(random[:]); err != nil {
		return "", fmt.Errorf("generate identity: %w", err)
	}

	return hex.EncodeToString(random[:]), nil
}

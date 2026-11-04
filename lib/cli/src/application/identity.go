package application

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// IDSource generates opaque identifiers. Use cases call it through id(), so a
// test injects a deterministic source and production uses the default below.
type IDSource func() (string, error)

func (source IDSource) id() (string, error) {
	if source != nil {
		return source()
	}

	return newID()
}

// Clock supplies the current time the same way.
type Clock func() time.Time

func (clock Clock) now() time.Time {
	if clock != nil {
		return clock()
	}

	return time.Now().UTC()
}

// newID returns an opaque 128-bit identifier in lowercase hex. Project IDs,
// ownership tokens, run IDs, and observation IDs all use it.
func newID() (string, error) {
	var random [16]byte
	if _, err := rand.Read(random[:]); err != nil {
		return "", fmt.Errorf("generate identity: %w", err)
	}

	return hex.EncodeToString(random[:]), nil
}

// ID is id() for the packages beside this one, so their use cases inherit the
// same deterministic injection.
func (source IDSource) ID() (string, error) { return source.id() }

// Now is now() for the packages beside this one.
func (clock Clock) Now() time.Time { return clock.now() }

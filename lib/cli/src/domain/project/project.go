// Package project is the domain model of a consumer workspace that EX manages.
package project

import (
	"errors"
	"fmt"
	"regexp"
)

// ErrInvalidID reports project metadata whose identity EX cannot trust.
var ErrInvalidID = errors.New("invalid project ID")

const clusterNamePrefixLength = 12

var idPattern = regexp.MustCompile(`^[0-9a-f]{32}$`)

// ID is the stable, opaque 128-bit identity of a project, in lowercase hex.
// Only ParseID produces a usable value; the zero ID means "not set".
type ID struct{ value string }

// ParseID validates text as a project ID.
func ParseID(text string) (ID, error) {
	if !idPattern.MatchString(text) {
		return ID{}, fmt.Errorf(
			"%w: %q must be 32 lowercase hexadecimal characters",
			ErrInvalidID, text,
		)
	}

	return ID{value: text}, nil
}

func (id ID) String() string { return id.value }

// Project is a consumer workspace: who it is and where it lives.
// Root is the workspace's canonical directory; the domain treats it as opaque.
type Project struct {
	ID   ID
	Root string
}

// ClusterName is the name of the one local cluster this project may own.
// It derives from the ID, so two projects never claim the same cluster.
// The project's ID must come from ParseID.
func (project Project) ClusterName() string {
	return "ex-" + project.ID.value[:clusterNamePrefixLength]
}

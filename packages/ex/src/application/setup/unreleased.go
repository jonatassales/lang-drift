package setup

import "strings"

// unreleased reports whether version is a build npm does not publish, which
// a workflow could not pin: an unversioned build, such as 0.0.0-dev.
func unreleased(version string) bool {
	return version == "" || strings.HasSuffix(version, "-dev")
}

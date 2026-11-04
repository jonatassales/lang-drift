package configuration

import "fmt"

// ImageTag is the image a component builds into: derived from the runtime,
// the component and the revision, never authored, so two plans of the same
// revision build the same tag.
func ImageTag(runtimeID, name string, revision int) string {
	return fmt.Sprintf("langdrift-%s-%s:r%d", runtimeID, name, revision)
}

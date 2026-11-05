package configuration

// FenceComponent is one component of a version 3 fence, as the Plan card's
// COMPONENTS rows show it: its path, stack and role, its port and resolved
// image, its proposed share of the bound and its replica count.
type FenceComponent struct {
	Path      string `json:"path"`
	Stack     string `json:"stack"`
	Role      string `json:"role"`
	Port      int    `json:"port"`
	Image     string `json:"image"`
	CPUMillis int    `json:"cpuMillis"`
	MemoryMiB int    `json:"memoryMiB"`
	Replicas  int    `json:"replicas"`
}

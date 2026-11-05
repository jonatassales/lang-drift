package configuration

// Fence is the runtime.md fence, version 3, carrying exactly what the
// simulation shows: the runtime's id and revision, the bound in millicores
// and MiB, the suggested model and effort, and the components.
type Fence struct {
	Version    int              `json:"version"`
	RuntimeID  string           `json:"runtimeId"`
	Revision   int              `json:"revision"`
	CPUMillis  int              `json:"cpuMillis"`
	MemoryMiB  int              `json:"memoryMiB"`
	Model      string           `json:"model"`
	Effort     string           `json:"effort"`
	Components []FenceComponent `json:"components"`
}

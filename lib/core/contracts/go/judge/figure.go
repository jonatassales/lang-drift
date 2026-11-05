package judge

// Figure is one figure the material gives for CPU and memory.
type Figure struct {
	Topic     string `json:"topic"`
	Fact      string `json:"fact"`
	Source    string `json:"source"`
	CPUMillis int    `json:"cpuMillis"`
	MemoryMiB int    `json:"memoryMiB"`
}

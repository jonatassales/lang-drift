package protocol

import "time"

// Observation contains measured usage only; nil values mean unavailable.
type Observation struct {
	// Nil until a consumer integration supplies measured cumulative execution cost.
	LLMCostUSD         *float64            `json:"llmCostUsd,omitempty"`
	State              string              `json:"-"`
	CPUUsed            *float64            `json:"cpuUsed"`
	MemoryUsedMiB      *float64            `json:"memoryUsedMiB"`
	ObservedAt         time.Time           `json:"observedAt"`
	MeasuredFrom       *time.Time          `json:"measuredFrom,omitempty"`
	MeasurementWindows []MeasurementWindow `json:"measurementWindows,omitempty"`
}

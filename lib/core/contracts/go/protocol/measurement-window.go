package protocol

import "time"

// MeasurementWindow preserves one native pod sample interval. Aggregates use
// the oldest sample timestamp, so a stale pod cannot make evidence look fresh.
type MeasurementWindow struct {
	Pod        string    `json:"pod"`
	ObservedAt time.Time `json:"observedAt"`
	Window     string    `json:"window"`
}

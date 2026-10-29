// Package harness holds the records of the team's delivery harness as the
// CLI mirrors them: the hardware envelope and the workflow stages, and the
// shapes every mirror shares.
package harness

import (
	"encoding/json"
	"fmt"
	"time"
)

// Timestamp is a point in time as the records write it: RFC 3339, in UTC, to
// the second, so a clock with finer resolution never changes a record's bytes.
type Timestamp struct{ time.Time }

// At makes the timestamp of a moment.
func At(moment time.Time) Timestamp {
	return Timestamp{moment.UTC().Truncate(time.Second)}
}

// MarshalJSON writes the moment as RFC 3339 in UTC.
func (stamp Timestamp) MarshalJSON() ([]byte, error) {
	return json.Marshal(stamp.UTC().Truncate(time.Second).Format(time.RFC3339))
}

// UnmarshalJSON reads an RFC 3339 moment.
func (stamp *Timestamp) UnmarshalJSON(data []byte) error {
	var text string
	if err := json.Unmarshal(data, &text); err != nil {
		return err
	}

	moment, err := time.Parse(time.RFC3339, text)
	if err != nil {
		return fmt.Errorf("timestamp %q: %w", text, err)
	}

	stamp.Time = moment.UTC()

	return nil
}

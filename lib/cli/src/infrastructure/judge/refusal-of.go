package judge

import (
	"encoding/json"
	"fmt"
)

// refusalOf reads a refused request's body as the judge route writes it.
func refusalOf(payload []byte) (routeRefusal, error) {
	var refusal routeRefusal
	if err := json.Unmarshal(payload, &refusal); err != nil {
		return routeRefusal{}, fmt.Errorf("the console's refusal is unreadable: %w", err)
	}

	return refusal, nil
}

package manifest

import (
	"fmt"
	"math"
	"time"
)

// seconds reads a Docker duration as whole seconds, rounded up, at least
// the minimum Kubernetes takes.
func seconds(value string, minimum int, field *int) error {
	duration, err := time.ParseDuration(value)
	if err != nil || duration < 0 {
		return fmt.Errorf("HEALTHCHECK duration %q is not a duration", value)
	}

	*field = max(minimum, int(math.Ceil(duration.Seconds())))

	return nil
}

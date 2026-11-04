package manifest

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jonatassales/tao/packages/ex/src/domain/microruntime"
)

// applyOption sets the probe field one HEALTHCHECK option names.
func applyOption(probe *microruntime.Probe, option string) error {
	name, value, found := strings.Cut(option, "=")
	if !found {
		return fmt.Errorf("HEALTHCHECK option %q needs =value", option)
	}

	switch name {
	case "--interval":
		return seconds(value, 1, &probe.Period)

	case "--timeout":
		return seconds(value, 1, &probe.Timeout)

	case "--start-period":
		return seconds(value, 0, &probe.InitialDelay)

	case "--start-interval":
		return nil

	case "--retries":
		retries, err := strconv.Atoi(value)
		if err != nil || retries < 1 {
			return fmt.Errorf("HEALTHCHECK --retries %q is not a positive number", value)
		}

		probe.Failures = retries

		return nil
	}

	return fmt.Errorf("HEALTHCHECK option %q is not one Docker defines", name)
}

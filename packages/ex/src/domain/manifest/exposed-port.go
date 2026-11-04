package manifest

import (
	"strconv"
	"strings"

	"github.com/jonatassales/tao/packages/ex/src/domain/microruntime"
)

// ExposedPort is the first port the Dockerfile's first EXPOSE names, or zero
// when it names none a number can read.
func ExposedPort(dockerfile []byte) int {
	for _, step := range instructions(dockerfile) {
		if step.keyword != "EXPOSE" {
			continue
		}

		fields := strings.Fields(step.arguments)
		if len(fields) == 0 {
			return 0
		}

		port, err := strconv.Atoi(strings.TrimSuffix(strings.TrimSuffix(fields[0], "/tcp"), "/udp"))
		if err != nil || !microruntime.ValidPort(port) {
			return 0
		}

		return port
	}

	return 0
}

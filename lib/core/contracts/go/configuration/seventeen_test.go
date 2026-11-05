package configuration_test

import (
	"strconv"

	"github.com/jonatassales/tao/lib/core/contracts/go/configuration"
)

// seventeen is one component more than a fence holds, each on its own path,
// port and image, with a share of one millicore and one MiB.
func seventeen(template configuration.FenceComponent) []configuration.FenceComponent {
	components := make([]configuration.FenceComponent, 0, 17)

	for index := range 17 {
		component := template
		component.Path = "services/api-" + strconv.Itoa(index)
		component.Port = 9000 + index
		component.Image = "book-a-demo-runtime/api-" + strconv.Itoa(index) + ":r1"
		component.CPUMillis = 1
		component.MemoryMiB = 1
		components = append(components, component)
	}

	return components
}

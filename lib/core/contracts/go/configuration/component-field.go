package configuration

import "strconv"

// componentField names a field of the component at an index, as a problem
// reports it: components[1].port.
func componentField(index int, name string) string {
	return "components[" + strconv.Itoa(index) + "]." + name
}

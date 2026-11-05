package configuration

// ValidPort reports whether a TCP port is usable.
func ValidPort(port int) bool {
	return port >= 1 && port <= 65535
}

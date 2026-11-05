package configuration

// The bounds of a version 3 fence: the per-execution limits a bound must
// stay within, and how many components one runtime holds.
const (
	maximumFenceCPUMillis  = 64000
	maximumFenceMemoryMiB  = 1048576
	maximumFenceComponents = 16
)

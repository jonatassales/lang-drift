package support

import "io"

// IsTerminal reports whether output goes to a terminal, the same test an
// interactive input passes.
func IsTerminal(output io.Writer) bool {
	stream, readable := output.(io.Reader)
	if !readable {
		return false
	}

	return Interactive(stream)
}

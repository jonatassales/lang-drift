package configuration

import "strings"

// IsPlatformPath reports whether a path lies inside the platform's own
// .langdrift directory, which component sources never may.
func IsPlatformPath(text string) bool {
	return text == ".langdrift" || strings.HasPrefix(text, ".langdrift/")
}

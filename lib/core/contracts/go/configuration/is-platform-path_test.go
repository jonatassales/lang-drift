package configuration

import "testing"

func TestIsPlatformPath(t *testing.T) {
	inside := IsPlatformPath(".langdrift/state")
	beside := IsPlatformPath(".langdriftx")

	if !inside || beside {
		t.Error("IsPlatformPath does not bound .langdrift exactly")
	}
}

package configuration

import "testing"

func TestValidPort(t *testing.T) {
	edgesAccepted := ValidPort(1) && ValidPort(65535)
	outsideRefused := !ValidPort(0) && !ValidPort(65536)

	if !edgesAccepted || !outsideRefused {
		t.Error("ValidPort does not bound 1 to 65535")
	}
}

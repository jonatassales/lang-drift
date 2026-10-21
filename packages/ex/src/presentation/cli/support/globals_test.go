package support_test

import (
	"context"
	"testing"
	"time"

	"github.com/exphere/ex/src/presentation/cli/support"
)

func TestGlobals_Bounded(t *testing.T) {
	t.Parallel()

	globals := support.Globals{Timeout: time.Minute}

	ctx, cancel := globals.Bounded(context.Background())
	defer cancel()

	if deadline, bounded := ctx.Deadline(); !bounded || time.Until(deadline) > time.Minute {
		t.Fatalf("Bounded() deadline = %v, %t; want at most one minute", deadline, bounded)
	}
}

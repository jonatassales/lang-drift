package connected

import "github.com/jonatassales/tao/lib/cli/src/domain/microruntime"

// applyEnvironment points an owner-edited container at the unit's
// environment Secret, when the component names one, and leaves it alone
// otherwise.
func applyEnvironment(container microruntime.Object, workload microruntime.Workload) {
	if workload.EnvSecret.IsZero() {
		return
	}

	container["envFrom"] = workload.EnvFrom()
}

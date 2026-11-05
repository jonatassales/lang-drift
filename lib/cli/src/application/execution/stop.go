package execution

import (
	"context"
	"errors"
	"fmt"

	"github.com/jonatassales/tao/lib/cli/src/application"
	"github.com/jonatassales/tao/lib/cli/src/application/inputs"
	"github.com/jonatassales/tao/lib/cli/src/application/plan"
	"github.com/jonatassales/tao/lib/cli/src/application/session"
	"github.com/jonatassales/tao/lib/cli/src/domain/harness"
	"github.com/jonatassales/tao/lib/cli/src/domain/microruntime"
	"github.com/jonatassales/tao/lib/cli/src/domain/run"
)

// StopRequest lets a waiting run go.
type StopRequest struct {
	Project application.ProjectSelection
	Runtime string
	By      string
	Session session.Sink
}

// StopResult says which lap was abandoned, by whom, what it released and
// what comes next.
type StopResult struct {
	Runtime       string `json:"runtime"`
	Lap           string `json:"lap"`
	By            string `json:"by"`
	ReleasedCents int    `json:"releasedCents"`
	Status        string `json:"status"`
}

// Stop abandons a run that waits at a stage: the runtime stops and is
// cleaned up as the end of a run does, the open lap closes as abandoned with
// who and when, the reservation is released as finops.Abandon settles it,
// and the run finishes abandoned, its evidence and memory kept, so the
// runtime can be planned again. A run that is not waiting is refused with its
// state.
func (service Service) Stop(ctx context.Context, request StopRequest) (StopResult, error) {
	owner, err := service.Projects.Open(request.Project)
	if err != nil {
		return StopResult{}, err
	}

	runtime, err := microruntime.ParseIdentifier(request.Runtime)
	if err != nil {
		return StopResult{}, fmt.Errorf("%w: %v", inputs.ErrInvalid, err)
	}

	unlock, err := service.Laps.Lock(owner, runtime)
	if err != nil {
		return StopResult{}, fmt.Errorf("%w: %s is busy: %v", ErrRefused, runtime, err)
	}
	defer unlock()

	record, err := service.Records.ReadRun(owner, runtime)
	if errors.Is(err, application.ErrNotFound) {
		return StopResult{}, application.NotFound("no such runtime: %s has no run to stop", runtime)
	}

	if err != nil {
		return StopResult{}, err
	}

	if record.State != run.StateWaiting {
		return StopResult{}, fmt.Errorf("%w: the run of %s is %s; only a run that waits at a stage can be stopped", ErrRefused, runtime, record.State)
	}

	lap, found, err := service.Laps.Open(owner, runtime)
	if err != nil {
		return StopResult{}, err
	}

	if !found {
		return StopResult{}, fmt.Errorf("%w: no lap of %s is open; nothing waits", ErrRefused, runtime)
	}

	guards, err := plan.LoadGuardrails(service.Inputs, owner)
	if err != nil {
		return StopResult{}, err
	}

	config, planned, err := service.Records.ReadPlan(owner, runtime)
	if err != nil {
		return StopResult{}, err
	}

	ledger, err := service.Records.ReadLedger(owner, runtime.String())
	if err != nil {
		return StopResult{}, err
	}

	sink := request.Session
	if sink == nil {
		sink = session.Discard{}
	}

	stopping := &playing{
		service: service,
		owner:   owner,
		runtime: runtime,
		guards:  guards,
		config:  config,
		planned: planned,
		record:  record,
		ledger:  ledger,
		lap:     lap,
		sink:    sink,
		by:      orDefault(request.By, inputs.DefaultBy),
		now:     func() harness.Timestamp { return harness.At(service.Clock.Now()) },
	}
	stopping.stopRequest = "requested by " + stopping.by + " with dft agent stop"

	emit(sink, session.StopCommand, session.Data{"runtime": runtime.String()})

	return stopping.abandonLap(ctx)
}

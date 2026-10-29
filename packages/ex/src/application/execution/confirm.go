package execution

import (
	"context"
	"fmt"

	"github.com/jonatassales/tao/packages/ex/src/application"
	"github.com/jonatassales/tao/packages/ex/src/application/inputs"
	"github.com/jonatassales/tao/packages/ex/src/application/session"
	"github.com/jonatassales/tao/packages/ex/src/domain/harness"
	"github.com/jonatassales/tao/packages/ex/src/domain/microruntime"
)

// ConfirmRequest records the owner's confirmation of a stage that waits.
type ConfirmRequest struct {
	Project application.ProjectSelection
	Runtime string
	Stage   string
	By      string
	Session session.Sink
}

// ConfirmResult says which stage was confirmed and what comes next.
type ConfirmResult struct {
	Runtime string `json:"runtime"`
	Stage   string `json:"stage"`
	By      string `json:"by"`
	Status  string `json:"status"`
}

// Confirm records a confirmation on the open lap; the next run plays on
// from the confirmed stage.
func (service Service) Confirm(_ context.Context, request ConfirmRequest) (ConfirmResult, error) {
	owner, err := service.Projects.Open(request.Project)
	if err != nil {
		return ConfirmResult{}, err
	}

	runtime, err := microruntime.ParseIdentifier(request.Runtime)
	if err != nil {
		return ConfirmResult{}, fmt.Errorf("%w: %v", inputs.ErrInvalid, err)
	}

	record, err := service.Records.ReadRun(owner, runtime)
	if err != nil {
		return ConfirmResult{}, err
	}

	lap, found, err := service.Laps.Open(owner, runtime)
	if err != nil {
		return ConfirmResult{}, err
	}
	if !found {
		return ConfirmResult{}, fmt.Errorf("%w: no lap of %s is open; nothing waits", ErrRefused, runtime)
	}

	state, known := stageOf(record, request.Stage)
	if !known {
		return ConfirmResult{}, fmt.Errorf("%w: %s has no stage %q", inputs.ErrInvalid, runtime, request.Stage)
	}
	if state.State != StageWaiting {
		return ConfirmResult{}, fmt.Errorf("%w: stage %s is %s, not waiting", ErrRefused, request.Stage, state.State)
	}

	by := orDefault(request.By, "owner")
	confirmStage(&record, &lap, request.Stage, by, harness.At(service.Clock.Now()))

	if err := service.Laps.Save(owner, runtime, lap); err != nil {
		return ConfirmResult{}, err
	}

	if err := service.Records.WriteLap(owner, runtime, lap); err != nil {
		return ConfirmResult{}, err
	}

	if err := service.Records.WriteRun(owner, runtime, record); err != nil {
		return ConfirmResult{}, err
	}

	sink := request.Session
	if sink == nil {
		sink = session.Discard{}
	}
	emit(sink, session.Confirmed, session.Data{"stage": request.Stage, "by": by})

	return ConfirmResult{Runtime: runtime.String(), Stage: request.Stage, By: by, Status: "confirmed; ex agent run " + runtime.String() + " plays on"}, nil
}

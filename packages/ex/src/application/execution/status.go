package execution

import (
	"context"
	"fmt"

	"github.com/jonatassales/tao/packages/ex/src/application"
	"github.com/jonatassales/tao/packages/ex/src/application/inputs"
	"github.com/jonatassales/tao/packages/ex/src/domain/microruntime"
	"github.com/jonatassales/tao/packages/ex/src/domain/run"
)

// StatusRequest asks where a runtime's run stands.
type StatusRequest struct {
	Project application.ProjectSelection
	Runtime string
}

// StatusResult is the run record and the lap that is open, if any.
type StatusResult struct {
	Run     run.Run `json:"run"`
	OpenLap *Lap    `json:"openLap"`
}

// Status reads the run record and the open lap.
func (service Service) Status(_ context.Context, request StatusRequest) (StatusResult, error) {
	owner, err := service.Projects.Open(request.Project)
	if err != nil {
		return StatusResult{}, err
	}

	runtime, err := microruntime.ParseIdentifier(request.Runtime)
	if err != nil {
		return StatusResult{}, fmt.Errorf("%w: %v", inputs.ErrInvalid, err)
	}

	record, err := service.Records.ReadRun(owner, runtime)
	if err != nil {
		return StatusResult{}, err
	}

	result := StatusResult{Run: record}

	open, found, err := service.Laps.Open(owner, runtime)
	if err != nil {
		return StatusResult{}, err
	}
	if found {
		result.OpenLap = &open
	}

	return result, nil
}

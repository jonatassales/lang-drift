package inputs_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/jonatassales/tao/packages/ex/src/application"
	"github.com/jonatassales/tao/packages/ex/src/application/inputs"
	"github.com/jonatassales/tao/packages/ex/src/domain/finops"
	"github.com/jonatassales/tao/packages/ex/src/domain/harness"
	"github.com/jonatassales/tao/packages/ex/src/domain/project"
)

type fakeProjects struct{}

func (fakeProjects) Open(selection application.ProjectSelection) (project.Project, error) {
	return project.Project{Root: selection.Directory}, nil
}

type fakeStore struct{ records map[string]inputs.Record }

func (store *fakeStore) Load(_ project.Project, kind string) (inputs.Record, error) {
	record, found := store.records[kind]
	if !found {
		return inputs.Record{}, application.NotFound("no %s input", kind)
	}

	return record, nil
}

func (store *fakeStore) Save(_ project.Project, record inputs.Record) error {
	if store.records == nil {
		store.records = map[string]inputs.Record{}
	}
	store.records[record.Kind] = record

	return nil
}

type fakeMirrors struct {
	written map[string]any
	next    string
}

func (mirrors *fakeMirrors) WriteMirror(_ project.Project, relative string, value any) error {
	if mirrors.written == nil {
		mirrors.written = map[string]any{}
	}
	mirrors.written[relative] = value

	return nil
}

func (mirrors *fakeMirrors) NextPriceNumber(project.Project) (string, error) { return mirrors.next, nil }

var moment = time.Date(2026, 10, 27, 18, 4, 0, 0, time.UTC)

func newService(store *fakeStore, mirrors *fakeMirrors, ids ...string) inputs.Service {
	remaining := append([]string{}, ids...)

	return inputs.Service{
		Projects: fakeProjects{},
		Store:    store,
		Mirrors:  mirrors,
		IDs: func() (string, error) {
			id := remaining[0]
			remaining = remaining[1:]

			return id, nil
		},
		Clock: func() time.Time { return moment },
	}
}

var selection = application.ProjectSelection{Directory: "/workspace", Explicit: true}

func TestService_SetEnvelopeRecordsAndMirrorsWithProvenance(t *testing.T) {
	t.Parallel()

	store, mirrors := &fakeStore{}, &fakeMirrors{}
	service := newService(store, mirrors, "7c2a0000000000000000000000000000", "9d4e0000000000000000000000000000")
	request := inputs.EnvelopeRequest{
		Project:          selection,
		PerExecution:     harness.Resources{CPUMillis: 4000, MemoryMiB: 6144},
		Default:          harness.Resources{CPUMillis: 2000, MemoryMiB: 3072},
		NotRuntimeFields: []string{"gpu", "disk"},
	}

	result, err := service.SetEnvelope(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}

	if result.Record != "envelope-7c2a" || result.Mirror != inputs.EnvelopeMirror || result.Status != "set" {
		t.Fatalf("result = %+v", result)
	}

	envelope, _ := mirrors.written[inputs.EnvelopeMirror].(harness.Envelope)
	wantMirror := harness.Mirror{Of: "cli", Record: "envelope-7c2a", SetBy: "owner", Via: "cli", SetAt: harness.At(moment), ReadOnly: true, Note: inputs.MirrorNote}
	if envelope.Mirror != wantMirror || envelope.Version != 1 || envelope.Status != "draft" || envelope.Owner != "team" || envelope.Default.CPUMillis != 2000 {
		t.Fatalf("mirror = %+v", envelope)
	}

	if record := store.records["envelope"]; record.Previous != nil || record.SetBy != "owner" || record.Via != "cli" {
		t.Fatalf("first record = %+v, want no previous one", record)
	}

	request.Default.CPUMillis = 1500
	if _, err := service.SetEnvelope(context.Background(), request); err != nil {
		t.Fatal(err)
	}

	record := store.records["envelope"]
	if record.ID != "envelope-9d4e" || record.Previous == nil || record.Previous.ID != "envelope-7c2a" {
		t.Fatalf("second record = %+v, want the first one as previous", record)
	}

	var value harness.Envelope
	if err := json.Unmarshal(record.Value, &value); err != nil || value.Default.CPUMillis != 1500 {
		t.Fatalf("record value = %s, %v", record.Value, err)
	}
}

func TestService_SetModelStoresTheChoiceWithoutAMirror(t *testing.T) {
	t.Parallel()

	store, mirrors := &fakeStore{}, &fakeMirrors{}
	service := newService(store, mirrors, "abcd0000000000000000000000000000")

	result, err := service.SetModel(context.Background(), inputs.ModelRequest{Project: selection, Model: "claude-sonnet-5-5", Effort: "medium"})
	if err != nil {
		t.Fatal(err)
	}

	if result.Mirror != "" || len(mirrors.written) != 0 {
		t.Fatalf("a choice was mirrored: %+v", result)
	}

	var choice harness.Choice
	if err := json.Unmarshal(store.records["choice"].Value, &choice); err != nil || choice.Planner != "claude-sonnet-5-5" || choice.Effort != "medium" {
		t.Fatalf("choice = %+v, %v; want the planner defaulting to the model", choice, err)
	}
}

func TestService_SetPricesTakesTheNextNumber(t *testing.T) {
	t.Parallel()

	store, mirrors := &fakeStore{}, &fakeMirrors{next: "0003"}
	service := newService(store, mirrors, "ffff0000000000000000000000000000")
	prices := finops.Prices{Unit: "cents-per-million-tokens", ReadOn: "2026-10-06", Models: []finops.ModelPrice{{ID: "claude-sonnet-5-5", Name: "Claude Sonnet 5.5", InputCents: 200, OutputCents: 1000}}}

	result, err := service.SetPrices(context.Background(), inputs.PricesRequest{Project: selection, Slug: "october", Prices: prices})
	if err != nil {
		t.Fatal(err)
	}

	if result.Record != "prices-0003" || result.Mirror != "finops/prices/0003-october.prices.json" {
		t.Fatalf("result = %+v", result)
	}

	written, _ := mirrors.written[result.Mirror].(finops.Prices)
	if written.Currency != "USD" || written.Cache != "uncached" || written.Mirror.Record != "prices-0003" {
		t.Fatalf("mirror = %+v", written)
	}
}

func TestService_RefusesInvalidInput(t *testing.T) {
	t.Parallel()

	validWorkflow := harness.Workflow{Stages: []harness.Stage{{ID: "planning", Label: "Planning", Human: "required"}, {ID: "merge", Label: "Merge", Human: "owner-only"}}}

	tests := []struct {
		name string
		call func(service inputs.Service) error
	}{
		{name: "default envelope above the ceiling", call: func(service inputs.Service) error {
			_, err := service.SetEnvelope(context.Background(), inputs.EnvelopeRequest{Project: selection, PerExecution: harness.Resources{CPUMillis: 1000, MemoryMiB: 1024}, Default: harness.Resources{CPUMillis: 2000, MemoryMiB: 512}})

			return err
		}},
		{name: "duplicate stage ids", call: func(service inputs.Service) error {
			workflow := harness.Workflow{Stages: []harness.Stage{{ID: "planning", Human: "required"}, {ID: "planning", Human: "optional"}}}
			_, err := service.SetWorkflow(context.Background(), inputs.WorkflowRequest{Project: selection, Workflow: workflow})

			return err
		}},
		{name: "unknown human marker", call: func(service inputs.Service) error {
			workflow := harness.Workflow{Stages: []harness.Stage{{ID: "planning", Human: "always"}}}
			_, err := service.SetWorkflow(context.Background(), inputs.WorkflowRequest{Project: selection, Workflow: workflow})

			return err
		}},
		{name: "a default human stage that is not a stage", call: func(service inputs.Service) error {
			workflow := validWorkflow
			workflow.Defaults = harness.Defaults{HumanAt: []string{"review"}}
			_, err := service.SetWorkflow(context.Background(), inputs.WorkflowRequest{Project: selection, Workflow: workflow})

			return err
		}},
		{name: "an effort that does not exist", call: func(service inputs.Service) error {
			_, err := service.SetModel(context.Background(), inputs.ModelRequest{Project: selection, Model: "claude-sonnet-5-5", Effort: "max"})

			return err
		}},
		{name: "allocations above the period", call: func(service inputs.Service) error {
			request := inputs.BudgetRequest{Project: selection, Period: "2026-10", Bounded: finops.BudgetBounds{PeriodCents: 1000, PerRunCents: 500, ReserveMarginPercent: 20}, Allocations: []finops.Allocation{{Ticket: "a", Cents: 600}, {Ticket: "b", Cents: 600}}}
			_, err := service.SetBudget(context.Background(), request)

			return err
		}},
		{name: "a budget period that is not a month", call: func(service inputs.Service) error {
			_, err := service.SetBudget(context.Background(), inputs.BudgetRequest{Project: selection, Period: "October"})

			return err
		}},
		{name: "prices in another unit", call: func(service inputs.Service) error {
			prices := finops.Prices{Unit: "dollars-per-million-tokens", ReadOn: "2026-10-06", Models: []finops.ModelPrice{{ID: "m"}}}
			_, err := service.SetPrices(context.Background(), inputs.PricesRequest{Project: selection, Slug: "x", Prices: prices})

			return err
		}},
		{name: "forecasts with a zero effort", call: func(service inputs.Service) error {
			forecasts := finops.Forecasts{Unit: "thousands-of-tokens", ByEffort: finops.ByEffort{Low: finops.TokensThousands{Input: 1, Output: 1}, Medium: finops.TokensThousands{Input: 1, Output: 1}}}
			_, err := service.SetForecasts(context.Background(), inputs.ForecastsRequest{Project: selection, Forecasts: forecasts})

			return err
		}},
		{name: "a driver without an executable", call: func(service inputs.Service) error {
			_, err := service.SetDriver(context.Background(), inputs.DriverRequest{Project: selection})

			return err
		}},
		{name: "gates without an action", call: func(service inputs.Service) error {
			_, err := service.SetGates(context.Background(), inputs.GatesRequest{Project: selection, Gates: []finops.Gate{{ID: "FG-01"}}})

			return err
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			store, mirrors := &fakeStore{}, &fakeMirrors{next: "0001"}
			err := tt.call(newService(store, mirrors, "abcd0000000000000000000000000000"))

			if !errors.Is(err, inputs.ErrInvalid) {
				t.Fatalf("error = %v, want ErrInvalid", err)
			}

			if len(store.records) != 0 || len(mirrors.written) != 0 {
				t.Fatal("invalid input was recorded")
			}
		})
	}
}

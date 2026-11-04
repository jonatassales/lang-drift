package application

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jonatassales/tao/apps/microservices/console-stubs/src/domain"
	"github.com/jonatassales/tao/apps/microservices/console-stubs/src/infrastructure/seed"
)

type blockingProbe struct {
	calls   atomic.Int32
	active  atomic.Int32
	ready   chan struct{}
	release chan struct{}
}

func (probe *blockingProbe) Check(ctx context.Context) error {
	probe.active.Add(1)
	defer probe.active.Add(-1)
	if probe.calls.Add(1) == 50 {
		close(probe.ready)
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-probe.release:
		return nil
	}
}
func TestLoadTesterBoundsOwnershipAndCooldown(t *testing.T) {
	probe := &blockingProbe{ready: make(chan struct{}), release: make(chan struct{})}
	tester := &LoadTester{Probe: probe}
	done := make(chan domain.LoadTestResult, 1)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	go func() { report, _ := tester.Run(ctx); done <- report }()
	select {
	case <-probe.ready:
	case <-ctx.Done():
		t.Fatal("workers did not start")
	}
	if _, err := tester.Run(ctx); !errors.Is(err, domain.ErrLoadTestBusy) {
		t.Fatal("overlapping burst accepted")
	}
	close(probe.release)
	report := <-done
	if probe.calls.Load() != 50 || probe.active.Load() != 0 || report.Succeeded != 50 || report.Failed != 0 || report.Concurrency != 50 || report.P50Ms > report.P95Ms || report.P95Ms > report.MaxMs {
		t.Fatalf("invalid report or worker lifecycle: %+v", report)
	}
	if _, err := tester.Run(ctx); !errors.Is(err, domain.ErrLoadTestBusy) {
		t.Fatal("cooldown bypassed")
	}
}
func TestCanceledLoadTestJoinsWorkersAndDoesNotPersist(t *testing.T) {
	probe := &blockingProbe{ready: make(chan struct{}), release: make(chan struct{})}
	store, err := seed.New(time.Now())
	if err != nil {
		t.Fatal(err)
	}
	service := Service{Store: store, Now: time.Now, LoadTests: &LoadTester{Probe: probe}}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() { _, err := service.RunLoadTest(ctx, "seed-runtime-1"); done <- err }()
	select {
	case <-probe.ready:
	case <-time.After(2 * time.Second):
		t.Fatal("workers did not start")
	}
	cancel()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("cancellation lost: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("workers leaked")
	}
	if probe.active.Load() != 0 || store.Snapshot().Runtimes[0].LoadTest != nil {
		t.Fatal("canceled report persisted or workers survived")
	}
}

func TestLoadDeadlineStopsEveryWorker(t *testing.T) {
	probe := &blockingProbe{ready: make(chan struct{}), release: make(chan struct{})}
	tester := &LoadTester{Probe: probe}
	started := time.Now()
	report, err := tester.Run(context.Background())
	if err != nil || report.Failed != 50 || report.Succeeded != 0 || probe.active.Load() != 0 || time.Since(started) > 4*time.Second {
		t.Fatalf("deadline failed to contain the burst: %+v, %v", report, err)
	}
}

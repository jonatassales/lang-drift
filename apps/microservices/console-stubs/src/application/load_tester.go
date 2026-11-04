package application

import (
	"context"
	"slices"
	"sync"
	"time"

	"github.com/jonatassales/tao/apps/microservices/console-stubs/src/domain"
)

// Probe is bound by composition to the local stub endpoint; callers cannot choose a target.
type Probe interface {
	Check(context.Context) error
}

// LoadTester admits one 50-attempt burst process-wide, with a cooldown after completion.
// Run owns and joins all workers; request cancellation and a three-second deadline stop I/O.
type LoadTester struct {
	Probe Probe
	mu    sync.Mutex
	busy  bool
	next  time.Time
}

func (tester *LoadTester) Run(ctx context.Context) (domain.LoadTestResult, error) {
	if err := ctx.Err(); err != nil {
		return domain.LoadTestResult{}, err
	}
	tester.mu.Lock()
	if tester.busy || time.Now().Before(tester.next) {
		tester.mu.Unlock()
		return domain.LoadTestResult{}, domain.ErrLoadTestBusy
	}
	tester.busy = true
	tester.mu.Unlock()
	defer func() {
		tester.mu.Lock()
		tester.busy = false
		tester.next = time.Now().Add(10 * time.Second)
		tester.mu.Unlock()
	}()
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	started := time.Now()
	const count = 50
	durations := make([]float64, count)
	success := make([]bool, count)
	start := make(chan struct{})
	var workers sync.WaitGroup
	for index := range count {
		workers.Add(1)
		go func() {
			defer workers.Done()
			select {
			case <-ctx.Done():
				return
			case <-start:
			}
			begin := time.Now()
			success[index] = tester.Probe.Check(ctx) == nil
			durations[index] = float64(time.Since(begin).Microseconds()) / 1000
		}()
	}
	close(start)
	workers.Wait()
	completed := time.Now()
	slices.Sort(durations)
	passed := 0
	for _, ok := range success {
		if ok {
			passed++
		}
	}
	return domain.LoadTestResult{
		Target: "extub-loopback", Requests: count, Concurrency: count,
		Succeeded: passed, Failed: count - passed,
		DurationMs: float64(completed.Sub(started).Microseconds()) / 1000,
		P50Ms:      durations[24], P95Ms: durations[47], MaxMs: durations[49],
		StartedAt: started.UTC().Format(time.RFC3339Nano), CompletedAt: completed.UTC().Format(time.RFC3339Nano),
	}, nil
}

func (service Service) RunLoadTest(ctx context.Context, id string) (domain.LoadTestResult, error) {
	catalog := service.Store.Snapshot()
	index := catalog.RuntimeIndex(id)
	if index < 0 {
		return domain.LoadTestResult{}, domain.ErrNotFound
	}
	runtime := catalog.Runtimes[index]
	if service.LoadTests == nil || runtime.DesiredState != "running" || runtime.CredentialRevoked {
		return domain.LoadTestResult{}, domain.ErrConflict
	}
	result, err := service.LoadTests.Run(ctx)
	if err != nil {
		return result, err
	}
	if err := ctx.Err(); err != nil {
		return result, err
	}
	err = service.Store.Update(func(current *domain.Catalog) error {
		index := current.RuntimeIndex(id)
		if index < 0 || current.Runtimes[index].DesiredState == "deleted" {
			return domain.ErrConflict
		}
		current.Runtimes[index].LoadTest = &result
		return nil
	})
	return result, err
}

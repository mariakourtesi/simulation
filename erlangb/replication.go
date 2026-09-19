package erlangb

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

// SimulationParams bundles the parameters shared by every replicated run.
type SimulationParams struct {
	ArrivalRate        float64
	ServiceRate        float64
	Capacity           int
	NumCallsToSimulate int
	WarmupFraction     float64
}

// SeedRun is the result of running the simulation once for a given seed.
type SeedRun struct {
	Seed    int64
	Result  RunResult
	Elapsed time.Duration
}

// ReplicationResult aggregates RunSimulation results across seeds.
type ReplicationResult struct {
	QMean         []float64
	Utilization   float64
	BlockingMean  float64
	BlockingStdev float64
	N             int
	// WallClockElapsed is how long Replicate actually took to run,
	// wall-clock time. With seeds run concurrently, this is typically
	// far less than the sum of each run's own Elapsed.
	WallClockElapsed time.Duration
}

// ProgressFunc is called as each seed finishes, for progress reporting.
// done is the 1-based count of seeds completed so far, total is len(seeds),
// and sinceStart is wall-clock time elapsed since Replicate began, useful
// for seeing the concurrency pay off as seeds land close together.
type ProgressFunc func(done, total int, run SeedRun, sinceStart time.Duration)

// Replicate runs the simulation once per seed, concurrently across
// goroutines, and aggregates the results
// workers caps how many seeds run at once; 0 defaults to runtime.NumCPU(),
// mirroring ProcessPoolExecutor's default. onProgress may be nil.
func Replicate(seeds []int64, params SimulationParams, workers int, onProgress ProgressFunc) (ReplicationResult, error) {
	if len(seeds) < 2 {
		return ReplicationResult{}, fmt.Errorf("erlangb: need at least 2 seeds to compute a standard deviation, got %d", len(seeds))
	}
	if workers <= 0 {
		workers = runtime.NumCPU()
	}

	type outcome struct {
		index int
		run   SeedRun
		err   error
	}

	overallStart := time.Now()

	sem := make(chan struct{}, workers)
	outcomes := make(chan outcome, len(seeds))
	var wg sync.WaitGroup

	for i, seed := range seeds {
		wg.Add(1)
		go func(index int, seed int64) {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			start := time.Now()
			rng := rand.New(rand.NewSource(seed))
			result, err := RunSimulation(rng, params.ArrivalRate, params.ServiceRate, params.Capacity, params.NumCallsToSimulate, params.WarmupFraction)
			elapsed := time.Since(start)

			outcomes <- outcome{index: index, run: SeedRun{Seed: seed, Result: result, Elapsed: elapsed}, err: err}
		}(i, seed)
	}

	go func() {
		wg.Wait()
		close(outcomes)
	}()

	runs := make([]RunResult, len(seeds))
	done := 0
	for o := range outcomes {
		if o.err != nil {
			return ReplicationResult{}, fmt.Errorf("erlangb: seed %d: %w", o.run.Seed, o.err)
		}
		runs[o.index] = o.run.Result
		done++
		if onProgress != nil {
			onProgress(done, len(seeds), o.run, time.Since(overallStart))
		}
	}

	capacity := params.Capacity
	qMean := make([]float64, capacity+1)
	for j := 0; j <= capacity; j++ {
		var sum float64
		for _, r := range runs {
			sum += r.Q[j]
		}
		qMean[j] = sum / float64(len(runs))
	}

	var utilSum float64
	blockingValues := make([]float64, len(runs))
	for i, r := range runs {
		utilSum += r.Utilization
		blockingValues[i] = r.CallBlocking
	}
	utilMean := utilSum / float64(len(runs))

	blockingMean, blockingStdev, err := SummaryStats(blockingValues)
	if err != nil {
		return ReplicationResult{}, err
	}

	return ReplicationResult{
		QMean:            qMean,
		Utilization:      utilMean,
		BlockingMean:     blockingMean,
		BlockingStdev:    blockingStdev,
		N:                len(seeds),
		WallClockElapsed: time.Since(overallStart),
	}, nil
}

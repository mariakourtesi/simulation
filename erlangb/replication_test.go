package erlangb

import (
	"math"
	"math/rand"
	"testing"
	"time"
)

func TestReplicate_TooFewSeeds(t *testing.T) {
	params := SimulationParams{ArrivalRate: 5, ServiceRate: 1, Capacity: 3, NumCallsToSimulate: 1000, WarmupFraction: 0.05}
	if _, err := Replicate([]int64{1}, params, 0, nil); err == nil {
		t.Error("expected an error for a single seed, got nil")
	}
	if _, err := Replicate(nil, params, 0, nil); err == nil {
		t.Error("expected an error for no seeds, got nil")
	}
}

func TestReplicate_MatchesRunningEachSeedDirectly(t *testing.T) {
	seeds := []int64{1, 2, 3, 4}
	params := SimulationParams{ArrivalRate: 5, ServiceRate: 1, Capacity: 3, NumCallsToSimulate: 20_000, WarmupFraction: 0.05}

	got, err := Replicate(seeds, params, 0, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Run each seed directly and aggregate by hand; Replicate should agree
	// exactly, since it's the same deterministic RunSimulation call, just
	// distributed across goroutines.
	var blockingValues []float64
	qSum := make([]float64, params.Capacity+1)
	var utilSum float64
	for _, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		result, err := RunSimulation(rng, params.ArrivalRate, params.ServiceRate, params.Capacity, params.NumCallsToSimulate, params.WarmupFraction)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		blockingValues = append(blockingValues, result.CallBlocking)
		utilSum += result.Utilization
		for j, q := range result.Q {
			qSum[j] += q
		}
	}
	wantBlockingMean, wantBlockingStdev, err := SummaryStats(blockingValues)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	wantUtilMean := utilSum / float64(len(seeds))

	if math.Abs(got.BlockingMean-wantBlockingMean) > 1e-12 {
		t.Errorf("BlockingMean = %v, want %v", got.BlockingMean, wantBlockingMean)
	}
	if math.Abs(got.BlockingStdev-wantBlockingStdev) > 1e-12 {
		t.Errorf("BlockingStdev = %v, want %v", got.BlockingStdev, wantBlockingStdev)
	}
	if math.Abs(got.Utilization-wantUtilMean) > 1e-12 {
		t.Errorf("Utilization = %v, want %v", got.Utilization, wantUtilMean)
	}
	for j := range qSum {
		wantQ := qSum[j] / float64(len(seeds))
		if math.Abs(got.QMean[j]-wantQ) > 1e-12 {
			t.Errorf("QMean[%d] = %v, want %v", j, got.QMean[j], wantQ)
		}
	}
	if got.N != len(seeds) {
		t.Errorf("N = %d, want %d", got.N, len(seeds))
	}
}

func TestReplicate_ProgressCallbackSeesEverySeed(t *testing.T) {
	seeds := []int64{10, 20, 30}
	params := SimulationParams{ArrivalRate: 5, ServiceRate: 1, Capacity: 3, NumCallsToSimulate: 5_000, WarmupFraction: 0.05}

	seen := make(map[int64]bool)
	var mu chan struct{} = make(chan struct{}, 1)
	mu <- struct{}{}

	onProgress := func(done, total int, run SeedRun, sinceStart time.Duration) {
		<-mu
		seen[run.Seed] = true
		mu <- struct{}{}
		if total != len(seeds) {
			t.Errorf("total = %d, want %d", total, len(seeds))
		}
		if sinceStart < 0 {
			t.Errorf("sinceStart = %v, want >= 0", sinceStart)
		}
	}

	if _, err := Replicate(seeds, params, 0, onProgress); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, s := range seeds {
		if !seen[s] {
			t.Errorf("progress callback never reported seed %d", s)
		}
	}
}

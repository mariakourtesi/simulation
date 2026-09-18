package erlangb

import (
	"math"
	"math/rand"
	"testing"
)

func TestRunSimulation_InvalidInputs(t *testing.T) {
	tests := []struct {
		name        string
		arrivalRate float64
		serviceRate float64
		capacity    int
		numCalls    int
		warmupFrac  float64
	}{
		{"zero arrival rate", 0, 1, 5, 1000, 0.05},
		{"negative arrival rate", -1, 1, 5, 1000, 0.05},
		{"zero service rate", 5, 0, 5, 1000, 0.05},
		{"zero capacity", 5, 1, 0, 1000, 0.05},
		{"negative capacity", 5, 1, -1, 1000, 0.05},
		{"zero calls", 5, 1, 5, 0, 0.05},
		{"negative warmup", 5, 1, 5, 1000, -0.1},
		{"warmup of 1", 5, 1, 5, 1000, 1.0},
	}

	rng := rand.New(rand.NewSource(1))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := RunSimulation(rng, tt.arrivalRate, tt.serviceRate, tt.capacity, tt.numCalls, tt.warmupFrac); err == nil {
				t.Errorf("expected an error, got nil")
			}
		})
	}
}

func TestRunSimulation_QSumsToOne(t *testing.T) {
	rng := rand.New(rand.NewSource(7))
	result, err := RunSimulation(rng, 5.0, 1.0, 3, 50_000, 0.05)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var sum float64
	for _, q := range result.Q {
		sum += q
	}
	if math.Abs(sum-1.0) > 1e-9 {
		t.Errorf("sum(Q) = %v, want 1.0", sum)
	}
}

func TestRunSimulation_MatchesAnalyticalBlocking(t *testing.T) {
	const (
		arrivalRate = 5.0
		serviceRate = 1.0
		capacity    = 5
		numCalls    = 300_000
	)

	rng := rand.New(rand.NewSource(42))
	result, err := RunSimulation(rng, arrivalRate, serviceRate, capacity, numCalls, 0.05)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	analytical := RecurrentErlangB(capacity, arrivalRate/serviceRate)

	// Generous tolerance: the simulated blocking probability is a sample
	// mean over ~300k counted calls, not an exact value, but it should
	// land close to the analytical Erlang-B blocking probability.
	tolerance := 0.01
	if math.Abs(result.CallBlocking-analytical) > tolerance {
		t.Errorf("simulated call blocking = %v, analytical = %v, diff exceeds tolerance %v", result.CallBlocking, analytical, tolerance)
	}

	// PASTA property: call blocking (an arrival-average) should match
	// q(capacity) (a time-average), since arrivals are Poisson.
	if math.Abs(result.CallBlocking-result.Q[capacity]) > tolerance {
		t.Errorf("call blocking = %v, q(capacity) = %v, PASTA mismatch exceeds tolerance %v", result.CallBlocking, result.Q[capacity], tolerance)
	}
}

func TestRunSimulation_DeterministicGivenSeed(t *testing.T) {
	run := func() RunResult {
		rng := rand.New(rand.NewSource(99))
		result, err := RunSimulation(rng, 5.0, 1.0, 4, 10_000, 0.05)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		return result
	}

	a := run()
	b := run()
	if a.CallBlocking != b.CallBlocking {
		t.Errorf("same seed produced different call blocking: %v vs %v", a.CallBlocking, b.CallBlocking)
	}
}

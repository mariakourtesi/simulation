package erlangb

import (
	"math"
	"math/rand"
	"testing"
)

func TestExponentialInterarrival_InvalidRate(t *testing.T) {
	rng := rand.New(rand.NewSource(1))

	for _, rate := range []float64{0, -1, -0.5} {
		if _, err := ExponentialInterarrival(rate, rng); err == nil {
			t.Errorf("ExponentialInterarrival(%v, rng) expected an error, got nil", rate)
		}
	}
}

func TestExponentialInterarrival_AlwaysPositive(t *testing.T) {
	rng := rand.New(rand.NewSource(42))

	for i := 0; i < 10_000; i++ {
		got, err := ExponentialInterarrival(5.0, rng)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got <= 0 {
			t.Fatalf("interarrival time must be > 0, got %v", got)
		}
	}
}

// TestExponentialInterarrival_MeanConvergesToOneOverRate checks the
// distribution itself: for an exponential distribution with rate lambda,
// the theoretical mean is 1/lambda. Over enough samples, the sample mean
// should land close to that value.
func TestExponentialInterarrival_MeanConvergesToOneOverRate(t *testing.T) {
	tests := []struct {
		name string
		rate float64
	}{
		{"rate=1", 1.0},
		{"rate=5", 5.0},
		{"rate=0.1", 0.1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rng := rand.New(rand.NewSource(7))

			const n = 200_000
			var sum float64
			for i := 0; i < n; i++ {
				v, err := ExponentialInterarrival(tt.rate, rng)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				sum += v
			}

			mean := sum / n
			want := 1 / tt.rate
			// Generous tolerance to keep the test stable across seeds/runs.
			tolerance := want * 0.02
			if math.Abs(mean-want) > tolerance {
				t.Errorf("sample mean = %v, want approx %v (+/- %v)", mean, want, tolerance)
			}
		})
	}
}

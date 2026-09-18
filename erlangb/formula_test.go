package erlangb

import (
	"math"
	"testing"
)

func TestRecurrentErlangB_KnownValues(t *testing.T) {
	tests := []struct {
		name        string
		capacity    int
		trafficLoad float64
		want        float64
	}{
		// B(0, A) = 1 for any load: with no servers, every call is blocked.
		{"zero capacity", 0, 5.0, 1.0},
		// B(1, A) = A / (1 + A), the single-server loss formula.
		{"one server, load 1", 1, 1.0, 0.5},
		{"one server, load 4", 1, 4.0, 0.8},
		// B(2, 2) worked out by hand from the recurrence: B1 = 2/3,
		// B2 = (2 * 2/3) / (2 + 2 * 2/3) = 0.4.
		{"two servers, load 2", 2, 2.0, 0.4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RecurrentErlangB(tt.capacity, tt.trafficLoad)
			if math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("RecurrentErlangB(%d, %v) = %v, want %v", tt.capacity, tt.trafficLoad, got, tt.want)
			}
		})
	}
}

func TestRecurrentErlangB_MonotonicInCapacity(t *testing.T) {
	// Adding servers can only reduce (or leave unchanged) the blocking
	// probability for a fixed offered load.
	trafficLoad := 5.0
	prev := RecurrentErlangB(0, trafficLoad)
	for c := 1; c <= 10; c++ {
		got := RecurrentErlangB(c, trafficLoad)
		if got > prev {
			t.Fatalf("RecurrentErlangB(%d, %v) = %v is greater than RecurrentErlangB(%d, ...) = %v", c, trafficLoad, got, c-1, prev)
		}
		prev = got
	}
}

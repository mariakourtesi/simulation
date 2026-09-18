package erlangb

import (
	"math"
	"testing"
)

func TestStateDistribution_SumsToOne(t *testing.T) {
	for _, capacity := range []int{0, 1, 5, 20} {
		q := StateDistribution(capacity, 5.0)
		var sum float64
		for _, v := range q {
			sum += v
		}
		if math.Abs(sum-1.0) > 1e-9 {
			t.Errorf("capacity=%d: sum(q) = %v, want 1.0", capacity, sum)
		}
	}
}

func TestStateDistribution_LastTermMatchesRecurrentErlangB(t *testing.T) {
	// q(capacity) is, by definition, the Erlang-B blocking probability.
	tests := []struct {
		capacity    int
		trafficLoad float64
	}{
		{1, 1.0},
		{2, 2.0},
		{5, 5.0},
		{10, 8.0},
	}

	for _, tt := range tests {
		q := StateDistribution(tt.capacity, tt.trafficLoad)
		want := RecurrentErlangB(tt.capacity, tt.trafficLoad)
		got := q[tt.capacity]
		if math.Abs(got-want) > 1e-9 {
			t.Errorf("capacity=%d load=%v: q(capacity) = %v, want %v (RecurrentErlangB)", tt.capacity, tt.trafficLoad, got, want)
		}
	}
}

func TestStateDistribution_KnownValues(t *testing.T) {
	// capacity=1, load=1: terms = [1, 1], sum = 2, q = [0.5, 0.5].
	got := StateDistribution(1, 1.0)
	want := []float64{0.5, 0.5}
	for i := range want {
		if math.Abs(got[i]-want[i]) > 1e-9 {
			t.Errorf("q[%d] = %v, want %v", i, got[i], want[i])
		}
	}
}

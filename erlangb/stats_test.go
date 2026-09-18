package erlangb

import (
	"math"
	"testing"
)

func TestSummaryStats_TooFewValues(t *testing.T) {
	if _, _, err := SummaryStats([]float64{1.0}); err == nil {
		t.Error("expected an error for a single value, got nil")
	}
	if _, _, err := SummaryStats(nil); err == nil {
		t.Error("expected an error for no values, got nil")
	}
}

func TestSummaryStats_KnownValues(t *testing.T) {
	// Mean of 2, 4, 4, 4, 5, 5, 7, 9 is 5; sample stdev (n-1) is 2.138...
	values := []float64{2, 4, 4, 4, 5, 5, 7, 9}
	mean, stdev, err := SummaryStats(values)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if math.Abs(mean-5.0) > 1e-9 {
		t.Errorf("mean = %v, want 5.0", mean)
	}
	wantStdev := 2.1380899352993950
	if math.Abs(stdev-wantStdev) > 1e-9 {
		t.Errorf("stdev = %v, want %v", stdev, wantStdev)
	}
}

func TestSummaryStats_IdenticalValuesHaveZeroStdev(t *testing.T) {
	mean, stdev, err := SummaryStats([]float64{3, 3, 3, 3})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mean != 3 {
		t.Errorf("mean = %v, want 3", mean)
	}
	if stdev != 0 {
		t.Errorf("stdev = %v, want 0", stdev)
	}
}

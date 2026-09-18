package erlangb

import (
	"fmt"
	"math"
)

// SummaryStats returns the sample mean and sample standard deviation
// (using Bessel's correction, n-1, matching Python's statistics.mean and
// statistics.stdev) of values. At least 2 values are required, a standard
// deviation is not defined for a single sample.
func SummaryStats(values []float64) (mean, stdev float64, err error) {
	if len(values) < 2 {
		return 0, 0, fmt.Errorf("erlangb: need at least 2 values to compute a standard deviation, got %d", len(values))
	}

	var sum float64
	for _, v := range values {
		sum += v
	}
	mean = sum / float64(len(values))

	var sumSquaredDiffs float64
	for _, v := range values {
		diff := v - mean
		sumSquaredDiffs += diff * diff
	}
	variance := sumSquaredDiffs / float64(len(values)-1)
	stdev = math.Sqrt(variance)

	return mean, stdev, nil
}

package erlangb

// StateDistribution computes the analytical time-stationary state
// distribution q(j) for an M/M/c/c (Erlang-B) loss system: the probability
// of exactly j busy servers, for j = 0..capacity.
//
//	q(j) = (A^j / j!) / sum_{k=0..c} (A^k / k!)
//
// Returns a slice of length capacity+1 summing to 1. q(capacity) equals the
// Erlang-B blocking probability B(capacity, offeredLoad).
func StateDistribution(capacity int, offeredLoad float64) []float64 {
	terms := make([]float64, capacity+1)
	terms[0] = 1.0 // A^0 / 0! = 1
	for j := 1; j <= capacity; j++ {
		terms[j] = terms[j-1] * offeredLoad / float64(j) // A^j / j!, built incrementally
	}

	var total float64
	for _, t := range terms {
		total += t
	}

	q := make([]float64, capacity+1)
	for i, t := range terms {
		q[i] = t / total
	}
	return q
}

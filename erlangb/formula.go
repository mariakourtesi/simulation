package erlangb

// RecurrentErlangB computes the Erlang-B blocking probability B(capacity,
// trafficLoad) for an M/M/c/c loss system, using the numerically stable
// recurrence:
//
//	B(0) = 1
//	B(c) = (trafficLoad * B(c-1)) / (c + trafficLoad * B(c-1))
//
// trafficLoad is the offered load in Erlangs (arrivalRate / serviceRate).
func RecurrentErlangB(capacity int, trafficLoad float64) float64 {
	b := 1.0
	for c := 1; c <= capacity; c++ {
		b = (trafficLoad * b) / (float64(c) + trafficLoad*b)
	}
	return b
}

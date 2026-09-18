package erlangb

import (
	"container/heap"
	"fmt"
	"math/rand"
)

// RunResult holds the results of a single Erlang-B simulation run.
type RunResult struct {
	// Q is the simulated fraction of time the system spends with exactly j
	// busy servers, for j = 0..capacity. len(Q) == capacity+1.
	Q []float64
	// CallBlocking is the fraction of arriving calls that were blocked
	// because all servers were busy.
	CallBlocking float64
	// Utilization is the time-average number of busy servers divided by
	// capacity.
	Utilization float64
}

// RunSimulation runs one Erlang-B (M/M/c/c) discrete-event simulation:
// calls arrive as a Poisson process, each is accepted if a server is free
// (and held for an exponential service time) or blocked if all servers are
// busy.
//
// The first warmupFraction of arrivals advance the system state but are
// excluded from every statistic; the loop generates enough extra arrivals
// up front so that exactly numCallsToSimulate calls are counted afterwards.
//
// rng is an injected source of randomness: seed it for reproducible runs
// and deterministic tests.
func RunSimulation(rng *rand.Rand, arrivalRate, serviceRate float64, capacity, numCallsToSimulate int, warmupFraction float64) (RunResult, error) {
	if arrivalRate <= 0 {
		return RunResult{}, fmt.Errorf("erlangb: arrivalRate must be positive, got %v", arrivalRate)
	}
	if serviceRate <= 0 {
		return RunResult{}, fmt.Errorf("erlangb: serviceRate must be positive, got %v", serviceRate)
	}
	if capacity < 1 {
		return RunResult{}, fmt.Errorf("erlangb: capacity must be >= 1, got %v", capacity)
	}
	if numCallsToSimulate <= 0 {
		return RunResult{}, fmt.Errorf("erlangb: numCallsToSimulate must be positive, got %v", numCallsToSimulate)
	}
	if warmupFraction < 0 || warmupFraction >= 1 {
		return RunResult{}, fmt.Errorf("erlangb: warmupFraction must be in [0, 1), got %v", warmupFraction)
	}

	busyServers := 0
	blockedCalls := 0
	acceptedCount := 0
	arrivalsGenerated := 0

	warmupCalls := int(warmupFraction * float64(numCallsToSimulate))
	totalArrivalsNeeded := numCallsToSimulate + warmupCalls

	timeInState := make([]float64, capacity+1)
	lastEventTime := 0.0

	events := &Events{}
	firstArrival, err := ExponentialInterarrival(arrivalRate, rng)
	if err != nil {
		return RunResult{}, err
	}
	heap.Push(events, Event{Time: firstArrival, Kind: Arrival})

	warmedUp := warmupCalls == 0

loop:
	for events.Len() > 0 {
		ev := heap.Pop(events).(Event)
		now := ev.Time

		if warmedUp {
			timeInState[busyServers] += now - lastEventTime
		}
		lastEventTime = now

		switch ev.Kind {
		case Arrival:
			arrivalsGenerated++
			if arrivalsGenerated > totalArrivalsNeeded {
				break loop
			}

			nextInterarrival, err := ExponentialInterarrival(arrivalRate, rng)
			if err != nil {
				return RunResult{}, err
			}
			heap.Push(events, Event{Time: now + nextInterarrival, Kind: Arrival})

			accepted := busyServers < capacity
			if accepted {
				busyServers++
				// Service time reuses ExponentialInterarrival with
				// serviceRate: it draws from the same exponential
				// distribution, just interpreted as a duration rather
				// than a gap between arrivals, exactly as the Python
				// version reuses expo(service_rate) for departures.
				serviceTime, err := ExponentialInterarrival(serviceRate, rng)
				if err != nil {
					return RunResult{}, err
				}
				heap.Push(events, Event{Time: now + serviceTime, Kind: Departure})
			}

			if !warmedUp && arrivalsGenerated > warmupCalls {
				warmedUp = true
			}

			if warmedUp {
				if accepted {
					acceptedCount++
				} else {
					blockedCalls++
				}
			}

		case Departure:
			busyServers--
		}
	}

	totalOfferedCalls := acceptedCount + blockedCalls
	if totalOfferedCalls == 0 {
		return RunResult{}, fmt.Errorf("erlangb: no calls were counted; check warmupFraction vs numCallsToSimulate")
	}

	var totalTime float64
	for _, ts := range timeInState {
		totalTime += ts
	}

	q := make([]float64, capacity+1)
	for i, ts := range timeInState {
		q[i] = ts / totalTime
	}

	callBlocking := float64(blockedCalls) / float64(totalOfferedCalls)

	var avgBusy float64
	for j, qj := range q {
		avgBusy += float64(j) * qj
	}
	utilization := avgBusy / float64(capacity)

	return RunResult{Q: q, CallBlocking: callBlocking, Utilization: utilization}, nil
}

package erlangb

import (
	"fmt"
	"math"
	"math/rand"
)

// ExponentialInterarrival draws a single interarrival time from an
// exponential distribution with the given arrival rate (lambda), i.e. the
// time until the next call arrives in a Poisson arrival process.
//
// It uses inverse transform sampling: if U ~ Uniform(0, 1), then
// -ln(1-U) / lambda is exponentially distributed with rate lambda.
//
// rate is the mean number of arrivals per unit time and must be > 0.
// rng is an injected source of randomness so callers can seed it for
// reproducible simulations and deterministic tests.
func ExponentialInterarrival(rate float64, rng *rand.Rand) (float64, error) {
	if rate <= 0 {
		return 0, fmt.Errorf("erlangb: rate must be positive, got %v", rate)
	}

	u := rng.Float64() // U ~ Uniform(0, 1)
	fmt.Println("what is U:", u)
	return -math.Log(1-u) / rate, nil
}

package main

import (
	"fmt"
	"math/rand"
	"simulation/erlangb"
)

func main() {
	rng := rand.New(rand.NewSource(1))
	rate := 5.0

	clock := 0.0
	count := 0

	for {
		gap, err := erlangb.ExponentialInterarrival(rate, rng)
		if err != nil {
			panic(err)
		}
		fmt.Println("GAP", gap)
		clock += gap // advance the clock by this gap

		if clock > 1.0 {
			break // this arrival falls beyond 1 minute, stop
		}

		count++
		fmt.Printf("call %d arrived at t = %.4f min (gap %.4f)\n", count, clock, gap)
	}

	fmt.Printf("\ntotal arrivals within 1 minute: %d (expected ~%.0f)\n", count, rate)
}

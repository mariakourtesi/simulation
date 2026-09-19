package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"simulation/erlangb"
)

func main() {
	arrivalRate := flag.Float64("arrival-rate", 5, "mean arrivals per unit time (lambda)")
	serviceRate := flag.Float64("service-rate", 1, "service completion rate per call (mu)")
	capacity := flag.Int("capacity", 5, "number of servers / capacity (C)")
	calls := flag.Int("calls", 10_000_000, "calls counted per run, after warm-up (warm-up calls run in addition, on top of this)")
	warmup := flag.Float64("warmup", 0.05, "fraction of --calls run first as warm-up and excluded from every statistic")
	seedsFlag := flag.String("seeds", "42,50,58,59,57,38,39,68,28,80", "comma-separated random seeds (need at least 2, to compute a standard deviation)")
	workers := flag.Int("workers", 0, "max seeds to run concurrently (0 = number of CPUs)")
	quiet := flag.Bool("quiet", false, "suppress per-seed progress logging")

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Erlang-B (M/M/c/c) loss-system simulation: runs one simulation per seed,")
		fmt.Fprintln(os.Stderr, "concurrently, and reports blocking probability against the analytical model.")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintf(os.Stderr, "Usage: %s [flags]\n\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Flags:")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Examples:")
		fmt.Fprintf(os.Stderr, "  %s\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --capacity 10 --arrival-rate 8 --service-rate 1\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --calls 1000000 --seeds 1,2,3,4,5 --quiet\n", os.Args[0])
	}
	flag.Parse()

	seeds, err := parseSeeds(*seedsFlag)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	params := erlangb.SimulationParams{
		ArrivalRate:        *arrivalRate,
		ServiceRate:        *serviceRate,
		Capacity:           *capacity,
		NumCallsToSimulate: *calls,
		WarmupFraction:     *warmup,
	}

	var onProgress erlangb.ProgressFunc
	if !*quiet {
		onProgress = func(done, total int, run erlangb.SeedRun, sinceStart time.Duration) {
			fmt.Printf("[%2d/%d] seed %3d done  blocking=%.7f  (%5.1fs this run, %6.1fs since start)\n",
				done, total, run.Seed, run.Result.CallBlocking, run.Elapsed.Seconds(), sinceStart.Seconds())
		}
	}

	results, err := erlangb.Replicate(seeds, params, *workers, onProgress)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	fmt.Println()
	erlangb.Report(results, *arrivalRate, *serviceRate, *capacity)
}

// parseSeeds turns a comma-separated flag value like "42,50,58" into seeds.
func parseSeeds(s string) ([]int64, error) {
	parts := strings.Split(s, ",")
	seeds := make([]int64, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		v, err := strconv.ParseInt(p, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid seed %q: %w", p, err)
		}
		seeds = append(seeds, v)
	}
	return seeds, nil
}

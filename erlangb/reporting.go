package erlangb

import (
	"fmt"
	"strings"
)

func row(left, right string, lw, rw int) string {
	return fmt.Sprintf("| %-*s | %*s |", lw, left, rw, right)
}

func rule(lw, rw int) string {
	return "+" + strings.Repeat("-", lw+2) + "+" + strings.Repeat("-", rw+2) + "+"
}

// Report prints the same three comparison tables as the Python version:
// the q(j) state distribution (simulated vs analytical), the PASTA
// property check, and summary statistics across seeds.
func Report(results ReplicationResult, arrivalRate, serviceRate float64, capacity int) {
	offeredLoad := arrivalRate / serviceRate
	analyticalQ := StateDistribution(capacity, offeredLoad)
	analyticalBlocking := RecurrentErlangB(capacity, offeredLoad)
	qMean := results.QMean

	// ---- q(j) state distribution table (three columns) ----
	lw, cw := 8, 14
	top := "+" + strings.Repeat("-", lw+2) + "+" + strings.Repeat("-", cw+2) + "+" + strings.Repeat("-", cw+2) + "+"
	fmt.Println("q(j): fraction of TIME the system spends with j servers busy")
	fmt.Println(top)
	fmt.Printf("| %-*s | %*s | %*s |\n", lw, "state j", cw, "simulated", cw, "analytical")
	fmt.Println(top)

	var qMeanSum, analyticalQSum float64
	for j := 0; j <= capacity; j++ {
		fmt.Printf("| %-*d | %*.7f | %*.7f |\n", lw, j, cw, qMean[j], cw, analyticalQ[j])
		qMeanSum += qMean[j]
		analyticalQSum += analyticalQ[j]
	}
	fmt.Println(top)
	fmt.Printf("| %-*s | %*.7f | %*.7f |\n", lw, "sum", cw, qMeanSum, cw, analyticalQSum)
	fmt.Println(top)
	fmt.Println()

	// ---- PASTA property table ----
	lw, rw := 40, 11
	fmt.Println("PASTA property: call blocking (arrival-avg) == q(C) (time-avg)")
	fmt.Println(rule(lw, rw))
	fmt.Println(row("quantity", "value", lw, rw))
	fmt.Println(rule(lw, rw))
	fmt.Println(row("call blocking (fraction of arrivals)", fmt.Sprintf("%.7f", results.BlockingMean), lw, rw))
	fmt.Println(row(fmt.Sprintf("q(%d) (fraction of time system full)", capacity), fmt.Sprintf("%.7f", qMean[capacity]), lw, rw))
	fmt.Println(row(fmt.Sprintf("analytical Erlang-B B(%d,%g)", capacity, offeredLoad), fmt.Sprintf("%.7f", analyticalBlocking), lw, rw))
	fmt.Println(rule(lw, rw))
	fmt.Println()

	// ---- summary statistics table ----
	lw, rw = 28, 24
	fmt.Printf("Summary across %d seeds\n", results.N)
	fmt.Println(rule(lw, rw))
	fmt.Println(row("metric", "value", lw, rw))
	fmt.Println(rule(lw, rw))
	fmt.Println(row("utilization (avg busy / C)", fmt.Sprintf("%.7f", results.Utilization), lw, rw))
	fmt.Println(row("mean call blocking", fmt.Sprintf("%.7f", results.BlockingMean), lw, rw))
	fmt.Println(row("standard deviation", fmt.Sprintf("%.7f", results.BlockingStdev), lw, rw))
	fmt.Println(row("wall-clock time (this run)", fmt.Sprintf("%.3fs", results.WallClockElapsed.Seconds()), lw, rw))
	fmt.Println(rule(lw, rw))
}

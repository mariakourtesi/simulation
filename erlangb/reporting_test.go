package erlangb

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestRow_PadsLeftAndRight(t *testing.T) {
	got := row("a", "b", 5, 3)
	want := "| a     |   b |"
	if got != want {
		t.Errorf("row(...) = %q, want %q", got, want)
	}
}

func TestRule_MatchesColumnWidths(t *testing.T) {
	got := rule(5, 3)
	want := "+-------+-----+"
	if got != want {
		t.Errorf("rule(...) = %q, want %q", got, want)
	}
}

func TestReport_ProducesExpectedSections(t *testing.T) {
	results := ReplicationResult{
		QMean:         []float64{0.1, 0.2, 0.3, 0.4},
		Utilization:   0.6,
		BlockingMean:  0.4,
		BlockingStdev: 0.01,
		N:             4,
	}

	output := captureStdout(t, func() {
		Report(results, 5, 1, 3)
	})

	for _, want := range []string{
		"q(j): fraction of TIME",
		"PASTA property",
		"Summary across 4 seeds",
	} {
		if !strings.Contains(output, want) {
			t.Errorf("report output missing %q\n---\n%s", want, output)
		}
	}
}

// captureStdout redirects os.Stdout for the duration of fn and returns
// everything written to it.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	original := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = original

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("io.Copy: %v", err)
	}
	return buf.String()
}

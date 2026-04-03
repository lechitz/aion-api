package main

import (
	"testing"
	"time"
)

func TestPercentileDuration(t *testing.T) {
	t.Parallel()

	sorted := []time.Duration{
		10 * time.Millisecond,
		20 * time.Millisecond,
		30 * time.Millisecond,
		40 * time.Millisecond,
		50 * time.Millisecond,
	}

	if got := percentileDuration(sorted, 0.50); got != 30*time.Millisecond {
		t.Fatalf("p50 = %s, want 30ms", got)
	}
	if got := percentileDuration(sorted, 0.95); got != 40*time.Millisecond {
		t.Fatalf("p95 = %s, want 40ms", got)
	}
}

func TestEvaluateThresholds(t *testing.T) {
	t.Parallel()

	result := scenarioResult{
		errorRatePct: 0,
		p50:          80 * time.Millisecond,
		p95:          150 * time.Millisecond,
	}

	if err := evaluateThresholds(result, threshold{
		MaxErrorRatePct: 0,
		MaxP50Ms:        100,
		MaxP95Ms:        200,
	}); err != nil {
		t.Fatalf("evaluateThresholds() returned error: %v", err)
	}

	if err := evaluateThresholds(result, threshold{
		MaxErrorRatePct: 0,
		MaxP50Ms:        50,
		MaxP95Ms:        120,
	}); err == nil {
		t.Fatal("expected threshold evaluation to fail")
	}
}

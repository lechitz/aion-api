package main

import (
	"testing"
	"time"
)

func TestBuildPayloadVariesByNonce(t *testing.T) {
	t.Parallel()

	first, err := buildPayload(time.Date(2026, time.March, 28, 16, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("buildPayload(first) returned error: %v", err)
	}

	second, err := buildPayload(time.Date(2026, time.March, 28, 16, 0, 1, 0, time.UTC))
	if err != nil {
		t.Fatalf("buildPayload(second) returned error: %v", err)
	}

	if string(first) == string(second) {
		t.Fatalf("expected payloads to differ when nonce source time differs")
	}
}


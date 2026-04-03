package main

import (
	"testing"
	"time"
)

func TestLoadConfigDefaults(t *testing.T) {
	t.Parallel()

	cfg := loadConfig(nil)

	if cfg.scenario != defaultScenario {
		t.Fatalf("scenario = %q, want %q", cfg.scenario, defaultScenario)
	}
	if cfg.requests != defaultRequests {
		t.Fatalf("requests = %d, want %d", cfg.requests, defaultRequests)
	}
	if cfg.concurrency != defaultConcurrency {
		t.Fatalf("concurrency = %d, want %d", cfg.concurrency, defaultConcurrency)
	}
	if cfg.timeout != defaultTimeout {
		t.Fatalf("timeout = %s, want %s", cfg.timeout, defaultTimeout)
	}
}

func TestLoadConfigOverrides(t *testing.T) {
	t.Parallel()

	cfg := loadConfig([]string{
		"--scenario", "record-projections-latest",
		"--requests", "80",
		"--concurrency", "8",
		"--warmup", "4",
		"--timeout", "3s",
		"--tag-id", "33",
		"--record-limit", "15",
		"--dashboard-date", "2026-03-28",
		"--dashboard-timezone", "America/Sao_Paulo",
	})

	if cfg.scenario != "record-projections-latest" {
		t.Fatalf("scenario = %q", cfg.scenario)
	}
	if cfg.requests != 80 || cfg.concurrency != 8 || cfg.warmupRequests != 4 || cfg.recordLimit != 15 {
		t.Fatalf("unexpected numeric overrides: %+v", cfg)
	}
	if cfg.timeout != 3*time.Second {
		t.Fatalf("timeout = %s, want 3s", cfg.timeout)
	}
	if cfg.tagID != "33" {
		t.Fatalf("tagID = %q, want 33", cfg.tagID)
	}
	if cfg.dashboardDate != "2026-03-28" || cfg.dashboardTZ != "America/Sao_Paulo" {
		t.Fatalf("unexpected dashboard overrides: %+v", cfg)
	}
}

func TestLoadConfigRealtimeScenarioDefaults(t *testing.T) {
	t.Parallel()

	cfg := loadConfig([]string{"--scenario", realtimeRecordScenario})

	if cfg.requests != defaultRealtimeRequests {
		t.Fatalf("requests = %d, want %d", cfg.requests, defaultRealtimeRequests)
	}
	if cfg.concurrency != defaultRealtimeConcurrency {
		t.Fatalf("concurrency = %d, want %d", cfg.concurrency, defaultRealtimeConcurrency)
	}
	if cfg.warmupRequests != defaultRealtimeWarmup {
		t.Fatalf("warmupRequests = %d, want %d", cfg.warmupRequests, defaultRealtimeWarmup)
	}
	if cfg.timeout != defaultRealtimeTimeout {
		t.Fatalf("timeout = %s, want %s", cfg.timeout, defaultRealtimeTimeout)
	}
}

func TestLoadConfigRealtimeScenarioExplicitOverridesWin(t *testing.T) {
	t.Parallel()

	cfg := loadConfig([]string{
		"--scenario", realtimeRecordScenario,
		"--requests", "12",
		"--concurrency", "3",
		"--warmup", "1",
		"--timeout", "12s",
	})

	if cfg.requests != 12 || cfg.concurrency != 3 || cfg.warmupRequests != 1 {
		t.Fatalf("unexpected realtime overrides: %+v", cfg)
	}
	if cfg.timeout != 12*time.Second {
		t.Fatalf("timeout = %s, want 12s", cfg.timeout)
	}
}

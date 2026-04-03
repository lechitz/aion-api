// Package main provides the versioned local load-test entrypoints for aion-api.
package main

import (
	"flag"
	"time"
)

const (
	defaultScenario            = "auth-login"
	defaultBaseURL             = "http://localhost:5001"
	defaultRequests            = 60
	defaultConcurrency         = 6
	defaultWarmupRequests      = 5
	defaultTimeout             = 5 * time.Second
	defaultRealtimeRequests    = 20
	defaultRealtimeConcurrency = 4
	defaultRealtimeWarmup      = 2
	defaultRealtimeTimeout     = 30 * time.Second
	defaultUsername            = "testuser"
	defaultPassword            = "Test@123"
	defaultTagID               = "24"
	defaultRecordLimit         = 20
	defaultDashboardTZ         = "UTC"
	defaultThresholdsFile      = "hack/tools/load-test/thresholds.json"
)

type config struct {
	scenario       string
	baseURL        string
	requests       int
	concurrency    int
	warmupRequests int
	timeout        time.Duration
	username       string
	password       string
	tagID          string
	recordLimit    int
	dashboardDate  string
	dashboardTZ    string
	thresholdsFile string
}

func loadConfig(args []string) config {
	fs := flag.NewFlagSet("load-test", flag.ContinueOnError)
	cfg := config{}

	fs.StringVar(&cfg.scenario, "scenario", defaultScenario, "load-test scenario to run")
	fs.StringVar(&cfg.baseURL, "base-url", defaultBaseURL, "aion-api base URL")
	fs.IntVar(&cfg.requests, "requests", defaultRequests, "number of measured requests")
	fs.IntVar(&cfg.concurrency, "concurrency", defaultConcurrency, "number of concurrent workers")
	fs.IntVar(&cfg.warmupRequests, "warmup", defaultWarmupRequests, "number of warmup requests before measurement")
	fs.DurationVar(&cfg.timeout, "timeout", defaultTimeout, "per-request timeout")
	fs.StringVar(&cfg.username, "username", defaultUsername, "seeded username for authenticated scenarios")
	fs.StringVar(&cfg.password, "password", defaultPassword, "seeded password for authenticated scenarios")
	fs.StringVar(&cfg.tagID, "tag-id", defaultTagID, "tag id used by record-creating scenarios")
	fs.IntVar(&cfg.recordLimit, "record-limit", defaultRecordLimit, "recordProjectionsLatest limit for GraphQL scenario")
	fs.StringVar(&cfg.dashboardDate, "dashboard-date", "", "dashboardSnapshot date in YYYY-MM-DD; defaults to current UTC date")
	fs.StringVar(&cfg.dashboardTZ, "dashboard-timezone", defaultDashboardTZ, "dashboardSnapshot timezone")
	fs.StringVar(&cfg.thresholdsFile, "thresholds-file", defaultThresholdsFile, "path to committed scenario thresholds")

	_ = fs.Parse(args)
	return applyScenarioDefaults(cfg)
}

func applyScenarioDefaults(cfg config) config {
	if cfg.scenario != realtimeRecordScenario {
		return cfg
	}

	if cfg.requests == defaultRequests {
		cfg.requests = defaultRealtimeRequests
	}
	if cfg.concurrency == defaultConcurrency {
		cfg.concurrency = defaultRealtimeConcurrency
	}
	if cfg.warmupRequests == defaultWarmupRequests {
		cfg.warmupRequests = defaultRealtimeWarmup
	}
	if cfg.timeout == defaultTimeout {
		cfg.timeout = defaultRealtimeTimeout
	}

	return cfg
}

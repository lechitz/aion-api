package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	authLoginScenario         = "auth-login"
	dashboardSnapshotScenario = "dashboard-snapshot"
	recordProjectionsScenario = "record-projections-latest"
	realtimeRecordScenario    = "realtime-record-created"
	authLoginPath             = "/aion/api/v1/auth/login"
	graphqlPath               = "/aion/api/v1/graphql"
	realtimeStreamPath        = "/aion/api/v1/realtime/events/stream"
	dashboardSnapshotQuery    = "contracts/graphql/queries/dashboard/snapshot.graphql"
	recordProjectionsQuery    = "contracts/graphql/queries/records/projections-latest.graphql"
	contentTypeJSON           = "application/json"
)

type threshold struct {
	MaxErrorRatePct float64 `json:"max_error_rate_pct"`
	MaxP50Ms        float64 `json:"max_p50_ms"`
	MaxP95Ms        float64 `json:"max_p95_ms"`
}

type thresholdsFile map[string]threshold

type scenarioDefinition struct {
	name    string
	prepare func(context.Context, *http.Client, config) (*scenarioState, error)
	execute func(context.Context, *http.Client, config, *scenarioState) error
}

type scenarioState struct {
	token         string
	query         string
	dashboardDate string
}

type loginEnvelope struct {
	Result struct {
		Token string `json:"token"`
	} `json:"result"`
}

type graphqlEnvelope struct {
	Data   map[string]json.RawMessage `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

type createRecordResult struct {
	ID string `json:"id"`
}

type sseEnvelope struct {
	Event string
	Data  string
}

type projectionChangedPayload struct {
	RecordID interface{} `json:"recordId"`
	Action   string      `json:"action"`
}

type sample struct {
	duration time.Duration
	err      error
}

type scenarioResult struct {
	name         string
	requests     int
	concurrency  int
	warmup       int
	successes    int
	failures     int
	totalRuntime time.Duration
	p50          time.Duration
	p95          time.Duration
	errorRatePct float64
	rps          float64
}

func run(ctx context.Context, cfg config) error {
	def, err := resolveScenario(cfg.scenario)
	if err != nil {
		return err
	}

	thresholds, err := loadThresholds(cfg.thresholdsFile)
	if err != nil {
		return err
	}

	target, ok := thresholds[def.name]
	if !ok {
		return fmt.Errorf("thresholds missing for scenario %q", def.name)
	}

	client := &http.Client{Timeout: cfg.timeout}

	state, err := def.prepare(ctx, client, cfg)
	if err != nil {
		return err
	}
	if err := warmupScenario(ctx, client, cfg, def, state); err != nil {
		return err
	}

	result, err := executeScenario(ctx, client, cfg, def, state)
	if err != nil {
		return err
	}

	printResult(result, target)
	return evaluateThresholds(result, target)
}

func resolveScenario(name string) (scenarioDefinition, error) {
	switch name {
	case authLoginScenario:
		return scenarioDefinition{
			name: authLoginScenario,
			prepare: func(context.Context, *http.Client, config) (*scenarioState, error) {
				return &scenarioState{}, nil
			},
			execute: func(ctx context.Context, client *http.Client, cfg config, _ *scenarioState) error {
				_, err := login(ctx, client, cfg)
				return err
			},
		}, nil
	case recordProjectionsScenario:
		return scenarioDefinition{
			name: recordProjectionsScenario,
			prepare: func(ctx context.Context, client *http.Client, cfg config) (*scenarioState, error) {
				token, err := login(ctx, client, cfg)
				if err != nil {
					return nil, err
				}
				query, err := os.ReadFile(recordProjectionsQuery)
				if err != nil {
					return nil, err
				}
				return &scenarioState{
					token: token,
					query: strings.TrimSpace(string(query)),
				}, nil
			},
			execute: fetchRecordProjectionsLatest,
		}, nil
	case dashboardSnapshotScenario:
		return scenarioDefinition{
			name: dashboardSnapshotScenario,
			prepare: func(ctx context.Context, client *http.Client, cfg config) (*scenarioState, error) {
				token, err := login(ctx, client, cfg)
				if err != nil {
					return nil, err
				}
				query, err := os.ReadFile(dashboardSnapshotQuery)
				if err != nil {
					return nil, err
				}

				dashboardDate := strings.TrimSpace(cfg.dashboardDate)
				if dashboardDate == "" {
					dashboardDate = time.Now().UTC().Format("2006-01-02")
				}

				return &scenarioState{
					token:         token,
					query:         strings.TrimSpace(string(query)),
					dashboardDate: dashboardDate,
				}, nil
			},
			execute: fetchDashboardSnapshot,
		}, nil
	case realtimeRecordScenario:
		return scenarioDefinition{
			name: realtimeRecordScenario,
			prepare: func(ctx context.Context, client *http.Client, cfg config) (*scenarioState, error) {
				token, err := login(ctx, client, cfg)
				if err != nil {
					return nil, err
				}
				return &scenarioState{token: token}, nil
			},
			execute: runRealtimeRecordScenario,
		}, nil
	default:
		return scenarioDefinition{}, fmt.Errorf("unsupported scenario %q", name)
	}
}

func loadThresholds(path string) (thresholdsFile, error) {
	if path != defaultThresholdsFile {
		return nil, fmt.Errorf("unsupported thresholds file %q", path)
	}

	raw, err := os.ReadFile(defaultThresholdsFile)
	if err != nil {
		return nil, err
	}

	var thresholds thresholdsFile
	if err := json.Unmarshal(raw, &thresholds); err != nil {
		return nil, err
	}
	return thresholds, nil
}

func warmupScenario(ctx context.Context, client *http.Client, cfg config, def scenarioDefinition, state *scenarioState) error {
	for i := range cfg.warmupRequests {
		if err := def.execute(ctx, client, cfg, state); err != nil {
			return fmt.Errorf("warmup failed on request %d: %w", i+1, err)
		}
	}
	return nil
}

func executeScenario(ctx context.Context, client *http.Client, cfg config, def scenarioDefinition, state *scenarioState) (scenarioResult, error) {
	start := time.Now()
	samples := make([]sample, 0, cfg.requests)
	jobs := make(chan struct{}, cfg.requests)
	results := make(chan sample, cfg.requests)

	for range cfg.requests {
		jobs <- struct{}{}
	}
	close(jobs)

	workerCount := cfg.concurrency
	if workerCount < 1 {
		workerCount = 1
	}

	var wg sync.WaitGroup
	for range workerCount {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range jobs {
				reqStart := time.Now()
				err := def.execute(ctx, client, cfg, state)
				results <- sample{
					duration: time.Since(reqStart),
					err:      err,
				}
			}
		}()
	}

	wg.Wait()
	close(results)

	for result := range results {
		samples = append(samples, result)
	}

	outcome := summarize(def.name, cfg, time.Since(start), samples)
	if outcome.requests == 0 {
		return outcome, errors.New("no load-test samples collected")
	}
	return outcome, nil
}

func summarize(name string, cfg config, totalRuntime time.Duration, samples []sample) scenarioResult {
	durations := make([]time.Duration, 0, len(samples))
	successes := 0
	failures := 0
	for _, sample := range samples {
		if sample.err == nil {
			successes++
		} else {
			failures++
		}
		durations = append(durations, sample.duration)
	}

	sort.Slice(durations, func(i int, j int) bool {
		return durations[i] < durations[j]
	})

	result := scenarioResult{
		name:         name,
		requests:     len(samples),
		concurrency:  cfg.concurrency,
		warmup:       cfg.warmupRequests,
		successes:    successes,
		failures:     failures,
		totalRuntime: totalRuntime,
	}
	if result.requests > 0 {
		result.errorRatePct = (float64(failures) / float64(result.requests)) * 100
	}
	if totalRuntime > 0 {
		result.rps = float64(successes) / totalRuntime.Seconds()
	}
	if len(durations) > 0 {
		result.p50 = percentileDuration(durations, 0.50)
		result.p95 = percentileDuration(durations, 0.95)
	}

	return result
}

func percentileDuration(sorted []time.Duration, percentile float64) time.Duration {
	if len(sorted) == 0 {
		return 0
	}
	if len(sorted) == 1 {
		return sorted[0]
	}

	index := int(percentile * float64(len(sorted)-1))
	if index < 0 {
		index = 0
	}
	if index >= len(sorted) {
		index = len(sorted) - 1
	}
	return sorted[index]
}

func printResult(result scenarioResult, target threshold) {
	_, _ = fmt.Fprintf(
		os.Stdout,
		"scenario=%s requests=%d concurrency=%d warmup=%d successes=%d failures=%d error_rate=%.2f%% rps=%.2f p50_ms=%.2f p95_ms=%.2f thresholds={max_error_rate_pct=%.2f,max_p50_ms=%.2f,max_p95_ms=%.2f}\n",
		result.name,
		result.requests,
		result.concurrency,
		result.warmup,
		result.successes,
		result.failures,
		result.errorRatePct,
		result.rps,
		durationToMS(result.p50),
		durationToMS(result.p95),
		target.MaxErrorRatePct,
		target.MaxP50Ms,
		target.MaxP95Ms,
	)
}

func evaluateThresholds(result scenarioResult, target threshold) error {
	var failed []string
	if result.errorRatePct > target.MaxErrorRatePct {
		failed = append(failed, fmt.Sprintf("error_rate %.2f%% > %.2f%%", result.errorRatePct, target.MaxErrorRatePct))
	}
	if durationToMS(result.p50) > target.MaxP50Ms {
		failed = append(failed, fmt.Sprintf("p50 %.2fms > %.2fms", durationToMS(result.p50), target.MaxP50Ms))
	}
	if durationToMS(result.p95) > target.MaxP95Ms {
		failed = append(failed, fmt.Sprintf("p95 %.2fms > %.2fms", durationToMS(result.p95), target.MaxP95Ms))
	}
	if len(failed) > 0 {
		return errors.New(strings.Join(failed, "; "))
	}
	return nil
}

func durationToMS(value time.Duration) float64 {
	return float64(value) / float64(time.Millisecond)
}

func login(ctx context.Context, client *http.Client, cfg config) (string, error) {
	body, err := json.Marshal(map[string]string{
		"username": cfg.username,
		"password": cfg.password,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(cfg.baseURL, "/")+authLoginPath, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", contentTypeJSON)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		payload, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("login returned status %d: %s", resp.StatusCode, string(payload))
	}

	var envelope loginEnvelope
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return "", err
	}
	if envelope.Result.Token == "" {
		return "", errors.New("login token missing")
	}
	return envelope.Result.Token, nil
}

func fetchRecordProjectionsLatest(ctx context.Context, client *http.Client, cfg config, state *scenarioState) error {
	body, err := json.Marshal(map[string]any{
		"query": state.query,
		"variables": map[string]any{
			"limit": cfg.recordLimit,
		},
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(cfg.baseURL, "/")+graphqlPath, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", contentTypeJSON)
	req.Header.Set("Authorization", "Bearer "+state.token)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		payload, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("record projections returned status %d: %s", resp.StatusCode, string(payload))
	}

	var envelope graphqlEnvelope
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return err
	}
	if len(envelope.Errors) > 0 {
		return errors.New(envelope.Errors[0].Message)
	}
	if _, ok := envelope.Data["recordProjectionsLatest"]; !ok {
		return errors.New("recordProjectionsLatest field missing")
	}
	return nil
}

func fetchDashboardSnapshot(ctx context.Context, client *http.Client, cfg config, state *scenarioState) error {
	body, err := json.Marshal(map[string]any{
		"query": state.query,
		"variables": map[string]any{
			"date":     state.dashboardDate,
			"timezone": cfg.dashboardTZ,
		},
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(cfg.baseURL, "/")+graphqlPath, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", contentTypeJSON)
	req.Header.Set("Authorization", "Bearer "+state.token)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		payload, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("dashboard snapshot returned status %d: %s", resp.StatusCode, string(payload))
	}

	var envelope graphqlEnvelope
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return err
	}
	if len(envelope.Errors) > 0 {
		return errors.New(envelope.Errors[0].Message)
	}
	if _, ok := envelope.Data["dashboardSnapshot"]; !ok {
		return errors.New("dashboardSnapshot field missing")
	}
	return nil
}

func runRealtimeRecordScenario(ctx context.Context, client *http.Client, cfg config, state *scenarioState) error {
	streamCtx, cancel := context.WithTimeout(ctx, cfg.timeout)
	defer cancel()

	events := make(chan sseEnvelope, 16)
	errs := make(chan error, 1)

	go consumeSSE(streamCtx, cfg, state.token, events, errs)

	if err := waitConnected(streamCtx, events, errs); err != nil {
		return err
	}

	recordID, err := createRecord(ctx, client, cfg, state.token, "codex realtime load")
	if err != nil {
		return err
	}
	defer func() {
		_ = deleteRecord(context.Background(), client, cfg, state.token, recordID)
	}()

	payload, err := waitProjectionChanged(streamCtx, events, errs, recordID)
	if err != nil {
		return err
	}
	if payload.Action != "created" {
		return fmt.Errorf("unexpected realtime action %q for record_id=%s", payload.Action, recordID)
	}
	return nil
}

func consumeSSE(ctx context.Context, cfg config, token string, events chan<- sseEnvelope, errs chan<- error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, strings.TrimRight(cfg.baseURL, "/")+realtimeStreamPath, nil)
	if err != nil {
		errs <- err
		return
	}
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := (&http.Client{Timeout: cfg.timeout}).Do(req)
	if err != nil {
		errs <- err
		return
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		errs <- fmt.Errorf("sse stream returned status=%d body=%s", resp.StatusCode, string(body))
		return
	}

	reader := bufio.NewScanner(resp.Body)
	var currentEvent string
	var dataLines []string

	for reader.Scan() {
		line := reader.Text()
		if line == "" {
			if currentEvent != "" {
				events <- sseEnvelope{
					Event: currentEvent,
					Data:  strings.Join(dataLines, "\n"),
				}
			}
			currentEvent = ""
			dataLines = dataLines[:0]
			continue
		}
		if strings.HasPrefix(line, ":") {
			continue
		}
		if strings.HasPrefix(line, "event:") {
			currentEvent = strings.TrimSpace(strings.TrimPrefix(line, "event:"))
			continue
		}
		if strings.HasPrefix(line, "data:") {
			dataLines = append(dataLines, strings.TrimSpace(strings.TrimPrefix(line, "data:")))
		}
	}

	if err := reader.Err(); err != nil && !errors.Is(err, context.Canceled) {
		errs <- err
	}
}

func waitConnected(ctx context.Context, events <-chan sseEnvelope, errs <-chan error) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("waiting connected event: %w", ctx.Err())
		case err := <-errs:
			return err
		case event := <-events:
			if event.Event == "connected" {
				return nil
			}
		}
	}
}

func waitProjectionChanged(ctx context.Context, events <-chan sseEnvelope, errs <-chan error, recordID string) (projectionChangedPayload, error) {
	for {
		select {
		case <-ctx.Done():
			return projectionChangedPayload{}, fmt.Errorf("waiting record_projection_changed: %w", ctx.Err())
		case err := <-errs:
			return projectionChangedPayload{}, err
		case event := <-events:
			if event.Event != "record_projection_changed" {
				continue
			}
			var payload projectionChangedPayload
			if err := json.Unmarshal([]byte(event.Data), &payload); err != nil {
				return projectionChangedPayload{}, err
			}
			if fmt.Sprint(payload.RecordID) == recordID {
				return payload, nil
			}
		}
	}
}

func createRecord(ctx context.Context, client *http.Client, cfg config, token string, description string) (string, error) {
	query := fmt.Sprintf(
		`mutation { createRecord(input: { tagId: %q, description: %q, source: "codex-load-test", status: "published" }) { id } }`,
		cfg.tagID,
		description,
	)
	var result createRecordResult
	if err := graphql(ctx, client, cfg.baseURL, token, query, map[string]any{}, "createRecord", &result); err != nil {
		return "", err
	}
	if result.ID == "" {
		return "", errors.New("createRecord returned empty id")
	}
	return result.ID, nil
}

func deleteRecord(ctx context.Context, client *http.Client, cfg config, token string, recordID string) error {
	query := fmt.Sprintf(`mutation { softDeleteRecord(input: { id: %q }) }`, recordID)
	var deleted bool
	return graphql(ctx, client, cfg.baseURL, token, query, map[string]any{}, "softDeleteRecord", &deleted)
}

func graphql(ctx context.Context, client *http.Client, host string, token string, query string, variables map[string]any, field string, out interface{}) error {
	if variables == nil {
		variables = map[string]any{}
	}
	body, err := json.Marshal(map[string]any{
		"query":     query,
		"variables": variables,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(host, "/")+graphqlPath, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", contentTypeJSON)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		payload, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("graphql returned status %d: %s", resp.StatusCode, string(payload))
	}

	var envelope graphqlEnvelope
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return err
	}
	if len(envelope.Errors) > 0 {
		return errors.New(envelope.Errors[0].Message)
	}
	raw, ok := envelope.Data[field]
	if !ok {
		return fmt.Errorf("graphql field %s missing", field)
	}
	return json.Unmarshal(raw, out)
}

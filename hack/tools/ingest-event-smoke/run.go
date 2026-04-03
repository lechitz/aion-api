package main

import (
	"bytes"
	"context"
	"crypto/sha1" //nolint:gosec // deterministic compatibility with current ingest event id algorithm
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	rawPayloadBucket = "aion-raw-events"
)

type ingestEnvelope struct {
	EventID       string          `json:"event_id"`
	EventType     string          `json:"event_type"`
	EventVersion  string          `json:"event_version"`
	OccurredAtUTC string          `json:"occurred_at_utc"`
	Source        string          `json:"source"`
	RawPayloadRef string          `json:"raw_payload_ref"`
	PayloadJSON   json.RawMessage `json:"payload_json"`
}

type smokePayload struct {
	Steps      int    `json:"steps"`
	CapturedAt string `json:"captured_at"`
	SmokeNonce string `json:"smoke_nonce"`
}

func run(ctx context.Context, cfg config) error {
	payload, err := buildPayload(time.Now().UTC())
	if err != nil {
		return err
	}
	expectedEventID := sha1Hex([]byte(strings.TrimSpace(cfg.source) + ":" + string(payload)))
	expectedRawRef := "/data/raw/" + rawPayloadBucket + "/" + sanitizeSource(cfg.source) + "/" + expectedEventID + ".json"

	startOffset, err := getLastOffset(cfg.kafkaBroker, cfg.topic)
	if err != nil {
		return err
	}

	if err := postWebhook(ctx, cfg, payload); err != nil {
		return err
	}

	envelope, err := waitIngestMessage(ctx, cfg, startOffset, expectedEventID)
	if err != nil {
		return err
	}

	if envelope.EventType != "ingest.normalized" || envelope.EventVersion != "v1" {
		return fmt.Errorf("unexpected ingest event metadata: %+v", envelope)
	}
	if envelope.Source != cfg.source {
		return fmt.Errorf("unexpected source: %s", envelope.Source)
	}
	if envelope.RawPayloadRef != expectedRawRef {
		return fmt.Errorf("unexpected raw payload ref: %s", envelope.RawPayloadRef)
	}
	if string(envelope.PayloadJSON) != string(payload) {
		return fmt.Errorf("unexpected payload json: %s", string(envelope.PayloadJSON))
	}

	_, _ = fmt.Fprintf(os.Stdout, "ingest event smoke passed for event_id=%s\n", envelope.EventID)
	return nil
}

func buildPayload(now time.Time) ([]byte, error) {
	payload := smokePayload{
		Steps:      42,
		CapturedAt: "2026-03-13T16:30:00Z",
		SmokeNonce: now.UTC().Format(time.RFC3339Nano),
	}

	return json.Marshal(payload)
}

func postWebhook(ctx context.Context, cfg config, payload []byte) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(cfg.ingestHost, "/")+"/webhooks/raw", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Aion-Source", cfg.source)

	client := &http.Client{Timeout: cfg.timeout}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected ingest status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func getLastOffset(broker string, topic string) (int64, error) {
	conn, err := kafka.DialLeader(context.Background(), "tcp", broker, topic, 0)
	if err != nil {
		return 0, err
	}
	defer func() { _ = conn.Close() }()

	return conn.ReadLastOffset()
}

func waitIngestMessage(ctx context.Context, cfg config, startOffset int64, expectedEventID string) (ingestEnvelope, error) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{cfg.kafkaBroker},
		Topic:     cfg.topic,
		Partition: 0,
		MinBytes:  1,
		MaxBytes:  10e6,
	})
	defer func() { _ = reader.Close() }()

	if err := reader.SetOffset(startOffset); err != nil {
		return ingestEnvelope{}, err
	}

	deadline := time.Now().Add(cfg.timeout)
	for time.Now().Before(deadline) {
		readCtx, cancel := context.WithTimeout(ctx, cfg.pollSleep)
		msg, err := reader.ReadMessage(readCtx)
		cancel()
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				continue
			}
			if strings.Contains(err.Error(), "context deadline exceeded") {
				continue
			}
			return ingestEnvelope{}, err
		}

		var envelope ingestEnvelope
		if err := json.Unmarshal(msg.Value, &envelope); err != nil {
			return ingestEnvelope{}, err
		}
		if envelope.EventID == expectedEventID {
			return envelope, nil
		}
	}

	return ingestEnvelope{}, fmt.Errorf("ingest event %s not observed on topic %s", expectedEventID, cfg.topic)
}

func sha1Hex(payload []byte) string {
	//nolint:gosec // deterministic compatibility with current ingest event id algorithm
	sum := sha1.Sum(payload)
	return hex.EncodeToString(sum[:])
}

func sanitizeSource(source string) string {
	replacer := strings.NewReplacer("/", "-", " ", "-", ":", "-", "\\", "-")
	return strings.Trim(strings.ToLower(replacer.Replace(source)), "-.")
}

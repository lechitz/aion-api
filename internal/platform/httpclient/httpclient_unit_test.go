package httpclient_test

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/lechitz/aion-api/internal/platform/httpclient"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type trackingRoundTripper struct {
	status int
	called bool
}

func (t *trackingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	t.called = true
	return &http.Response{
		StatusCode: t.status,
		Body:       io.NopCloser(strings.NewReader("")),
		Header:     make(http.Header),
		Request:    req,
	}, nil
}

func TestNewInstrumentedClient_DefaultTimeoutAndTransport(t *testing.T) {
	client := httpclient.NewInstrumentedClient(httpclient.Options{})
	if client == nil {
		t.Fatal("expected non-nil client")
		return
	}
	if client.Timeout != 15*time.Second {
		t.Fatalf("unexpected default timeout: %s", client.Timeout)
	}
	if client.Transport == nil {
		t.Fatal("expected non-nil transport")
	}
}

func TestNewInstrumentedClient_DisableInstrumentationUsesBaseTransport(t *testing.T) {
	base := &trackingRoundTripper{status: http.StatusNoContent}

	client := httpclient.NewInstrumentedClient(httpclient.Options{
		Timeout:                2 * time.Second,
		BaseTransport:          base,
		DisableInstrumentation: true,
	})

	if client.Timeout != 2*time.Second {
		t.Fatalf("unexpected timeout: %s", client.Timeout)
	}
	if client.Transport != base {
		t.Fatal("expected base transport when instrumentation is disabled")
	}
}

func TestNewInstrumentedClient_InstrumentedTransportCallsBaseRoundTripper(t *testing.T) {
	called := false
	base := roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		called = true
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("ok")),
			Header:     make(http.Header),
			Request:    req,
		}, nil
	})

	client := httpclient.NewInstrumentedClient(httpclient.Options{
		BaseTransport: base,
	})

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com", nil)
	if err != nil {
		t.Fatalf("request build failed: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if !called {
		t.Fatal("expected wrapped transport to call base round tripper")
	}
}

func TestClientAdapterDo(t *testing.T) {
	rt := roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusAccepted,
			Body:       io.NopCloser(strings.NewReader("accepted")),
			Header:     make(http.Header),
			Request:    req,
		}, nil
	})
	httpClient := &http.Client{Transport: rt}
	client := httpclient.NewClient(httpClient)

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com", nil)
	if err != nil {
		t.Fatalf("request build failed: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Do failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusAccepted {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
}

package observability_test

import (
	"testing"

	"github.com/lechitz/aion-api/internal/platform/observability"
	"github.com/stretchr/testify/require"
)

func TestParseHeaders(t *testing.T) {
	require.Empty(t, observability.ParseHeaders(""))

	headers := observability.ParseHeaders("a=1, b=2, malformed, c = 3 ")
	require.Equal(t, "1", headers["a"])
	require.Equal(t, "2", headers["b"])
	require.Equal(t, "3", headers["c"])
}

func TestNormalizeEndpoint(t *testing.T) {
	normalized, err := observability.NormalizeEndpoint("aion-dev-otel-collector:4318")
	require.NoError(t, err)
	require.Equal(t, "http://aion-dev-otel-collector:4318", normalized)

	normalized, err = observability.NormalizeEndpoint("https://otel:4318/path")
	require.NoError(t, err)
	require.Equal(t, "https://otel:4318/path", normalized)

	normalized, err = observability.NormalizeEndpoint(" ")
	require.NoError(t, err)
	require.Empty(t, normalized)

	_, err = observability.NormalizeEndpoint("http://[::1")
	require.Error(t, err)
}

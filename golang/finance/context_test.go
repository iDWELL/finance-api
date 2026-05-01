package finance

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBaseURLFromCtx(t *testing.T) {
	t.Parallel()

	t.Run("falls back to BaseURL when not set", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, BaseURL, baseURLFromCtx(t.Context()))
	})

	t.Run("falls back to BaseURL when empty string", func(t *testing.T) {
		t.Parallel()
		ctx := WithBaseURL(t.Context(), "")
		assert.Equal(t, BaseURL, baseURLFromCtx(ctx))
	})

	t.Run("returns custom URL", func(t *testing.T) {
		t.Parallel()

		expected := "http://localhost:9999"
		ctx := WithBaseURL(t.Context(), expected)
		assert.Equal(t, expected, baseURLFromCtx(ctx))
	})
}

func TestHTTPClientFromCtx(t *testing.T) {
	t.Parallel()

	t.Run("returns DefaultClient when not set", func(t *testing.T) {
		t.Parallel()
		assert.Same(t, http.DefaultClient, httpClientFromCtx(t.Context()))
	})

	t.Run("returns DefaultClient when nil set", func(t *testing.T) {
		t.Parallel()
		ctx := WithHTTPClient(t.Context(), nil)
		assert.Same(t, http.DefaultClient, httpClientFromCtx(ctx))
	})

	t.Run("returns custom client", func(t *testing.T) {
		t.Parallel()

		custom := &http.Client{}
		ctx := WithHTTPClient(t.Context(), custom)
		assert.Same(t, custom, httpClientFromCtx(ctx))
	})
}

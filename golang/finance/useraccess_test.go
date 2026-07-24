package finance

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRevokeUserAccess_Success(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/user/access/revoke", r.URL.Path)
		assert.Contains(t, r.Header.Get("Authorization"), "Bearer ")

		var requests []UserAccessRevokeRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&requests))
		require.Len(t, requests, 2)

		results := []UserAccessRevokeResult{
			{Email: requests[0].Email, Success: true},
			{Email: requests[1].Email, Success: false, Error: "user not found"},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(results)
	}))
	defer srv.Close()

	ctx := WithBaseURL(t.Context(), srv.URL)
	got, err := RevokeUserAccess(ctx, "key", []UserAccessRevokeRequest{
		{Email: "a@example.com"},
		{Email: "missing@example.com"},
	})
	require.NoError(t, err)
	require.Len(t, got, 2)
	assert.True(t, got[0].Success)
	assert.Equal(t, "a@example.com", got[0].Email)
	assert.False(t, got[1].Success)
	assert.Equal(t, "user not found", got[1].Error)
}

func TestRevokeUserAccess_ServerError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(errorHandler(http.StatusInternalServerError))
	defer srv.Close()

	ctx := WithBaseURL(t.Context(), srv.URL)
	_, err := RevokeUserAccess(ctx, "key", []UserAccessRevokeRequest{{Email: "a@example.com"}})
	require.Error(t, err)
}

func TestRevokeUserAccess_NetworkError(t *testing.T) {
	t.Parallel()

	ctx := WithBaseURL(t.Context(), "http://127.0.0.1:1")
	_, err := RevokeUserAccess(ctx, "key", []UserAccessRevokeRequest{{Email: "a@example.com"}})
	require.Error(t, err)
}

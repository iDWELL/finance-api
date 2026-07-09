package finance

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateClientCompanyCostCenter_Success(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.Header.Get("Authorization"), "Bearer ")
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var req struct {
			Query     string            `json:"query"`
			Variables map[string]string `json:"variables"`
		}
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&req))

		// the API schema requires NonEmptyText!, String! is rejected before execution
		assert.Contains(t, req.Query, "$number: NonEmptyText!")
		assert.Contains(t, req.Query, "$description: NonEmptyText!")
		assert.Equal(t, "1234/5678", req.Variables["number"])
		assert.Equal(t, "Heating", req.Variables["description"])

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"createClientCompanyCostCenter":{"clientCompanyCostCenter":{"rowId":"aabbccdd-1234-5678-abcd-ef0123456789","number":"1234/5678"}}}}`))
	}))
	defer srv.Close()

	ctx := WithBaseURL(t.Context(), srv.URL)
	cc, err := CreateClientCompanyCostCenter(ctx, "key", "1234/5678", "Heating")
	require.NoError(t, err)
	assert.Equal(t, "1234/5678", cc.Number)
	assert.Equal(t, "aabbccdd-1234-5678-abcd-ef0123456789", cc.ID.String())
}

func TestCreateClientCompanyCostCenter_GraphQLError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(okJSONHandler(t, map[string]any{
		"errors": []map[string]string{{"message": "some graphql error"}},
	}))
	defer srv.Close()

	ctx := WithBaseURL(t.Context(), srv.URL)
	_, err := CreateClientCompanyCostCenter(ctx, "key", "1234/5678", "Heating")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "some graphql error")
}

func TestCreateClientCompanyCostCenter_ServerError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(errorHandler(http.StatusBadRequest))
	defer srv.Close()

	ctx := WithBaseURL(t.Context(), srv.URL)
	_, err := CreateClientCompanyCostCenter(ctx, "key", "1234/5678", "Heating")
	require.Error(t, err)
}

func TestCreateClientCompanyCostCenter_NetworkError(t *testing.T) {
	t.Parallel()

	ctx := WithBaseURL(t.Context(), "http://127.0.0.1:1")
	_, err := CreateClientCompanyCostCenter(ctx, "key", "1234/5678", "Heating")
	require.Error(t, err)
}

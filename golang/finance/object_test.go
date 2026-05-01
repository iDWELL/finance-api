package finance

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostObjectInstancesWithIDProp_ValidationErrors(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		className   string
		idPropName  string
		objects     []map[string]any
		expectedErr string
	}{
		"empty className": {
			className:   "",
			idPropName:  "id",
			objects:     []map[string]any{{"id": "1"}},
			expectedErr: "className is required",
		},
		"empty idPropName": {
			className:   "Asset",
			idPropName:  "",
			objects:     []map[string]any{{"id": "1"}},
			expectedErr: "idPropName is required",
		},
		"invalid className characters": {
			className:   "my class!",
			idPropName:  "id",
			objects:     []map[string]any{{"id": "1"}},
			expectedErr: "className contains invalid characters",
		},
		"invalid idPropName characters": {
			className:   "Asset",
			idPropName:  "my prop!",
			objects:     []map[string]any{{"id": "1"}},
			expectedErr: "idPropName contains invalid characters",
		},
		"object missing ID prop": {
			className:   "Asset",
			idPropName:  "assetId",
			objects:     []map[string]any{{"name": "thing"}},
			expectedErr: `object at index 0 has no ID prop "assetId"`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := WithBaseURL(t.Context(), "http://127.0.0.1:1")
			err := PostObjectInstancesWithIDProp(ctx, "key", tc.className, tc.idPropName, tc.objects, "")
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.expectedErr)
		})
	}
}

func TestPostObjectInstancesWithIDProp_Success(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/masterdata/upsert-objects/Asset/id-prop/assetId", r.URL.Path)
		assert.Equal(t, "myapp", r.URL.Query().Get("source"))
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	ctx := WithBaseURL(t.Context(), srv.URL)
	err := PostObjectInstancesWithIDProp(ctx, "key", "Asset", "assetId",
		[]map[string]any{{"assetId": "A001", "name": "Main Building"}},
		"myapp",
	)
	require.NoError(t, err)
}

func TestPostObjectInstancesWithIDProp_ServerError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(errorHandler(http.StatusInternalServerError))
	defer srv.Close()

	ctx := WithBaseURL(t.Context(), srv.URL)
	err := PostObjectInstancesWithIDProp(ctx, "key", "Asset", "assetId",
		[]map[string]any{{"assetId": "A001"}},
		"",
	)
	require.Error(t, err)
}

func TestPostObjectInstancesWithIDProp_NetworkError(t *testing.T) {
	t.Parallel()

	ctx := WithBaseURL(t.Context(), "http://127.0.0.1:1")
	err := PostObjectInstancesWithIDProp(ctx, "key", "Asset", "assetId",
		[]map[string]any{{"assetId": "A001"}},
		"",
	)
	require.Error(t, err)
}

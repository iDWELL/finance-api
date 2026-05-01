package finance

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/domonda/go-types/account"
	"github.com/domonda/go-types/bank"
	"github.com/domonda/go-types/money"
	"github.com/domonda/go-types/notnull"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func okJSONHandler(t *testing.T, body any) http.HandlerFunc {
	t.Helper()

	return func(w http.ResponseWriter, r *http.Request) {
		t.Helper()
		assert.Contains(t, r.Header.Get("Authorization"), "Bearer ")
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(body)
	}
}

func errorHandler(status int) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(status)
	}
}

// ---- PostGLAccounts ----

func TestPostGLAccounts_Success(t *testing.T) {
	t.Parallel()

	results := []*ImportGLAccountResult{{NormalizedNumber: "4200", State: ImportStateCreated}}

	srv := httptest.NewServer(okJSONHandler(t, results))
	defer srv.Close()

	ctx := WithBaseURL(t.Context(), srv.URL)
	got, err := PostGLAccounts(ctx, "key", []*GLAccount{{Number: account.Number("4200")}}, false, false, false, false, "test")
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, ImportStateCreated, got[0].State)
}

func TestPostGLAccounts_ServerError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(errorHandler(http.StatusInternalServerError))
	defer srv.Close()

	ctx := WithBaseURL(t.Context(), srv.URL)
	_, err := PostGLAccounts(ctx, "key", []*GLAccount{{Number: "4200"}}, false, false, false, false, "")
	require.Error(t, err)
}

func TestPostGLAccounts_ValidationError(t *testing.T) {
	t.Parallel()

	ctx := WithBaseURL(t.Context(), "http://127.0.0.1:1")
	_, err := PostGLAccounts(ctx, "key", []*GLAccount{{Number: ""}}, false, false, true, false, "")
	require.Error(t, err)
}

// ---- PostBankAccounts ----

func TestPostBankAccounts_Success(t *testing.T) {
	t.Parallel()

	results := []*ImportBankAccountResult{{State: ImportStateCreated}}

	srv := httptest.NewServer(okJSONHandler(t, results))
	defer srv.Close()

	ctx := WithBaseURL(t.Context(), srv.URL)
	got, err := PostBankAccounts(ctx, "key", []*BankAccount{{
		IBAN:     bank.IBAN("AT611904300234573201"),
		BIC:      bank.BIC("OPSKATWW"),
		Currency: money.Currency("EUR"),
		Holder:   notnull.TrimmedString("Holder"),
	}}, false, false, "test")
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, ImportStateCreated, got[0].State)
}

func TestPostBankAccounts_ServerError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(errorHandler(http.StatusBadRequest))
	defer srv.Close()

	ctx := WithBaseURL(t.Context(), srv.URL)
	_, err := PostBankAccounts(ctx, "key", []*BankAccount{{
		IBAN: "AT611904300234573201", BIC: "OPSKATWW", Currency: "EUR", Holder: "H",
	}}, false, false, "")
	require.Error(t, err)
}

// ---- PostObjectGroups ----

func TestPostObjectGroups_Success(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/masterdata/real-estate-object-groups", r.URL.Path)
		assert.Equal(t, "myapp", r.URL.Query().Get("source"))
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	ctx := WithBaseURL(t.Context(), srv.URL)
	err := PostObjectGroups(ctx, "key", []*ObjectGroup{{ExternalID: "g1", Name: "G1", ObjectIDs: []string{"001"}}}, "myapp")
	require.NoError(t, err)
}

func TestPostObjectGroups_ServerError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(errorHandler(http.StatusInternalServerError))
	defer srv.Close()

	ctx := WithBaseURL(t.Context(), srv.URL)
	err := PostObjectGroups(ctx, "key", []*ObjectGroup{{ExternalID: "g1", Name: "G1"}}, "")
	require.Error(t, err)
}

func TestPostObjectGroups_NetworkError(t *testing.T) {
	t.Parallel()

	ctx := WithBaseURL(t.Context(), "http://127.0.0.1:1")
	err := PostObjectGroups(ctx, "key", []*ObjectGroup{{ExternalID: "g1", Name: "G1"}}, "")
	require.Error(t, err)
}

// ---- PostObjectRoles ----

func TestPostObjectRoles_Success(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/masterdata/real-estate-object-roles", r.URL.Path)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	ctx := WithBaseURL(t.Context(), srv.URL)
	err := PostObjectRoles(ctx, "key", []*ObjectRole{{Name: "Manager", ObjectID: "001", Email: "m@example.com"}}, "myapp")
	require.NoError(t, err)
}

func TestPostObjectRoles_ServerError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(errorHandler(http.StatusBadRequest))
	defer srv.Close()

	ctx := WithBaseURL(t.Context(), srv.URL)
	err := PostObjectRoles(ctx, "key", []*ObjectRole{{Name: "M", ObjectID: "001", Email: "m@example.com"}}, "")
	require.Error(t, err)
}

func TestPostObjectRoles_NetworkError(t *testing.T) {
	t.Parallel()

	ctx := WithBaseURL(t.Context(), "http://127.0.0.1:1")
	err := PostObjectRoles(ctx, "key", []*ObjectRole{{Name: "M", ObjectID: "001", Email: "m@example.com"}}, "")
	require.Error(t, err)
}

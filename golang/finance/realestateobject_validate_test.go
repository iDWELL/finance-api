package finance

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/domonda/go-types/account"
	"github.com/domonda/go-types/bank"
	"github.com/domonda/go-types/country"
	"github.com/domonda/go-types/notnull"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validRealEstateObject(t *testing.T) *RealEstateObject {
	t.Helper()

	return &RealEstateObject{
		Type:          RealEstateObjectTypeWEG,
		Number:        account.Number("4200"),
		Country:       country.Code("AT"),
		StreetAddress: notnull.TrimmedString("Main St 1"),
	}
}

func TestRealEstateObjectValidate(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		build       func(*testing.T) *RealEstateObject
		expectedErr string
	}{
		"valid object": {
			build: func(t *testing.T) *RealEstateObject { t.Helper(); return validRealEstateObject(t) },
		},
		"invalid Type": {
			build: func(t *testing.T) *RealEstateObject {
				t.Helper()

				o := validRealEstateObject(t)
				o.Type = "BOGUS"

				return o
			},
			expectedErr: "RealEstateObject.Type",
		},
		"invalid Number": {
			build: func(t *testing.T) *RealEstateObject {
				t.Helper()

				o := validRealEstateObject(t)
				o.Number = ""

				return o
			},
			expectedErr: "RealEstateObject.Number",
		},
		"invalid AccountingArea": {
			build: func(t *testing.T) *RealEstateObject {
				t.Helper()

				o := validRealEstateObject(t)
				o.AccountingArea = "acc no!" // spaces and special chars are invalid

				return o
			},
			expectedErr: "RealEstateObject.AccountingArea",
		},
		"invalid UserAccount": {
			build: func(t *testing.T) *RealEstateObject {
				t.Helper()

				o := validRealEstateObject(t)
				o.UserAccount = "user no!" // spaces and special chars are invalid

				return o
			},
			expectedErr: "RealEstateObject.UserAccount",
		},
		"invalid Country": {
			build: func(t *testing.T) *RealEstateObject {
				t.Helper()

				o := validRealEstateObject(t)
				o.Country = "INVALID"

				return o
			},
			expectedErr: "RealEstateObject.Country",
		},
		"empty StreetAddress": {
			build: func(t *testing.T) *RealEstateObject {
				t.Helper()

				o := validRealEstateObject(t)
				o.StreetAddress = ""

				return o
			},
			expectedErr: "RealEstateObject.StreetAddress",
		},
		"invalid BankAccount IBAN": {
			build: func(t *testing.T) *RealEstateObject {
				t.Helper()

				o := validRealEstateObject(t)
				o.BankAccounts = []bank.Account{{IBAN: "NOTANIBAN"}}

				return o
			},
			expectedErr: "RealEstateObject.BankAccounts[0]",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := tc.build(t).Validate()
			if tc.expectedErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestPostRealEstateObjects_ValidationError(t *testing.T) {
	t.Parallel()

	ctx := WithBaseURL(t.Context(), "http://127.0.0.1:1")
	invalid := &RealEstateObject{Type: "BOGUS", Number: "", StreetAddress: ""}
	err := PostRealEstateObjects(ctx, "key", []*RealEstateObject{invalid}, "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "RealEstateObject at index 0")
}

func TestPostRealEstateObjects_Success(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/masterdata/real-estate-objects", r.URL.Path)
		assert.Contains(t, r.Header.Get("Authorization"), "Bearer ")
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	ctx := WithBaseURL(t.Context(), srv.URL)
	err := PostRealEstateObjects(ctx, "key", []*RealEstateObject{validRealEstateObject(t)}, "myapp")
	require.NoError(t, err)
}

func TestPostRealEstateObjects_ServerError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(errorHandler(http.StatusInternalServerError))
	defer srv.Close()

	ctx := WithBaseURL(t.Context(), srv.URL)
	err := PostRealEstateObjects(ctx, "key", []*RealEstateObject{validRealEstateObject(t)}, "")
	require.Error(t, err)
}

func TestPostRealEstateObjects_NetworkError(t *testing.T) {
	t.Parallel()

	ctx := WithBaseURL(t.Context(), "http://127.0.0.1:1")
	err := PostRealEstateObjects(ctx, "key", []*RealEstateObject{validRealEstateObject(t)}, "")
	require.Error(t, err)
}

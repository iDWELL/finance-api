package finance

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/domonda/go-types/account"
	"github.com/domonda/go-types/notnull"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validObjectTenantOwner(t *testing.T) *ObjectTenantOwner {
	t.Helper()

	return &ObjectTenantOwner{
		ObjectNo:      account.Number("4200"),
		TenantOwnerID: 1,
		TenantOwnerNo: 1,
		UnitNo:        1,
		Unit:          notnull.TrimmedString("Unit 1"),
		OwnerLinkNo:   1,
		Owner:         notnull.TrimmedString("John Doe"),
	}
}

func TestObjectTenantOwnerValidate(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		build       func(*testing.T) *ObjectTenantOwner
		expectedErr string
	}{
		"valid tenant owner": {
			build: func(t *testing.T) *ObjectTenantOwner { t.Helper(); return validObjectTenantOwner(t) },
		},
		"invalid ObjectNo": {
			build: func(t *testing.T) *ObjectTenantOwner {
				t.Helper()

				o := validObjectTenantOwner(t)
				o.ObjectNo = ""

				return o
			},
			expectedErr: "ObjectTenantOwner.ObjectNo",
		},
		"empty Unit": {
			build: func(t *testing.T) *ObjectTenantOwner {
				t.Helper()

				o := validObjectTenantOwner(t)
				o.Unit = ""

				return o
			},
			expectedErr: "ObjectTenantOwner.Unit",
		},
		"empty Owner": {
			build: func(t *testing.T) *ObjectTenantOwner {
				t.Helper()

				o := validObjectTenantOwner(t)
				o.Owner = ""

				return o
			},
			expectedErr: "ObjectTenantOwner.Owner",
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

func TestPostObjectTenantOwners_ValidationError(t *testing.T) {
	t.Parallel()

	ctx := WithBaseURL(t.Context(), "http://127.0.0.1:1")
	invalid := &ObjectTenantOwner{ObjectNo: "", Unit: "", Owner: ""}
	err := PostObjectTenantOwners(ctx, "key", []*ObjectTenantOwner{invalid}, "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "ObjectTenantOwner at index 0")
}

func TestPostObjectTenantOwners_Success(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/masterdata/real-estate-object-tenant-owners", r.URL.Path)
		assert.Contains(t, r.Header.Get("Authorization"), "Bearer ")
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	ctx := WithBaseURL(t.Context(), srv.URL)
	err := PostObjectTenantOwners(ctx, "key", []*ObjectTenantOwner{validObjectTenantOwner(t)}, "myapp")
	require.NoError(t, err)
}

func TestPostObjectTenantOwners_ServerError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(errorHandler(http.StatusInternalServerError))
	defer srv.Close()

	ctx := WithBaseURL(t.Context(), srv.URL)
	err := PostObjectTenantOwners(ctx, "key", []*ObjectTenantOwner{validObjectTenantOwner(t)}, "")
	require.Error(t, err)
}

func TestPostObjectTenantOwners_NetworkError(t *testing.T) {
	t.Parallel()

	ctx := WithBaseURL(t.Context(), "http://127.0.0.1:1")
	err := PostObjectTenantOwners(ctx, "key", []*ObjectTenantOwner{validObjectTenantOwner(t)}, "")
	require.Error(t, err)
}

package finance

import (
	"testing"

	"github.com/domonda/go-types/account"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGLAccountValidate(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		build       func(*testing.T) *GLAccount
		expectedErr string
	}{
		"valid GL account": {
			build: func(t *testing.T) *GLAccount { t.Helper(); return &GLAccount{Number: account.Number("4200")} },
		},
		"empty number": {
			build:       func(t *testing.T) *GLAccount { t.Helper(); return &GLAccount{Number: ""} },
			expectedErr: "GLAccount.Number",
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

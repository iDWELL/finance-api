package finance

import (
	"testing"

	"github.com/domonda/go-types/bank"
	"github.com/domonda/go-types/notnull"
	"github.com/domonda/go-types/nullable"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validPartner(t *testing.T) *Partner {
	t.Helper()

	return &Partner{Name: notnull.TrimmedString("ACME Corp")}
}

func TestPartnerValidate(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		build       func(*testing.T) *Partner
		expectedErr string
	}{
		"valid minimal partner": {
			build: func(t *testing.T) *Partner { t.Helper(); return validPartner(t) },
		},
		"empty name": {
			build:       func(t *testing.T) *Partner { t.Helper(); p := validPartner(t); p.Name = ""; return p },
			expectedErr: "empty Partner.Name",
		},
		"invalid IBAN": {
			build: func(t *testing.T) *Partner {
				t.Helper()

				p := validPartner(t)
				p.IBAN = bank.NullableIBAN("NOTANIBAN")

				return p
			},
			expectedErr: "invalid Partner.IBAN",
		},
		"invalid BIC": {
			build: func(t *testing.T) *Partner {
				t.Helper()

				p := validPartner(t)
				p.BIC = bank.NullableBIC("!!")

				return p
			},
			expectedErr: "invalid Partner.BIC",
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

func TestPartnerNormalizedAlternativeNames(t *testing.T) {
	t.Parallel()

	t.Run("trims whitespace, removes empty, sorts", func(t *testing.T) {
		t.Parallel()

		p := &Partner{
			Name:             "ACME",
			AlternativeNames: notnull.StringArray{"  Zebra Co  ", "", "Alpha Ltd", "  "},
		}
		assert.Equal(t, []string{"Alpha Ltd", "Zebra Co"}, p.NormalizedAlternativeNames())
	})

	t.Run("nil returns empty", func(t *testing.T) {
		t.Parallel()

		p := &Partner{Name: "ACME"}
		assert.Empty(t, p.NormalizedAlternativeNames())
	})
}

func TestPartnerEqualAlternativeNames(t *testing.T) {
	t.Parallel()

	p := &Partner{
		Name:             "ACME",
		AlternativeNames: notnull.StringArray{"Alpha", "Beta"},
	}
	assert.True(t, p.EqualAlternativeNames([]string{"Beta", "Alpha"}))
	assert.True(t, p.EqualAlternativeNames([]string{"Alpha", "Beta"}))
	assert.False(t, p.EqualAlternativeNames([]string{"Alpha"}))
	assert.False(t, p.EqualAlternativeNames([]string{"Alpha", "Gamma"}))
}

func TestPartnerHasLocation(t *testing.T) {
	t.Parallel()

	p := validPartner(t)
	assert.False(t, p.HasLocation())

	p.City = nullable.TrimmedString("Vienna")
	assert.True(t, p.HasLocation())
}

func TestPartnerNormalize_Valid(t *testing.T) {
	t.Parallel()

	p := validPartner(t)
	errs := p.Normalize(false)
	assert.Empty(t, errs)
}

func TestPartnerNormalize_EmptyName(t *testing.T) {
	t.Parallel()

	p := &Partner{}
	errs := p.Normalize(false)
	require.Len(t, errs, 1)
	assert.Contains(t, errs[0].Error(), "name is empty")
}

func TestPartnerString(t *testing.T) {
	t.Parallel()

	p := validPartner(t)
	assert.Contains(t, p.String(), "ACME Corp")
}

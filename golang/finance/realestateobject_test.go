package finance

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRealEstateObjectTypeValid(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		typ      RealEstateObjectType
		expected bool
	}{
		"WEG valid":       {RealEstateObjectTypeWEG, true},
		"HI valid":        {RealEstateObjectTypeHI, true},
		"SUB valid":       {RealEstateObjectTypeSUB, true},
		"KREIS valid":     {RealEstateObjectTypeKREIS, true},
		"MANDANT valid":   {RealEstateObjectTypeMANDANT, true},
		"MRG valid":       {RealEstateObjectTypeMRG, true},
		"HBH valid":       {RealEstateObjectTypeHBH, true},
		"empty invalid":   {"", false},
		"lower invalid":   {"weg", false},
		"unknown invalid": {"UNKNOWN", false},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, tc.typ.Valid())
		})
	}
}

func TestRealEstateObjectTypeValidate(t *testing.T) {
	t.Parallel()

	t.Run("valid type returns nil", func(t *testing.T) {
		t.Parallel()
		require.NoError(t, RealEstateObjectTypeWEG.Validate())
	})

	t.Run("invalid type returns error mentioning value", func(t *testing.T) {
		t.Parallel()

		err := RealEstateObjectType("BOGUS").Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "BOGUS")
	})
}

func TestRealEstateObjectTypeIsVirtual(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		typ     RealEstateObjectType
		virtual bool
	}{
		"KREIS is virtual":   {RealEstateObjectTypeKREIS, true},
		"MANDANT is virtual": {RealEstateObjectTypeMANDANT, true},
		"WEG not virtual":    {RealEstateObjectTypeWEG, false},
		"HI not virtual":     {RealEstateObjectTypeHI, false},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.virtual, tc.typ.IsVirtual())
		})
	}
}

func TestRealEstateObjectTypeString(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "WEG", RealEstateObjectTypeWEG.String())
	assert.Equal(t, "HBH", RealEstateObjectTypeHBH.String())
}

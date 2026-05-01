package finance

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImportStateValid(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		state    ImportState
		expected bool
	}{
		"UNCHANGED valid":   {ImportStateUnchanged, true},
		"UPDATED valid":     {ImportStateUpdated, true},
		"CREATED valid":     {ImportStateCreated, true},
		"ERROR valid":       {ImportStateError, true},
		"empty invalid":     {"", false},
		"lowercase invalid": {"unchanged", false},
		"unknown invalid":   {"PENDING", false},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, tc.state.Valid())
		})
	}
}

func TestImportStateValidate(t *testing.T) {
	t.Parallel()

	t.Run("valid state returns nil", func(t *testing.T) {
		t.Parallel()
		require.NoError(t, ImportStateCreated.Validate())
	})

	t.Run("invalid state returns error mentioning value", func(t *testing.T) {
		t.Parallel()

		err := ImportState("BOGUS").Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "BOGUS")
	})
}

func TestImportStateEnums(t *testing.T) {
	t.Parallel()

	enums := ImportState("").Enums()
	assert.Len(t, enums, 4)
	assert.Contains(t, enums, ImportStateUnchanged)
	assert.Contains(t, enums, ImportStateUpdated)
	assert.Contains(t, enums, ImportStateCreated)
	assert.Contains(t, enums, ImportStateError)
}

func TestImportStateEnumStrings(t *testing.T) {
	t.Parallel()
	assert.Equal(t, []string{"UNCHANGED", "UPDATED", "CREATED", "ERROR"}, ImportState("").EnumStrings())
}

func TestImportStateString(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "UNCHANGED", ImportStateUnchanged.String())
	assert.Equal(t, "ERROR", ImportStateError.String())
}

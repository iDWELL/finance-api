package finance

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestObjectTenantOwnerRoleValid(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		role     ObjectTenantOwnerRole
		expected bool
	}{
		"UNSPECIFIED valid": {ObjectTenantOwnerRoleUnspecified, true},
		"TENANT valid":      {ObjectTenantOwnerRoleTenant, true},
		"OWNER valid":       {ObjectTenantOwnerRoleOwner, true},
		"empty invalid":     {"", false},
		"lowercase invalid": {"tenant", false},
		"German invalid":    {"Mieter", false},
		"unknown invalid":   {"VERWALTER", false},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.expected, tc.role.Valid())
			if tc.expected {
				require.NoError(t, tc.role.Validate())
			} else {
				require.Error(t, tc.role.Validate())
			}
		})
	}
}

func TestObjectTenantOwnerRoleEnums(t *testing.T) {
	t.Parallel()

	roles := ObjectTenantOwnerRole("").Enums()
	strs := ObjectTenantOwnerRole("").EnumStrings()

	require.Len(t, roles, 3)
	require.Len(t, strs, len(roles))
	for i, role := range roles {
		assert.Equal(t, strs[i], role.String())
		assert.True(t, role.Valid())
	}
}

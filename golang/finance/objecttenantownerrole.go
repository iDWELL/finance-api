package finance

import "fmt"

//go:generate go tool go-enum $GOFILE

// ObjectTenantOwnerRole is the role a person has on a unit of a real estate object.
// A row of the tenant/owner master data identifies one person, on one unit, in one
// role, so the role is part of that identity: the same person can be the tenant of
// one unit and the owner of another, or both on the same unit.
type ObjectTenantOwnerRole string //#enum

const (
	// ObjectTenantOwnerRoleUnspecified is the role of a record whose source lists
	// tenants and owners together without telling them apart. It is the zero value,
	// so a payload that omits the role keeps its previous meaning.
	ObjectTenantOwnerRoleUnspecified ObjectTenantOwnerRole = "UNSPECIFIED"

	// ObjectTenantOwnerRoleTenant is a tenant (Mieter)
	ObjectTenantOwnerRoleTenant ObjectTenantOwnerRole = "TENANT"

	// ObjectTenantOwnerRoleOwner is an owner (Eigentümer)
	ObjectTenantOwnerRoleOwner ObjectTenantOwnerRole = "OWNER"
)

// Valid indicates if r is any of the valid values for ObjectTenantOwnerRole
func (r ObjectTenantOwnerRole) Valid() bool {
	switch r {
	case
		ObjectTenantOwnerRoleUnspecified,
		ObjectTenantOwnerRoleTenant,
		ObjectTenantOwnerRoleOwner:
		return true
	}

	return false
}

// Validate returns an error if r is none of the valid values for ObjectTenantOwnerRole
func (r ObjectTenantOwnerRole) Validate() error {
	if !r.Valid() {
		return fmt.Errorf("invalid value %#v for type finance.ObjectTenantOwnerRole", r)
	}

	return nil
}

// Enums returns all valid values for ObjectTenantOwnerRole
func (ObjectTenantOwnerRole) Enums() []ObjectTenantOwnerRole {
	return []ObjectTenantOwnerRole{
		ObjectTenantOwnerRoleUnspecified,
		ObjectTenantOwnerRoleTenant,
		ObjectTenantOwnerRoleOwner,
	}
}

// EnumStrings returns all valid values for ObjectTenantOwnerRole as strings
func (ObjectTenantOwnerRole) EnumStrings() []string {
	return []string{
		"UNSPECIFIED",
		"TENANT",
		"OWNER",
	}
}

// String implements the fmt.Stringer interface for ObjectTenantOwnerRole
func (r ObjectTenantOwnerRole) String() string {
	return string(r)
}

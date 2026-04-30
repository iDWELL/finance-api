package finance

import (
	"context"
	"fmt"
	"net/url"
)

// ObjectRole represents a role assignment linking a user to a real estate object.
type ObjectRole struct {
	Name     string `json:"Name"`
	ObjectID string `json:"ObjectID"`
	Email    string `json:"Email"`
}

// PostObjectRoles upserts role assignments for real estate objects via the iDWELL Finance API.
//
// Arguments:
//   - ctx:    Context for the HTTP request (for cancellation and timeouts)
//   - apiKey: API key (bearer token) for the iDWELL Finance API
//   - roles:  Role assignments to insert or update
//   - source: String describing the source of the data; use your company or service name
//
// API endpoint: https://idwell.ai/api/public/masterdata/real-estate-object-roles
func PostObjectRoles(ctx context.Context, apiKey string, roles []*ObjectRole, source string) error {
	vals := make(url.Values)
	if source != "" {
		vals.Set("source", source)
	}
	response, err := postJSON(ctx, apiKey, "/masterdata/real-estate-object-roles", vals, roles)
	if err != nil {
		return err
	}
	defer func() { _ = response.Body.Close() }()
	if response.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}
	return nil
}

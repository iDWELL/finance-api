package finance

import (
	"context"
	"fmt"
	"net/url"
)

// ObjectGroup represents a group of real estate objects identified by an external ID.
type ObjectGroup struct {
	ExternalID string   `json:"ExternalID"`
	Name       string   `json:"Name"`
	AreaCode   *string  `json:"AreaCode"`
	ObjectIDs  []string `json:"ObjectIDs"`
}

// PostObjectGroups upserts real estate object groups via the iDWELL Finance API.
//
// Arguments:
//   - ctx:    Context for the HTTP request (for cancellation and timeouts)
//   - apiKey: API key (bearer token) for the iDWELL Finance API
//   - groups: Object groups to insert or update
//   - source: String describing the source of the data; use your company or service name
//
// API endpoint: https://idwell.ai/api/public/masterdata/real-estate-object-groups
func PostObjectGroups(ctx context.Context, apiKey string, groups []ObjectGroup, source string) error {
	vals := make(url.Values)
	if source != "" {
		vals.Set("source", source)
	}
	response, err := postJSON(ctx, apiKey, "/masterdata/real-estate-object-groups", vals, groups)
	if err != nil {
		return err
	}
	defer func() { _ = response.Body.Close() }()
	if response.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}
	return nil
}

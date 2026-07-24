package finance

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// UserAccessRevokeRequest identifies a user by email
// whose portal access should be revoked.
type UserAccessRevokeRequest struct {
	Email string `json:"email"`
}

// UserAccessRevokeResult reports the outcome of revoking
// portal access from a single user.
type UserAccessRevokeResult struct {
	Email   string `json:"email"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// RevokeUserAccess revokes portal access from users via the iDWELL Finance API.
// The portal (client company) is determined by the apiKey.
// Only the access assignment is removed, the user records themselves are kept.
// Revoking access from a user that has no access is not an error,
// but a user that can't be found by email is reported as a per-entry error.
//
// Arguments:
//   - ctx:      Context for the HTTP request (for cancellation and timeouts)
//   - apiKey:   API key (bearer token) for the iDWELL Finance API
//   - requests: Users to revoke access from, identified by email
//
// API endpoint: https://idwell.ai/api/public/user/access/revoke
func RevokeUserAccess(ctx context.Context, apiKey string, requests []UserAccessRevokeRequest) ([]UserAccessRevokeResult, error) {
	response, err := postJSON(ctx, apiKey, "/user/access/revoke", nil, requests)
	if err != nil {
		return nil, err
	}

	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	var results []UserAccessRevokeResult
	if err := json.NewDecoder(response.Body).Decode(&results); err != nil {
		return nil, err
	}

	return results, nil
}

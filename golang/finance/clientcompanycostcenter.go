package finance

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/domonda/go-types/uu"
	"github.com/pkg/errors"
)

// ClientCompanyCostCenter is a cost center belonging to a `ClientCompany`.
type ClientCompanyCostCenter struct {
	ID     uu.ID  `json:"rowId"`
	Number string `json:"number"`
}

type createClientCompanyCostCenterResponse struct {
	Data struct {
		CreateClientCompanyCostCenter struct {
			ClientCompanyCostCenter ClientCompanyCostCenter `json:"clientCompanyCostCenter"`
		} `json:"createClientCompanyCostCenter"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

// CreateClientCompanyCostCenter creates a cost-center master for the client
// company scoped by apiKey via the `createClientCompanyCostCenter` mutation.
// The currency defaults to EUR on the API side. Returns the new cost center.
func CreateClientCompanyCostCenter(ctx context.Context, apiKey, number, description string) (ClientCompanyCostCenter, error) {
	query := `mutation($number: String!, $description: String!) {
  createClientCompanyCostCenter(input: {
    number: $number,
    description: $description
  }) {
    clientCompanyCostCenter {
      rowId
      number
    }
  }
}`

	body, err := json.Marshal(map[string]any{
		"query": query,
		"variables": map[string]string{
			"number":      number,
			"description": description,
		},
	})
	if err != nil {
		return ClientCompanyCostCenter{}, errors.Wrap(err, "failed to marshal json data")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURLFromCtx(ctx)+GraphqlEndpoint, bytes.NewBuffer(body))
	if err != nil {
		return ClientCompanyCostCenter{}, errors.Wrap(err, "failed to create request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := httpClientFromCtx(ctx).Do(req) //nolint:gosec // intentional HTTP call to API URL from context
	if err != nil {
		return ClientCompanyCostCenter{}, errors.Wrap(err, "failed to execute request")
	}

	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return ClientCompanyCostCenter{}, errors.Wrap(err, "failed to read response body")
	}

	if resp.StatusCode != http.StatusOK {
		return ClientCompanyCostCenter{}, fmt.Errorf("create client company cost center: unexpected status code: %d, body: %s", resp.StatusCode, string(respBody))
	}

	var result createClientCompanyCostCenterResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return ClientCompanyCostCenter{}, errors.Wrap(err, "failed to unmarshal response body")
	}

	if len(result.Errors) > 0 {
		return ClientCompanyCostCenter{}, fmt.Errorf("create client company cost center: %s", result.Errors[0].Message)
	}

	return result.Data.CreateClientCompanyCostCenter.ClientCompanyCostCenter, nil
}

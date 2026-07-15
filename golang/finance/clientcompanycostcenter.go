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
	query := `mutation($number: NonEmptyText!, $description: NonEmptyText!) {
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

type updateClientCompanyCostCenterResponse struct {
	Data struct {
		UpdateClientCompanyCostCenter struct {
			ClientCompanyCostCenter ClientCompanyCostCenter `json:"clientCompanyCostCenter"`
		} `json:"updateClientCompanyCostCenter"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

// UpdateClientCompanyCostCenter updates the cost-center master identified by id
// for the client company scoped by apiKey via the `updateClientCompanyCostCenter`
// mutation. If currency is omitted, the cost center's current currency is kept
// unchanged; passing currency sets it to currency[0]. Returns the updated cost center.
func UpdateClientCompanyCostCenter(ctx context.Context, apiKey string, id uu.ID, number, description string, currency ...string) (ClientCompanyCostCenter, error) {
	query := `mutation($id: UUID!, $number: NonEmptyText!, $description: NonEmptyText!, $currency: CurrencyCode) {
  updateClientCompanyCostCenter(input: {
    id: $id,
    number: $number,
    description: $description,
    currency: $currency
  }) {
    clientCompanyCostCenter {
      rowId
      number
    }
  }
}`

	variables := map[string]any{
		"id":          id.String(),
		"number":      number,
		"description": description,
		"currency":    nil,
	}
	if len(currency) > 0 {
		variables["currency"] = currency[0]
	}

	body, err := json.Marshal(map[string]any{
		"query":     query,
		"variables": variables,
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
		return ClientCompanyCostCenter{}, fmt.Errorf("update client company cost center: unexpected status code: %d, body: %s", resp.StatusCode, string(respBody))
	}

	var result updateClientCompanyCostCenterResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return ClientCompanyCostCenter{}, errors.Wrap(err, "failed to unmarshal response body")
	}

	if len(result.Errors) > 0 {
		return ClientCompanyCostCenter{}, fmt.Errorf("update client company cost center: %s", result.Errors[0].Message)
	}

	return result.Data.UpdateClientCompanyCostCenter.ClientCompanyCostCenter, nil
}

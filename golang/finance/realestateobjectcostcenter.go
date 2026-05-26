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

// RealEstateObjectCostCenter is the association between a `RealEstateObject`
// and a `ClientCompanyCostCenter`. A cost center can be linked to many real
// estate objects and vice versa.
type RealEstateObjectCostCenter struct {
	ID                        uu.ID `json:"rowId"`
	ObjectInstanceID          uu.ID `json:"objectInstanceId"`
	ClientCompanyCostCenterID uu.ID `json:"clientCompanyCostCenterId"`
}

type createRealEstateObjectCostCenterResponse struct {
	Data struct {
		CreateRealEstateObjectCostCenter struct {
			RealEstateObjectCostCenter RealEstateObjectCostCenter `json:"realEstateObjectCostCenter"`
		} `json:"createRealEstateObjectCostCenter"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

// CreateRealEstateObjectCostCenter associates a `ClientCompanyCostCenter`
// with a `RealEstateObject` via the GraphQL `createRealEstateObjectCostCenter`
// mutation. Returns the new association row.
func CreateRealEstateObjectCostCenter(ctx context.Context, apiKey string, objectInstanceID, clientCompanyCostCenterID uu.ID) (RealEstateObjectCostCenter, error) {
	query := `mutation($objectInstanceId: UUID!, $clientCompanyCostCenterId: UUID!) {
  createRealEstateObjectCostCenter(input: {
    objectInstanceId: $objectInstanceId,
    clientCompanyCostCenterId: $clientCompanyCostCenterId
  }) {
    realEstateObjectCostCenter {
      rowId
      objectInstanceId
      clientCompanyCostCenterId
    }
  }
}`

	body, err := json.Marshal(map[string]any{
		"query": query,
		"variables": map[string]string{
			"objectInstanceId":          objectInstanceID.String(),
			"clientCompanyCostCenterId": clientCompanyCostCenterID.String(),
		},
	})
	if err != nil {
		return RealEstateObjectCostCenter{}, errors.Wrap(err, "failed to marshal json data")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURLFromCtx(ctx)+GraphqlEndpoint, bytes.NewBuffer(body))
	if err != nil {
		return RealEstateObjectCostCenter{}, errors.Wrap(err, "failed to create request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := httpClientFromCtx(ctx).Do(req) //nolint:gosec // intentional HTTP call to API URL from context
	if err != nil {
		return RealEstateObjectCostCenter{}, errors.Wrap(err, "failed to execute request")
	}

	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return RealEstateObjectCostCenter{}, errors.Wrap(err, "failed to read response body")
	}

	if resp.StatusCode != http.StatusOK {
		return RealEstateObjectCostCenter{}, fmt.Errorf("create real estate object cost center: unexpected status code: %d, body: %s", resp.StatusCode, string(respBody))
	}

	var result createRealEstateObjectCostCenterResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return RealEstateObjectCostCenter{}, errors.Wrap(err, "failed to unmarshal response body")
	}

	if len(result.Errors) > 0 {
		return RealEstateObjectCostCenter{}, fmt.Errorf("create real estate object cost center: %s", result.Errors[0].Message)
	}

	return result.Data.CreateRealEstateObjectCostCenter.RealEstateObjectCostCenter, nil
}

// DeleteRealEstateObjectCostCenter removes a `ClientCompanyCostCenter`
// association from a `RealEstateObject` by the association's row ID via the
// GraphQL `deleteRealEstateObjectCostCenter` mutation.
func DeleteRealEstateObjectCostCenter(ctx context.Context, apiKey string, id uu.ID) error {
	query := `mutation($id: UUID!) {
  deleteRealEstateObjectCostCenter(input: { id: $id }) {
    realEstateObjectCostCenter { rowId }
  }
}`

	body, err := json.Marshal(map[string]any{
		"query":     query,
		"variables": map[string]string{"id": id.String()},
	})
	if err != nil {
		return errors.Wrap(err, "failed to marshal json data")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURLFromCtx(ctx)+GraphqlEndpoint, bytes.NewBuffer(body))
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := httpClientFromCtx(ctx).Do(req) //nolint:gosec // intentional HTTP call to API URL from context
	if err != nil {
		return errors.Wrap(err, "failed to execute request")
	}

	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read response body")
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("delete real estate object cost center: unexpected status code: %d, body: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return errors.Wrap(err, "failed to unmarshal response body")
	}

	if len(result.Errors) > 0 {
		return fmt.Errorf("delete real estate object cost center: %s", result.Errors[0].Message)
	}

	return nil
}

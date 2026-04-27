package finance

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/domonda/go-types/country"
	"github.com/domonda/go-types/email"
	"github.com/domonda/go-types/notnull"
	"github.com/domonda/go-types/nullable"
	"github.com/domonda/go-types/uu"
	"github.com/domonda/go-types/vat"
	"github.com/pkg/errors"
)

type MainLocation struct {
	ID      uu.ID                  `json:"rowId"`
	Country country.Code           `json:"country"`
	Zip     nullable.TrimmedString `json:"zip"`
	City    nullable.TrimmedString `json:"city"`
	Street  nullable.TrimmedString `json:"street"`
	Phone   nullable.TrimmedString `json:"phone"`
	Email   email.NullableAddress  `json:"email"`
	VatNo   vat.NullableID         `json:"vatNo"`
	TaxNo   nullable.TrimmedString `json:"taxNo"`
}

type Company struct {
	ID               uu.ID                  `json:"rowId"`
	Name             notnull.TrimmedString  `json:"name"`
	AlternativeNames nullable.StringArray   `json:"alternativeNames"`
	BrandName        nullable.TrimmedString `json:"brandName"`
	LegalForm        nullable.TrimmedString `json:"legalForm"`
	MainLocation     MainLocation           `json:"mainLocation"`
}

type CompanyResponse struct {
	Data struct {
		CurrentClientCompany struct {
			CompanyByCompanyRowID Company `json:"companyByCompanyRowId"`
		} `json:"currentClientCompany"`
	}
}

func GetCurrentCompany(ctx context.Context, apiKey string) (Company, error) {
	query := `{
  currentClientCompany{
    companyByCompanyRowId{
      rowId
      name
      alternativeNames
      brandName
      legalForm
      mainLocation{
        rowId
        zip
        country
        city
        street
        phone
        email
        vatNo
        taxNo
      }
    }
  }
}`

	jsonData := map[string]string{
		"query": query,
	}

	jsonValue, err := json.Marshal(jsonData)
	if err != nil {
		return Company{}, errors.Wrap(err, "failed to marshal json data")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURLFromCtx(ctx)+GraphqlEndpoint, bytes.NewBuffer(jsonValue))
	if err != nil {
		return Company{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	response, err := httpClientFromCtx(ctx).Do(req)
	if err != nil {
		return Company{}, fmt.Errorf("the HTTP request failed with error %w", err)
	}

	defer func() {
		if err := response.Body.Close(); err != nil {
			fmt.Printf("failed to close response body: %v\n", err)
		}
	}()

	if response.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(response.Body)
		return Company{}, fmt.Errorf("non-200 status code: %d, response: %s", response.StatusCode, string(bodyBytes))
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return Company{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var resp CompanyResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return Company{}, fmt.Errorf("JSON unmarshal failed: %w", err)
	}

	return resp.Data.CurrentClientCompany.CompanyByCompanyRowID, nil
}

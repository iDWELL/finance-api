package finance

import (
	"testing"

	"github.com/domonda/go-types/bank"
	"github.com/domonda/go-types/date"
	"github.com/domonda/go-types/money"
	"github.com/domonda/go-types/nullable"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ptrAmount(v float64) *money.Amount { a := money.Amount(v); return &a }
func ptrRate(v float64) *money.Rate     { r := money.Rate(v); return &r }

func validInvoice(t *testing.T) *Invoice {
	t.Helper()
	return &Invoice{Currency: "EUR"}
}

func TestInvoiceValidate_NilReceiver(t *testing.T) {
	t.Parallel()

	var inv *Invoice

	err := inv.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "<nil>")
}

func TestInvoiceValidate(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		build       func(*testing.T) *Invoice
		expectedErr string
	}{
		"minimal valid invoice": {
			build: func(t *testing.T) *Invoice { t.Helper(); return validInvoice(t) },
		},
		"invalid currency": {
			build:       func(t *testing.T) *Invoice { t.Helper(); inv := validInvoice(t); inv.Currency = "NOTACCY"; return inv },
			expectedErr: "invalid currency",
		},
		"negative net": {
			build:       func(t *testing.T) *Invoice { t.Helper(); inv := validInvoice(t); inv.Net = ptrAmount(-1); return inv },
			expectedErr: "net amount must not be negative",
		},
		"negative total": {
			build: func(t *testing.T) *Invoice {
				t.Helper()
				inv := validInvoice(t)
				inv.Total = ptrAmount(-0.01)

				return inv
			},
			expectedErr: "total amount must not be negative",
		},
		"total less than net": {
			build: func(t *testing.T) *Invoice {
				t.Helper()
				inv := validInvoice(t)
				inv.Net = ptrAmount(100)
				inv.Total = ptrAmount(50)

				return inv
			},
			expectedErr: "must not be smaller than net",
		},
		"total equal net is valid": {
			build: func(t *testing.T) *Invoice {
				t.Helper()
				inv := validInvoice(t)
				inv.Net = ptrAmount(100)
				inv.Total = ptrAmount(100)

				return inv
			},
		},
		"VAT percent negative": {
			build: func(t *testing.T) *Invoice {
				t.Helper()
				inv := validInvoice(t)
				inv.VATPercent = ptrRate(-1)

				return inv
			},
			expectedErr: "vat percent",
		},
		"VAT percent above 100": {
			build: func(t *testing.T) *Invoice {
				t.Helper()
				inv := validInvoice(t)
				inv.VATPercent = ptrRate(101)

				return inv
			},
			expectedErr: "vat percent",
		},
		"VAT percent 0 valid": {
			build: func(t *testing.T) *Invoice {
				t.Helper()
				inv := validInvoice(t)
				inv.VATPercent = ptrRate(0)

				return inv
			},
		},
		"VAT percentage out of range": {
			build: func(t *testing.T) *Invoice {
				t.Helper()
				inv := validInvoice(t)
				inv.VATPercentages = nullable.FloatArray{20, -1}

				return inv
			},
			expectedErr: "vat percentage[1]",
		},
		"VAT amount negative": {
			build: func(t *testing.T) *Invoice {
				t.Helper()
				inv := validInvoice(t)
				inv.VATAmounts = nullable.FloatArray{-5}

				return inv
			},
			expectedErr: "vat amount[0]",
		},
		"discount percent negative": {
			build: func(t *testing.T) *Invoice {
				t.Helper()
				inv := validInvoice(t)
				inv.DiscountPercent = ptrRate(-1)

				return inv
			},
			expectedErr: "discount percent",
		},
		"discount percent above 100": {
			build: func(t *testing.T) *Invoice {
				t.Helper()
				inv := validInvoice(t)
				inv.DiscountPercent = ptrRate(200)

				return inv
			},
			expectedErr: "discount percent",
		},
		"conversion rate zero": {
			build: func(t *testing.T) *Invoice {
				t.Helper()
				inv := validInvoice(t)
				inv.ConversionRate = ptrRate(0)

				return inv
			},
			expectedErr: "conversion rate must be greater zero",
		},
		"conversion rate negative": {
			build: func(t *testing.T) *Invoice {
				t.Helper()
				inv := validInvoice(t)
				inv.ConversionRate = ptrRate(-1)

				return inv
			},
			expectedErr: "conversion rate must be greater zero",
		},
		"conversion rate positive valid": {
			build: func(t *testing.T) *Invoice {
				t.Helper()
				inv := validInvoice(t)
				inv.ConversionRate = ptrRate(1.2)

				return inv
			},
		},
		"deliveredFrom without deliveredUntil": {
			build: func(t *testing.T) *Invoice {
				t.Helper()
				inv := validInvoice(t)
				inv.DeliveredFrom = date.NullableDate("2024-01-01")

				return inv
			},
			expectedErr: "deliveredFrom date needs deliveredUntil",
		},
		"deliveredFrom after deliveredUntil": {
			build: func(t *testing.T) *Invoice {
				t.Helper()
				inv := validInvoice(t)
				inv.DeliveredFrom = date.NullableDate("2024-01-10")
				inv.DeliveredUntil = date.NullableDate("2024-01-01")

				return inv
			},
			expectedErr: "must not be after",
		},
		"deliveredFrom same as deliveredUntil valid": {
			build: func(t *testing.T) *Invoice {
				t.Helper()
				inv := validInvoice(t)
				inv.DeliveredFrom = date.NullableDate("2024-01-01")
				inv.DeliveredUntil = date.NullableDate("2024-01-01")

				return inv
			},
		},
		"invalid IBAN": {
			build: func(t *testing.T) *Invoice {
				t.Helper()
				inv := validInvoice(t)
				inv.IBAN = bank.NullableIBAN("NOTANIBAN")

				return inv
			},
			expectedErr: "invalid invoice IBAN",
		},
		"invalid BIC": {
			build: func(t *testing.T) *Invoice {
				t.Helper()
				inv := validInvoice(t)
				inv.BIC = bank.NullableBIC("!!")

				return inv
			},
			expectedErr: "invalid invoice BIC",
		},
		"empty cost center key": {
			build: func(t *testing.T) *Invoice {
				t.Helper()
				inv := validInvoice(t)
				inv.CostCenters = map[string]money.Amount{"": 100}

				return inv
			},
			expectedErr: "empty costCenter string",
		},
		"cost center zero amount": {
			build: func(t *testing.T) *Invoice {
				t.Helper()
				inv := validInvoice(t)
				inv.CostCenters = map[string]money.Amount{"CC1": 0}

				return inv
			},
			expectedErr: "must not be zero",
		},
		"cost centers sum exceeds net": {
			build: func(t *testing.T) *Invoice {
				t.Helper()
				inv := validInvoice(t)
				inv.Net = ptrAmount(100)
				inv.Total = ptrAmount(120)
				inv.CostCenters = map[string]money.Amount{"CC1": 101}

				return inv
			},
			expectedErr: "greater than invoice net sum",
		},
		"cost centers sum equal net valid": {
			build: func(t *testing.T) *Invoice {
				t.Helper()
				inv := validInvoice(t)
				inv.Net = ptrAmount(100)
				inv.Total = ptrAmount(120)
				inv.CostCenters = map[string]money.Amount{"CC1": 50, "CC2": 50}

				return inv
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := tc.build(t).Validate()
			if tc.expectedErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestInvoiceValidate_DeliveryNoteNumbersTrimmed(t *testing.T) {
	t.Parallel()

	inv := validInvoice(t)
	inv.DeliveryNoteNumbers = []string{"  DN-001  ", "  ", "DN-002", ""}
	require.NoError(t, inv.Validate())
	assert.Equal(t, []string{"DN-001", "DN-002"}, inv.DeliveryNoteNumbers)
}

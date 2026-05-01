package finance

import (
	"testing"

	"github.com/domonda/go-types/bank"
	"github.com/domonda/go-types/money"
	"github.com/domonda/go-types/notnull"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validBankAccount(t *testing.T) *BankAccount {
	t.Helper()

	return &BankAccount{
		IBAN:     bank.IBAN("AT611904300234573201"),
		BIC:      bank.BIC("OPSKATWW"),
		Currency: money.Currency("EUR"),
		Holder:   notnull.TrimmedString("Test Holder"),
	}
}

func TestBankAccountValidate(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		build       func(*testing.T) *BankAccount
		expectedErr string
	}{
		"valid bank account": {
			build: func(t *testing.T) *BankAccount { t.Helper(); return validBankAccount(t) },
		},
		"invalid IBAN": {
			build:       func(t *testing.T) *BankAccount { t.Helper(); a := validBankAccount(t); a.IBAN = "NOTANIBAN"; return a },
			expectedErr: "BankAccount.IBAN",
		},
		"invalid BIC": {
			build:       func(t *testing.T) *BankAccount { t.Helper(); a := validBankAccount(t); a.BIC = "!!"; return a },
			expectedErr: "BankAccount.BIC",
		},
		"invalid currency": {
			build: func(t *testing.T) *BankAccount {
				t.Helper()
				a := validBankAccount(t)
				a.Currency = "NOTACCY"

				return a
			},
			expectedErr: "BankAccount.Currency",
		},
		"empty holder": {
			build:       func(t *testing.T) *BankAccount { t.Helper(); a := validBankAccount(t); a.Holder = ""; return a },
			expectedErr: "BankAccount.Holder",
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

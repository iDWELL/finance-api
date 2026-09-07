package finance

import (
	"encoding/json"
	"testing"

	"github.com/domonda/go-types/account"
	"github.com/domonda/go-types/bank"
	"github.com/domonda/go-types/country"
	"github.com/domonda/go-types/date"
	"github.com/domonda/go-types/notnull"
	"github.com/domonda/go-types/nullable"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRealEstateObjectUnsentFields(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		jsonObject string
		expected   []string
	}{
		"all fields sent": {
			jsonObject: `{
				"Type": "WEG",
				"Number": "4200",
				"AccountingArea": "10",
				"UserAccount": "20",
				"Description": "note",
				"StreetAddress": "Hauptstraße 1",
				"AlternativeAddresses": ["Nebenstraße 2"],
				"ZipCode": "1010",
				"City": "Wien",
				"Country": "AT",
				"BankAccounts": [{"IBAN": "AT611904300234573201"}],
				"Active": true,
				"ManagementEnd": "2026-12-31"
			}`,
			expected: nil,
		},
		"accounting fields not sent": {
			jsonObject: `{
				"Type": "WEG",
				"Number": "4200",
				"Description": "note",
				"StreetAddress": "Hauptstraße 1",
				"AlternativeAddresses": ["Nebenstraße 2"],
				"ZipCode": "1010",
				"City": "Wien",
				"Country": "AT",
				"BankAccounts": [{"IBAN": "AT611904300234573201"}],
				"Active": true,
				"ManagementEnd": "2026-12-31"
			}`,
			expected: []string{"AccountingArea", "UserAccount"},
		},
		"accounting fields sent as null are sent": {
			jsonObject: `{
				"Type": "WEG",
				"Number": "4200",
				"AccountingArea": null,
				"UserAccount": null,
				"Description": "note",
				"StreetAddress": "Hauptstraße 1",
				"AlternativeAddresses": ["Nebenstraße 2"],
				"ZipCode": "1010",
				"City": "Wien",
				"Country": "AT",
				"BankAccounts": [{"IBAN": "AT611904300234573201"}],
				"Active": true,
				"ManagementEnd": "2026-12-31"
			}`,
			expected: nil,
		},
		"lower camel case field names are sent": {
			jsonObject: `{
				"type": "WEG",
				"number": "4200",
				"accountingArea": "10",
				"userAccount": "20",
				"description": "note",
				"streetAddress": "Hauptstraße 1",
				"alternativeAddresses": ["Nebenstraße 2"],
				"zipCode": "1010",
				"city": "Wien",
				"country": "AT",
				"bankAccounts": [{"IBAN": "AT611904300234573201"}],
				"active": true,
				"managementEnd": "2026-12-31"
			}`,
			expected: nil,
		},
		"only required fields sent": {
			jsonObject: `{
				"Type": "WEG",
				"Number": "4200",
				"StreetAddress": "Hauptstraße 1",
				"Country": "AT",
				"Active": true
			}`,
			expected: realEstateObjectOptionalFields,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var object RealEstateObject
			require.NoError(t, json.Unmarshal([]byte(tc.jsonObject), &object))
			require.NoError(t, object.Validate())

			assert.Equal(t, tc.expected, object.UnsentFields())
		})
	}
}

// TestRealEstateObjectUnsentFieldsWithoutJSON covers an object that was assembled
// in Go instead of being unmarshalled from JSON: every field of it is authoritative,
// so none of them counts as unsent even when it is empty.
func TestRealEstateObjectUnsentFieldsWithoutJSON(t *testing.T) {
	t.Parallel()

	object := RealEstateObject{
		Type:          RealEstateObjectTypeWEG,
		Number:        account.Number("4200"),
		StreetAddress: notnull.TrimmedString("Hauptstraße 1"),
		Country:       country.AT,
		Active:        true,
	}

	assert.Nil(t, object.UnsentFields())
}

// TestRealEstateObjectMarshalOmitsUnsetOptionalFields covers the sender side of a
// partial payload: a sender that leaves an optional field at its zero value must not
// send the field at all, because a field sent as null clears it for the receiver.
func TestRealEstateObjectMarshalOmitsUnsetOptionalFields(t *testing.T) {
	t.Parallel()

	object := RealEstateObject{
		Type:          RealEstateObjectTypeWEG,
		Number:        account.Number("4200"),
		Description:   nullable.TrimmedString("note"),
		StreetAddress: notnull.TrimmedString("Hauptstraße 1"),
		ZipCode:       nullable.TrimmedString("1010"),
		City:          nullable.TrimmedString("Wien"),
		Country:       country.AT,
		Active:        true,
	}

	jsonObject, err := json.Marshal(object)
	require.NoError(t, err)

	var fields map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(jsonObject, &fields))

	assert.NotContains(t, fields, "AccountingArea", "unset optional field must not be sent: %s", jsonObject)
	assert.NotContains(t, fields, "UserAccount", "unset optional field must not be sent: %s", jsonObject)
	assert.NotContains(t, fields, "AlternativeAddresses")
	assert.NotContains(t, fields, "BankAccounts")
	assert.NotContains(t, fields, "ManagementEnd")

	assert.Contains(t, fields, "Type")
	assert.Contains(t, fields, "Number")
	assert.Contains(t, fields, "StreetAddress")
	assert.Contains(t, fields, "Country")
	assert.Contains(t, fields, "Active")
	assert.Contains(t, fields, "Description")
	assert.Contains(t, fields, "ZipCode")
	assert.Contains(t, fields, "City")
}

// TestRealEstateObjectSendReceiveWithoutAccountingFields covers the round trip of
// the sync service payload: a sender that has no accounting data leaves the fields
// at their zero value and the receiver sees them as unsent instead of as cleared.
func TestRealEstateObjectSendReceiveWithoutAccountingFields(t *testing.T) {
	t.Parallel()

	sent := RealEstateObject{
		Type:          RealEstateObjectTypeWEG,
		Number:        account.Number("4200"),
		StreetAddress: notnull.TrimmedString("Hauptstraße 1"),
		ZipCode:       nullable.TrimmedString("1010"),
		City:          nullable.TrimmedString("Wien"),
		Country:       country.AT,
		Active:        true,
	}

	jsonObjects, err := json.Marshal([]*RealEstateObject{&sent})
	require.NoError(t, err)

	var received []RealEstateObject
	require.NoError(t, json.Unmarshal(jsonObjects, &received))
	require.Len(t, received, 1)

	assert.Equal(t,
		[]string{"AccountingArea", "UserAccount", "Description", "AlternativeAddresses", "BankAccounts", "ManagementEnd"},
		received[0].UnsentFields(),
	)
}

// TestRealEstateObjectMarshalSendsSetOptionalFields guards the opposite direction of
// the omitzero struct tags: a sender that owns an optional field still sends it,
// including the empty bank account slice that clears the bank accounts.
func TestRealEstateObjectMarshalSendsSetOptionalFields(t *testing.T) {
	t.Parallel()

	object := RealEstateObject{
		Type:           RealEstateObjectTypeWEG,
		Number:         account.Number("4200"),
		AccountingArea: account.NullableNumber("10"),
		UserAccount:    account.NullableNumber("20"),
		StreetAddress:  notnull.TrimmedString("Hauptstraße 1"),
		Country:        country.AT,
		BankAccounts:   []bank.Account{},
		Active:         true,
		ManagementEnd:  date.NullableDate("2026-12-31"),
	}

	jsonObject, err := json.Marshal(object)
	require.NoError(t, err)

	var fields map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(jsonObject, &fields))

	assert.Contains(t, fields, "AccountingArea")
	assert.Contains(t, fields, "UserAccount")
	assert.Contains(t, fields, "BankAccounts")
	assert.Contains(t, fields, "ManagementEnd")
}

package bill_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvoiceTaxTagsMigration(t *testing.T) {
	// Sample document taken from spanish examples.
	in := `{
		"$schema": "https://gobl.org/draft-0/bill/invoice",
		"uuid": "3aea7b56-59d8-4beb-90bd-f8f280d852a0",
		"type": "standard",
		"code": "SAMPLE-001",
		"issue_date": "2022-02-01",
		"currency": "EUR",
		"tax": {
			"tags": [
				"simplified"
			]
		},
		"supplier": {
			"name": "Provide One S.L.",
			"tax_id": {
				"country": "ES",
				"code": "B98602642"
			},
			"addresses": [
				{
					"num": "42",
					"street": "Calle Pradillo",
					"locality": "Madrid",
					"region": "Madrid",
					"code": "28002",
					"country": "ES"
				}
			],
			"emails": [
				{
					"addr": "billing@example.com"
				}
			]
		},
		"lines": [
			{
				"i": 1,
				"quantity": "20",
				"item": {
					"name": "Main product",
					"price": "90.00"
				},
				"sum": "1800.00",
				"discounts": [
					{
						"percent": "10%",
						"amount": "180.00",
						"reason": "Special discount"
					}
				],
				"taxes": [
					{
						"cat": "VAT",
						"rate": "standard",
						"percent": "21.0%"
					}
				],
				"total": "1620.00"
			},
			{
				"i": 2,
				"quantity": "1",
				"item": {
					"name": "Something else",
					"price": "10.00"
				},
				"sum": "10.00",
				"taxes": [
					{
						"cat": "VAT",
						"rate": "standard",
						"percent": "21.0%"
					}
				],
				"total": "10.00"
			}
		],
		"totals": {
			"sum": "1630.00",
			"total": "1630.00",
			"taxes": {
				"categories": [
					{
						"code": "VAT",
						"rates": [
							{
								"key": "standard",
								"base": "1630.00",
								"percent": "21.0%",
								"amount": "342.30"
							}
						],
						"amount": "342.30"
					}
				],
				"sum": "342.30"
			},
			"tax": "342.30",
			"total_with_tax": "1972.30",
			"payable": "1972.30"
		}
	}`
	inv := new(bill.Invoice)
	require.NoError(t, json.Unmarshal([]byte(in), inv))

	assert.Equal(t, "simplified", inv.GetTags()[0].String())
}

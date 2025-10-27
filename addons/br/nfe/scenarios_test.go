package nfe_test

import (
	"testing"

	"github.com/invopop/gobl/addons/br/nfe"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvoiceTypeScenarios(t *testing.T) {
	tests := []struct {
		name    string
		invType cbc.Key
		tags    []cbc.Key
		out     cbc.Code
	}{
		{
			name:    "Standard Invoice (NFe)",
			invType: bill.InvoiceTypeStandard,
			tags:    []cbc.Key{},
			out:     nfe.ModelNFe,
		},
		{
			name:    "Credit Note (NFe)",
			invType: bill.InvoiceTypeCreditNote,
			tags:    []cbc.Key{},
			out:     nfe.ModelNFe,
		},
		{
			name:    "Debit Note (NFe)",
			invType: bill.InvoiceTypeDebitNote,
			tags:    []cbc.Key{},
			out:     nfe.ModelNFe,
		},
		{
			name:    "Simplified Invoice (NFCe)",
			invType: bill.InvoiceTypeStandard,
			tags:    []cbc.Key{tax.TagSimplified},
			out:     nfe.ModelNFCe,
		},
		{
			name:    "Simplified Credit Note (NFCe)",
			invType: bill.InvoiceTypeCreditNote,
			tags:    []cbc.Key{tax.TagSimplified},
			out:     nfe.ModelNFCe,
		},
		{
			name:    "Simplified Debit Note (NFCe)",
			invType: bill.InvoiceTypeDebitNote,
			tags:    []cbc.Key{tax.TagSimplified},
			out:     nfe.ModelNFCe,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := validInvoice()
			i.Type = tt.invType
			i.SetTags(tt.tags...)
			require.NoError(t, i.Calculate())
			require.NotNil(t, i.Tax)
			require.NotNil(t, i.Tax.Ext)
			assert.Equal(t, tt.out, i.Tax.Ext[nfe.ExtKeyModel])
		})
	}
}

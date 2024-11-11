package mydata_test

import (
	"testing"

	"github.com/invopop/gobl/addons/gr/mydata"
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
		out     string
	}{
		{
			name:    "Default standard invoice",
			invType: bill.InvoiceTypeStandard,
			tags:    []cbc.Key{},
			out:     "2.1",
		},
		{
			name:    "Sales Invoice",
			invType: bill.InvoiceTypeStandard,
			tags:    []cbc.Key{mydata.TagGoods},
			out:     "1.1",
		},
		{
			name:    "Sales Invoice/Third Country Supplies",
			invType: bill.InvoiceTypeStandard,
			tags:    []cbc.Key{mydata.TagGoods, mydata.TagExport},
			out:     "1.3",
		},
		{
			name:    "Sales Invoice/Intra-community Supplies",
			invType: bill.InvoiceTypeStandard,
			tags:    []cbc.Key{mydata.TagGoods, mydata.TagExport, mydata.TagEU},
			out:     "1.2",
		},
		{
			name:    "Sales Invoice/Sale on Behalf of Third Parties",
			invType: bill.InvoiceTypeStandard,
			tags:    []cbc.Key{mydata.TagGoods, tax.TagSelfBilled},
			out:     "1.4",
		},
		{
			name:    "Service Rendered Invoice",
			invType: bill.InvoiceTypeStandard,
			tags:    []cbc.Key{mydata.TagServices},
			out:     "2.1",
		},
		{
			name:    "Third Country Service Rendered Invoice",
			invType: bill.InvoiceTypeStandard,
			tags:    []cbc.Key{mydata.TagServices, mydata.TagExport},
			out:     "2.3",
		},
		{
			name:    "Intra-community Service Rendered Invoice",
			invType: bill.InvoiceTypeStandard,
			tags:    []cbc.Key{mydata.TagServices, mydata.TagExport, mydata.TagEU},
			out:     "2.2",
		},
		{
			name:    "Credit Invoice/Associated",
			invType: bill.InvoiceTypeCreditNote,
			tags:    []cbc.Key{},
			out:     "5.1",
		},
		{
			name:    "Simplified Invoice",
			invType: bill.InvoiceTypeStandard,
			tags:    []cbc.Key{tax.TagSimplified},
			out:     "11.3",
		},
		{
			name:    "Retail Sales Receipt",
			invType: bill.InvoiceTypeStandard,
			tags:    []cbc.Key{mydata.TagGoods, tax.TagSimplified},
			out:     "11.1",
		},
		{
			name:    "Service Rendered Receipt",
			invType: bill.InvoiceTypeStandard,
			tags:    []cbc.Key{mydata.TagServices, tax.TagSimplified},
			out:     "11.2",
		},
		{
			name:    "Retail Sales Credit Note",
			invType: bill.InvoiceTypeCreditNote,
			tags:    []cbc.Key{tax.TagSimplified},
			out:     "11.4",
		},
		{
			name:    "Retail Sales Receipt on Behalf of Third Parties",
			invType: bill.InvoiceTypeCreditNote,
			tags:    []cbc.Key{mydata.TagGoods, tax.TagSimplified, tax.TagSelfBilled},
			out:     "11.5",
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
			assert.Equal(t, tt.out, i.Tax.Ext[mydata.ExtKeyInvoiceType].String())
		})
	}
}

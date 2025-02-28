package tax_test

import (
	"testing"

	"github.com/invopop/gobl/addons/es/facturae"
	"github.com/invopop/gobl/addons/es/tbai"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestCorrections(t *testing.T) {
	cds1 := tax.CorrectionSet{
		{
			Schema: bill.ShortSchemaInvoice,
			Extensions: []cbc.Key{
				facturae.ExtKeyCorrection,
			},
		},
		{
			Schema: "note/message",
			Extensions: []cbc.Key{
				"note-correction",
			},
			CopyTax: true,
		},
	}
	cds2 := tax.CorrectionSet{
		{
			Schema: bill.ShortSchemaInvoice,
			Extensions: []cbc.Key{
				tbai.ExtKeyCorrection,
			},
			CopyTax: true,
		},
	}
	cd1 := cds1.Def(bill.ShortSchemaInvoice)
	assert.Equal(t, facturae.ExtKeyCorrection, cd1.Extensions[0])
	assert.NotNil(t, cds1.Def("note/message"))
	assert.Nil(t, cds2.Def("note/message"))
	assert.False(t, cd1.CopyTax)

	cd2 := cds2.Def(bill.ShortSchemaInvoice)
	cd3 := cd1.Merge(cd2)
	assert.Len(t, cd3.Extensions, 2)
	assert.Contains(t, cd3.Extensions, facturae.ExtKeyCorrection)
	assert.Contains(t, cd3.Extensions, tbai.ExtKeyCorrection)
	assert.True(t, cd3.CopyTax)
}

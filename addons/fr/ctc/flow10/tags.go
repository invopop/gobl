package flow10

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

// French CTC Flow 10 tag keys. Absence of a tag means the document is B2B.
const (
	// TagB2C marks a document as reporting a business-to-consumer transaction.
	// Applied to both bill.Invoice and bill.Payment, it switches Flow 10
	// validation into the B2C rule set (no customer SIREN required, etc.).
	TagB2C cbc.Key = "b2c"
)

var b2cTagDef = &cbc.Definition{
	Key: TagB2C,
	Name: i18n.String{
		i18n.EN: "B2C",
		i18n.FR: "B2C",
	},
}

var invoiceTags = &tax.TagSet{
	Schema: bill.ShortSchemaInvoice,
	List:   []*cbc.Definition{b2cTagDef},
}

var paymentTags = &tax.TagSet{
	Schema: bill.ShortSchemaPayment,
	List:   []*cbc.Definition{b2cTagDef},
}

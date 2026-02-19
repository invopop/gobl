package en16931

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
)

const (
	// ExtKeyExemptionReason is used to provide a free-text reason for a VAT
	// exemption when no CEF VATEX code is available.
	ExtKeyExemptionReason cbc.Key = "en16931-exemption-reason"
)

var extensions = []*cbc.Definition{
	{
		Key: ExtKeyExemptionReason,
		Name: i18n.String{
			i18n.EN: "Exemption Reason",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Free-text description of the reason a line item is exempt from VAT (BT-120, BR-E-10).

				Use this extension when no specific CEF VATEX code (` + "`cef-vatex`" + `) applies.
				Exactly one of ` + "`en16931-exemption-reason`" + ` or ` + "`cef-vatex`" + ` must be present
				on every tax combo whose UNTDID tax category is **E** (Exempt).
			`),
		},
	},
}

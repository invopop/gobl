package bis

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
)

// IdentityKeyGreekMARK identifies a Greek MARK (Mydata Auto-Registration Key)
// number issued by the Greek Independent Authority for Public Revenue.
const IdentityKeyGreekMARK cbc.Key = "gr-mark"

// identities lists the party identity types recognised by the Peppol addon.
var identities = []*cbc.Definition{
	{
		Key: IdentityKeyGreekMARK,
		Name: i18n.String{
			i18n.EN: "Greek MARK Number",
		},
		Desc: i18n.String{
			i18n.EN: "Mydata Auto-Registration Key assigned by the Greek Independent " +
				"Authority for Public Revenue (IAPR) to each registered invoice.",
		},
	},
}

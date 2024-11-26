package in

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
)

// Extension keys used in India.
const (
	ExtKeySupplyPlace cbc.Key = "in-supply-place"
)

var extensions = []*cbc.KeyDefinition{
	{
		Key: ExtKeySupplyPlace,
		Name: i18n.String{
			i18n.EN: "Place of Supply",
			i18n.HI: "आपूर्ति का स्थान",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				The location to which the goods or services are supplied. In GST, this is referred to as the 'Place of Supply'.
			`),
			i18n.HI: here.Doc(`
				वह स्थान जहां वस्तुएं या सेवाएं प्रदान की जाती हैं। GST में इसे 'आपूर्ति का स्थान' कहा जाता है।
			`),
		},
	},
}

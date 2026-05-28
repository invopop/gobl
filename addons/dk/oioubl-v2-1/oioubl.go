// Package oioubl provides extensions and validations for the Danish OIOUBL 2.1
// standard used on the NemHandel e-invoicing network.
package oioubl

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

const (
	// Key identifies the Danish OIOUBL addon family. Individual versions
	// append a suffix; the family key is used as the fault-code namespace
	// so that rules carrying across versions keep stable codes.
	Key cbc.Key = "dk-oioubl"

	// V2_1 is the key for OIOUBL version 2.1, the current production version
	// used on the NemHandel network.
	V2_1 cbc.Key = Key + "-v2-1"
)

func init() {
	tax.RegisterAddonDef(newAddon())
	rules.RegisterWithGuard(
		Key.String(),
		rules.GOBL.Add("DK-OIOUBL"),
		is.InContext(tax.AddonIn(V2_1)),
		billInvoiceRules(),
		billStatusRules(),
	)
}

func newAddon() *tax.AddonDef {
	return &tax.AddonDef{
		Key: V2_1,
		Name: i18n.String{
			i18n.EN: "Danish OIOUBL 2.1",
			i18n.DA: "Dansk OIOUBL 2.1",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Support for the Danish OIOUBL 2.1 standard used on the NemHandel
				e-invoicing network, mandatory for business-to-government (B2G)
				transactions in Denmark since 2005.

				OIOUBL 2.1 is a national profile of UBL 2.1, maintained by
				Erhvervsstyrelsen (the Danish Business Authority). Unlike many
				European profiles it predates and does not extend EN 16931.

				This addon translates the OIOUBL Schematron rules (v1.17.1, live
				since 2026-05-18) into GOBL validations. OIOUBL 2.1 is scheduled
				to be replaced by NemHandel BIS 4 starting in 2028.
			`),
			i18n.DA: here.Doc(`
				Understøttelse af den danske OIOUBL 2.1-standard, som anvendes på
				NemHandel-netværket og har været obligatorisk for offentlige
				indkøb (B2G) i Danmark siden 2005.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "OIOUBL 2.1 documentation overview",
				},
				URL: "https://nemhandel.dk/vejledning-oioubl-dokumentationsoversigt",
			},
			{
				Title: i18n.String{
					i18n.EN: "OIOUBL Schematron v1.17.1 (released 2026-02-19)",
				},
				URL: "https://nemhandel.dk/oioubl-21-schematron-version-1171",
			},
		},
	}
}

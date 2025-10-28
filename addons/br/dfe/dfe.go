// Package dfe provides common extensions and identities for Brazilian fiscal documents.
// This addon contains the base functionality shared by both NF-e and NFS-e.
package dfe

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

const (
	// V1 is the key for the DFE addon version 1.
	V1 cbc.Key = "br-dfe-v1"
)

func init() {
	tax.RegisterAddonDef(newAddon())
}

func newAddon() *tax.AddonDef {
	return &tax.AddonDef{
		Key: V1,
		Name: i18n.String{
			i18n.EN: "Brazilian Electronic Fiscal Documents",
			i18n.PT: "Documentos Fiscais Eletrônicos Brasileiros",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Common base addon for all Brazilian electronic fiscal documents.
			`),
		},
		Extensions: extensions,
		Identities: identities,
		Normalizer: normalize,
	}
}

func normalize(doc any) {
	switch obj := doc.(type) {
	case *org.Identity:
		normalizeIdentity(obj)
	case *org.Party:
		normalizeParty(obj)
	}
}

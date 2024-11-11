// Package nfse handles extensions and validation rules to issue NFS-e in
// Brazil.
package nfse

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

const (
	// V1 identifies the NFS-e addon version
	V1 cbc.Key = "br-nfse-v1"
)

func init() {
	tax.RegisterAddonDef(newAddon())
}

func newAddon() *tax.AddonDef {
	return &tax.AddonDef{
		Key: V1,
		Name: i18n.String{
			i18n.EN: "Brazil NFS-e 1.X",
		},
		Extensions: extensions,
		Validator:  validate,
	}
}

func validate(doc any) error {
	switch obj := doc.(type) {
	case *bill.Invoice:
		return validateInvoice(obj)
	case *bill.Line:
		return validateLine(obj)
	case *org.Item:
		return validateItem(obj)
	}
	return nil
}

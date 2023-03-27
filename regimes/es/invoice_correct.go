package es

import "github.com/invopop/gobl/bill"

type invoiceCorrector struct{}

func (c invoiceCorrector) Correct(inv *bill.Invoice, opts *bill.Options) error {
	inv.Type = bill.InvoiceTypeCorrective
	if opts.Refund {
		// Spain doesn't support credit-notes, so we invert each
		// of the values.
		inv.Invert()
	}
	return nil
}

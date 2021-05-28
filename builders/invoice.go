package builders

import "github.com/invopop/gobl"

type Invoice interface {
	SetSupplier(p *gobl.Party) Invoice
	SetCustomer(p *gobl.Party) Invoice
	AddLine(l *gobl.InvoiceLine) Invoice
	CalculateTotals() Invoice
	Validate() error
	Document() *gobl.Document
}

type invoiceBuilder struct {
	doc *gobl.Document
}

func NewInvoice(doc *gobl.Document) Invoice {
	b := new(invoiceBuilder)
	b.doc = doc
	return b
}

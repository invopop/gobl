// Package bill provides models for dealing with Billing and specifically invoicing.
package bill

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/schema"
	"github.com/invopop/gobl/tax"
)

func init() {
	schema.Register(schema.GOBL.Add("bill"),
		// Primary schemas
		CorrectionOptions{},
		Delivery{},
		Invoice{},
		Order{},
		Payment{},
		// Sub-schemas - used by primaries
		Charge{},
		Discount{},
		Line{},
		Ordering{},
		PaymentDetails{},
		DeliveryDetails{},
		Tax{},
		Totals{},
	)
}

// Constants used to help identify document schemas
const (
	ShortSchemaOrder    = "bill/order"
	ShortSchemaDelivery = "bill/delivery"
	ShortSchemaInvoice  = "bill/invoice"
	ShortSchemaPayment  = "bill/payment"
)

// billable defines a contract for document types within the bill package that
// share common structure and behavior (e.g. invoice, payment, delivery or
// order). It defines the set of methods these types must implement to enable
// generic handling.
//
// The methods are exported to allow for generic handling from outside the
// package. However, the interface itself is kept private to promote the use of
// small, purpose-built interfaces for each use case.
type billable interface {
	// From tax.Regime
	RegimeDef() *tax.RegimeDef

	// From tax.Tags
	HasTags(tags ...cbc.Key) bool
	GetTags() []cbc.Key

	GetIssueDate() cal.Date
	GetIssueTime() *cal.Time
	GetValueDate() *cal.Date
	GetTax() *Tax
	GetPreceding() []*org.DocumentRef
	GetSupplier() *org.Party
	GetCustomer() *org.Party
	GetCurrency() currency.Code
	GetExchangeRates() []*currency.ExchangeRate
	GetLines() []*Line
	GetDiscounts() []*Discount
	GetCharges() []*Charge
	GetPaymentDetails() *PaymentDetails
	GetTotals() *Totals
	GetComplements() []*schema.Object

	SetCode(cbc.Code)
	SetIssueDate(cal.Date)
	SetIssueTime(*cal.Time)
	SetCurrency(currency.Code)
	SetTotals(*Totals)
}

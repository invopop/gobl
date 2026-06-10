package bill

import (
	"github.com/invopop/gobl/norm"
)

// init registers the intrinsic normalizers for the bill package's document
// types. Each function performs only the type's own normalization; the norm
// engine handles recursion into nested values, the global cleaning of codes
// and extensions, and the application of regime and addon normalizers.
//
// Types whose only normalization was code/extension cleaning (Status,
// StatusLine, Reason, Action, Fault, Ordering, LineCharge, LineDiscount) have
// no normalizer here: the global cbc.Code and tax.Extensions normalizers cover
// them automatically.
func init() {
	norm.Register(
		norm.For(normalizeInvoice),
		norm.For(normalizeDelivery),
		norm.For(normalizePayment),
		norm.For(normalizeOrder),
		norm.For(normalizeLine),
		norm.For(normalizeSubLine),
		norm.For(normalizeBillTax),
		norm.For(normalizeCharge),
		norm.For(normalizeDiscount),
	)
}

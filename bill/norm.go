package bill

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/tax"
)

// init registers the intrinsic normalizers for the bill package's document
// types. Each function performs only the type's own normalization; the norm
// engine handles recursion into nested values and the application of regime
// and addon normalizers.
func init() {
	norm.Register("bill",
		norm.For(normalizeInvoice),
		norm.For(normalizeDelivery),
		norm.For(normalizePayment),
		norm.For(normalizeOrder),
		norm.For(normalizeStatus),
		norm.For(normalizeStatusLine),
		norm.For(normalizeReason),
		norm.For(normalizeAction),
		norm.For(normalizeFault),
		norm.For(normalizeOrdering),
		norm.For(normalizeLine),
		norm.For(normalizeSubLine),
		norm.For(normalizeBillTax),
		norm.For(normalizeCharge),
		norm.For(normalizeDiscount),
		norm.For(normalizeLineCharge),
		norm.For(normalizeLineDiscount),
	)
}

func normalizeInvoice(inv *Invoice) {
	if inv.Type == cbc.KeyEmpty {
		inv.Type = InvoiceTypeStandard
	}
	inv.Series = cbc.NormalizeCode(inv.Series)
	inv.Code = cbc.NormalizeCode(inv.Code)
}

func normalizeDelivery(dlv *Delivery) {
	if dlv.Type == cbc.KeyEmpty {
		dlv.Type = DeliveryTypeAdvice
	}
	dlv.Series = cbc.NormalizeCode(dlv.Series)
	dlv.Code = cbc.NormalizeCode(dlv.Code)
}

func normalizePayment(pmt *Payment) {
	if pmt.Type == cbc.KeyEmpty {
		pmt.Type = PaymentTypeReceipt
	}
	pmt.Series = cbc.NormalizeCode(pmt.Series)
	pmt.Code = cbc.NormalizeCode(pmt.Code)
}

func normalizeOrder(ord *Order) {
	if ord.Type == cbc.KeyEmpty {
		ord.Type = OrderTypePurchase
	}
	ord.Series = cbc.NormalizeCode(ord.Series)
	ord.Code = cbc.NormalizeCode(ord.Code)
}

func normalizeStatus(st *Status) {
	st.Series = cbc.NormalizeCode(st.Series)
	st.Code = cbc.NormalizeCode(st.Code)
	st.Ext = st.Ext.Clean()
}

func normalizeStatusLine(sl *StatusLine) {
	if sl == nil {
		return
	}
	sl.Ext = sl.Ext.Clean()
}

func normalizeReason(r *Reason) {
	if r == nil {
		return
	}
	r.Ext = r.Ext.Clean()
}

func normalizeAction(a *Action) {
	if a == nil {
		return
	}
	a.Ext = a.Ext.Clean()
}

func normalizeFault(f *Fault) {
	if f == nil {
		return
	}
	f.Code = cbc.NormalizeCode(f.Code)
}

func normalizeOrdering(o *Ordering) {
	if o == nil {
		return
	}
	o.Code = cbc.NormalizeCode(o.Code)
	o.Cost = cbc.NormalizeCode(o.Cost)
}

func normalizeLine(l *Line) {
	if l == nil {
		return
	}
	normalizeLineItemPrice(l)
	l.Taxes = tax.CleanSet(l.Taxes)
	l.Discounts = CleanLineDiscounts(l.Discounts)
	l.Charges = CleanLineCharges(l.Charges)
	l.Breakdown = CleanSubLines(l.Breakdown)
}

func normalizeSubLine(sl *SubLine) {
	if sl == nil {
		return
	}
	normalizeSubLineItemPrice(sl)
	sl.Discounts = CleanLineDiscounts(sl.Discounts)
	sl.Charges = CleanLineCharges(sl.Charges)
}

func normalizeBillTax(t *Tax) {
	if t == nil {
		return
	}
	// migration for old rounding rules
	switch t.Rounding {
	case "sum-then-round":
		t.Rounding = tax.RoundingRulePrecise
	case "round-then-sum":
		t.Rounding = tax.RoundingRuleCurrency
	}
	t.Ext = t.Ext.Clean()
}

func normalizeCharge(m *Charge) {
	if m == nil {
		return
	}
	m.Code = cbc.NormalizeCode(m.Code)
	m.Taxes = tax.CleanSet(m.Taxes)
	m.Ext = m.Ext.Clean()
}

func normalizeDiscount(m *Discount) {
	if m == nil {
		return
	}
	m.Code = cbc.NormalizeCode(m.Code)
	m.Taxes = tax.CleanSet(m.Taxes)
	m.Ext = m.Ext.Clean()
}

func normalizeLineCharge(lc *LineCharge) {
	lc.Code = cbc.NormalizeCode(lc.Code)
	lc.Ext = lc.Ext.Clean()
}

func normalizeLineDiscount(ld *LineDiscount) {
	ld.Code = cbc.NormalizeCode(ld.Code)
	ld.Ext = ld.Ext.Clean()
}

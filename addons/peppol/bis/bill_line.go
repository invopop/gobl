package bis

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
)

func billLineRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.Assert("07", "line period start must be within invoice period (PEPPOL-EN16931-R110)",
			is.Func("line period start", lineStartsWithinInvoice),
		),
		rules.Assert("08", "line period end must be within invoice period (PEPPOL-EN16931-R111)",
			is.Func("line period end", lineEndsWithinInvoice),
		),
	)
}

func lineStartsWithinInvoice(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil || inv.Ordering == nil || inv.Ordering.Period == nil {
		return true
	}
	invStart := inv.Ordering.Period.Start
	if invStart.IsZero() {
		return true
	}
	for _, line := range inv.Lines {
		if line == nil || line.Period == nil || line.Period.Start.IsZero() {
			continue
		}
		if line.Period.Start.Before(invStart.Date) {
			return false
		}
	}
	return true
}

func lineEndsWithinInvoice(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil || inv.Ordering == nil || inv.Ordering.Period == nil {
		return true
	}
	invEnd := inv.Ordering.Period.End
	if invEnd.IsZero() {
		return true
	}
	for _, line := range inv.Lines {
		if line == nil || line.Period == nil || line.Period.End.IsZero() {
			continue
		}
		if line.Period.End.After(invEnd.Date) {
			return false
		}
	}
	return true
}

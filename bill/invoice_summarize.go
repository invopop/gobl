package bill

import (
	"encoding/json"
	"fmt"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

// SummarizeLines replaces the invoice's lines with one line per group of lines
// that are taxed and classified in the same way, for documents with too many
// lines to be handled in full. The original lines are discarded, so the source
// document should be kept by the caller when the detail is still needed.
//
// Lines are grouped by their taxes, extensions, and seller, together with the
// item's key and extensions, so that the summary carries the same details that
// regimes and addons use to classify the original lines. Each summary line has
// a quantity of one and a price equal to the sum of the group's line totals.
// Document-level discounts, charges, and payment details are kept unchanged.
//
// The invoice is tagged with `summarized` and given a note with the number of
// lines replaced. After summarizing the invoice is recalculated, and any
// difference in the totals will raise an error and restore the original lines.
func (inv *Invoice) SummarizeLines() error {
	if inv.HasTags(tax.TagBypass) {
		return fmt.Errorf("cannot summarize an invoice with tag bypass")
	}
	if inv.HasTags(tax.TagSummarized) {
		return fmt.Errorf("invoice lines already summarized")
	}
	if err := ensureCalculated(inv); err != nil {
		return err
	}
	if inv.Totals == nil {
		return fmt.Errorf("cannot summarize an invoice without totals")
	}
	orig := *inv.Totals
	lines := inv.Lines
	tags := inv.GetTags()

	summary, count, err := summarizeLines(inv)
	if err != nil {
		inv.Lines = lines
		return restoreSummarized(inv, err)
	}
	inv.Lines = summary
	inv.SetTags(append(tags, tax.TagSummarized)...)
	if err := inv.Calculate(); err != nil {
		inv.Lines = lines
		inv.SetTags(tags...)
		return restoreSummarized(inv, err)
	}
	if err := compareSummarizedTotals(&orig, inv.Totals); err != nil {
		inv.Lines = lines
		inv.SetTags(tags...)
		return restoreSummarized(inv, err)
	}

	inv.Notes = append(inv.Notes, &org.Note{
		Key:  org.NoteKeyGeneral,
		Text: fmt.Sprintf("Lines summarized from %d original lines.", count),
	})
	return nil
}

// restoreSummarized recalculates the invoice after its original lines have
// been put back, as summarizing replaces their amounts, and returns the error
// that prevented the summary.
func restoreSummarized(inv *Invoice, err error) error {
	if cerr := inv.Calculate(); cerr != nil {
		return fmt.Errorf("%w (restoring: %s)", err, cerr.Error())
	}
	return err
}

// lineSummary accumulates the totals of a group of lines.
type lineSummary struct {
	first *Line
	count int
	total num.Amount
}

// summarizeLines groups the invoice's lines, returning a line for each group in
// the order they first appear, and the number of lines summarized. Line totals
// are recalculated in place without the final rounding, so that the sums keep
// the precision used to calculate the original taxes.
func summarizeLines(inv *Invoice) ([]*Line, int, error) {
	cur := inv.GetCurrency()
	if err := calculateLines(inv.Lines, cur, inv.GetExchangeRates(), roundingRule(inv)); err != nil {
		return nil, 0, fmt.Errorf("lines: %w", err)
	}
	groups := make(map[string]*lineSummary)
	order := make([]*lineSummary, 0)
	count := 0
	for i, l := range inv.Lines {
		if l == nil || l.Item == nil || l.Total == nil {
			// lines without a price have no effect on the totals
			continue
		}
		key, err := lineSummaryKey(l)
		if err != nil {
			return nil, 0, fmt.Errorf("lines: %d: %w", i, err)
		}
		g, ok := groups[key]
		if !ok {
			g = &lineSummary{first: l, total: cur.Def().Zero()}
			groups[key] = g
			order = append(order, g)
		}
		g.total = g.total.MatchPrecision(*l.Total).Add(*l.Total)
		g.count++
		count++
	}

	summary := make([]*Line, len(order))
	for i, g := range order {
		price := g.total
		summary[i] = &Line{
			Quantity: num.MakeAmount(1, 0),
			Item: &org.Item{
				Key:   g.first.Item.Key,
				Name:  fmt.Sprintf("Summary of %d lines", g.count),
				Price: &price,
				Ext:   g.first.Item.Ext,
			},
			Taxes:  g.first.Taxes,
			Seller: g.first.Seller,
			Ext:    g.first.Ext,
		}
	}
	return summary, count, nil
}

// lineSummaryKey identifies the group a line belongs to from the fields that
// are copied to the summary line. Any field later added to tax combos will
// be included automatically.
func lineSummaryKey(l *Line) (string, error) {
	data, err := json.Marshal(struct {
		Taxes   tax.Set        `json:"taxes,omitempty"`
		Ext     tax.Extensions `json:"ext,omitzero"`
		Seller  *org.Party     `json:"seller,omitempty"`
		ItemKey cbc.Key        `json:"item_key,omitempty"`
		ItemExt tax.Extensions `json:"item_ext,omitzero"`
	}{
		Taxes:   l.Taxes,
		Ext:     l.Ext,
		Seller:  l.Seller,
		ItemKey: l.Item.Key,
		ItemExt: l.Item.Ext,
	})
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// compareSummarizedTotals ensures the summarized invoice has the same totals
// as the original.
func compareSummarizedTotals(orig, t *Totals) error {
	pairs := []struct {
		name string
		a, b num.Amount
	}{
		{"sum", orig.Sum, t.Sum},
		{"total", orig.Total, t.Total},
		{"tax", orig.Tax, t.Tax},
		{"total with tax", orig.TotalWithTax, t.TotalWithTax},
		{"payable", orig.Payable, t.Payable},
	}
	for _, p := range pairs {
		if !p.a.Equals(p.b) {
			return fmt.Errorf("summarized invoice %s does not match %s != %s", p.name, p.a.String(), p.b.String())
		}
	}
	return nil
}

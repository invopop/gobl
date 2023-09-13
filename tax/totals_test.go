package tax_test

import (
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

// taxableLine is a very simple implementation of what the totals calculator requires.
type taxableLine struct {
	taxes  tax.Set
	amount num.Amount
}

func (tl *taxableLine) GetTaxes() tax.Set {
	return tl.taxes
}

func (tl *taxableLine) GetTotal() num.Amount {
	return tl.amount
}

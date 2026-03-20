package nfe_test

import (
"testing"

"github.com/invopop/gobl/addons/br/nfe"
"github.com/invopop/gobl/bill"
"github.com/invopop/gobl/num"
"github.com/invopop/gobl/org"
"github.com/invopop/gobl/regimes/br"
"github.com/invopop/gobl/rules"
"github.com/invopop/gobl/tax"
"github.com/stretchr/testify/assert"
)

func TestLineValidation(t *testing.T) {
tests := []struct {
name string
line *bill.Line
err  string
}{
{
name: "valid line with all required taxes",
line: &bill.Line{
Index:    1,
Quantity: num.MakeAmount(1, 0),
Item:     &org.Item{Name: "Test Item"},
Taxes: tax.Set{
{Category: br.TaxCategoryICMS},
{Category: br.TaxCategoryPIS},
{Category: br.TaxCategoryCOFINS},
},
},
},
{
name: "nil line",
line: nil,
},
{
name: "missing taxes",
line: &bill.Line{},
err:  "ICMS tax category is required",
},
{
name: "empty taxes",
line: &bill.Line{
Taxes: tax.Set{},
},
err: "ICMS tax category is required",
},
{
name: "missing ICMS tax",
line: &bill.Line{
Taxes: tax.Set{
{Category: br.TaxCategoryPIS},
{Category: br.TaxCategoryCOFINS},
},
},
err: "ICMS tax category is required",
},
{
name: "missing PIS tax",
line: &bill.Line{
Taxes: tax.Set{
{Category: br.TaxCategoryICMS},
{Category: br.TaxCategoryCOFINS},
},
},
err: "PIS tax category is required",
},
{
name: "missing COFINS tax",
line: &bill.Line{
Taxes: tax.Set{
{Category: br.TaxCategoryICMS},
{Category: br.TaxCategoryPIS},
},
},
err: "COFINS tax category is required",
},
}

for _, ts := range tests {
t.Run(ts.name, func(t *testing.T) {
err := rules.Validate(ts.line, tax.AddonContext(nfe.V4))
if ts.err == "" {
assert.NoError(t, err)
} else {
if assert.Error(t, err) {
assert.Contains(t, err.Error(), ts.err)
}
}
})
}
}

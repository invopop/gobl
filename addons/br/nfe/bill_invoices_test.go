package nfe_test

import (
"fmt"
"testing"

"github.com/invopop/gobl/addons/br/nfe"
"github.com/invopop/gobl/bill"
"github.com/invopop/gobl/cbc"
"github.com/invopop/gobl/num"
"github.com/invopop/gobl/org"
"github.com/invopop/gobl/pay"
"github.com/invopop/gobl/regimes/br"
"github.com/invopop/gobl/rules"
"github.com/invopop/gobl/tax"
"github.com/stretchr/testify/assert"
"github.com/stretchr/testify/require"
)

func TestInvoicesValidation(t *testing.T) {
t.Run("validates tax extensions", func(t *testing.T) {
inv := validInvoice()
require.NoError(t, inv.Calculate())
inv.Tax = nil
err := rules.Validate(inv)
assert.ErrorContains(t, err, "tax details are required")

inv.Tax = &bill.Tax{}
err = rules.Validate(inv)
assert.ErrorContains(t, err, "model extension is required")
assert.ErrorContains(t, err, "presence extension is required")

inv.Tax.Ext = tax.Extensions{
nfe.ExtKeyModel:    nfe.ModelNFe,
nfe.ExtKeyPresence: nfe.PresenceDelivery,
}
err = rules.Validate(inv)
assert.ErrorContains(t, err, "delivery presence not allowed for NFe")

inv.Tax.Ext[nfe.ExtKeyPresence] = nfe.PresenceInPerson
err = rules.Validate(inv)
assert.NoError(t, err)
})

t.Run("validates required notes", func(t *testing.T) {
inv := validInvoice()
require.NoError(t, inv.Calculate())
inv.Notes = nil
err := rules.Validate(inv)
assert.ErrorContains(t, err, "a reason note is required")

inv.Notes = []*org.Note{nil}
err = rules.Validate(inv)
assert.ErrorContains(t, err, "a reason note is required")

inv.Notes[0] = &org.Note{
Key:  org.NoteKeyGeneral,
Text: "General note",
}
err = rules.Validate(inv)
assert.ErrorContains(t, err, "a reason note is required")

inv.Notes[0].Key = org.NoteKeyReason
inv.Notes[0].Text = "1234567890123456789012345678901234567890123456789012345678901" // 61 chars
err = rules.Validate(inv)
assert.ErrorContains(t, err, "reason note text must be between 1 and 60 characters")

inv.Notes[0].Text = "123456789012345678901234567890123456789012345678901234567890" // 60 chars
err = rules.Validate(inv)
assert.NoError(t, err)
})

t.Run("validates payment when invoice is due", func(t *testing.T) {
inv := validInvoice()
require.NoError(t, inv.Calculate())
inv.Totals = &bill.Totals{}
inv.Payment = nil

inv.Totals.Due = &num.AmountZero
err := rules.Validate(inv)
assert.NoError(t, err)

inv.Totals.Due = nil
err = rules.Validate(inv)
assert.ErrorContains(t, err, "payment details are required")

inv.Totals.Due = num.NewAmount(1, 2)
err = rules.Validate(inv)
assert.ErrorContains(t, err, "payment details are required")

inv.Payment = &bill.PaymentDetails{}
err = rules.Validate(inv)
assert.ErrorContains(t, err, "payment instructions are required")

inv.Payment.Instructions = &pay.Instructions{Key: pay.MeansKeyCash, Ext: tax.Extensions{nfe.ExtKeyPaymentMeans: "01"}}
err = rules.Validate(inv)
assert.NoError(t, err)
})

t.Run("validates invoice totals due field", func(t *testing.T) {
inv := validInvoice()
require.NoError(t, inv.Calculate())
inv.Totals = &bill.Totals{}

inv.Totals.Due = num.NewAmount(-1, 2)
err := rules.Validate(inv)
assert.ErrorContains(t, err, "due amount must be zero or positive")

inv.Totals.Due = &num.AmountZero
err = rules.Validate(inv)
assert.NoError(t, err)

inv.Totals.Due = num.NewAmount(1, 2)
err = rules.Validate(inv)
assert.NoError(t, err)
})

t.Run("validates NFe presence when model is NFe", func(t *testing.T) {
inv := validInvoice()
require.NoError(t, inv.Calculate())
inv.Tax.Ext[nfe.ExtKeyModel] = nfe.ModelNFe
inv.Tax.Ext[nfe.ExtKeyPresence] = nfe.PresenceDelivery
err := rules.Validate(inv)
assert.ErrorContains(t, err, "delivery presence not allowed for NFe")

inv.Tax.Ext[nfe.ExtKeyPresence] = nfe.PresenceInPerson
err = rules.Validate(inv)
assert.NoError(t, err)
})

t.Run("validates NFCe presence when model is NFCe", func(t *testing.T) {
inv := validInvoice()
inv.Customer = nil // For NFCe, customer is optional, so remove it before Calculate
require.NoError(t, inv.Calculate())

inv.Tax.Ext[nfe.ExtKeyModel] = nfe.ModelNFCe
inv.Tax.Ext[nfe.ExtKeyPresence] = nfe.PresenceNotApplicable
err := rules.Validate(inv)
assert.ErrorContains(t, err, "NFCe presence must be in-person or delivery")

inv.Tax.Ext[nfe.ExtKeyPresence] = nfe.PresenceInPerson
err = rules.Validate(inv)
assert.NoError(t, err)
})
}

func TestInvoiceSeriesValidation(t *testing.T) {
tests := []struct {
series cbc.Code
err    string
}{
{series: "0"},
{series: "1"},
{series: "12"},
{series: "123"},
{series: "999"},
{series: "", err: "series is required"},
{series: "1000", err: "series format is invalid"},
{series: "abc", err: "series format is invalid"},
{series: "012", err: "series format is invalid"},
{series: "00", err: "series format is invalid"},
{series: "-3", err: "series format is invalid"},
}

for _, tt := range tests {
name := fmt.Sprintf("validates series %s", tt.series)
t.Run(name, func(t *testing.T) {
inv := validInvoice()
require.NoError(t, inv.Calculate())
inv.Series = tt.series
err := rules.Validate(inv)
if tt.err != "" {
assert.ErrorContains(t, err, tt.err)
} else {
assert.NoError(t, err)
}
})
}
}

func TestSupplierValidation(t *testing.T) {
t.Run("nil supplier", func(t *testing.T) {
inv := validInvoice()
require.NoError(t, inv.Calculate())
inv.Supplier = nil
err := rules.Validate(inv)
// NFe addon does not add a supplier-required rule; GOBL handles that
// When supplier is nil, NFe-specific supplier field rules are skipped
if err != nil {
assert.NotContains(t, err.Error(), "supplier name is required")
assert.NotContains(t, err.Error(), "supplier state registration identity is required")
}
})

t.Run("validates supplier name", func(t *testing.T) {
inv := validInvoice()
require.NoError(t, inv.Calculate())
inv.Supplier.Name = ""
err := rules.Validate(inv)
assert.ErrorContains(t, err, "supplier name is required")

inv.Supplier.Name = "Test Company"
err = rules.Validate(inv)
assert.NoError(t, err)
})

t.Run("validates supplier addresses required", func(t *testing.T) {
inv := validInvoice()
require.NoError(t, inv.Calculate())
inv.Supplier.Addresses = nil
err := rules.Validate(inv)
assert.ErrorContains(t, err, "supplier addresses are required")

inv.Supplier.Addresses = []*org.Address{}
err = rules.Validate(inv)
assert.ErrorContains(t, err, "supplier addresses are required")

inv.Supplier.Addresses = []*org.Address{nil}
err = rules.Validate(inv)
assert.ErrorContains(t, err, "address must not be empty")

inv.Supplier.Addresses = []*org.Address{
{
Street:   "Rua Test",
Number:   "100",
Locality: "São Paulo",
State:    "SP",
Code:     "01310100",
},
}
err = rules.Validate(inv)
assert.NoError(t, err)
})

t.Run("validates supplier state registration identity", func(t *testing.T) {
inv := validInvoice()
require.NoError(t, inv.Calculate())
inv.Supplier.Identities = nil
err := rules.Validate(inv)
assert.ErrorContains(t, err, "supplier state registration identity is required")

inv.Supplier.Identities = []*org.Identity{
{
Key:  nfe.IdentityKeyStateReg,
Code: "35503304557308",
},
}
err = rules.Validate(inv)
assert.NoError(t, err)
})

t.Run("validates supplier tax ID required", func(t *testing.T) {
inv := validInvoice()
require.NoError(t, inv.Calculate())
inv.Supplier.TaxID = nil
err := rules.Validate(inv)
assert.ErrorContains(t, err, "supplier tax ID is required")

inv.Supplier.TaxID = &tax.Identity{Country: "BR"}
err = rules.Validate(inv)
assert.ErrorContains(t, err, "supplier tax ID code is required")

inv.Supplier.TaxID.Code = "55263640000186"
err = rules.Validate(inv)
assert.NoError(t, err)
})

t.Run("validates supplier municipality extension when addresses exist", func(t *testing.T) {
inv := validInvoice()
require.NoError(t, inv.Calculate())
inv.Supplier.Ext = nil
err := rules.Validate(inv)
assert.ErrorContains(t, err, "municipality extension is required")

inv.Supplier.Ext = tax.Extensions{
"br-ibge-municipality": "3304557",
}
err = rules.Validate(inv)
assert.NoError(t, err)
})

t.Run("validates supplier address fields", func(t *testing.T) {
inv := validInvoice()
require.NoError(t, inv.Calculate())
inv.Supplier.Addresses = []*org.Address{
{
Street:   "",
Number:   "100",
Locality: "São Paulo",
State:    "SP",
Code:     "01310100",
},
}
err := rules.Validate(inv)
assert.ErrorContains(t, err, "street is required")

inv.Supplier.Addresses[0].Street = "Rua Test"
inv.Supplier.Addresses[0].Number = ""
err = rules.Validate(inv)
assert.ErrorContains(t, err, "number is required")

inv.Supplier.Addresses[0].Number = "100"
inv.Supplier.Addresses[0].Locality = ""
err = rules.Validate(inv)
assert.ErrorContains(t, err, "locality is required")

inv.Supplier.Addresses[0].Locality = "São Paulo"
inv.Supplier.Addresses[0].State = ""
err = rules.Validate(inv)
assert.ErrorContains(t, err, "state is required")

inv.Supplier.Addresses[0].State = "SP"
inv.Supplier.Addresses[0].Code = ""
err = rules.Validate(inv)
assert.ErrorContains(t, err, "address code is required")
})
}

func TestCustomerValidation(t *testing.T) {
t.Run("validates customer required for NFe", func(t *testing.T) {
inv := validInvoice()
inv.Tax.Ext[nfe.ExtKeyModel] = nfe.ModelNFe
require.NoError(t, inv.Calculate())
inv.Customer = nil
err := rules.Validate(inv)
assert.ErrorContains(t, err, "customer is required for NFe")

inv.Customer = &org.Party{
Name: "Test Customer",
TaxID: &tax.Identity{
Country: "BR",
Code:    "05700736000196",
},
Addresses: []*org.Address{
{
Street:   "Rua das Flores",
Number:   "123",
Locality: "São Paulo",
State:    "SP",
Code:     "01310000",
},
},
Ext: tax.Extensions{
"br-ibge-municipality": "3550308",
},
}
err = rules.Validate(inv)
assert.NoError(t, err)
})

t.Run("validates customer addresses required for NFe", func(t *testing.T) {
inv := validInvoice()
inv.Tax.Ext[nfe.ExtKeyModel] = nfe.ModelNFe
require.NoError(t, inv.Calculate())
inv.Customer.Addresses = nil
err := rules.Validate(inv)
assert.ErrorContains(t, err, "customer addresses are required for NFe")

inv.Customer.Addresses = []*org.Address{}
err = rules.Validate(inv)
assert.ErrorContains(t, err, "customer addresses are required for NFe")

inv.Customer.Addresses = []*org.Address{nil}
err = rules.Validate(inv)
assert.ErrorContains(t, err, "address must not be empty")

inv.Customer.Addresses = []*org.Address{
{
Street:   "Rua das Flores",
Number:   "123",
Locality: "São Paulo",
State:    "SP",
Code:     "01310000",
},
}
err = rules.Validate(inv)
assert.NoError(t, err)
})

t.Run("customer not required for NFCe", func(t *testing.T) {
inv := validInvoice()
inv.Customer = nil
require.NoError(t, inv.Calculate())
// Override model to NFCe after Calculate to bypass scenario normalization
inv.Tax.Ext[nfe.ExtKeyModel] = nfe.ModelNFCe
inv.Tax.Ext[nfe.ExtKeyPresence] = nfe.PresenceInPerson
err := rules.Validate(inv)
assert.NoError(t, err)
})

t.Run("validates customer tax ID required", func(t *testing.T) {
inv := validInvoice()
require.NoError(t, inv.Calculate())
inv.Customer.TaxID = nil
err := rules.Validate(inv)
assert.ErrorContains(t, err, "customer tax ID is required")

inv.Customer.TaxID = &tax.Identity{Country: "BR"}
err = rules.Validate(inv)
assert.ErrorContains(t, err, "customer tax ID code is required")

inv.Customer.TaxID.Code = "05700736000196"
err = rules.Validate(inv)
assert.NoError(t, err)
})

t.Run("validates customer municipality when addresses exist", func(t *testing.T) {
inv := validInvoice()
require.NoError(t, inv.Calculate())
inv.Customer.Ext = nil
err := rules.Validate(inv)
assert.ErrorContains(t, err, "municipality extension is required")

inv.Customer.Ext = tax.Extensions{
"br-ibge-municipality": "3550308",
}
err = rules.Validate(inv)
assert.NoError(t, err)
})

t.Run("validates customer address fields", func(t *testing.T) {
inv := validInvoice()
require.NoError(t, inv.Calculate())
inv.Customer.Addresses = []*org.Address{
{
Street:   "",
Number:   "123",
Locality: "São Paulo",
State:    "SP",
Code:     "01310000",
},
}
err := rules.Validate(inv)
assert.ErrorContains(t, err, "street is required")

inv.Customer.Addresses[0].Street = "Rua das Flores"
inv.Customer.Addresses[0].Number = ""
err = rules.Validate(inv)
assert.ErrorContains(t, err, "number is required")

inv.Customer.Addresses[0].Number = "123"
inv.Customer.Addresses[0].Locality = ""
err = rules.Validate(inv)
assert.ErrorContains(t, err, "locality is required")

inv.Customer.Addresses[0].Locality = "São Paulo"
inv.Customer.Addresses[0].State = ""
err = rules.Validate(inv)
assert.ErrorContains(t, err, "state is required")

inv.Customer.Addresses[0].State = "SP"
inv.Customer.Addresses[0].Code = ""
err = rules.Validate(inv)
assert.ErrorContains(t, err, "address code is required")
})
}

func validInvoice() *bill.Invoice {
return &bill.Invoice{
Regime:   tax.WithRegime("BR"),
Addons:   tax.WithAddons(nfe.V4),
Currency: "BRL",
Type:     bill.InvoiceTypeStandard,
Series:   cbc.Code("123"),
Supplier: &org.Party{
Name: "Test Supplier LTDA",
TaxID: &tax.Identity{
Country: "BR",
Code:    "55263640000186",
},
Identities: []*org.Identity{
{
Key:  nfe.IdentityKeyStateReg,
Code: "35503304557308",
},
},
Addresses: []*org.Address{
{
Street:   "Av Paulista",
Number:   "1578",
Locality: "São Paulo",
State:    "SP",
Code:     "01310100",
},
},
Ext: tax.Extensions{
"br-ibge-municipality": "3304557",
},
},
Tax: &bill.Tax{
Ext: tax.Extensions{
nfe.ExtKeyModel:    nfe.ModelNFe,
nfe.ExtKeyPresence: nfe.PresenceInPerson,
},
},
Customer: &org.Party{
Name: "Test Customer LTDA",
TaxID: &tax.Identity{
Country: "BR",
Code:    "05700736000196",
},
Addresses: []*org.Address{
{
Street:   "Rua das Flores",
Number:   "123",
Locality: "São Paulo",
State:    "SP",
Code:     "01310000",
},
},
Ext: tax.Extensions{
"br-ibge-municipality": "3550308",
},
},
Notes: []*org.Note{
{
Key:  org.NoteKeyReason,
Text: "VENDA DE MERCADORIA",
},
},
Payment: &bill.PaymentDetails{
Instructions: &pay.Instructions{
Key: pay.MeansKeyCash,
},
},
Lines: []*bill.Line{
{
Quantity: num.MakeAmount(1, 0),
Item: &org.Item{
Name:  "Test Item",
Price: num.NewAmount(100, 2),
},
Taxes: tax.Set{
{Category: br.TaxCategoryICMS, Percent: num.NewPercentage(12, 2)},
{Category: br.TaxCategoryPIS, Percent: num.NewPercentage(165, 5)},
{Category: br.TaxCategoryCOFINS, Percent: num.NewPercentage(76, 4)},
},
},
},
}
}

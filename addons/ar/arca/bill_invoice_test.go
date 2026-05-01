package arca_test

import (
	"testing"

	"github.com/invopop/gobl/addons/ar/arca"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvoiceCustomerVATStatusNormalization(t *testing.T) {
	ad := tax.AddonForKey(arca.V4)

	t.Run("customer without tax ID sets final consumer", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Customer.TaxID = nil
		inv.Customer.Ext = tax.Extensions{}
		inv.Customer.Identities = []*org.Identity{
			{
				Code: "12345678",
				Ext: tax.ExtensionsOf(tax.ExtMap{
					arca.ExtKeyIdentityType: "96", // DNI
				}),
			},
		}
		ad.Normalizer(inv)
		assert.Equal(t, "5", inv.Customer.Ext.Get(arca.ExtKeyVATStatus).String())
	})

	t.Run("customer with AR tax ID sets registered VAT responsible", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Customer.Ext = tax.Extensions{}
		ad.Normalizer(inv)
		assert.Equal(t, "1", inv.Customer.Ext.Get(arca.ExtKeyVATStatus).String())
	})

	t.Run("customer with foreign tax ID sets foreign customer", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Customer.TaxID = &tax.Identity{
			Country: "US",
			Code:    "123456789",
		}
		inv.Customer.Ext = tax.Extensions{}
		ad.Normalizer(inv)
		assert.Equal(t, "9", inv.Customer.Ext.Get(arca.ExtKeyVATStatus).String())
	})

	t.Run("customer with existing valid VAT status is preserved", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Customer.Ext = tax.ExtensionsOf(tax.ExtMap{
			arca.ExtKeyVATStatus: "6", // Monotributo Responsible - valid for AR tax ID
		})
		ad.Normalizer(inv)
		assert.Equal(t, "6", inv.Customer.Ext.Get(arca.ExtKeyVATStatus).String())
	})

	t.Run("customer with invalid VAT status for AR tax ID is corrected", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		// "9" (Foreign Customer) is not valid for AR tax ID
		inv.Customer.Ext = tax.ExtensionsOf(tax.ExtMap{
			arca.ExtKeyVATStatus: "9",
		})
		ad.Normalizer(inv)
		// Should be corrected to default "1" (Responsable Inscripto)
		assert.Equal(t, "1", inv.Customer.Ext.Get(arca.ExtKeyVATStatus).String())
	})

	t.Run("AR customer with final consumer status is corrected", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		// "5" (Consumidor Final) is not valid for AR tax ID - final consumers don't have CUIT
		inv.Customer.Ext = tax.ExtensionsOf(tax.ExtMap{
			arca.ExtKeyVATStatus: "5",
		})
		ad.Normalizer(inv)
		// Should be corrected to default "1" (Responsable Inscripto)
		assert.Equal(t, "1", inv.Customer.Ext.Get(arca.ExtKeyVATStatus).String())
	})

	t.Run("AR customer with VAT exempt status is preserved", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		// "4" (IVA Sujeto Exento) is valid for AR tax ID - exempt organizations have CUIT
		inv.Customer.Ext = tax.ExtensionsOf(tax.ExtMap{
			arca.ExtKeyVATStatus: "4",
		})
		ad.Normalizer(inv)
		assert.Equal(t, "4", inv.Customer.Ext.Get(arca.ExtKeyVATStatus).String())
	})

	t.Run("customer with invalid VAT status for foreign tax ID is corrected", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Customer.TaxID = &tax.Identity{
			Country: "US",
			Code:    "123456789",
		}
		// "1" (Responsable Inscripto) is not valid for foreign tax ID
		inv.Customer.Ext = tax.ExtensionsOf(tax.ExtMap{
			arca.ExtKeyVATStatus: "1",
		})
		ad.Normalizer(inv)
		// Should be corrected to default "9" (Foreign Customer)
		assert.Equal(t, "9", inv.Customer.Ext.Get(arca.ExtKeyVATStatus).String())
	})

	t.Run("customer without tax ID preserves valid uncategorized status", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Customer.TaxID = nil
		inv.Customer.Identities = []*org.Identity{
			{
				Code: "12345678",
				Ext: tax.ExtensionsOf(tax.ExtMap{
					arca.ExtKeyIdentityType: "96", // DNI
				}),
			},
		}
		// "7" (Sujeto No Categorizado) is valid for no tax ID
		inv.Customer.Ext = tax.ExtensionsOf(tax.ExtMap{
			arca.ExtKeyVATStatus: "7",
		})
		ad.Normalizer(inv)
		assert.Equal(t, "7", inv.Customer.Ext.Get(arca.ExtKeyVATStatus).String())
	})

	t.Run("foreign customer can be set as foreign supplier", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Customer.TaxID = &tax.Identity{
			Country: "US",
			Code:    "123456789",
		}
		// "8" (Proveedor del Exterior) is valid for foreign tax ID
		inv.Customer.Ext = tax.ExtensionsOf(tax.ExtMap{
			arca.ExtKeyVATStatus: "8",
		})
		ad.Normalizer(inv)
		assert.Equal(t, "8", inv.Customer.Ext.Get(arca.ExtKeyVATStatus).String())
	})

	t.Run("nil customer does not panic", func(t *testing.T) {
		inv := testInvoiceSimplified(t)
		inv.Customer = nil
		assert.NotPanics(t, func() {
			ad.Normalizer(inv)
		})
	})
}

func TestInvoiceConceptNormalization(t *testing.T) {
	t.Run("only goods sets transaction type to products", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "1", inv.Tax.Ext.Get(arca.ExtKeyConcept).String())
	})

	t.Run("only services sets transaction type to services", func(t *testing.T) {
		inv := testInvoiceWithServices(t)
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "2", inv.Tax.Ext.Get(arca.ExtKeyConcept).String())
	})

	t.Run("default item key (empty) treated as services", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		// Item.Key is empty by default, treated as services
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "2", inv.Tax.Ext.Get(arca.ExtKeyConcept).String())
	})

	t.Run("mixed goods and services sets transaction type to products and services", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Lines = append(inv.Lines, &bill.Line{
			Quantity: num.MakeAmount(1, 0),
			Item: &org.Item{
				Name:  "Service Item",
				Price: num.NewAmount(5000, 2),
				Key:   org.ItemKeyServices,
			},
			Taxes: tax.Set{
				{
					Category: "VAT",
					Rate:     "standard",
				},
			},
		})
		// For mixed goods/services, we need ordering and payment
		inv.Ordering = testOrdering()
		inv.Payment = testPayment()
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "3", inv.Tax.Ext.Get(arca.ExtKeyConcept).String())
	})

	t.Run("nil item defaults to services", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Lines = append(inv.Lines, &bill.Line{
			Quantity: num.MakeAmount(1, 0),
			Item:     nil, // nil item defaults to services
		})
		// With goods + nil (service), we need ordering and payment for mixed concept
		inv.Ordering = testOrdering()
		inv.Payment = testPayment()
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "3", inv.Tax.Ext.Get(arca.ExtKeyConcept).String()) // mixed goods and services
	})

	t.Run("only nil items sets transaction type to services", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Lines = []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item:     nil, // nil item defaults to services
			},
		}
		inv.Ordering = testOrdering()
		inv.Payment = testPayment()
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "2", inv.Tax.Ext.Get(arca.ExtKeyConcept).String()) // services
	})

	t.Run("existing tax extensions are merged", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Tax = &bill.Tax{
			Ext: tax.ExtensionsOf(tax.ExtMap{
				arca.ExtKeyDocType: "1",
			}),
		}
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "1", inv.Tax.Ext.Get(arca.ExtKeyConcept).String())
		assert.Equal(t, "1", inv.Tax.Ext.Get(arca.ExtKeyDocType).String())
	})

	t.Run("empty lines does not set concept", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Lines = nil
		require.NoError(t, inv.Calculate())
		assert.Empty(t, inv.Tax.Ext.Get(arca.ExtKeyConcept))
	})

	t.Run("nil tax is initialized", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Tax = nil
		require.NoError(t, inv.Calculate())
		require.NotNil(t, inv.Tax)
		assert.Equal(t, "1", inv.Tax.Ext.Get(arca.ExtKeyConcept).String())
	})
}

func TestNormalizeBillInvoiceTaxDocType(t *testing.T) {
	ad := tax.AddonForKey(arca.V4)

	t.Run("nil tax is initialized", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Type = bill.InvoiceTypeStandard
		inv.Tax = nil
		ad.Normalizer(inv)
		require.NotNil(t, inv.Tax)
		// With AR customer (VAT status 1 after normalization), should get type A doc type
		assert.Equal(t, "1", inv.Tax.Ext.Get(arca.ExtKeyDocType).String())
	})

	t.Run("predefined doc type is not modified", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Type = bill.InvoiceTypeStandard
		inv.Tax.Ext = tax.ExtensionsOf(tax.ExtMap{
			arca.ExtKeyDocType: "51", // Invoice A with withholding legend
		})
		inv.Customer.Ext = tax.ExtensionsOf(tax.ExtMap{
			arca.ExtKeyVATStatus: "10", // VAT Exempt (would normally trigger type B)
		})
		ad.Normalizer(inv)
		// Should keep the existing doc type and not change it
		assert.Equal(t, "51", inv.Tax.Ext.Get(arca.ExtKeyDocType).String())
	})

	t.Run("proforma invoice type does not set doc type", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Type = bill.InvoiceTypeProforma
		inv.Tax.Ext = tax.Extensions{}
		inv.Customer.Ext = tax.ExtensionsOf(tax.ExtMap{
			arca.ExtKeyVATStatus: "1", // Would normally trigger type A
		})
		ad.Normalizer(inv)
		// Proforma is not supported, so no doc type should be set
		assert.Empty(t, inv.Tax.GetExt(arca.ExtKeyDocType))
	})

	t.Run("corrective invoice type does not set doc type", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Type = bill.InvoiceTypeCorrective
		inv.Tax.Ext = tax.Extensions{}
		inv.Customer.Ext = tax.ExtensionsOf(tax.ExtMap{
			arca.ExtKeyVATStatus: "1", // Would normally trigger type A
		})
		ad.Normalizer(inv)
		// Corrective is not supported, so no doc type should be set
		assert.Empty(t, inv.Tax.GetExt(arca.ExtKeyDocType))
	})

	t.Run("other invoice type does not set doc type", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Type = bill.InvoiceTypeOther
		inv.Tax.Ext = tax.Extensions{}
		inv.Customer.Ext = tax.ExtensionsOf(tax.ExtMap{
			arca.ExtKeyVATStatus: "1", // Would normally trigger type A
		})
		ad.Normalizer(inv)
		// Other is not supported, so no doc type should be set
		assert.Empty(t, inv.Tax.GetExt(arca.ExtKeyDocType))
	})

	t.Run("type B proforma does not set doc type", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Type = bill.InvoiceTypeProforma
		inv.Tax.Ext = tax.Extensions{}
		inv.Customer.Ext = tax.ExtensionsOf(tax.ExtMap{
			arca.ExtKeyVATStatus: "5", // Final Consumer - would normally trigger type B
		})
		ad.Normalizer(inv)
		// Proforma is not supported, so no doc type should be set
		assert.Empty(t, inv.Tax.GetExt(arca.ExtKeyDocType))
	})

	t.Run("type C proforma does not set doc type", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Type = bill.InvoiceTypeProforma
		inv.Tax.Ext = tax.Extensions{}
		inv.Tags = tax.Tags{List: []cbc.Key{arca.TagMonotax}}
		ad.Normalizer(inv)
		// Proforma is not supported, so no doc type should be set
		assert.Empty(t, inv.Tax.GetExt(arca.ExtKeyDocType))
	})

	t.Run("empty invoice type does not set doc type", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		// Empty type is not "standard"
		inv.Type = ""
		inv.Tax.Ext = tax.Extensions{}
		inv.Customer.Ext = tax.ExtensionsOf(tax.ExtMap{
			arca.ExtKeyVATStatus: "1",
		})
		ad.Normalizer(inv)
		// Empty type is not supported, so no doc type should be set
		assert.Empty(t, inv.Tax.GetExt(arca.ExtKeyDocType))
	})

	t.Run("customer with nil ext gets normalized then doc type set", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Type = bill.InvoiceTypeStandard
		inv.Tax.Ext = tax.Extensions{}
		inv.Customer.Ext = tax.Extensions{} // nil ext, customer exists with AR tax ID
		ad.Normalizer(inv)
		// Customer normalization runs first and sets VAT status to "1" (AR tax ID)
		// Then doc type normalization sets type A
		assert.Equal(t, "1", inv.Tax.Ext.Get(arca.ExtKeyDocType).String())
	})

	t.Run("empty string doc type is treated as not set", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Type = bill.InvoiceTypeStandard
		inv.Tax.Ext = tax.ExtensionsOf(tax.ExtMap{
			arca.ExtKeyDocType: "", // Empty string
		})
		inv.Customer.Ext = tax.ExtensionsOf(tax.ExtMap{
			arca.ExtKeyVATStatus: "1",
		})
		ad.Normalizer(inv)
		// Empty string is treated as not set, should get doc type
		assert.Equal(t, "1", inv.Tax.Ext.Get(arca.ExtKeyDocType).String())
	})

	t.Run("standard type A invoice sets doc type 1", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Type = bill.InvoiceTypeStandard
		inv.Tax.Ext = tax.Extensions{}
		inv.Customer.Ext = tax.ExtensionsOf(tax.ExtMap{
			arca.ExtKeyVATStatus: "1", // Registered Company - type A
		})
		ad.Normalizer(inv)
		assert.Equal(t, "1", inv.Tax.Ext.Get(arca.ExtKeyDocType).String())
	})

	t.Run("credit note type A sets doc type 3", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Type = bill.InvoiceTypeCreditNote
		inv.Tax.Ext = tax.Extensions{}
		inv.Customer.Ext = tax.ExtensionsOf(tax.ExtMap{
			arca.ExtKeyVATStatus: "1",
		})
		ad.Normalizer(inv)
		assert.Equal(t, "3", inv.Tax.Ext.Get(arca.ExtKeyDocType).String())
	})

	t.Run("debit note type A sets doc type 2", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Type = bill.InvoiceTypeDebitNote
		inv.Tax.Ext = tax.Extensions{}
		inv.Customer.Ext = tax.ExtensionsOf(tax.ExtMap{
			arca.ExtKeyVATStatus: "1",
		})
		ad.Normalizer(inv)
		assert.Equal(t, "2", inv.Tax.Ext.Get(arca.ExtKeyDocType).String())
	})

	t.Run("standard type B invoice sets doc type 6", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Type = bill.InvoiceTypeStandard
		inv.Tax.Ext = tax.Extensions{}
		inv.Customer.Ext = tax.ExtensionsOf(tax.ExtMap{
			arca.ExtKeyVATStatus: "10", // VAT Exempt Law 19640 - valid for AR, triggers type B
		})
		ad.Normalizer(inv)
		assert.Equal(t, "6", inv.Tax.Ext.Get(arca.ExtKeyDocType).String())
	})

	t.Run("credit note type B sets doc type 8", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Type = bill.InvoiceTypeCreditNote
		inv.Tax.Ext = tax.Extensions{}
		inv.Customer.Ext = tax.ExtensionsOf(tax.ExtMap{
			arca.ExtKeyVATStatus: "10", // VAT Exempt Law 19640 - valid for AR, triggers type B
		})
		ad.Normalizer(inv)
		assert.Equal(t, "8", inv.Tax.Ext.Get(arca.ExtKeyDocType).String())
	})

	t.Run("debit note type B sets doc type 7", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Type = bill.InvoiceTypeDebitNote
		inv.Tax.Ext = tax.Extensions{}
		inv.Customer.Ext = tax.ExtensionsOf(tax.ExtMap{
			arca.ExtKeyVATStatus: "10", // VAT Exempt Law 19640 - valid for AR, triggers type B
		})
		ad.Normalizer(inv)
		assert.Equal(t, "7", inv.Tax.Ext.Get(arca.ExtKeyDocType).String())
	})

	t.Run("standard type C invoice sets doc type 11", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Type = bill.InvoiceTypeStandard
		inv.Tax.Ext = tax.Extensions{}
		inv.Tags = tax.Tags{List: []cbc.Key{arca.TagMonotax}}
		ad.Normalizer(inv)
		assert.Equal(t, "11", inv.Tax.Ext.Get(arca.ExtKeyDocType).String())
	})

	t.Run("credit note type C sets doc type 13", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Type = bill.InvoiceTypeCreditNote
		inv.Tax.Ext = tax.Extensions{}
		inv.Tags = tax.Tags{List: []cbc.Key{arca.TagMonotax}}
		ad.Normalizer(inv)
		assert.Equal(t, "13", inv.Tax.Ext.Get(arca.ExtKeyDocType).String())
	})

	t.Run("debit note type C sets doc type 12", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Type = bill.InvoiceTypeDebitNote
		inv.Tax.Ext = tax.Extensions{}
		inv.Tags = tax.Tags{List: []cbc.Key{arca.TagMonotax}}
		ad.Normalizer(inv)
		assert.Equal(t, "12", inv.Tax.Ext.Get(arca.ExtKeyDocType).String())
	})

	t.Run("no customer defaults to type B", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Type = bill.InvoiceTypeStandard
		inv.Tax.Ext = tax.Extensions{}
		inv.Customer = nil
		ad.Normalizer(inv)
		assert.Equal(t, "6", inv.Tax.Ext.Get(arca.ExtKeyDocType).String())
	})

	t.Run("monotax tag takes precedence over VAT status", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Type = bill.InvoiceTypeStandard
		inv.Tax.Ext = tax.Extensions{}
		inv.Tags = tax.Tags{List: []cbc.Key{arca.TagMonotax}}
		inv.Customer.Ext = tax.ExtensionsOf(tax.ExtMap{
			arca.ExtKeyVATStatus: "1", // Would normally trigger type A
		})
		ad.Normalizer(inv)
		// Tag takes precedence, should get type C
		assert.Equal(t, "11", inv.Tax.Ext.Get(arca.ExtKeyDocType).String())
	})
}

func TestInvoiceDocTypeNormalization(t *testing.T) {
	ad := tax.AddonForKey(arca.V4)

	t.Run("monotax tag sets type C invoice", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Tax.Ext = tax.Extensions{} // Clear doc type
		inv.Tags = tax.Tags{List: []cbc.Key{arca.TagMonotax}}
		inv.Customer.Ext = tax.ExtensionsOf(tax.ExtMap{
			arca.ExtKeyVATStatus: "1", // Even with type A VAT status, tag takes precedence
		})
		ad.Normalizer(inv)
		assert.Equal(t, "11", inv.Tax.Ext.Get(arca.ExtKeyDocType).String()) // Invoice C
	})

	t.Run("monotax tag sets type C credit note", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Type = bill.InvoiceTypeCreditNote
		inv.Tax.Ext = tax.Extensions{}
		inv.Tags = tax.Tags{List: []cbc.Key{arca.TagMonotax}}
		ad.Normalizer(inv)
		assert.Equal(t, "13", inv.Tax.Ext.Get(arca.ExtKeyDocType).String()) // Credit Note C
	})

	t.Run("monotax tag sets type C debit note", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Type = bill.InvoiceTypeDebitNote
		inv.Tax.Ext = tax.Extensions{}
		inv.Tags = tax.Tags{List: []cbc.Key{arca.TagMonotax}}
		ad.Normalizer(inv)
		assert.Equal(t, "12", inv.Tax.Ext.Get(arca.ExtKeyDocType).String()) // Debit Note C
	})

	t.Run("VAT status 1 sets type A invoice", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Tax.Ext = tax.Extensions{}
		inv.Customer.Ext = tax.ExtensionsOf(tax.ExtMap{
			arca.ExtKeyVATStatus: "1", // Registered Company
		})
		ad.Normalizer(inv)
		assert.Equal(t, "1", inv.Tax.Ext.Get(arca.ExtKeyDocType).String()) // Invoice A
	})

	t.Run("VAT status 6 sets type A invoice", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Tax.Ext = tax.Extensions{}
		inv.Customer.Ext = tax.ExtensionsOf(tax.ExtMap{
			arca.ExtKeyVATStatus: "6", // Monotributo Responsible
		})
		ad.Normalizer(inv)
		assert.Equal(t, "1", inv.Tax.Ext.Get(arca.ExtKeyDocType).String()) // Invoice A
	})

	t.Run("VAT status 13 sets type A invoice", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Tax.Ext = tax.Extensions{}
		inv.Customer.Ext = tax.ExtensionsOf(tax.ExtMap{
			arca.ExtKeyVATStatus: "13", // Social Monotributista
		})
		ad.Normalizer(inv)
		assert.Equal(t, "1", inv.Tax.Ext.Get(arca.ExtKeyDocType).String()) // Invoice A
	})

	t.Run("VAT status 16 sets type A invoice", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Tax.Ext = tax.Extensions{}
		inv.Customer.Ext = tax.ExtensionsOf(tax.ExtMap{
			arca.ExtKeyVATStatus: "16", // Promoted Independent Worker Monotributista
		})
		ad.Normalizer(inv)
		assert.Equal(t, "1", inv.Tax.Ext.Get(arca.ExtKeyDocType).String()) // Invoice A
	})

	t.Run("VAT status 10 sets type B invoice", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Tax.Ext = tax.Extensions{}
		inv.Customer.Ext = tax.ExtensionsOf(tax.ExtMap{
			arca.ExtKeyVATStatus: "10", // VAT Exempt Law 19640
		})
		ad.Normalizer(inv)
		assert.Equal(t, "6", inv.Tax.Ext.Get(arca.ExtKeyDocType).String()) // Invoice B
	})

	t.Run("VAT status 15 sets type B invoice", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Tax.Ext = tax.Extensions{}
		inv.Customer.Ext = tax.ExtensionsOf(tax.ExtMap{
			arca.ExtKeyVATStatus: "15", // VAT Not Applicable
		})
		ad.Normalizer(inv)
		assert.Equal(t, "6", inv.Tax.Ext.Get(arca.ExtKeyDocType).String()) // Invoice B
	})

	t.Run("VAT status 5 with no tax ID sets type B invoice", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Tax.Ext = tax.Extensions{}
		inv.Customer.TaxID = nil
		inv.Customer.Identities = []*org.Identity{
			{
				Code: "12345678",
				Ext: tax.ExtensionsOf(tax.ExtMap{
					arca.ExtKeyIdentityType: "96", // DNI
				}),
			},
		}
		inv.Customer.Ext = tax.ExtensionsOf(tax.ExtMap{
			arca.ExtKeyVATStatus: "5", // Final Consumer
		})
		ad.Normalizer(inv)
		assert.Equal(t, "6", inv.Tax.Ext.Get(arca.ExtKeyDocType).String()) // Invoice B
	})

	t.Run("VAT status 9 with foreign tax ID sets type B invoice", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Tax.Ext = tax.Extensions{}
		inv.Customer.TaxID = &tax.Identity{
			Country: "US",
			Code:    "123456789",
		}
		inv.Customer.Ext = tax.ExtensionsOf(tax.ExtMap{
			arca.ExtKeyVATStatus: "9", // Foreign Customer
		})
		ad.Normalizer(inv)
		assert.Equal(t, "6", inv.Tax.Ext.Get(arca.ExtKeyDocType).String()) // Invoice B
	})

	t.Run("type A credit note with VAT status 1", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Type = bill.InvoiceTypeCreditNote
		inv.Tax.Ext = tax.Extensions{}
		inv.Customer.Ext = tax.ExtensionsOf(tax.ExtMap{
			arca.ExtKeyVATStatus: "1",
		})
		ad.Normalizer(inv)
		assert.Equal(t, "3", inv.Tax.Ext.Get(arca.ExtKeyDocType).String()) // Credit Note A
	})

	t.Run("type A debit note with VAT status 1", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Type = bill.InvoiceTypeDebitNote
		inv.Tax.Ext = tax.Extensions{}
		inv.Customer.Ext = tax.ExtensionsOf(tax.ExtMap{
			arca.ExtKeyVATStatus: "1",
		})
		ad.Normalizer(inv)
		assert.Equal(t, "2", inv.Tax.Ext.Get(arca.ExtKeyDocType).String()) // Debit Note A
	})

	t.Run("type B credit note with VAT status 10", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Type = bill.InvoiceTypeCreditNote
		inv.Tax.Ext = tax.Extensions{}
		inv.Customer.Ext = tax.ExtensionsOf(tax.ExtMap{
			arca.ExtKeyVATStatus: "10", // VAT Exempt Law 19640
		})
		ad.Normalizer(inv)
		assert.Equal(t, "8", inv.Tax.Ext.Get(arca.ExtKeyDocType).String()) // Credit Note B
	})

	t.Run("type B debit note with VAT status 15", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Type = bill.InvoiceTypeDebitNote
		inv.Tax.Ext = tax.Extensions{}
		inv.Customer.Ext = tax.ExtensionsOf(tax.ExtMap{
			arca.ExtKeyVATStatus: "15", // VAT Not Applicable
		})
		ad.Normalizer(inv)
		assert.Equal(t, "7", inv.Tax.Ext.Get(arca.ExtKeyDocType).String()) // Debit Note B
	})

	t.Run("existing doc type is preserved", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Tax.Ext = tax.ExtensionsOf(tax.ExtMap{
			arca.ExtKeyDocType: "51", // Invoice A with withholding legend
		})
		inv.Customer.Ext = tax.ExtensionsOf(tax.ExtMap{
			arca.ExtKeyVATStatus: "10", // VAT Exempt (would normally trigger type B)
		})
		ad.Normalizer(inv)
		// Should keep the existing doc type and not change it
		assert.Equal(t, "51", inv.Tax.Ext.Get(arca.ExtKeyDocType).String())
	})

	t.Run("no customer defaults to type B", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Tax.Ext = tax.Extensions{}
		inv.Customer = nil
		ad.Normalizer(inv)
		assert.Equal(t, "6", inv.Tax.Ext.Get(arca.ExtKeyDocType).String()) // Invoice B
	})

	t.Run("customer without extensions gets normalized to type A", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Tax.Ext = tax.Extensions{}
		inv.Customer.Ext = tax.Extensions{}
		ad.Normalizer(inv)
		// Customer has AR tax ID, so it will be normalized to VAT status "1" (Registered Company)
		// which triggers type A doc type
		assert.Equal(t, "1", inv.Tax.Ext.Get(arca.ExtKeyDocType).String()) // Invoice A
	})
}

func TestInvoiceSeriesValidation(t *testing.T) {
	t.Run("missing series fails", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Series = ""
		assertValidationError(t, inv, "series is required")
	})

	t.Run("non-numeric series fails", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Series = "ABC"
		assertValidationError(t, inv, "series must be a valid number between 1 and 99998")
	})

	t.Run("series with leading zeros is valid", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Series = "00001"
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("series below range fails", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Series = "0"
		assertValidationError(t, inv, "series must be a valid number between 1 and 99998")
	})

	t.Run("series above range fails", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Series = "99999"
		assertValidationError(t, inv, "series must be a valid number between 1 and 99998")
	})

	t.Run("series at lower bound is valid", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Series = "1"
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("series at upper bound is valid", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Series = "99998"
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("series with very large number overflows", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Series = "9999999999999999999999999" // number too large for strconv
		assertValidationError(t, inv, "series must be a valid number between 1 and 99998")
	})
}

func TestInvoiceTaxValidation(t *testing.T) {
	t.Run("valid standard invoice", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("missing tax fails", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		require.NoError(t, inv.Calculate())
		inv.Tax = nil
		require.ErrorContains(t, rules.Validate(inv), "tax is required")
	})

	t.Run("missing doc type fails", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		require.NoError(t, inv.Calculate())
		inv.Tax.Ext = tax.Extensions{}
		require.ErrorContains(t, rules.Validate(inv), "tax requires 'ar-arca-doc-type' extension")
	})
}

func TestInvoiceCustomerValidation(t *testing.T) {
	t.Run("B2B invoice automatically gets type A with customer", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		// With AR customer, should automatically become type A
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "1", inv.Tax.Ext.Get(arca.ExtKeyDocType).String())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("B2C invoice automatically gets type B without customer", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		// Without customer, should automatically become type B
		inv.Customer = nil
		inv.Tax.Ext = tax.Extensions{} // Clear doc type so it can be auto-detected
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "6", inv.Tax.Ext.Get(arca.ExtKeyDocType).String())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("customer not required for invoice type B (006)", func(t *testing.T) {
		inv := testInvoiceSimplified(t)
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("customer not required for debit note type B (007)", func(t *testing.T) {
		inv := testInvoiceSimplified(t)
		inv.Type = bill.InvoiceTypeDebitNote
		inv.Tax.Ext = inv.Tax.Ext.Set(arca.ExtKeyDocType, "7")
		inv.Preceding = testPreceding()
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("customer not required for credit note type B (008)", func(t *testing.T) {
		inv := testInvoiceSimplified(t)
		inv.Type = bill.InvoiceTypeCreditNote
		inv.Tax.Ext = inv.Tax.Ext.Set(arca.ExtKeyDocType, "8")
		inv.Preceding = testPreceding()
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("customer with tax ID is valid", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("customer with identity and ext is valid for type B", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Tax.Ext = inv.Tax.Ext.Set(arca.ExtKeyDocType, "6") // Type B
		inv.Customer.TaxID = nil
		inv.Customer.Identities = []*org.Identity{
			{
				Code: "12345678",
				Ext: tax.ExtensionsOf(tax.ExtMap{
					arca.ExtKeyIdentityType: "96", // DNI
				}),
			},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("customer without tax ID or identity ext fails", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Customer.TaxID = nil
		inv.Customer.Identities = nil
		assertValidationError(t, inv, "must have a tax_id, or an identity with ext 'ar-arca-identity-type'")
	})

	t.Run("customer with identity but no ext fails", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Customer.TaxID = nil
		inv.Customer.Identities = []*org.Identity{
			{
				Code: "12345678", // No ext
			},
		}
		assertValidationError(t, inv, "must have a tax_id, or an identity with ext 'ar-arca-identity-type'")
	})

	t.Run("customer with tax ID but missing code passes addon validation", func(t *testing.T) {
		// The ARCA addon does not enforce tax ID code presence;
		// that is handled by core tax.Identity rules.
		inv := testInvoiceWithGoods(t)
		inv.Customer.TaxID.Code = ""
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("doc type 49 with final consumer VAT status passes", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Tax.Ext = inv.Tax.Ext.Set(arca.ExtKeyDocType, arca.TypeUsedGoodsPurchaseInvoice) // 49
		inv.Customer.TaxID = nil
		inv.Customer.Identities = []*org.Identity{
			{
				Code: "12345678",
				Ext: tax.ExtensionsOf(tax.ExtMap{
					arca.ExtKeyIdentityType: "96", // DNI
				}),
			},
		}
		inv.Customer.Ext = tax.ExtensionsOf(tax.ExtMap{
			arca.ExtKeyVATStatus: "5", // Final Consumer
		})
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("doc type 49 with registered company VAT status fails", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Tax.Ext = inv.Tax.Ext.Set(arca.ExtKeyDocType, arca.TypeUsedGoodsPurchaseInvoice) // 49
		inv.Customer.Ext = tax.ExtensionsOf(tax.ExtMap{
			arca.ExtKeyVATStatus: "1", // Registered VAT Company
		})
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "document type 49 (Used Goods Purchase Invoice) requires customer VAT status to be 5 (Final Consumer)")
	})

	t.Run("doc type 49 with monotributo VAT status fails", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Tax.Ext = inv.Tax.Ext.Set(arca.ExtKeyDocType, arca.TypeUsedGoodsPurchaseInvoice) // 49
		inv.Customer.Ext = tax.ExtensionsOf(tax.ExtMap{
			arca.ExtKeyVATStatus: "6", // Monotributo Responsible
		})
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "document type 49 (Used Goods Purchase Invoice) requires customer VAT status to be 5 (Final Consumer)")
	})

	t.Run("doc type 49 with foreign customer VAT status fails", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Tax.Ext = inv.Tax.Ext.Set(arca.ExtKeyDocType, arca.TypeUsedGoodsPurchaseInvoice) // 49
		inv.Customer.TaxID = &tax.Identity{
			Country: "US",
			Code:    "123456789",
		}
		inv.Customer.Ext = tax.ExtensionsOf(tax.ExtMap{
			arca.ExtKeyVATStatus: "9", // Foreign Customer
		})
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "document type 49 (Used Goods Purchase Invoice) requires customer VAT status to be 5 (Final Consumer)")
	})

	t.Run("doc type 49 without customer passes", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Tax.Ext = inv.Tax.Ext.Set(arca.ExtKeyDocType, arca.TypeUsedGoodsPurchaseInvoice) // 49
		inv.Customer = nil
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})
}

func TestInvoiceServiceRequirements(t *testing.T) {
	t.Run("services require ordering with period", func(t *testing.T) {
		inv := testInvoiceWithServices(t)
		inv.Ordering = nil
		assertValidationError(t, inv, "ordering is required for services")
	})

	t.Run("services require ordering period", func(t *testing.T) {
		inv := testInvoiceWithServices(t)
		inv.Ordering = &bill.Ordering{}
		assertValidationError(t, inv, "ordering period is required for services")
	})

	t.Run("services with valid ordering passes", func(t *testing.T) {
		inv := testInvoiceWithServices(t)
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("services require payment terms", func(t *testing.T) {
		inv := testInvoiceWithServices(t)
		inv.Payment = nil
		assertValidationError(t, inv, "payment is required for services")
	})

	t.Run("services require payment terms object", func(t *testing.T) {
		inv := testInvoiceWithServices(t)
		inv.Payment = &bill.PaymentDetails{}
		assertValidationError(t, inv, "payment terms are required for services")
	})

	t.Run("services require payment due dates", func(t *testing.T) {
		inv := testInvoiceWithServices(t)
		inv.Payment = &bill.PaymentDetails{
			Terms: &pay.Terms{},
		}
		assertValidationError(t, inv, "payment due dates are required for services")
	})

	t.Run("products do not require ordering or payment", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Ordering = nil
		inv.Payment = nil
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("products with payment due dates fails", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Payment = testPayment()
		assertValidationError(t, inv, "payment due dates must not be set for goods")
	})

	t.Run("products with payment but no due dates passes", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Payment = &bill.PaymentDetails{
			Terms: &pay.Terms{
				Notes: "Payment on delivery",
			},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("products with empty payment passes", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Payment = &bill.PaymentDetails{}
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("mixed goods and services requires ordering and payment", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Lines = append(inv.Lines, &bill.Line{
			Quantity: num.MakeAmount(1, 0),
			Item: &org.Item{
				Name:  "Service Item",
				Price: num.NewAmount(5000, 2),
				Key:   org.ItemKeyServices,
			},
			Taxes: tax.Set{
				{
					Category: "VAT",
					Rate:     "standard",
				},
			},
		})
		inv.Ordering = nil
		assertValidationError(t, inv, "ordering is required for services")
	})
}

func TestCreditNoteValidation(t *testing.T) {
	t.Run("valid credit note type A", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Type = bill.InvoiceTypeCreditNote
		inv.Tax.Ext = tax.Extensions{} // Clear doc type so it can be auto-detected
		inv.Preceding = testPreceding()
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
		assert.Equal(t, "3", inv.Tax.Ext.Get(arca.ExtKeyDocType).String())
	})

	t.Run("valid credit note type B", func(t *testing.T) {
		inv := testInvoiceSimplified(t)
		inv.Type = bill.InvoiceTypeCreditNote
		inv.Tax.Ext = tax.Extensions{} // Clear doc type so it can be auto-detected
		inv.Preceding = testPreceding()
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
		assert.Equal(t, "8", inv.Tax.Ext.Get(arca.ExtKeyDocType).String())
	})
}

func TestDebitNoteValidation(t *testing.T) {
	t.Run("valid debit note type A", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Type = bill.InvoiceTypeDebitNote
		inv.Tax.Ext = tax.Extensions{} // Clear doc type so it can be auto-detected
		inv.Preceding = testPreceding()
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
		assert.Equal(t, "2", inv.Tax.Ext.Get(arca.ExtKeyDocType).String())
	})

	t.Run("valid debit note type B", func(t *testing.T) {
		inv := testInvoiceSimplified(t)
		inv.Type = bill.InvoiceTypeDebitNote
		inv.Tax.Ext = tax.Extensions{} // Clear doc type so it can be auto-detected
		inv.Preceding = testPreceding()
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
		assert.Equal(t, "7", inv.Tax.Ext.Get(arca.ExtKeyDocType).String())
	})
}

func TestInvoiceTypeDocTypeValidation(t *testing.T) {
	t.Run("credit note with credit note doc type passes", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Type = bill.InvoiceTypeCreditNote
		inv.Tax.Ext = inv.Tax.Ext.Set(arca.ExtKeyDocType, "3") // Credit Note A
		inv.Preceding = testPreceding()
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("credit note with standard doc type fails", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Type = bill.InvoiceTypeCreditNote
		inv.Tax.Ext = inv.Tax.Ext.Set(arca.ExtKeyDocType, "1") // Standard Invoice A
		inv.Preceding = testPreceding()
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invoice type is credit-note but ar-arca-doc-type is not a credit note")
	})

	t.Run("debit note with debit note doc type passes", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Type = bill.InvoiceTypeDebitNote
		inv.Tax.Ext = inv.Tax.Ext.Set(arca.ExtKeyDocType, "2") // Debit Note A
		inv.Preceding = testPreceding()
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("debit note with standard doc type fails", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Type = bill.InvoiceTypeDebitNote
		inv.Tax.Ext = inv.Tax.Ext.Set(arca.ExtKeyDocType, "1") // Standard Invoice A
		inv.Preceding = testPreceding()
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invoice type is debit-note but ar-arca-doc-type is not a debit note")
	})

	t.Run("standard invoice with credit note doc type fails", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Type = bill.InvoiceTypeStandard
		inv.Tax.Ext = inv.Tax.Ext.Set(arca.ExtKeyDocType, "3") // Credit Note A
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "ar-arca-doc-type is a credit note but invoice type is not credit-note")
	})

	t.Run("standard invoice with debit note doc type fails", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Type = bill.InvoiceTypeStandard
		inv.Tax.Ext = inv.Tax.Ext.Set(arca.ExtKeyDocType, "2") // Debit Note A
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "doc type is a debit note but invoice type is not debit-note")
	})

	t.Run("credit note with FCE credit note doc type passes", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Type = bill.InvoiceTypeCreditNote
		inv.Tax.Ext = inv.Tax.Ext.Set(arca.ExtKeyDocType, "203") // MiPyMEs Electronic Credit Note (FCE A)
		inv.Preceding = testPreceding()
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("debit note with FCE debit note doc type passes", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Type = bill.InvoiceTypeDebitNote
		inv.Tax.Ext = inv.Tax.Ext.Set(arca.ExtKeyDocType, "202") // MiPyMEs Electronic Debit Note (FCE A)
		inv.Preceding = testPreceding()
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})
}

func TestInvoicePrecedingValidation(t *testing.T) {
	t.Run("credit note requires preceding", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Type = bill.InvoiceTypeCreditNote
		inv.Preceding = nil
		assertValidationError(t, inv, "preceding documents are required for credit/debit notes")
	})

	t.Run("debit note requires preceding", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Type = bill.InvoiceTypeDebitNote
		inv.Preceding = nil
		assertValidationError(t, inv, "preceding documents are required for credit/debit notes")
	})

	t.Run("standard invoice does not require preceding", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Type = bill.InvoiceTypeStandard
		inv.Preceding = nil
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("credit note with valid preceding passes", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Type = bill.InvoiceTypeCreditNote
		inv.Tax.Ext = tax.Extensions{} // Clear doc type so normalization sets it to credit note
		inv.Preceding = []*org.DocumentRef{
			{
				Series: "1",
				Code:   "100",
				Ext: tax.ExtensionsOf(tax.ExtMap{
					arca.ExtKeyDocType: "1",
				}),
			},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("debit note with valid preceding passes", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Type = bill.InvoiceTypeDebitNote
		inv.Tax.Ext = tax.Extensions{} // Clear doc type so normalization sets it to debit note
		inv.Preceding = []*org.DocumentRef{
			{
				Series: "1",
				Code:   "100",
				Ext: tax.ExtensionsOf(tax.ExtMap{
					arca.ExtKeyDocType: "1",
				}),
			},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("credit note with preceding missing doc type fails", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Type = bill.InvoiceTypeCreditNote
		inv.Tax.Ext = tax.Extensions{} // Clear doc type so normalization sets it to credit note
		inv.Preceding = []*org.DocumentRef{
			{
				Series: "1",
				Code:   "100",
			},
		}
		assertValidationError(t, inv, "preceding document requires 'ar-arca-doc-type' extension")
	})

	t.Run("debit note with preceding missing doc type fails", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Type = bill.InvoiceTypeDebitNote
		inv.Tax.Ext = tax.Extensions{} // Clear doc type so normalization sets it to debit note
		inv.Preceding = []*org.DocumentRef{
			{
				Series: "1",
				Code:   "100",
			},
		}
		assertValidationError(t, inv, "preceding document requires 'ar-arca-doc-type' extension")
	})

	t.Run("credit note with multiple preceding validates all", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Type = bill.InvoiceTypeCreditNote
		inv.Tax.Ext = tax.Extensions{} // Clear doc type so normalization sets it to credit note
		inv.Preceding = []*org.DocumentRef{
			{
				Series: "1",
				Code:   "100",
				Ext: tax.ExtensionsOf(tax.ExtMap{
					arca.ExtKeyDocType: "1",
				}),
			},
			{
				Series: "1",
				Code:   "101",
				// Missing doc type
			},
		}
		assertValidationError(t, inv, "preceding document requires 'ar-arca-doc-type' extension")
	})

	t.Run("standard invoice with preceding validates doc type", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Type = bill.InvoiceTypeStandard
		inv.Preceding = []*org.DocumentRef{
			{
				Series: "1",
				Code:   "100",
				// Missing doc type
			},
		}
		assertValidationError(t, inv, "preceding document requires 'ar-arca-doc-type' extension")
	})
}

func TestValidateFunctionsWithNilValues(t *testing.T) {
	t.Run("validate with nil tax", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Tax = nil
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "tax is required")
	})

	t.Run("validate ordering with nil for services", func(t *testing.T) {
		inv := testInvoiceWithServices(t)
		// Need to set concept to "2" (services) to trigger ordering validation
		inv.Tax.Ext = inv.Tax.Ext.Set(arca.ExtKeyConcept, "2")
		inv.Ordering = nil
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "ordering is required for services")
	})

	t.Run("validate payment with nil for services", func(t *testing.T) {
		inv := testInvoiceWithServices(t)
		// Need to set concept to "2" (services) to trigger payment validation
		inv.Tax.Ext = inv.Tax.Ext.Set(arca.ExtKeyConcept, "2")
		inv.Payment = nil
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "payment is required for services")
	})
}

func TestCorrectionDefinitions(t *testing.T) {
	t.Run("correction definitions exist for credit and debit notes", func(t *testing.T) {
		ad := tax.AddonForKey(arca.V4)
		require.NotNil(t, ad.Corrections)
		// Check that invoice correction definition exists
		def := ad.Corrections.Def(bill.ShortSchemaInvoice)
		require.NotNil(t, def)
		assert.True(t, def.HasType(bill.InvoiceTypeCreditNote))
		assert.True(t, def.HasType(bill.InvoiceTypeDebitNote))
		assert.True(t, def.HasExtension(arca.ExtKeyDocType))
	})
}

func TestTypeCInvoiceLineTaxesValidation(t *testing.T) {
	t.Run("type C invoice without taxes on lines passes", func(t *testing.T) {
		inv := testInvoiceTypeC(t)
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("type C invoice with taxes on lines fails", func(t *testing.T) {
		inv := testInvoiceTypeC(t)
		inv.Lines[0].Taxes = tax.Set{
			{
				Category: "VAT",
				Rate:     "standard",
			},
		}
		assertValidationError(t, inv, "type C invoices (simplified tax scheme) must not have taxes on lines")
	})

	t.Run("type C debit note without taxes on lines passes", func(t *testing.T) {
		inv := testInvoiceTypeC(t)
		inv.Type = bill.InvoiceTypeDebitNote
		inv.Tax.Ext = inv.Tax.Ext.Set(arca.ExtKeyDocType, "12") // Debit Note C
		inv.Preceding = testPreceding()
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("type C debit note with taxes on lines fails", func(t *testing.T) {
		inv := testInvoiceTypeC(t)
		inv.Type = bill.InvoiceTypeDebitNote
		inv.Tax.Ext = inv.Tax.Ext.Set(arca.ExtKeyDocType, "12") // Debit Note C
		inv.Preceding = testPreceding()
		inv.Lines[0].Taxes = tax.Set{
			{
				Category: "VAT",
				Rate:     "standard",
			},
		}
		assertValidationError(t, inv, "type C invoices (simplified tax scheme) must not have taxes on lines")
	})

	t.Run("type C credit note without taxes on lines passes", func(t *testing.T) {
		inv := testInvoiceTypeC(t)
		inv.Type = bill.InvoiceTypeCreditNote
		inv.Tax.Ext = inv.Tax.Ext.Set(arca.ExtKeyDocType, "13") // Credit Note C
		inv.Preceding = testPreceding()
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("type C credit note with taxes on lines fails", func(t *testing.T) {
		inv := testInvoiceTypeC(t)
		inv.Type = bill.InvoiceTypeCreditNote
		inv.Tax.Ext = inv.Tax.Ext.Set(arca.ExtKeyDocType, "13") // Credit Note C
		inv.Preceding = testPreceding()
		inv.Lines[0].Taxes = tax.Set{
			{
				Category: "VAT",
				Rate:     "standard",
			},
		}
		assertValidationError(t, inv, "type C invoices (simplified tax scheme) must not have taxes on lines")
	})

	t.Run("FCE type C invoice without taxes on lines passes", func(t *testing.T) {
		inv := testInvoiceTypeC(t)
		inv.Tax.Ext = inv.Tax.Ext.Set(arca.ExtKeyDocType, "211") // FCE Invoice C
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("FCE type C invoice with taxes on lines fails", func(t *testing.T) {
		inv := testInvoiceTypeC(t)
		inv.Tax.Ext = inv.Tax.Ext.Set(arca.ExtKeyDocType, "211") // FCE Invoice C
		inv.Lines[0].Taxes = tax.Set{
			{
				Category: "VAT",
				Rate:     "standard",
			},
		}
		assertValidationError(t, inv, "type C invoices (simplified tax scheme) must not have taxes on lines")
	})

	t.Run("type A invoice with taxes on lines passes", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		// Type A invoice (doc type "1") should allow taxes
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("type B invoice with taxes on lines passes", func(t *testing.T) {
		inv := testInvoiceSimplified(t)
		// Type B invoice (doc type "6") should allow taxes
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("type C invoice with multiple lines without taxes passes", func(t *testing.T) {
		inv := testInvoiceTypeC(t)
		inv.Lines = append(inv.Lines, &bill.Line{
			Quantity: num.MakeAmount(2, 0),
			Item: &org.Item{
				Name:  "Another Service",
				Price: num.NewAmount(5000, 2),
				Key:   org.ItemKeyServices,
			},
		})
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("type C invoice with taxes on second line fails", func(t *testing.T) {
		inv := testInvoiceTypeC(t)
		inv.Lines = append(inv.Lines, &bill.Line{
			Quantity: num.MakeAmount(2, 0),
			Item: &org.Item{
				Name:  "Another Service",
				Price: num.NewAmount(5000, 2),
				Key:   org.ItemKeyServices,
			},
			Taxes: tax.Set{
				{
					Category: "VAT",
					Rate:     "standard",
				},
			},
		})
		assertValidationError(t, inv, "type C invoices (simplified tax scheme) must not have taxes on lines")
	})
}

func TestInvoiceCurrencyValidation(t *testing.T) {
	t.Run("non-ARS currency without exchange rates", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Currency = "USD"
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "[GOBL-AR-ARCA-BILL-INVOICE-24] invoice must be in ARS or provide exchange rate for conversion")
	})

	t.Run("non-ARS currency with exchange rates", func(t *testing.T) {
		inv := testInvoiceWithGoods(t)
		inv.Currency = "USD"
		inv.ExchangeRates = []*currency.ExchangeRate{
			{
				From:   "USD",
				To:     "ARS",
				Amount: num.MakeAmount(1050, 0),
			},
		}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.NoError(t, err)
	})
}

func TestTourismInvoiceTypeT(t *testing.T) {
	t.Run("valid type T invoice (195) passes", func(t *testing.T) {
		inv := testInvoiceTourism(t)
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("type T invoice without tourism relation fails", func(t *testing.T) {
		inv := testInvoiceTourism(t)
		inv.Tax.Ext = inv.Tax.Ext.Delete(arca.ExtKeyTourismRelation)
		assertValidationError(t, inv, "tourism invoice requires 'ar-arca-tourism-relation' extension")
	})

	t.Run("type T invoice with services does not require ordering or payment", func(t *testing.T) {
		inv := testInvoiceTourism(t)
		inv.Lines[0].Item.Key = org.ItemKeyServices
		inv.Ordering = nil
		inv.Payment = nil
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("type T debit note (196) is recognized as debit note", func(t *testing.T) {
		inv := testInvoiceTourism(t)
		inv.Type = bill.InvoiceTypeDebitNote
		inv.Tax.Ext = inv.Tax.Ext.Set(arca.ExtKeyDocType, "196")
		inv.Preceding = testPreceding()
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("type T credit note (197) is recognized as credit note", func(t *testing.T) {
		inv := testInvoiceTourism(t)
		inv.Type = bill.InvoiceTypeCreditNote
		inv.Tax.Ext = inv.Tax.Ext.Set(arca.ExtKeyDocType, "197")
		inv.Preceding = testPreceding()
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("type T invoice with standard doc type fails type check", func(t *testing.T) {
		inv := testInvoiceTourism(t)
		inv.Type = bill.InvoiceTypeCreditNote
		inv.Tax.Ext = inv.Tax.Ext.Set(arca.ExtKeyDocType, "195") // 195 is an invoice, not a credit note
		inv.Preceding = testPreceding()
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invoice type is credit-note but ar-arca-doc-type is not a credit note")
	})

	t.Run("type T invoice (196) as standard fails type check", func(t *testing.T) {
		inv := testInvoiceTourism(t)
		inv.Type = bill.InvoiceTypeStandard
		inv.Tax.Ext = inv.Tax.Ext.Set(arca.ExtKeyDocType, "196") // 196 is debit note
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "doc type is a debit note but invoice type is not debit-note")
	})

	t.Run("type T invoice without tourism code on lines fails", func(t *testing.T) {
		inv := testInvoiceTourism(t)
		inv.Lines[0].Taxes[0].Ext = tax.Extensions{}
		assertValidationError(t, inv, "tourism invoice line requires 'ar-arca-tourism-code' extension")
	})

	t.Run("type T invoice without customer addresses fails", func(t *testing.T) {
		inv := testInvoiceTourism(t)
		inv.Customer.Addresses = nil
		assertValidationError(t, inv, "tourism invoice customer requires an address")
	})

	t.Run("type T invoice with non-21% VAT rate fails", func(t *testing.T) {
		inv := testInvoiceTourism(t)
		inv.Lines[0].Taxes[0].Rate = "reduced"
		assertValidationError(t, inv, "tourism invoice line VAT rate must be '5'")
	})

	t.Run("type T invoice with empty line taxes fails", func(t *testing.T) {
		inv := testInvoiceTourism(t)
		inv.Lines[0].Taxes = nil
		assertValidationError(t, inv, "tourism invoice line requires taxes")
	})
}

// Helper functions

func assertValidationError(t *testing.T, inv *bill.Invoice, expected string) {
	t.Helper()
	require.NoError(t, inv.Calculate())
	err := rules.Validate(inv)
	require.ErrorContains(t, err, expected)
}

func testInvoiceStandard(t *testing.T) *bill.Invoice {
	t.Helper()
	inv := &bill.Invoice{
		Addons: tax.WithAddons(arca.V4),
		Type:   bill.InvoiceTypeStandard,
		Series: "1",
		Code:   "123",
		Tax: &bill.Tax{
			Ext: tax.ExtensionsOf(tax.ExtMap{
				arca.ExtKeyDocType: "1",
			}),
		},
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "AR",
				Code:    "30500010912", // Valid company CUIT
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "AR",
				Code:    "20172543597", // Valid individual CUIL
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.NewAmount(10000, 2),
					Unit:  org.UnitPackage,
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     "standard",
					},
				},
			},
		},
	}
	return inv
}

func testInvoiceWithGoods(t *testing.T) *bill.Invoice {
	t.Helper()
	inv := testInvoiceStandard(t)
	inv.Lines[0].Item.Key = org.ItemKeyGoods
	return inv
}

func testInvoiceWithServices(t *testing.T) *bill.Invoice {
	t.Helper()
	inv := testInvoiceStandard(t)
	inv.Lines[0].Item.Key = org.ItemKeyServices
	inv.Ordering = testOrdering()
	inv.Payment = testPayment()
	return inv
}

func testInvoiceSimplified(t *testing.T) *bill.Invoice {
	t.Helper()
	inv := testInvoiceWithGoods(t)
	// B2C invoice: no customer (doc type will be set automatically by scenario)
	inv.Customer = nil
	inv.Tax.Ext = inv.Tax.Ext.Set(arca.ExtKeyDocType, "6")
	return inv
}

func testInvoiceTypeC(t *testing.T) *bill.Invoice {
	t.Helper()
	inv := &bill.Invoice{
		Addons: tax.WithAddons(arca.V4),
		Series: "1",
		Code:   "123",
		Tax: &bill.Tax{
			Ext: tax.ExtensionsOf(tax.ExtMap{
				arca.ExtKeyDocType: "11", // Invoice C (monotributo)
			}),
		},
		Supplier: &org.Party{
			Name: "Test Supplier Monotributo",
			TaxID: &tax.Identity{
				Country: "AR",
				Code:    "20172543597",
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "AR",
				Code:    "30500010912",
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Service Item",
					Price: num.NewAmount(10000, 2),
					Key:   org.ItemKeyServices,
				},
				// No taxes for type C invoices
			},
		},
		Ordering: testOrdering(),
		Payment:  testPayment(),
	}
	inv.SetTags(arca.TagMonotax) // Type C uses monotax tag
	return inv
}

func testOrdering() *bill.Ordering {
	return &bill.Ordering{
		Period: &cal.Period{
			Start: cal.MakeDate(2024, 1, 1),
			End:   cal.MakeDate(2024, 1, 31),
		},
	}
}

func testPayment() *bill.PaymentDetails {
	return &bill.PaymentDetails{
		Terms: &pay.Terms{
			DueDates: []*pay.DueDate{
				{
					Date:   cal.NewDate(2024, 2, 15),
					Amount: num.MakeAmount(10000, 2),
				},
			},
		},
	}
}

func testInvoiceTourism(t *testing.T) *bill.Invoice {
	t.Helper()
	return &bill.Invoice{
		Addons: tax.WithAddons(arca.V4),
		Type:   bill.InvoiceTypeStandard,
		Series: "1",
		Code:   "123",
		Tax: &bill.Tax{
			Ext: tax.ExtensionsOf(tax.ExtMap{
				arca.ExtKeyDocType:         "195",
				arca.ExtKeyTourismRelation: "1",
			}),
		},
		Supplier: &org.Party{
			Name: "Test Hotel",
			TaxID: &tax.Identity{
				Country: "AR",
				Code:    "30500010912",
			},
		},
		Customer: &org.Party{
			Name: "Foreign Tourist",
			TaxID: &tax.Identity{
				Country: "US",
				Code:    "123456789",
			},
			Addresses: []*org.Address{
				{
					Street:   "5th Avenue 100",
					Locality: "New York",
					Country:  "US",
				},
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Hotel Room",
					Price: num.NewAmount(10000, 2),
					Key:   org.ItemKeyGoods,
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     "standard",
						Ext: tax.ExtensionsOf(tax.ExtMap{
							arca.ExtKeyTourismCode: "1",
						}),
					},
				},
			},
		},
	}
}

func testPreceding() []*org.DocumentRef {
	return []*org.DocumentRef{
		{
			Series: "1",
			Code:   "100",
			Ext: tax.ExtensionsOf(tax.ExtMap{
				arca.ExtKeyDocType: "1",
			}),
		},
	}
}

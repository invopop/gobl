package favat_test

import (
	"testing"

	"github.com/invopop/gobl/addons/pl/favat"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func creditNote() *bill.Invoice {
	inv := &bill.Invoice{
		Regime:    tax.WithRegime("PL"),
		Addons:    tax.WithAddons(favat.V3),
		Currency:  currency.PLN,
		Code:      "TEST",
		Type:      bill.InvoiceTypeCreditNote,
		IssueDate: cal.MakeDate(2022, 12, 29),
		Preceding: []*org.DocumentRef{
			{
				Code:      "TEST",
				IssueDate: cal.NewDate(2022, 12, 27),
				Ext: tax.Extensions{
					favat.ExtKeyEffectiveDate: "1",
				},
			},
		},
		Supplier: &org.Party{
			Name: "Test Party",
			TaxID: &tax.Identity{
				Country: "PL",
				Code:    "1111111111",
			},
			Addresses: []*org.Address{
				{
					Street:   "ul. Testowa 1",
					Locality: "Warsaw",
					Country:  "PL",
				},
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "PL",
				Code:    "2222222222",
			},
			Addresses: []*org.Address{
				{
					Locality: "Foo",
				},
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 3),
				Item: &org.Item{
					Name:  "bogus",
					Price: num.NewAmount(1000, 3),
				},
			},
		},
	}
	return inv
}

func standardInvoice() *bill.Invoice {
	inv := &bill.Invoice{
		Regime:    tax.WithRegime("PL"),
		Addons:    tax.WithAddons(favat.V3),
		Currency:  currency.PLN,
		Code:      "STD",
		Type:      bill.InvoiceTypeStandard,
		IssueDate: cal.MakeDate(2022, 12, 29),
		Supplier: &org.Party{
			Name: "Test Party",
			TaxID: &tax.Identity{
				Country: "PL",
				Code:    "1111111111",
			},
			Addresses: []*org.Address{
				{
					Street:   "ul. Testowa 1",
					Locality: "Warsaw",
					Country:  "PL",
				},
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "PL",
				Code:    "2222222222",
			},
			Addresses: []*org.Address{
				{
					Locality: "Foo",
				},
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 3),
				Item: &org.Item{
					Name:  "standard-bogus",
					Price: num.NewAmount(1000, 3),
				},
			},
		},
	}
	return inv
}

func TestBasicCreditNoteValidation(t *testing.T) {
	inv := creditNote()
	inv.Preceding[0].Reason = "Correcting an error"
	err := inv.Calculate()
	require.NoError(t, err)
	err = inv.Validate()
	assert.NoError(t, err)
	assert.Equal(t, inv.Preceding[0].Ext[favat.ExtKeyEffectiveDate], cbc.Code("1"))

	inv.Preceding[0].Ext["foo"] = "bar"
	err = inv.Validate()
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "preceding: (0: (ext: (foo: undefined.).).)")
	}
}

func TestBasicStandardInvoiceValidation(t *testing.T) {
	inv := standardInvoice()
	err := inv.Calculate()
	require.NoError(t, err)
	err = inv.Validate()
	assert.NoError(t, err)
}

func TestExemptStandardInvoiceValidation(t *testing.T) {
	inv := standardInvoice()
	inv.Tax = &bill.Tax{
		Ext: tax.Extensions{
			favat.ExtKeyExemption: "A",
		},
	}
	inv.Notes = []*org.Note{
		{
			Key:  org.NoteKeyLegal,
			Code: "A",
			Src:  favat.ExtKeyExemption,
			Text: "Art. 25a ust. 1 pkt 9 ustawy o VAT",
		},
	}

	err := inv.Calculate()
	require.NoError(t, err)
	err = inv.Validate()
	assert.NoError(t, err)
}

func TestExemptStandardInvoiceValidationFailsWithoutNote(t *testing.T) {
	inv := standardInvoice()
	inv.Tax = &bill.Tax{
		// valid exemption set but no matching note
		Ext: tax.Extensions{
			favat.ExtKeyExemption: "A",
		},
	}

	err := inv.Calculate()
	require.NoError(t, err)

	err = inv.Validate()
	assert.Error(t, err)
}

func TestExemptStandardInvoiceValidationFailsWithTooManyNotes(t *testing.T) {
	inv := standardInvoice()
	inv.Tax = &bill.Tax{
		Ext: tax.Extensions{
			favat.ExtKeyExemption: "A",
		},
	}
	inv.Notes = []*org.Note{
		{
			Key:  org.NoteKeyLegal,
			Code: "A",
			Src:  favat.ExtKeyExemption,
			Text: "Art. 25a ust. 1 pkt 9 ustawy o VAT",
		},
		{
			Key:  org.NoteKeyLegal,
			Code: "A",
			Src:  favat.ExtKeyExemption,
			Text: "Another exemption note",
		},
	}

	err := inv.Calculate()
	require.NoError(t, err)

	err = inv.Validate()
	assert.ErrorContains(t, err, "too many exemption notes")
}

func TestSupplierValidation(t *testing.T) {
	t.Run("valid supplier", func(t *testing.T) {
		inv := standardInvoice()
		inv.Supplier.Name = "Test Supplier"
		inv.Supplier.Addresses = []*org.Address{
			{
				Street:  "ul. Testowa 1",
				Country: "PL",
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("missing supplier name", func(t *testing.T) {
		inv := standardInvoice()
		inv.Supplier.Name = ""
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "supplier: (name: cannot be blank.)")
	})

	t.Run("missing addresses", func(t *testing.T) {
		inv := standardInvoice()
		inv.Supplier.Addresses = nil
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "supplier: (addresses: cannot be blank.)")
	})

	t.Run("empty addresses array", func(t *testing.T) {
		inv := standardInvoice()
		inv.Supplier.Addresses = []*org.Address{}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "supplier: (addresses: cannot be blank.)")
	})

	t.Run("missing country in first address", func(t *testing.T) {
		inv := standardInvoice()
		inv.Supplier.Addresses = []*org.Address{
			{
				Street:  "ul. Testowa 1",
				Country: "",
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "supplier: (addresses: (country: cannot be blank.).).")
	})

	t.Run("missing street in first address", func(t *testing.T) {
		inv := standardInvoice()
		inv.Supplier.Addresses = []*org.Address{
			{
				Street:  "",
				Country: "PL",
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "supplier: (addresses: (street: cannot be blank.).).")
	})

	t.Run("missing both country and street in first address", func(t *testing.T) {
		inv := standardInvoice()
		inv.Supplier.Addresses = []*org.Address{
			{
				Locality: "Warsaw",
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "country: cannot be blank")
		assert.ErrorContains(t, err, "street: cannot be blank")
	})
}

func TestCustomerJSTValidation(t *testing.T) {
	t.Run("valid JST customer with LGU recipient identity", func(t *testing.T) {
		inv := standardInvoice()
		inv.Customer.Ext = tax.Extensions{
			favat.ExtKeyJST: "1", // Customer is a Subordinate Local Government Unit
		}
		inv.Customer.Identities = []*org.Identity{
			{
				Code: "JST-12345",
				Ext: tax.Extensions{
					favat.ExtKeyThirdPartyRole: "8", // Local Government Unit (LGU) - recipient
				},
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("JST customer without required identity", func(t *testing.T) {
		inv := standardInvoice()
		inv.Customer.Ext = tax.Extensions{
			favat.ExtKeyJST: "1", // Customer is a Subordinate Local Government Unit
		}
		// No identities provided
		inv.Customer.Identities = nil
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "customer: (identities: missing identity with role '8' and code.)")
	})

	t.Run("JST customer with identity missing code", func(t *testing.T) {
		inv := standardInvoice()
		inv.Customer.Ext = tax.Extensions{
			favat.ExtKeyJST: "1", // Customer is a Subordinate Local Government Unit
		}
		inv.Customer.Identities = []*org.Identity{
			{
				Code: "", // Empty code
				Ext: tax.Extensions{
					favat.ExtKeyThirdPartyRole: "8", // Local Government Unit (LGU) - recipient
				},
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		// Identity's own validation catches empty code first
		assert.ErrorContains(t, err, "code: cannot be blank")
	})

	t.Run("JST customer with wrong role identity", func(t *testing.T) {
		inv := standardInvoice()
		inv.Customer.Ext = tax.Extensions{
			favat.ExtKeyJST: "1", // Customer is a Subordinate Local Government Unit
		}
		inv.Customer.Identities = []*org.Identity{
			{
				Code: "SOME-ID",
				Ext: tax.Extensions{
					favat.ExtKeyThirdPartyRole: "10", // Wrong role (GV member instead of LGU)
				},
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "customer: (identities: missing identity with role '8' and code.)")
	})

	t.Run("non-JST customer does not require identity", func(t *testing.T) {
		inv := standardInvoice()
		inv.Customer.Ext = tax.Extensions{
			favat.ExtKeyJST: "2", // Customer is NOT a Subordinate Local Government Unit
		}
		// No identities needed
		inv.Customer.Identities = nil
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
	})
}

func TestCustomerGroupVATValidation(t *testing.T) {
	t.Run("valid GroupVAT customer with GV member identity", func(t *testing.T) {
		inv := standardInvoice()
		inv.Customer.Ext = tax.Extensions{
			favat.ExtKeyGroupVAT: "1", // Customer is a Group VAT member
		}
		inv.Customer.Identities = []*org.Identity{
			{
				Code: "GV-67890",
				Ext: tax.Extensions{
					favat.ExtKeyThirdPartyRole: "10", // GV member - recipient
				},
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("GroupVAT customer without required identity", func(t *testing.T) {
		inv := standardInvoice()
		inv.Customer.Ext = tax.Extensions{
			favat.ExtKeyGroupVAT: "1", // Customer is a Group VAT member
		}
		// No identities provided
		inv.Customer.Identities = nil
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "customer: (identities: missing identity with role '10' and code.)")
	})

	t.Run("GroupVAT customer with identity missing code", func(t *testing.T) {
		inv := standardInvoice()
		inv.Customer.Ext = tax.Extensions{
			favat.ExtKeyGroupVAT: "1", // Customer is a Group VAT member
		}
		inv.Customer.Identities = []*org.Identity{
			{
				Code: "", // Empty code
				Ext: tax.Extensions{
					favat.ExtKeyThirdPartyRole: "10", // GV member - recipient
				},
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		// Identity's own validation catches empty code first
		assert.ErrorContains(t, err, "code: cannot be blank")
	})

	t.Run("GroupVAT customer with wrong role identity", func(t *testing.T) {
		inv := standardInvoice()
		inv.Customer.Ext = tax.Extensions{
			favat.ExtKeyGroupVAT: "1", // Customer is a Group VAT member
		}
		inv.Customer.Identities = []*org.Identity{
			{
				Code: "SOME-ID",
				Ext: tax.Extensions{
					favat.ExtKeyThirdPartyRole: "8", // Wrong role (LGU instead of GV member)
				},
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "customer: (identities: missing identity with role '10' and code.)")
	})

	t.Run("non-GroupVAT customer does not require identity", func(t *testing.T) {
		inv := standardInvoice()
		inv.Customer.Ext = tax.Extensions{
			favat.ExtKeyGroupVAT: "2", // Customer is NOT a Group VAT member
		}
		// No identities needed
		inv.Customer.Identities = nil
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
	})
}

func TestCustomerJSTAndGroupVATCombined(t *testing.T) {
	t.Run("customer with both JST and GroupVAT needs both identities", func(t *testing.T) {
		inv := standardInvoice()
		inv.Customer.Ext = tax.Extensions{
			favat.ExtKeyJST:      "1", // Customer is a Subordinate Local Government Unit
			favat.ExtKeyGroupVAT: "1", // Customer is also a Group VAT member
		}
		inv.Customer.Identities = []*org.Identity{
			{
				Code: "JST-12345",
				Ext: tax.Extensions{
					favat.ExtKeyThirdPartyRole: "8", // Local Government Unit (LGU) - recipient
				},
			},
			{
				Code: "GV-67890",
				Ext: tax.Extensions{
					favat.ExtKeyThirdPartyRole: "10", // GV member - recipient
				},
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("customer with both JST and GroupVAT missing JST identity", func(t *testing.T) {
		inv := standardInvoice()
		inv.Customer.Ext = tax.Extensions{
			favat.ExtKeyJST:      "1", // Customer is a Subordinate Local Government Unit
			favat.ExtKeyGroupVAT: "1", // Customer is also a Group VAT member
		}
		inv.Customer.Identities = []*org.Identity{
			{
				Code: "GV-67890",
				Ext: tax.Extensions{
					favat.ExtKeyThirdPartyRole: "10", // Only GV member identity
				},
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "customer: (identities: missing identity with role '8' and code.)")
	})

	t.Run("customer with both JST and GroupVAT missing GroupVAT identity", func(t *testing.T) {
		inv := standardInvoice()
		inv.Customer.Ext = tax.Extensions{
			favat.ExtKeyJST:      "1", // Customer is a Subordinate Local Government Unit
			favat.ExtKeyGroupVAT: "1", // Customer is also a Group VAT member
		}
		inv.Customer.Identities = []*org.Identity{
			{
				Code: "JST-12345",
				Ext: tax.Extensions{
					favat.ExtKeyThirdPartyRole: "8", // Only LGU identity
				},
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "customer: (identities: missing identity with role '10' and code.)")
	})

	t.Run("customer without JST and GroupVAT does not require identities", func(t *testing.T) {
		inv := standardInvoice()
		// No JST or GroupVAT extensions set
		inv.Customer.Ext = nil
		inv.Customer.Identities = nil
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
	})
}

func TestNormalizeInvoice(t *testing.T) {
	ad := tax.AddonForKey(favat.V3)

	t.Run("self-billed invoice sets extension", func(t *testing.T) {
		inv := standardInvoice()
		inv.SetTags(tax.TagSelfBilled)
		require.NoError(t, inv.Calculate())
		ad.Normalizer(inv)
		assert.Equal(t, "1", inv.Tax.Ext.Get(favat.ExtKeySelfBilling).String())
	})

	t.Run("reverse charge invoice sets extension", func(t *testing.T) {
		inv := standardInvoice()
		inv.SetTags(tax.TagReverseCharge)
		require.NoError(t, inv.Calculate())
		ad.Normalizer(inv)
		assert.Equal(t, "1", inv.Tax.Ext.Get(favat.ExtKeyReverseCharge).String())
	})

	t.Run("both self-billed and reverse charge", func(t *testing.T) {
		inv := standardInvoice()
		inv.SetTags(tax.TagSelfBilled, tax.TagReverseCharge)
		require.NoError(t, inv.Calculate())
		ad.Normalizer(inv)
		assert.Equal(t, "1", inv.Tax.Ext.Get(favat.ExtKeySelfBilling).String())
		assert.Equal(t, "1", inv.Tax.Ext.Get(favat.ExtKeyReverseCharge).String())
	})

	t.Run("regular invoice does not set extensions", func(t *testing.T) {
		inv := standardInvoice()
		require.NoError(t, inv.Calculate())
		ad.Normalizer(inv)
		assert.Equal(t, "", inv.Tax.Ext.Get(favat.ExtKeySelfBilling).String())
		assert.Equal(t, "", inv.Tax.Ext.Get(favat.ExtKeyReverseCharge).String())
	})
}

func TestCreditNoteValidation(t *testing.T) {
	t.Run("credit note without preceding", func(t *testing.T) {
		inv := creditNote()
		inv.Preceding = nil
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})
}

func TestSimplifiedInvoiceCustomerValidation(t *testing.T) {
	t.Run("simplified invoice without customer is valid", func(t *testing.T) {
		inv := standardInvoice()
		inv.SetTags(tax.TagSimplified)
		inv.Customer = nil
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("standard invoice without customer is invalid", func(t *testing.T) {
		inv := standardInvoice()
		inv.Customer = nil
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "customer: cannot be blank")
	})
}

func TestCustomerTaxIDValidation(t *testing.T) {
	t.Run("customer without tax ID", func(t *testing.T) {
		inv := standardInvoice()
		inv.Customer.TaxID = nil
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "customer: (tax_id: cannot be blank.)")
	})
}

func TestPrecedingValidation(t *testing.T) {
	t.Run("preceding without issue date", func(t *testing.T) {
		inv := creditNote()
		inv.Preceding[0].IssueDate = nil
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "preceding: (0: (issue_date: cannot be blank.).")
	})

	t.Run("preceding without code", func(t *testing.T) {
		inv := creditNote()
		inv.Preceding[0].Code = ""
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "preceding: (0: (code: cannot be blank.).")
	})
}

func TestNilValidation(t *testing.T) {
	ad := tax.AddonForKey(favat.V3)

	t.Run("nil supplier", func(t *testing.T) {
		inv := standardInvoice()
		inv.Supplier = nil
		require.NoError(t, inv.Calculate())
		err := ad.Validator(inv)
		assert.ErrorContains(t, err, "supplier: cannot be blank")
	})
}

func TestValidationEdgeCases(t *testing.T) {
	t.Run("customer with identity having different role", func(t *testing.T) {
		inv := standardInvoice()
		inv.Customer.Ext = tax.Extensions{
			favat.ExtKeyJST: "1",
		}
		inv.Customer.Identities = []*org.Identity{
			{
				Code: "TEST-123",
				Ext: tax.Extensions{
					favat.ExtKeyThirdPartyRole: "5", // Different role
				},
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "missing identity with role '8' and code")
	})

	t.Run("customer with identity matching role but no code", func(t *testing.T) {
		inv := standardInvoice()
		inv.Customer.Ext = tax.Extensions{
			favat.ExtKeyGroupVAT: "1",
		}
		inv.Customer.Identities = []*org.Identity{
			{
				Code: "", // Empty code
				Ext: tax.Extensions{
					favat.ExtKeyThirdPartyRole: "10",
				},
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		// Standard validation catches empty code
		assert.Error(t, err)
	})

	t.Run("exemption note without exemption extension is valid", func(t *testing.T) {
		inv := standardInvoice()
		inv.Notes = []*org.Note{
			{
				Key:  org.NoteKeyLegal,
				Code: "A",
				Src:  favat.ExtKeyExemption,
				Text: "Some note",
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		// No exemption extension set, so note is allowed but not required
		assert.NoError(t, err)
	})
}

func TestInvoiceWithNotes(t *testing.T) {
	t.Run("invoice with regular note (non-exemption)", func(t *testing.T) {
		inv := standardInvoice()
		inv.Notes = []*org.Note{
			{
				Key:  org.NoteKeyGeneral,
				Text: "Regular note",
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("invoice with legal note but different src", func(t *testing.T) {
		inv := standardInvoice()
		inv.Notes = []*org.Note{
			{
				Key:  org.NoteKeyLegal,
				Src:  "other-source",
				Text: "Legal note with different source",
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
	})
}

package oioubl_test

import (
	"testing"

	_ "github.com/invopop/gobl"
	oioubl "github.com/invopop/gobl/addons/dk/oioubl-v2-1"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testStatusResponse(t *testing.T) *bill.Status {
	t.Helper()
	return &bill.Status{
		Regime:    tax.WithRegime("DK"),
		Addons:    tax.WithAddons(oioubl.V2_1),
		Type:      bill.StatusTypeResponse,
		Code:      "RESP001",
		IssueDate: cal.MakeDate(2026, 1, 15),
		Supplier: &org.Party{
			Name:    "Kunde ApS",
			TaxID:   &tax.Identity{Country: "DK", Code: "88146328"},
			Inboxes: []*org.Inbox{{Scheme: "0184", Code: "88146328"}},
		},
		Customer: &org.Party{
			Name:    "Eksempel A/S",
			TaxID:   &tax.Identity{Country: "DK", Code: "12345674"},
			Inboxes: []*org.Inbox{{Scheme: "0184", Code: "12345674"}},
		},
		Lines: []*bill.StatusLine{
			{
				Key: bill.StatusEventRejected,
				Doc: &org.DocumentRef{Code: "INV1000"},
			},
		},
	}
}

func TestStatusValidation(t *testing.T) {
	t.Run("standard response", func(t *testing.T) {
		st := testStatusResponse(t)
		require.NoError(t, st.Calculate())
		require.NoError(t, rules.Validate(st))
	})

	t.Run("missing code (F-APR005)", func(t *testing.T) {
		st := testStatusResponse(t)
		st.Code = ""
		require.NoError(t, st.Calculate())
		err := rules.Validate(st)
		assert.ErrorContains(t, err, "F-APR005")
	})

	t.Run("missing supplier inboxes (F-APR012)", func(t *testing.T) {
		st := testStatusResponse(t)
		st.Supplier.Inboxes = nil
		require.NoError(t, st.Calculate())
		err := rules.Validate(st)
		assert.ErrorContains(t, err, "F-APR012")
	})

	t.Run("supplier without tax ID or identities (F-APR041)", func(t *testing.T) {
		st := testStatusResponse(t)
		st.Supplier.TaxID = nil
		st.Supplier.Identities = nil
		require.NoError(t, st.Calculate())
		err := rules.Validate(st)
		assert.ErrorContains(t, err, "F-APR041")
	})

	t.Run("supplier with identities and no tax ID passes", func(t *testing.T) {
		st := testStatusResponse(t)
		st.Supplier.TaxID = nil
		st.Supplier.Identities = []*org.Identity{{Key: "dk-cvr", Code: "88146328"}}
		require.NoError(t, st.Calculate())
		assert.NoError(t, rules.Validate(st))
	})

	t.Run("supplier without name or identities (F-LIB022)", func(t *testing.T) {
		st := testStatusResponse(t)
		st.Supplier.Name = ""
		st.Supplier.Identities = nil
		require.NoError(t, st.Calculate())
		err := rules.Validate(st)
		assert.ErrorContains(t, err, "F-LIB022")
	})

	t.Run("name-less supplier with a non-legal identity fails (F-LIB022)", func(t *testing.T) {
		// A tax-scope-only identity yields no PartyLegalEntity, so it cannot
		// produce valid OIOUBL and must be rejected.
		st := testStatusResponse(t)
		st.Supplier.Name = ""
		st.Supplier.Identities = []*org.Identity{{Scope: org.IdentityScopeTax, Code: "88146328"}}
		require.NoError(t, st.Calculate())
		err := rules.Validate(st)
		assert.ErrorContains(t, err, "F-LIB022")
	})

	t.Run("name-less supplier with a legal identity passes", func(t *testing.T) {
		st := testStatusResponse(t)
		st.Supplier.Name = ""
		st.Supplier.Identities = []*org.Identity{{Scope: org.IdentityScopeLegal, Code: "88146328"}}
		require.NoError(t, st.Calculate())
		assert.NoError(t, rules.Validate(st))
	})

	t.Run("missing customer", func(t *testing.T) {
		st := testStatusResponse(t)
		st.Customer = nil
		require.NoError(t, st.Calculate())
		err := rules.Validate(st)
		assert.ErrorContains(t, err, "customer is required")
	})

	t.Run("missing customer inboxes (F-APR008)", func(t *testing.T) {
		st := testStatusResponse(t)
		st.Customer.Inboxes = nil
		require.NoError(t, st.Calculate())
		err := rules.Validate(st)
		assert.ErrorContains(t, err, "F-APR008")
	})

	t.Run("customer without name or identities (F-LIB022)", func(t *testing.T) {
		st := testStatusResponse(t)
		st.Customer.Name = ""
		st.Customer.Identities = nil
		require.NoError(t, st.Calculate())
		err := rules.Validate(st)
		assert.ErrorContains(t, err, "F-LIB022")
	})

	t.Run("issuer absent is allowed", func(t *testing.T) {
		st := testStatusResponse(t)
		require.NoError(t, st.Calculate())
		assert.NoError(t, rules.Validate(st))
	})

	t.Run("issuer set without inboxes (F-APR008)", func(t *testing.T) {
		st := testStatusResponse(t)
		st.Issuer = &org.Party{
			Name:  "Invopop",
			TaxID: &tax.Identity{Country: "DK", Code: "12345674"},
		}
		require.NoError(t, st.Calculate())
		err := rules.Validate(st)
		assert.ErrorContains(t, err, "F-APR008")
	})

	t.Run("issuer set without tax ID or identities (F-APR040)", func(t *testing.T) {
		st := testStatusResponse(t)
		st.Issuer = &org.Party{
			Name:    "Invopop",
			Inboxes: []*org.Inbox{{Scheme: "0184", Code: "12345674"}},
		}
		require.NoError(t, st.Calculate())
		err := rules.Validate(st)
		assert.ErrorContains(t, err, "F-APR040")
	})

	t.Run("issuer set without name or identities (F-LIB022)", func(t *testing.T) {
		st := testStatusResponse(t)
		st.Issuer = &org.Party{
			TaxID:   &tax.Identity{Country: "DK", Code: "12345674"},
			Inboxes: []*org.Inbox{{Scheme: "0184", Code: "12345674"}},
		}
		require.NoError(t, st.Calculate())
		err := rules.Validate(st)
		assert.ErrorContains(t, err, "F-LIB022")
	})

	t.Run("issuer fully populated passes", func(t *testing.T) {
		st := testStatusResponse(t)
		st.Issuer = &org.Party{
			Name:    "Invopop",
			TaxID:   &tax.Identity{Country: "DK", Code: "12345674"},
			Inboxes: []*org.Inbox{{Scheme: "0184", Code: "12345674"}},
		}
		require.NoError(t, st.Calculate())
		assert.NoError(t, rules.Validate(st))
	})

	t.Run("recipient absent is allowed", func(t *testing.T) {
		st := testStatusResponse(t)
		require.NoError(t, st.Calculate())
		assert.NoError(t, rules.Validate(st))
	})

	t.Run("recipient set without inboxes (F-APR012)", func(t *testing.T) {
		st := testStatusResponse(t)
		st.Recipient = &org.Party{
			Name:  "Forsendelses Hub A/S",
			TaxID: &tax.Identity{Country: "DK", Code: "12345674"},
		}
		require.NoError(t, st.Calculate())
		err := rules.Validate(st)
		assert.ErrorContains(t, err, "F-APR012")
	})

	t.Run("recipient fully populated passes", func(t *testing.T) {
		st := testStatusResponse(t)
		st.Recipient = &org.Party{
			Name:    "Forsendelses Hub A/S",
			TaxID:   &tax.Identity{Country: "DK", Code: "12345674"},
			Inboxes: []*org.Inbox{{Scheme: "0184", Code: "12345674"}},
		}
		require.NoError(t, st.Calculate())
		assert.NoError(t, rules.Validate(st))
	})

	t.Run("missing line doc", func(t *testing.T) {
		st := testStatusResponse(t)
		st.Lines[0].Doc = nil
		require.NoError(t, st.Calculate())
		err := rules.Validate(st)
		assert.ErrorContains(t, err, "line document reference is required")
	})

	t.Run("non-response status skips F-APR rules", func(t *testing.T) {
		st := testStatusResponse(t)
		st.Type = bill.StatusTypeSystem
		st.Supplier.Inboxes = nil
		st.Customer = nil
		st.Lines[0].Doc = nil
		require.NoError(t, st.Calculate())
		assert.NoError(t, rules.Validate(st))
	})
}

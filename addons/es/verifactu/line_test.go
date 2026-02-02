package verifactu_test

import (
	"testing"

	_ "github.com/invopop/gobl"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/invopop/gobl/addons/es/verifactu"
)

func TestLineNormalization(t *testing.T) {
	addon := tax.AddonForKey(verifactu.V1)
	require.NotNil(t, addon)

	t.Run("nil line", func(t *testing.T) {
		assert.NotPanics(t, func() {
			var line *bill.Line
			addon.Normalizer(line)
		})
	})

	t.Run("line with no taxes", func(t *testing.T) {
		line := new(bill.Line)
		addon.Normalizer(line)
		assert.Nil(t, line.Notes)
	})

	t.Run("line with nil tax", func(t *testing.T) {
		line := &bill.Line{
			Taxes: tax.Set{nil},
		}
		addon.Normalizer(line)
		assert.Nil(t, line.Notes)
	})

	t.Run("line with taxes but no exemption", func(t *testing.T) {
		line := &bill.Line{
			Taxes: tax.Set{
				{
					Category: tax.CategoryVAT,
					Rate:     tax.RateGeneral,
					Percent:  num.NewPercentage(210, 3),
				},
			},
		}

		addon.Normalizer(line)
		assert.Nil(t, line.Notes)
	})

	t.Run("line with exemption extension E1", func(t *testing.T) {
		line := &bill.Line{
			Taxes: tax.Set{
				{
					Category: tax.CategoryVAT,
					Ext: tax.Extensions{
						verifactu.ExtKeyExempt: "E1",
					},
				},
			},
		}

		addon.Normalizer(line)
		require.Len(t, line.Notes, 1)

		note := line.Notes[0]
		assert.Equal(t, org.NoteKeyLegal, note.Key)
		assert.Equal(t, "E1", note.Code.String())
		assert.Equal(t, verifactu.ExtKeyExempt, note.Src)
		assert.Equal(t, "Exenta: por el artículo 20. Exenciones en operaciones interiores.", note.Text)
	})

	t.Run("line with exemption extension E5", func(t *testing.T) {
		line := &bill.Line{
			Taxes: tax.Set{
				{
					Category: tax.CategoryVAT,
					Ext: tax.Extensions{
						verifactu.ExtKeyExempt: "E5",
					},
				},
			},
		}

		addon.Normalizer(line)
		require.Len(t, line.Notes, 1)

		note := line.Notes[0]
		assert.Equal(t, org.NoteKeyLegal, note.Key)
		assert.Equal(t, "E5", note.Code.String())
		assert.Equal(t, verifactu.ExtKeyExempt, note.Src)
		assert.Equal(t, "Exenta: por el artículo 25. Exenciones en las entregas de bienes destinados a otro Estado miembro.", note.Text)
	})

	t.Run("line with existing notes", func(t *testing.T) {
		line := &bill.Line{
			Notes: []*org.Note{
				{
					Key:  org.NoteKeyGeneral,
					Text: "Existing note",
				},
			},
			Taxes: tax.Set{
				{
					Category: tax.CategoryVAT,
					Ext: tax.Extensions{
						verifactu.ExtKeyExempt: "E1",
					},
				},
			},
		}

		addon.Normalizer(line)
		require.Len(t, line.Notes, 2)

		// The existing note is preserved
		assert.Equal(t, "Existing note", line.Notes[0].Text)

		// The new exemption note is added
		exemptionNote := line.Notes[1]
		assert.Equal(t, org.NoteKeyLegal, exemptionNote.Key)
		assert.Equal(t, "E1", exemptionNote.Code.String())
		assert.Equal(t, verifactu.ExtKeyExempt, exemptionNote.Src)
		assert.Equal(t, "Exenta: por el artículo 20. Exenciones en operaciones interiores.", exemptionNote.Text)
	})

	t.Run("duplicate exemption prevention", func(t *testing.T) {
		line := &bill.Line{
			Notes: []*org.Note{
				{
					Key:  org.NoteKeyLegal,
					Code: "E2", // Code doesn't need to match. Validation will check this.
					Src:  verifactu.ExtKeyExempt,
					Text: "Exenta: por el artículo 21. Exenciones en las exportaciones de bienes.",
				},
			},
			Taxes: tax.Set{
				{
					Category: tax.CategoryVAT,
					Ext: tax.Extensions{
						verifactu.ExtKeyExempt: "E1",
					},
				},
			},
		}

		addon.Normalizer(line)
		require.Len(t, line.Notes, 1) // Should not add duplicate

		note := line.Notes[0]
		assert.Equal(t, "E2", note.Code.String())
		assert.Equal(t, "Exenta: por el artículo 21. Exenciones en las exportaciones de bienes.", note.Text)
	})

	t.Run("invalid exemption code", func(t *testing.T) {
		line := &bill.Line{
			Taxes: tax.Set{
				{
					Category: tax.CategoryVAT,
					Ext: tax.Extensions{
						verifactu.ExtKeyExempt: "INVALID",
					},
				},
			},
		}

		addon.Normalizer(line)
		assert.Nil(t, line.Notes) // Should not add note for invalid code
	})
}

func TestLineValidation(t *testing.T) {
	addon := tax.AddonForKey(verifactu.V1)
	require.NotNil(t, addon)

	t.Run("nil line", func(t *testing.T) {
		var line *bill.Line
		assert.NoError(t, addon.Validator(line))
	})

	t.Run("valid line", func(t *testing.T) {
		line := validLine()
		assert.NoError(t, addon.Validator(line))
	})

	t.Run("line with valid exemption note", func(t *testing.T) {
		line := validLine()
		line.Taxes = tax.Set{
			{
				Category: tax.CategoryVAT,
				Ext: tax.Extensions{
					verifactu.ExtKeyExempt: "E1",
				},
			},
		}
		line.Notes = []*org.Note{
			{
				Key:  org.NoteKeyLegal,
				Src:  verifactu.ExtKeyExempt,
				Code: "E1",
				Text: "Exenta: por el artículo 20. Exenciones en operaciones interiores.",
			},
		}
		assert.NoError(t, addon.Validator(line))
	})

	t.Run("line missing exemption note", func(t *testing.T) {
		line := validLine()
		line.Taxes = tax.Set{
			{
				Category: tax.CategoryVAT,
				Ext: tax.Extensions{
					verifactu.ExtKeyExempt: "E1",
				},
			},
		}
		err := addon.Validator(line)
		assert.ErrorContains(t, err, "notes: missing exemption note for code E1")
	})

	t.Run("line with unexpected exemption note", func(t *testing.T) {
		line := validLine()
		line.Notes = []*org.Note{
			{
				Key:  org.NoteKeyLegal,
				Src:  verifactu.ExtKeyExempt,
				Code: "E1",
				Text: "Exenta: por el artículo 20. Exenciones en operaciones interiores.",
			},
		}
		err := addon.Validator(line)
		assert.ErrorContains(t, err, "notes: (0: unexpected exemption note")
	})

	t.Run("line with mismatched exemption note code", func(t *testing.T) {
		line := validLine()
		line.Taxes = tax.Set{
			{
				Category: tax.CategoryVAT,
				Ext: tax.Extensions{
					verifactu.ExtKeyExempt: "E1",
				},
			},
		}
		line.Notes = []*org.Note{
			{
				Key:  org.NoteKeyLegal,
				Src:  verifactu.ExtKeyExempt,
				Code: "E2", // Different code than extension
				Text: "Exenta: por el artículo 21. Exenciones en las exportaciones de bienes.",
			},
		}
		err := addon.Validator(line)
		assert.ErrorContains(t, err, "notes: (0: note code E2 must match extension E1)")
	})

	t.Run("line with too many exemption notes", func(t *testing.T) {
		line := validLine()
		line.Taxes = tax.Set{
			{
				Category: tax.CategoryVAT,
				Ext: tax.Extensions{
					verifactu.ExtKeyExempt: "E1",
				},
			},
		}
		line.Notes = []*org.Note{
			{
				Key:  org.NoteKeyLegal,
				Src:  verifactu.ExtKeyExempt,
				Code: "E1",
				Text: "Exenta: por el artículo 20. Exenciones en operaciones interiores.",
			},
			{
				Key:  org.NoteKeyLegal,
				Src:  verifactu.ExtKeyExempt,
				Code: "E1",
				Text: "Duplicate exemption note",
			},
		}
		err := addon.Validator(line)
		assert.ErrorContains(t, err, "notes: (1: too many exemption notes)")
	})
}

func validLine() *bill.Line {
	return &bill.Line{
		Quantity: num.MakeAmount(1, 0),
		Taxes: tax.Set{
			{
				Category: tax.CategoryVAT,
				Rate:     tax.RateGeneral,
				Percent:  num.NewPercentage(210, 3),
			},
		},
	}
}

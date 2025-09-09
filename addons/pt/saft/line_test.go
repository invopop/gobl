package saft_test

import (
	"testing"

	_ "github.com/invopop/gobl"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/invopop/gobl/addons/pt/saft"
)

func TestLineNormalization(t *testing.T) {
	addon := tax.AddonForKey(saft.V1)
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

	t.Run("line with exemption extension", func(t *testing.T) {
		line := &bill.Line{
			Taxes: tax.Set{
				{
					Category: tax.CategoryVAT,
					Ext: tax.Extensions{
						saft.ExtKeyExemption: "M04",
					},
				},
			},
		}

		addon.Normalizer(line)
		require.Len(t, line.Notes, 1)

		note := line.Notes[0]
		assert.Equal(t, org.NoteKeyLegal, note.Key)
		assert.Equal(t, "M04", note.Code.String())
		assert.Equal(t, saft.ExtKeyExemption, note.Src)
		assert.Equal(t, "Artigo 13.º do CIVA", note.Text)
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
						saft.ExtKeyExemption: "M04",
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
		assert.Equal(t, "M04", exemptionNote.Code.String())
		assert.Equal(t, saft.ExtKeyExemption, exemptionNote.Src)
		assert.Equal(t, "Artigo 13.º do CIVA", exemptionNote.Text)
	})

	t.Run("duplicate exemption prevention", func(t *testing.T) {
		line := &bill.Line{
			Notes: []*org.Note{
				{
					Key:  org.NoteKeyLegal,
					Code: "M03", // Code doesn't need to match. Validation will check this.
					Src:  saft.ExtKeyExemption,
					Text: "Artigo 13.º do CIVA",
				},
			},
			Taxes: tax.Set{
				{
					Category: tax.CategoryVAT,
					Ext: tax.Extensions{
						saft.ExtKeyExemption: "M04",
					},
				},
			},
		}

		addon.Normalizer(line)
		require.Len(t, line.Notes, 1) // Should not add duplicate

		note := line.Notes[0]
		assert.Equal(t, "M03", note.Code.String())
		assert.Equal(t, "Artigo 13.º do CIVA", note.Text)
	})

	t.Run("invalid exemption code", func(t *testing.T) {
		line := &bill.Line{
			Taxes: tax.Set{
				{
					Category: tax.CategoryVAT,
					Ext: tax.Extensions{
						saft.ExtKeyExemption: "INVALID",
					},
				},
			},
		}

		addon.Normalizer(line)
		assert.Nil(t, line.Notes) // Should not add note for invalid code
	})
}

func TestLineValidation(t *testing.T) {
	addon := tax.AddonForKey(saft.V1)
	require.NotNil(t, addon)

	t.Run("nil line", func(t *testing.T) {
		var line *bill.Line
		assert.NoError(t, addon.Validator(line))
	})

	t.Run("line with no notes", func(t *testing.T) {
		line := new(bill.Line)
		assert.NoError(t, addon.Validator(line))
	})

	t.Run("line with valid exemption note", func(t *testing.T) {
		line := &bill.Line{
			Taxes: tax.Set{
				{
					Category: tax.CategoryVAT,
					Ext: tax.Extensions{
						saft.ExtKeyExemption: "M04",
					},
				},
			},
			Notes: []*org.Note{
				{
					Key:  org.NoteKeyLegal,
					Src:  saft.ExtKeyExemption,
					Code: "M04",
					Text: "Artigo 13.º do CIVA",
				},
			},
		}
		assert.NoError(t, addon.Validator(line))
	})

	t.Run("line missing exemption note", func(t *testing.T) {
		line := &bill.Line{
			Taxes: tax.Set{
				{
					Category: tax.CategoryVAT,
					Ext: tax.Extensions{
						saft.ExtKeyExemption: "M04",
					},
				},
			},
		}
		err := addon.Validator(line)
		assert.ErrorContains(t, err, "notes: missing exemption note for code M04")
	})

	t.Run("line with unexpected exemption note", func(t *testing.T) {
		line := &bill.Line{
			Notes: []*org.Note{
				{
					Key:  org.NoteKeyLegal,
					Src:  saft.ExtKeyExemption,
					Code: "M04",
					Text: "Artigo 13.º do CIVA",
				},
			},
		}
		err := addon.Validator(line)
		assert.ErrorContains(t, err, "notes: (0: unexpected exemption note")
	})

	t.Run("line with mismatched exemption note code", func(t *testing.T) {
		line := &bill.Line{
			Taxes: tax.Set{
				{
					Category: tax.CategoryVAT,
					Ext: tax.Extensions{
						saft.ExtKeyExemption: "M04",
					},
				},
			},
			Notes: []*org.Note{
				{
					Key:  org.NoteKeyLegal,
					Src:  saft.ExtKeyExemption,
					Code: "M01", // Different code than extension
					Text: "Artigo 13.º do CIVA",
				},
			},
		}
		err := addon.Validator(line)
		assert.ErrorContains(t, err, "notes: (0: note code M01 must match extension M04)")
	})

	t.Run("line with too many exemption notes", func(t *testing.T) {
		line := &bill.Line{
			Taxes: tax.Set{
				{
					Category: tax.CategoryVAT,
					Ext: tax.Extensions{
						saft.ExtKeyExemption: "M04",
					},
				},
			},
			Notes: []*org.Note{
				{
					Key:  org.NoteKeyLegal,
					Src:  saft.ExtKeyExemption,
					Code: "M04",
					Text: "Artigo 13.º do CIVA",
				},
				{
					Key:  org.NoteKeyLegal,
					Src:  saft.ExtKeyExemption,
					Code: "M04",
					Text: "Duplicate exemption note",
				},
			},
		}
		err := addon.Validator(line)
		assert.ErrorContains(t, err, "notes: (1: too many exemption notes)")
	})
}

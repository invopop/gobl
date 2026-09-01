package legal_test

import (
	"testing"

	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/legal"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/schema"
	"github.com/invopop/gobl/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContractNormalization(t *testing.T) {
	contract := &legal.Contract{
		Title:       " Contract title \x00 ",
		SubTitle:    " Subtitle ",
		Description: " Description ",
		Chapters: []*legal.Chapter{
			nil,
			{
				Anchor: " chapter-one ",
				Title:  " Chapter one ",
				Sections: []*legal.Section{
					nil,
					{
						Anchor:  " section-one ",
						Title:   " Section one ",
						Content: " Some terms. ",
						Sections: []*legal.Section{
							{Content: " Nested terms. "},
						},
					},
				},
			},
			{
				Ref:   " #chapter-one ",
				Title: " Chapter two ",
			},
		},
	}

	norm.Normalize(contract)

	assert.Equal(t, "Contract title", contract.Title)
	assert.Equal(t, "Subtitle", contract.SubTitle)
	assert.Equal(t, "Description", contract.Description)
	require.Len(t, contract.Chapters, 2)
	assert.Equal(t, int32(1), contract.Chapters[0].Index)
	assert.Equal(t, int32(2), contract.Chapters[1].Index)
	assert.Equal(t, "#chapter-one", contract.Chapters[1].Ref)
	require.Len(t, contract.Chapters[0].Sections, 1)
	assert.Equal(t, int32(1), contract.Chapters[0].Sections[0].Index)
	assert.Equal(t, "Some terms.", contract.Chapters[0].Sections[0].Content)
	assert.Equal(t, int32(1), contract.Chapters[0].Sections[0].Sections[0].Index)
}

func TestContractCalculate(t *testing.T) {
	contract := validContract()
	require.NoError(t, contract.Calculate())
	assert.Equal(t, int32(1), contract.Chapters[0].Index)
	assert.Equal(t, int32(1), contract.Chapters[0].Sections[0].Index)
}

func TestContractValidation(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		contract := validContract()
		norm.Normalize(contract)
		assert.NoError(t, rules.Validate(contract))
		assert.NoError(t, rules.Validate(*contract))
	})

	t.Run("title required", func(t *testing.T) {
		contract := validContract()
		contract.Title = ""
		norm.Normalize(contract)
		faults := rules.Validate(contract)
		require.Error(t, faults)
		assert.True(t, faults.HasCode("GOBL-LEGAL-CONTRACT-03"))
		assert.True(t, faults.HasPath("$.title"))
	})

	t.Run("section body required", func(t *testing.T) {
		contract := validContract()
		contract.Chapters[0].Sections[0].Content = ""
		norm.Normalize(contract)
		faults := rules.Validate(contract)
		require.Error(t, faults)
		assert.True(t, faults.HasCode("GOBL-LEGAL-SECTION-02"))
		assert.True(t, faults.HasPath("$.chapters[0].sections[0]"))
	})

	t.Run("anchors unique", func(t *testing.T) {
		contract := validContract()
		contract.Chapters[0].Sections = append(contract.Chapters[0].Sections,
			&legal.Section{Anchor: "terms", Content: "More terms."},
		)
		norm.Normalize(contract)
		faults := rules.Validate(contract)
		require.Error(t, faults)
		assert.True(t, faults.HasCode("GOBL-LEGAL-CONTRACT-10"))
	})

	t.Run("local reference resolves", func(t *testing.T) {
		contract := validContract()
		contract.Chapters[0].Sections[0].Ref = "#missing"
		norm.Normalize(contract)
		faults := rules.Validate(contract)
		require.Error(t, faults)
		assert.True(t, faults.HasCode("GOBL-LEGAL-CONTRACT-11"))
	})

	t.Run("anchor format", func(t *testing.T) {
		contract := validContract()
		contract.Chapters[0].Anchor = "not an anchor"
		norm.Normalize(contract)
		faults := rules.Validate(contract)
		require.Error(t, faults)
		assert.True(t, faults.HasCode("GOBL-LEGAL-CHAPTER-04"))
		assert.True(t, faults.HasPath("$.chapters[0].$anchor"))
	})

	t.Run("reference format", func(t *testing.T) {
		contract := validContract()
		contract.Chapters[0].Ref = "not a reference"
		norm.Normalize(contract)
		faults := rules.Validate(contract)
		require.Error(t, faults)
		assert.True(t, faults.HasCode("GOBL-LEGAL-CHAPTER-05"))
	})
}

func TestContractSchemaRegistration(t *testing.T) {
	contract := validContract()
	assert.Equal(t, schema.GOBL.Add(legal.ShortSchemaContract), schema.Lookup(contract))

	obj, err := schema.NewObject(contract)
	require.NoError(t, err)
	require.NoError(t, obj.Calculate())
	assert.False(t, contract.UUID.IsZero())
	assert.Equal(t, int32(1), contract.Chapters[0].Index)
}

func validContract() *legal.Contract {
	id := uuid.V7()
	return &legal.Contract{
		Identify:  uuid.Identify{UUID: id},
		Agreement: id,
		Language:  i18n.EN,
		Title:     "Service agreement",
		Parties: []*legal.Party{
			{
				Key:       "provider",
				Role:      "service-provider",
				DefinedAs: "Provider",
				Entity:    &org.Party{Name: "Example Provider Ltd"},
			},
		},
		Chapters: []*legal.Chapter{
			{
				Anchor: "scope",
				Title:  "Scope",
				Sections: []*legal.Section{
					{
						Anchor:  "terms",
						Title:   "Terms",
						Content: "The parties agree to the following terms.",
					},
				},
			},
		},
	}
}
